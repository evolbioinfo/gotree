package tests

import (
	"fmt"
	"github.com/evolbioinfo/gotree/io/newick"
	"strings"
	"testing"
)

var roottree1 string = "((1:1,2:1):1,(3:1,4:1):1);"
var roottree2 string = "(1:1,(2:1,(3:1,4:1):1):1);"
var roottree3 string = "(((3:1,4:1):1,2:1):1,1:1);"

func TestUnroot(t *testing.T) {
	tree, err := newick.NewParser(strings.NewReader(roottree1)).Parse()
	if err != nil {
		t.Error(err)
	}

	edges := tree.Edges()
	sumlen := tree.SumBranchLengths()
	if len(edges) != 6 {
		t.Error(fmt.Sprintf("The number of edges before unroot is not 6 (%d)", len(edges)))
	}
	if sumlen != 6 {
		t.Error(fmt.Sprintf("The sum of branch lengths before unroot is not 6  (%f)", sumlen))
	}

	tree.UnRoot()
	edges = tree.Edges()
	sumlen = tree.SumBranchLengths()

	if len(edges) != 5 {
		t.Error(fmt.Sprintf("The number of edges after unroot is not 6 (%d)", len(edges)))
	}

	if sumlen != 6 {
		t.Error(fmt.Sprintf("The sum of branch lengths after unroot is not 6 (%f)", sumlen))
	}
}

func TestUnroot2(t *testing.T) {
	tree, err := newick.NewParser(strings.NewReader(roottree2)).Parse()
	if err != nil {
		t.Error(err)
	}

	edges := tree.Edges()
	sumlen := tree.SumBranchLengths()
	if len(edges) != 6 {
		t.Error(fmt.Sprintf("The number of edges before unroot is not 6 (%d)", len(edges)))
	}
	if sumlen != 6 {
		t.Error(fmt.Sprintf("The sum of branch lengths before unroot is not 6  (%f)", sumlen))
	}

	tree.UnRoot()
	edges = tree.Edges()
	sumlen = tree.SumBranchLengths()

	if len(edges) != 5 {
		t.Error(fmt.Sprintf("The number of edges after unroot is not 6 (%d)", len(edges)))
	}

	if sumlen != 6 {
		t.Error(fmt.Sprintf("The sum of branch lengths after unroot is not 6 (%f)", sumlen))
	}
}

func TestUnroot3(t *testing.T) {
	tree, err := newick.NewParser(strings.NewReader(roottree3)).Parse()
	if err != nil {
		t.Error(err)
	}

	edges := tree.Edges()
	sumlen := tree.SumBranchLengths()
	if len(edges) != 6 {
		t.Error(fmt.Sprintf("The number of edges before unroot is not 6 (%d)", len(edges)))
	}
	if sumlen != 6 {
		t.Error(fmt.Sprintf("The sum of branch lengths before unroot is not 6  (%f)", sumlen))
	}

	tree.UnRoot()
	edges = tree.Edges()
	sumlen = tree.SumBranchLengths()

	if len(edges) != 5 {
		t.Error(fmt.Sprintf("The number of edges after unroot is not 6 (%d)", len(edges)))
	}

	if sumlen != 6 {
		t.Error(fmt.Sprintf("The sum of branch lengths after unroot is not 6 (%f)", sumlen))
	}
}
