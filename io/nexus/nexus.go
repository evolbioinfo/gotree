package nexus

import (
	"bytes"
	"fmt"
	"sort"
	"strconv"

	"github.com/evolbioinfo/goalign/align"
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

// WriteNexus generates a Nexus string from a tree channel
// if translate is true, then it replaces all tip names by indices and
// generates a translate block.
func WriteNexus(tchan <-chan tree.Trees, translate bool) (string, error) {
	var treeBuffer bytes.Buffer
	var fullBuffer bytes.Buffer

	var taxLabelsMap map[string]string = make(map[string]string)
	var taxLabelsSlice []string = make([]string, 0)
	var nbTax int = 0

	for t := range tchan {
		if t.Err != nil {
			return "", t.Err
		}

		for _, tip := range t.Tree.AllTipNames() {
			if _, ok := taxLabelsMap[tip]; !ok {
				taxLabelsMap[tip] = fmt.Sprintf("%d", nbTax)
				taxLabelsSlice = append(taxLabelsSlice, tip)
				nbTax++
			}
		}
		sort.Strings(taxLabelsSlice)

		renameTree := t.Tree
		if translate {
			renameTree = t.Tree.Clone()
			renameTree.Rename(taxLabelsMap)
		}
		treeBuffer.WriteString("  TREE tree")
		treeBuffer.WriteString(strconv.Itoa(t.Id))
		treeBuffer.WriteString(" = ")
		treeBuffer.WriteString(renameTree.Newick())
		treeBuffer.WriteString("\n")
	}

	fullBuffer.WriteString("#NEXUS\n")
	fullBuffer.WriteString("BEGIN TAXA;\n")
	fullBuffer.WriteString(" DIMENSIONS NTAX=")
	fullBuffer.WriteString(strconv.Itoa(len(taxLabelsMap)))
	fullBuffer.WriteString(";\n")
	fullBuffer.WriteString(" TAXLABELS")

	for _, tip := range taxLabelsSlice {
		fullBuffer.WriteString(" " + tip)
	}

	fullBuffer.WriteString(";\n")
	fullBuffer.WriteString("END;\n")
	fullBuffer.WriteString("BEGIN TREES;\n")

	if translate {
		fullBuffer.WriteString("  TRANSLATE\n")
		for _, tip := range taxLabelsSlice {
			fullBuffer.WriteString(fmt.Sprintf("   %s %s\n", taxLabelsMap[tip], tip))
		}
		fullBuffer.WriteString("  ;\n")
	}

	fullBuffer.Write(treeBuffer.Bytes())
	fullBuffer.WriteString("END;\n")

	return fullBuffer.String(), nil
}
