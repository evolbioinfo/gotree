package phyloxml

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"

	"github.com/evolbioinfo/gotree/tree"
)

// Structs for representation of the PhyloXML tree
type PhyloXML struct {
	XMLName     xml.Name    `xml:"phyloxml"`
	Phylogenies []Phylogeny `xml:"phylogeny"`
}

type Phylogeny struct {
	XMLName xml.Name `xml:"phylogeny"`
	Rooted  bool     `xml:"rooted,attr"`
	Root    Clade    `xml:"clade"`
}

type Clade struct {
	XMLName      xml.Name `xml:"clade"`
	Clades       []Clade  `xml:"clade"`
	BranchLength *float64 `xml:"branch_length"`
	Confidence   *float64 `xml:"confidence"`
	Name         string   `xml:"name"`
	Tax          Taxonomy `xml:"taxonomy"`
}

type Taxonomy struct {
	XMLName        xml.Name   `xml:"taxonomy"`
	TaxId          TaxonomyId `xml:"id"`
	ScientificName string     `xml:"scientific_name"`
	Code           string     `xml:"code"`
}

type TaxonomyId struct {
	Id       int
	Provider string `xml:"provider,attr"`
}

// Parser represents a parser.
type Parser struct {
	reader io.Reader
}

// NewParser returns a new instance of Parser.
func NewParser(r io.Reader) *Parser {
	return &Parser{reader: r}
}

func (p *Parser) Parse() (px *PhyloXML, err error) {
	px = &PhyloXML{}
	buf := new(bytes.Buffer)
	buf.ReadFrom(p.reader)
	err = xml.Unmarshal(buf.Bytes(), px)
	return
}

func (p *PhyloXML) IterateTrees(it func(*tree.Tree, error)) {
	for _, phylo := range p.Phylogenies {
		t := tree.NewTree()
		err := phylogenyToTree(&phylo, t)
		it(t, err)
	}
}

func (p *PhyloXML) FirstTree() (t *tree.Tree, err error) {
	for _, phylo := range p.Phylogenies {
		t := tree.NewTree()
		err = phylogenyToTree(&phylo, t)
		break
	}
	return
}

func phylogenyToTree(p *Phylogeny, t *tree.Tree) (err error) {
	err = cladeToTree(&p.Root, t, nil)
	return
}

func cladeToTree(c *Clade, t *tree.Tree, parent *tree.Node) (err error) {
	newNode := t.NewNode()
	if parent == nil {
		t.SetRoot(newNode)
	} else {
		e := t.ConnectNodes(parent, newNode)
		if c.BranchLength != nil {
			e.SetLength(*(c.BranchLength))
		}
		if len(c.Clades) > 0 {
			if c.Confidence != nil {
				e.SetSupport(*(c.Confidence))
			}
		}
	}
	if c.Name != "" {
		newNode.SetName(c.Name)
	} else if c.Tax.ScientificName != "" {
		newNode.SetName(c.Tax.ScientificName)
	} else if c.Tax.Code != "" {
		newNode.SetName(c.Tax.Code)
	}
	for _, cl := range c.Clades {
		err = cladeToTree(&cl, t, newNode)
		if err != nil {
			return
		}
	}
	if len(c.Clades) == 0 && newNode.Name() == "" {
		err = fmt.Errorf("One tip has no name")
	}
	return
}

// func printTaxonomy(t Taxonomy, level int) {
// 	tab := ""
// 	for i := 0; i < level; i++ {
// 		tab += "\t"
// 	}
// 	fmt.Printf("%ssciname:%s\n", tab, t.ScientificName)
// 	fmt.Printf("%scode:%s\n", tab, t.Code)
// 	fmt.Printf("%sid:%d\n", tab, t.TaxId.Id)
// 	fmt.Printf("%sprovider:%s\n", tab, t.TaxId.Provider)
// }

func WritePhyloXML(tchan <-chan tree.Trees) (string, error) {
	var buffer bytes.Buffer
	buffer.WriteString(
		`<?xml version="1.0" encoding="UTF-8"?>
<phyloxml xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" 
          xsi:schemaLocation="http://www.phyloxml.org http://www.phyloxml.org/1.10/phyloxml.xsd"
          xmlns="http://www.phyloxml.org">
`)
	for t := range tchan {
		if t.Err != nil {
			return "", t.Err
		}
		writePhylogeny(t.Tree, &buffer)
	}
	buffer.WriteString("</phyloxml>\n")
	return buffer.String(), nil
}

func writePhylogeny(t *tree.Tree, buf *bytes.Buffer) {
	buf.WriteString(fmt.Sprintf("  <phylogeny rooted=\"%t\">\n", t.Rooted()))
	writeClade(t.Root(), nil, nil, buf, 1)
	buf.WriteString("  </phylogeny>\n")
}

func writeClade(n *tree.Node, prev *tree.Node, e *tree.Edge, buf *bytes.Buffer, level int) {
	tab := "  "
	for i := 0; i < level; i++ {
		tab += "  "
	}

	buf.WriteString(tab + "<clade>\n")
	if n.Name() != "" {
		buf.WriteString(fmt.Sprintf("%s<name>%s</name>\n", tab, n.Name()))
	}
	if prev != nil && e != nil {
		if e.Length() != tree.NIL_LENGTH {
			buf.WriteString(fmt.Sprintf("%s<branch_length>%s</branch_length>\n", tab, e.LengthString()))
		}
		if !n.Tip() && e.Support() != tree.NIL_SUPPORT {
			buf.WriteString(fmt.Sprintf("%s<confidence type=\"bootstrap\">%s</confidence>\n", tab, e.SupportString()))
		}
	}
	for i, child := range n.Neigh() {
		if child != prev {
			nextedge := n.Edges()[i]
			writeClade(child, n, nextedge, buf, level+1)
		}
	}
	buf.WriteString(tab + "</clade>\n")
}
