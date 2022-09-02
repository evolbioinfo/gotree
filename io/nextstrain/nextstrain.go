package nextstrain

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/evolbioinfo/gotree/tree"
)

// Structs for representation of the PhyloXML tree
type Nextstrain struct {
	Tree    NsNode `json:"tree"`
	Version string `json:"version"`
}

type NsNode struct {
	BranchAttr NSBranchAttributes `json:"branch_attrs"`
	Children   []NsNode           `json:"children"`
	Name       string             `json:"name"`
	Attributes NsNodeAttributes   `json:"node_attrs"`
}

type NsNodeAttributes struct {
	Divergence          float64 `json:"div"`
	LocalBranchingIndex struct {
		Value float64 `json:"value"`
	} `json:"lbi"`
	Date      NsDate   `json:"num_date"`
	Region    NsRegion `json:"region"`
	Accession string   `json:"accession"`
	Age       struct {
		Value string `json:"value"`
	} `json:"age"`
	CladeMembership struct {
		Value string `json:"value"`
	} `json:"clade_membership"`
	Country struct {
		Value string `json:"value"`
	} `json:"country"`
	Division struct {
		Value string `json:"value"`
	} `json:"division"`
	Epiweek struct {
		Value string `json:"value"`
	} `json:"epiweek"`
	Gender struct {
		Value string `json:"value"`
	} `json:"gender"`
	OriginatingLab struct {
		Value string `json:"value"`
	} `json:"originating_lab"`
	Recency struct {
		Value string `json:"value"`
	} `json:"recency"`
	SubmittingLab struct {
		Value string `json:"value"`
	} `json:"submitting_lab"`
}

type NSBranchAttributes struct {
	Labels struct {
		Aa string `json:"aa"`
	} `json:"labels"`
	Mutations map[string][]string `json:"mutations"`
}

type NsDate struct {
	Value      float64   `json:"value"`
	Confidence []float64 `json:"confidence"`
}

type NsRegion struct {
	Value      string             `json:"value"`
	Entropy    float64            `json:"entropy"`
	Confidence map[string]float64 `json:"confidence"`
}

// Parser represents a parser.
type Parser struct {
	reader io.Reader
}

// NewParser returns a new instance of Parser.
func NewParser(r io.Reader) *Parser {
	return &Parser{reader: r}
}

func (p *Parser) Parse() (ns *Nextstrain, err error) {
	ns = &Nextstrain{}
	buf := new(bytes.Buffer)
	buf.ReadFrom(p.reader)
	err = json.Unmarshal(buf.Bytes(), ns)
	if err == nil && ns.Version != "v2" {
		err = fmt.Errorf("format error : gotree only supports nextstrain v2 format")
	}
	return
}

func (n *Nextstrain) IterateTrees(it func(*tree.Tree, error)) {
	phylo := n.Tree
	t := tree.NewTree()
	err := phylogenyToTree(&phylo, t)
	it(t, err)
}

func (n *Nextstrain) FirstTree() (t *tree.Tree, err error) {
	phylo := n.Tree
	t = tree.NewTree()
	err = phylogenyToTree(&phylo, t)
	return
}

func phylogenyToTree(r *NsNode, t *tree.Tree) (err error) {
	var nedges, nnodes int = 0, 0
	err = cladeToTree(r, t, nil, &nedges, &nnodes, r.Attributes.Divergence)
	return
}

func cladeToTree(c *NsNode, t *tree.Tree, parent *tree.Node, nedges, nnodes *int, prevdiv float64) (err error) {
	newNode := t.NewNode()
	newNode.SetId(*nnodes)
	(*nnodes)++
	if parent == nil {
		t.SetRoot(newNode)
	} else {
		e := t.ConnectNodes(parent, newNode)
		e.SetId(*nedges)
		(*nedges)++
		e.SetLength(c.Attributes.Divergence - prevdiv)
	}

	firstannot := true
	comment := ""
	if c.BranchAttr.Labels.Aa != "" {
		mut := strings.Replace(c.BranchAttr.Labels.Aa, ":", ".", -1)
		mut = strings.Replace(mut, " ", "", -1)
		mut = strings.Replace(mut, ",", "-", -1)
		if firstannot {
			mut = "&mutations=" + mut
		} else {
			mut = ",mutations=" + mut
		}
		comment += mut
		firstannot = false
	}

	if c.Attributes.Accession != "" {
		acc := strings.Replace(c.Attributes.Accession, ":", ".", -1)
		acc = strings.Replace(acc, " ", "", -1)
		acc = strings.Replace(acc, ",", "-", -1)
		if firstannot {
			acc = "&accession=" + acc
		} else {
			acc = ",accession=" + acc
		}
		comment += acc
		firstannot = false
	}

	if c.Attributes.Country.Value != "" {
		country := strings.Replace(c.Attributes.Country.Value, ":", ".", -1)
		country = strings.Replace(country, " ", "", -1)
		country = strings.Replace(country, ",", "-", -1)
		if firstannot {
			country = "&country=" + country
		} else {
			country = ",country=" + country
		}
		comment += country
		firstannot = false
	}

	if c.Attributes.Date.Value != 0.0 {
		date := fmt.Sprintf("%f", c.Attributes.Date.Value)
		if firstannot {
			date = "&date=" + date
		} else {
			date = ",date=" + date
		}
		comment += date
		firstannot = false
	}

	if !firstannot {
		newNode.AddComment(comment)
	}

	if c.Name != "" {
		newNode.SetName(c.Name)
	}
	for _, cl := range c.Children {
		err = cladeToTree(&cl, t, newNode, nedges, nnodes, c.Attributes.Divergence)
		if err != nil {
			return
		}
	}
	if len(c.Children) == 0 && newNode.Name() == "" {
		err = fmt.Errorf("one tip has no name")
	}
	return
}
