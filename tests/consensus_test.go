package tests

import (
	"github.com/fredericlemoine/gotree/io/utils"
	"github.com/fredericlemoine/gotree/tree"
	"testing"
)

/*
 Function to test consensus tree generation
 It compares majority and strict consensus to a given already computed
 consensus (from phylip consense)
*/
func TestConsensus(t *testing.T) {
	trees := make(chan tree.Trees)
	trees2 := make(chan tree.Trees)
	go func() {
		utils.ReadCompTrees("data/bootstrap_trees.nw.gz", trees)
	}()
	majority := tree.Consensus(trees, 0.5)
	edgeindex1 := tree.NewEdgeIndex(128, .75)
	for _, e := range majority.Edges() {
		edgeindex1.PutEdgeValue(e, 1, e.Length())
	}
	go func() {
		utils.ReadCompTrees("data/bootstrap_trees.nw.gz", trees2)
	}()
	strict := tree.Consensus(trees2, 1)
	edgeindex2 := tree.NewEdgeIndex(128, .75)
	for _, e := range strict.Edges() {
		edgeindex2.PutEdgeValue(e, 1, e.Length())
	}

	expected_majority, _ := utils.ReadRefTree("data/bootstrap_majority.nw.gz")
	edgeindex3 := tree.NewEdgeIndex(128, .75)
	for _, e := range expected_majority.Edges() {
		edgeindex3.PutEdgeValue(e, 1, e.Length())
		_, ok := edgeindex1.Value(e)
		if !ok {
			t.Error("One edge of the majority consensus is present in the expected tree but absent from the consensus")
		}
	}

	expected_strict, _ := utils.ReadRefTree("data/bootstrap_strict.nw.gz")
	edgeindex4 := tree.NewEdgeIndex(128, .75)
	for _, e := range expected_strict.Edges() {
		edgeindex4.PutEdgeValue(e, 1, e.Length())
		_, ok := edgeindex2.Value(e)
		if !ok {
			t.Error("One edge of the strict consensus is present in the expected tree but absent from the consensus")
		}
	}

	for _, e := range majority.Edges() {
		_, ok := edgeindex3.Value(e)
		if !ok {
			t.Error("One edge of the majority consensus is present in the consensus but absent from the expected tree")
		}
	}

	for _, e := range strict.Edges() {
		_, ok := edgeindex4.Value(e)
		if !ok {
			t.Error("One edge of the strict consensus is present in the consensus but absent from the expected tree")
		}
	}

}

// We generate a random tree and the consensus
func TestConsensus2(t *testing.T) {
	trees := make(chan tree.Trees, 4)
	trees2 := make(chan tree.Trees, 4)
	var randtree1, randtree2, randtree3 *tree.Tree
	var err error

	if randtree1, err = tree.RandomUniformBinaryTree(1000, false); err != nil {
		t.Error(err)
	}
	if randtree2, err = tree.RandomUniformBinaryTree(1000, false); err != nil {
		t.Error(err)
	}
	if randtree3, err = tree.RandomUniformBinaryTree(1000, false); err != nil {
		t.Error(err)
	}

	go func() {
		trees <- tree.Trees{randtree1, 1}
		close(trees)
	}()
	consens := tree.Consensus(trees, 0.5)
	edgeindex := tree.NewEdgeIndex(128, .75)
	if len(consens.Edges()) != len(randtree1.Edges()) {
		t.Error("Consensus and Initial trees have different number of edges")
	}
	for i := 0; i < 100; i++ {
		for _, e := range randtree1.Edges() {
			edgeindex.PutEdgeValue(e, 1, e.Length())
		}
	}
	for _, e := range consens.Edges() {
		if _, ok := edgeindex.Value(e); !ok {
			t.Error("Edge is not present in the consensus, and should be")
		}
	}

	go func() {
		trees2 <- tree.Trees{randtree1, 1}
		trees2 <- tree.Trees{randtree2, 2}
		trees2 <- tree.Trees{randtree3, 3}
		close(trees2)
	}()
	consens = tree.Consensus(trees2, 1)
	ntips, _ := consens.NbTips()
	if len(consens.Edges()) > ntips {
		t.Error("Strict Consensus of 3 random binary trees (1000 tips) should strongly probably be a star tree")
	}
}
