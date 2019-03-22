package download

import (
	"archive/tar"
	"bufio"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	gtio "github.com/evolbioinfo/gotree/io"
	"github.com/evolbioinfo/gotree/io/fileutils"
	"github.com/evolbioinfo/gotree/tree"
	"github.com/jlaffaye/ftp"
)

type NcbiTreeDownloader struct {
	server      string
	path        string
	tipstaxids  bool // If taxids only must be printed as tips of the output tree
	nodestaxids bool // If taxids only must be printed as internal nodes of the output tree
	mapfile     string
}

/* NCBI taxonomy downloader */
func NewNcbiTreeDownloader() *NcbiTreeDownloader {
	return &NcbiTreeDownloader{"ftp.ncbi.nih.gov:21", "/pub/taxonomy/taxdump.tar.gz", false, false, ""}
}

// If taxids only must be printed as internal nodes of the output tree
// Otherwise taxa names
func (d *NcbiTreeDownloader) SetInternalNodesTaxId(val bool) {
	d.nodestaxids = val
}

// If taxids only must be printed as tips of the output tree
// Otherwise taxa names
func (d *NcbiTreeDownloader) SetTipsTaxId(val bool) {
	d.tipstaxids = val
}

// If called, the downloader will also write mapping between
// ncbi taxid and tax name in the given output file
func (d *NcbiTreeDownloader) SetMapFileOutputPath(output string) {
	d.mapfile = output
}

// Download the NCBI taxonomy as a tree.Tree
func (d *NcbiTreeDownloader) Download(id string) (*tree.Tree, error) {
	var client *ftp.ServerConn
	var err error
	var reader *ftp.Response
	var gzreader *gzip.Reader
	var tarreader *tar.Reader
	var t *tree.Tree              // tree structure of the ncbi taxo
	var namemap map[string]string // map between node ids and node names

	// Connect to NCBI FTP Server
	client, err = ftp.Dial(d.server)
	defer client.Quit()
	if err != nil {
		return nil, err
	}
	if err = client.Login("anonymous", "anonymous@domain.com"); err != nil {
		return nil, err
	}
	// Retrieve the "taxdump.tar.gz" file
	gtio.LogInfo("Downloading from NCBI ftp")
	reader, err = client.Retr(d.path)
	if err != nil {
		return nil, err
	}

	// Reading tar gz and processing nodes.dmp and names.dmp
	gtio.LogInfo("Extracting files from archive")
	if gzreader, err = gzip.NewReader(reader); err != nil {
		return nil, err
	}
	tarreader = tar.NewReader(gzreader)

	for {
		header, err := tarreader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		switch header.Typeflag {
		case tar.TypeDir:
			continue
		case tar.TypeReg:
			// We handle names of ncbi taxonomy nodes
			if header.Name == "names.dmp" {
				gtio.LogInfo("Parsing name file")
				namemap, err = d.parseNcbiNames(tarreader)
				if err != nil {
					return nil, err
				}
			}
			// We handle the tree
			if header.Name == "nodes.dmp" {
				gtio.LogInfo("Parsing node file")
				t, err = d.parseNcbiTree(tarreader)
				if err != nil {
					return nil, err
				}
			}
		default:
			return nil, errors.New("Problem with tar archive")
		}
	}
	gtio.LogInfo("Removing single nodes")
	d.addSpeciesTips(t)
	t.RemoveSingleNodes()
	gtio.LogInfo("Renaming taxid -> taxnames")
	d.renameTreeNodes(t, namemap, d.tipstaxids, d.nodestaxids)
	err = d.writeMapfile(namemap)
	return t, err
}

func (d *NcbiTreeDownloader) renameTreeNodes(t *tree.Tree, namemap map[string]string, tipstaxids bool, nodestaxids bool) {
	for _, n := range t.Nodes() {
		if (n.Tip() && !tipstaxids) || (!n.Tip() && !nodestaxids) {
			if name, ok := namemap[n.Name()]; ok {
				n.SetName(name)
			}
		}
	}
}

/*
Parse name file and convert names with the following rules:
Special characters are replaces with "_" ->  '(', ')', '[', ']', ':', ',', ' ', ';'
*/
func (d *NcbiTreeDownloader) parseNcbiNames(reader io.Reader) (map[string]string, error) {
	r := bufio.NewReader(reader)
	l, err := fileutils.Readln(r)
	namemap := make(map[string]string)
	sciname := make(map[string]bool)

	re := regexp.MustCompile("[\\[\\]\\(\\)\\:\\,\\s\\;]")
	for err == nil {
		cols := regexp.MustCompile("\t*\\|\t*").Split(l, -1)
		tax := cols[0]
		name := cols[1]
		tpe := cols[3]
		clean := re.ReplaceAllString(name, "_")
		if tpe == "scientific name" {
			namemap[tax] = clean
			sciname[tax] = true
		} else if tpe == "synonym" {
			_, ok := sciname[tax]
			if !ok {
				namemap[tax] = clean
			}
		}
		l, err = fileutils.Readln(r)
	}
	return namemap, nil
}

// Build a gotree.tree.Tree
func (d *NcbiTreeDownloader) parseNcbiTree(reader io.Reader) (*tree.Tree, error) {
	r := bufio.NewReader(reader)
	l, err := fileutils.Readln(r)
	t := tree.NewTree()
	var root *tree.Node
	nodes := make(map[string]*tree.Node)
	for err == nil {
		cols := strings.Split(l, "\t|\t")
		tax := cols[0]
		parent := cols[1]
		rank := cols[2]
		n1, ok1 := nodes[tax]
		if !ok1 {
			n1 = t.NewNode()
			n1.SetName(tax)
			nodes[tax] = n1
		}
		// We add the rank in order to be able
		// resolve cases were a species have also children
		n1.AddComment(rank)
		n2, ok2 := nodes[parent]
		if !ok2 {
			n2 = t.NewNode()
			n2.SetName(parent)
			nodes[parent] = n2
		}
		if tax == parent {
			root = n1
			t.SetRoot(root)
			n1.SetName("")
		} else {
			t.ConnectNodes(n2, n1)
		}

		l, err = fileutils.Readln(r)
	}
	if root == nil {
		return nil, errors.New("No root found in the NCBI Taxonomy")
	}
	return t, nil
}

// if an internal node is a species, then we add a new tip
func (d *NcbiTreeDownloader) addSpeciesTips(t *tree.Tree) {
	for _, n := range t.Nodes() {
		if len(n.Neigh()) > 1 &&
			len(n.Comments()) > 0 &&
			n.Comments()[0] == "species" {
			tip := t.NewNode()
			tip.SetName(n.Name())
			tip.AddComment(n.Comments()[0])
			t.ConnectNodes(n, tip)
		}
	}
}

// if an internal node is a species, then we add a new tip
func (d *NcbiTreeDownloader) writeMapfile(namemap map[string]string) error {
	if d.mapfile != "" {
		f, err := os.Create(d.mapfile)
		if err != nil {
			return err
		}
		for k, v := range namemap {
			fmt.Fprintf(f, "%s\t%s\n", k, v)
		}
		f.Close()
	}
	return nil
}
