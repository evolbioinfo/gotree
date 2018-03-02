package tests

import (
	"fmt"
	"github.com/fredericlemoine/gotree/io/newick"
	"strings"
	"testing"
)

var treestring string = "(Tip2:1.00000,Node0:1.0000,((Tip7:1.00000,((Tip9:1.00000,Tip6:1.0000):1.0000,(Tip5:1.00000,Tip3:1.0000):1.0000):1.00):1.00,(Tip4:1.00000,(Tip8:1.00000,Tip1:1.000):0.126):0.127):0.125);"

func TestCollapse(t *testing.T) {
	tree, err2 := newick.NewParser(strings.NewReader(treestring)).Parse()

	if err2 != nil {
		t.Error(err2)
	}

	edges := tree.Edges()
	sumlen := tree.SumBranchLengths()
	if len(edges) != 17 {
		t.Error(fmt.Sprintf("The number of edges before collapse is not 17 (%d)", len(edges)))
	}
	if sumlen != 14.378 {
		t.Error(fmt.Sprintf("The sum of branch lengths before collapse is not  (%f)", sumlen))
	}

	tree.CollapseShortBranches(0.126)
	edges = tree.Edges()
	sumlen = tree.SumBranchLengths()

	if len(edges) != 15 {
		t.Error(fmt.Sprintf("The number of edges after collapse is not 15 (%d)", len(edges)))
	}

	if sumlen != 14.127 {
		t.Error(fmt.Sprintf("The sum of branch lengths after collapse is not 14.127 (%f)", sumlen))
	}
}

var treestring2 string = "(A:1,(B:1):1,C:1);"

func TestCollapseSingle(t *testing.T) {
	tree, err2 := newick.NewParser(strings.NewReader(treestring2)).Parse()
	if err2 != nil {
		t.Error(err2)
	}
	nbranches := len(tree.Edges())
	nnodes := len(tree.Nodes())
	sumlen := tree.SumBranchLengths()
	if nbranches != 4 {
		t.Error(fmt.Sprintf("The number of edges before collapse is not 4 (%d)", nbranches))
	}
	if sumlen != 4.0 {
		t.Error(fmt.Sprintf("The sum of branch lengths before collapse is not 4.0  (%f)", sumlen))
	}
	if nnodes != 5 {
		t.Error(fmt.Sprintf("The number of nodes before collapse is not 5  (%d)", nnodes))
	}

	tree.RemoveSingleNodes()
	nbranches = len(tree.Edges())
	sumlen = tree.SumBranchLengths()
	nnodes = len(tree.Nodes())

	if nbranches != 3 {
		t.Error(fmt.Sprintf("The number of edges after collapse is not 3 (%d)", nbranches))
	}

	if sumlen != 4.0 {
		t.Error(fmt.Sprintf("The sum of branch lengths after collapse is not 3.0 (%f)", sumlen))
	}
	if nnodes != 4 {
		t.Error(fmt.Sprintf("The number of nodes after collapse is not 4  (%d)", nnodes))
	}
}
