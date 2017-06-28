package download

import (
	"archive/tar"
	"bufio"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/fredericlemoine/gotree/io/utils"
	"github.com/fredericlemoine/gotree/tree"
	"github.com/jlaffaye/ftp"
)

type NcbiTreeDownloader struct {
	server string
	path   string
}

/* NCBI taxonomy downloader */
func NewNcbiTreeDownloader() *NcbiTreeDownloader {
	return &NcbiTreeDownloader{"ftp.ncbi.nih.gov:21", "/pub/taxonomy/taxdump.tar.gz"}
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
	reader, err = client.Retr(d.path)
	if err != nil {
		return nil, err
	}

	// Reading tar gz and processing nodes.dmp and names.dmp
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
			fmt.Println("Name: ", header.Name)
			// We handle names of ncbi taxonomy nodes
			if header.Name == "names.dmp" {
				namemap, err = ParseNcbiNames(tarreader)
				if err != nil {
					return nil, err
				}
			}
			// We handle the tree
			if header.Name == "nodes.dmp" {
				t, err = ParseNcbiTree(tarreader)
				if err != nil {
					return nil, err
				}
			}
		default:
			return nil, errors.New("Problem with tar archive")
		}
	}

	RenameTreeNodes(t, namemap)

	return t, err
}

func RenameTreeNodes(t *tree.Tree, namemap map[string]string) {
	for _, n := range t.Nodes() {
		if name, ok := namemap[n.Name()]; ok {
			n.SetName(name)
		}
	}
}

func ParseNcbiNames(reader io.Reader) (map[string]string, error) {
	r := bufio.NewReader(reader)
	l, err := utils.Readln(r)
	namemap := make(map[string]string)
	for err == nil {
		cols := regexp.MustCompile("\t*\\|\t*").Split(l, -1)
		tax := cols[0]
		name := cols[1]
		tpe := cols[3]
		if tpe == "scientific name" || tpe == "synonym" {
			fmt.Println(tax)
			fmt.Println(name)
			namemap[tax] = name
		}
		l, err = utils.Readln(r)
	}
	return namemap, nil
}

func ParseNcbiTree(reader io.Reader) (*tree.Tree, error) {
	r := bufio.NewReader(reader)
	l, err := utils.Readln(r)
	t := tree.NewTree()
	var root *tree.Node
	nodes := make(map[string]*tree.Node)
	for err == nil {
		cols := strings.Split(l, "\t|\t")
		tax := cols[0]
		parent := cols[1]
		n1, ok1 := nodes[tax]
		if !ok1 {
			n1 = t.NewNode()
			n1.SetName(tax)
			nodes[tax] = n1
		}
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

		l, err = utils.Readln(r)
	}
	if root == nil {
		return nil, errors.New("No root found in the NCBI Taxonomy")
	}
	// t.UpdateTipIndex()
	// t.ClearBitSets()
	// t.UpdateBitSet()
	// t.ComputeDepths()
	return t, nil
}
