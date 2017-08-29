package phyloxml

import (
	"bytes"
	"encoding/xml"
	"io"

	"github.com/fredericlemoine/gotree/tree"
)

// Structs for representation of the PhyloXML tree
type PhyloXML struct {
	XMLName xml.Name  `xml:"phyloxml"`
	Phylo   Phylogeny `xml:phylogeny`
}

type Phylogeny struct {
	XMLName xml.Name `xml:"phylogeny"`
	Rooted  bool     `xml:"rooted,attr"`
	Root    Clade    `xml:"clade"`
}

type Clade struct {
	XMLName      xml.Name `xml:"clade"`
	Clades       []Clade  `xml:"clade"`
	BranchLength float64  `xml:"branch_length"`
	Confidence   float64  `xml:"confidence"`
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

func (p *Parser) Parse() (t *tree.Tree, err error) {
	var q PhyloXML
	buf := new(bytes.Buffer)
	buf.ReadFrom(p.reader)
	err = xml.Unmarshal(buf.Bytes(), &q)
	if err == nil {
		t = phyloXMLToTree(q)
	}
	return
}

func phyloXMLToTree(q PhyloXML) (t *tree.Tree) {
	t = tree.NewTree()
	phylogenyToTree(q.Phylo, t)
	return
}

func phylogenyToTree(p Phylogeny, t *tree.Tree) {
	cladeToTree(p.Root, t, nil)
}

func cladeToTree(c Clade, t *tree.Tree, parent *tree.Node) {
	newNode := t.NewNode()
	if parent == nil {
		t.SetRoot(newNode)
	} else {
		e := t.ConnectNodes(parent, newNode)
		e.SetLength(c.BranchLength)
		e.SetSupport(c.Confidence)
	}
	if c.Name != "" {
		newNode.SetName(c.Name)
	} else if c.Tax.ScientificName != "" {
		newNode.SetName(c.Tax.ScientificName)
	} else if c.Tax.Code != "" {
		newNode.SetName(c.Tax.Code)
	}
	for _, cl := range c.Clades {
		cladeToTree(cl, t, newNode)
	}
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
