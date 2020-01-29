package download

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/evolbioinfo/gotree/tree"
)

// PantherTreeDownloader allows to download trees from Panther
// using a family ID
type PantherTreeDownloader struct {
	server string
	path   string
}

// NewPantherTreeDownloader initializes a new Panther tree downloader
func NewPantherTreeDownloader() *PantherTreeDownloader {
	return &PantherTreeDownloader{
		server: "http://pantherdb.org/",
		path:   "services/oai/pantherdb/treeinfo",
	}
}

// Download a tree from Panther
func (p *PantherTreeDownloader) Download(id string) (t *tree.Tree, err error) {
	geturl := fmt.Sprintf("%s/%s?family=%s", p.server, p.path, id)
	var getresponse *http.Response
	var responsebody []byte
	var answer PantherAnswer

	if getresponse, err = http.Get(geturl); err != nil {
		return
	}
	defer getresponse.Body.Close()

	if responsebody, err = ioutil.ReadAll(getresponse.Body); err != nil {
		return
	}

	if err = json.Unmarshal(responsebody, &answer); err != nil {
		err = fmt.Errorf("%s (%s)", err.Error(), string(responsebody))
		return
	}

	if answer.Search.Error != "" {
		err = errors.New(string(answer.Search.Error))
		return
	}

	if t, err = p.treeFromPantherAnswer(&answer); err != nil {
		return
	}

	return
}

// PantherAnswer is the root of Panther JSON answer
type PantherAnswer struct {
	Search PantherAnswerSearch `json:"search"`
}

// PantherAnswerSearch defines information on answer search
type PantherAnswerSearch struct {
	Product      PantherAnswerProduct      `json:"product"`
	SearchType   string                    `json:"search_type"`
	Parameters   PantherAnswerParameters   `json:"parameters"`
	TreeTopology PantherAnswerTreeTopology `json:"tree_topology"`
	Error        string                    `json:"error"`
}

// PantherAnswerProduct defines information version and source of the answer
type PantherAnswerProduct struct {
	Source  string  `json:"source"`
	Version float64 `json:"version"`
}

// PantherAnswerParameters defines information about family
type PantherAnswerParameters struct {
	Family string `json:"family"`
}

// PantherAnswerTreeTopology defines information on node
type PantherAnswerTreeTopology struct {
	AnnotationNode PantherAnswerAnnotationNode `json:"annotation_node"`
}

// PantherAnswerAnnotationNode defines information about the node
type PantherAnswerAnnotationNode struct {
	SfID           string                `json:"sf_id"`
	PersistentID   string                `json:"persistent_id"`
	BranchLength   float64               `json:"branch_length"`
	PropSfID       string                `json:"prop_sf_id"`
	EventType      string                `json:"event_type"`
	Species        string                `json:"species"`
	TreeNodeType   string                `json:"tree_node_type"`
	TaxonomicRange string                `json:"taxonomic_range"`
	SfName         string                `json:"sf_name"`
	GeneSymbol     string                `json:"gene_symbol"`
	NodeName       string                `json:"node_name"`
	Organism       string                `json:"organism"`
	Children       PantherAnswerChildren `json:"children"`
}

// PantherAnswerChildren defines the children type
type PantherAnswerChildren struct {
	TreeNodeType   string                        `json:"tree_node_type"`
	TaxonomicRange string                        `json:"taxonomic_range"`
	SfName         string                        `json:"sf_name"`
	AnnotationNode []PantherAnswerAnnotationNode `json:"annotation_node"`
}

func (p *PantherTreeDownloader) treeFromPantherAnswer(answer *PantherAnswer) (t *tree.Tree, err error) {
	var nedges, nnodes int = 0, 0

	t = tree.NewTree()
	err = p.annotationNodeToTree(&answer.Search.TreeTopology.AnnotationNode, t, nil, &nedges, &nnodes)
	return
}

func (p *PantherTreeDownloader) annotationNodeToTree(an *PantherAnswerAnnotationNode, t *tree.Tree, parent *tree.Node, nedges, nnodes *int) (err error) {
	newNode := t.NewNode()
	newNode.SetId(*nnodes)
	(*nnodes)++
	if parent == nil {
		t.SetRoot(newNode)
	} else {
		e := t.ConnectNodes(parent, newNode)
		e.SetId(*nedges)
		(*nedges)++
		if an.BranchLength != 0 {
			e.SetLength(an.BranchLength)
		}
	}
	if an.TreeNodeType == "LEAF" {
		if an.NodeName != "" {
			newNode.SetName(an.NodeName)
		} else {
			err = fmt.Errorf("One tip has no name -%s-", an.SfID)
		}
	} else {
		if len(an.Children.AnnotationNode) == 0 {
			err = fmt.Errorf("One internal node has 0 children -%s-", an.SfID)
		}
		if an.Species != "" {
			newNode.SetName(an.Species)
		}
	}

	if an.EventType != "" {
		newNode.AddComment(an.EventType)
	}

	if an.GeneSymbol != "" {
		newNode.AddComment(an.GeneSymbol)
	}

	if an.Organism != "" {
		newNode.AddComment(an.Organism)
	}

	if an.PersistentID != "" {
		newNode.AddComment(an.PersistentID)
	}

	for _, cl := range an.Children.AnnotationNode {
		if err = p.annotationNodeToTree(&cl, t, newNode, nedges, nnodes); err != nil {
			return
		}
	}

	return
}
