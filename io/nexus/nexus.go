package nexus

import (
	"bytes"
	"strconv"

	"github.com/fredericlemoine/goalign/align"
	"github.com/evolbioinfo/gotree/tree"
)

// The nexus structure, with several trees (gotree) and one alignment (goalign)
type Nexus struct {
	HasAlignment bool            // If the Nexus structure has contains an Alignment
	HasTrees     bool            // If the Nexus structure has contains a Tree
	GapChar      rune            // Gap character in the alignment
	MissingChar  rune            // Missing character in the alignment
	trees        []*tree.Tree    // Set of trees
	treeNames    []string        // Set of tree names
	align        align.Alignment // Alignment
}

func NewNexus() *Nexus {
	return &Nexus{
		HasAlignment: false,
		HasTrees:     false,
		GapChar:      '-',
		MissingChar:  '*',
		trees:        make([]*tree.Tree, 0),
		treeNames:    make([]string, 0),
		align:        nil,
	}
}

func (n *Nexus) AddTree(name string, t *tree.Tree) {
	n.trees = append(n.trees, t)
	n.treeNames = append(n.treeNames, name)
	n.HasTrees = true
}

func (n *Nexus) SetAlignment(align align.Alignment) {
	n.align = align
	n.HasAlignment = true
}
func (n *Nexus) Alignment() align.Alignment {
	return n.align
}

func (n *Nexus) IterateTrees(it func(string, *tree.Tree)) {
	for i, t := range n.trees {
		it(n.treeNames[i], t)
	}
}

// returns the first tree of the nexus data structure
// If no tree is present, then returns nil
func (n *Nexus) FirstTree() *tree.Tree {
	if len(n.trees) > 0 {
		return n.trees[0]
	}
	return nil
}

// returns the first tree of the nexus data structure
// If no tree is present, then returns nil
func (n *Nexus) NTrees() int {
	return len(n.trees)
}

// Generates a Nexus string from a tree channel.
func WriteNexus(tchan <-chan tree.Trees) (string, error) {
	var buffer bytes.Buffer
	buffer.WriteString("#NEXUS\n")
	taxlabels := false
	for t := range tchan {
		if t.Err != nil {
			return "", t.Err
		}

		if !taxlabels {
			buffer.WriteString("BEGIN TAXA;\n")
			tips := t.Tree.Tips()
			buffer.WriteString(" DIMENSIONS NTAX=")
			buffer.WriteString(strconv.Itoa(len(tips)))
			buffer.WriteString(";\n")
			buffer.WriteString(" TAXLABELS")
			for _, tip := range tips {
				buffer.WriteString(" " + tip.Name())
			}
			buffer.WriteString(";\n")
			buffer.WriteString("END;\n")
			buffer.WriteString("BEGIN TREES;\n")
			taxlabels = true
		}
		buffer.WriteString("  TREE tree")
		buffer.WriteString(strconv.Itoa(t.Id))
		buffer.WriteString(" = ")
		buffer.WriteString(t.Tree.Newick())
		buffer.WriteString("\n")
	}
	buffer.WriteString("END;\n")
	return buffer.String(), nil
}
