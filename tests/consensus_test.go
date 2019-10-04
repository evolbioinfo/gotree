package tests

import (
	"bufio"
	"io"
	"testing"

	"github.com/evolbioinfo/gotree/io/newick"
	"github.com/evolbioinfo/gotree/io/utils"
	"github.com/evolbioinfo/gotree/tree"
)

/*
 Function to test consensus tree generation
 It compares majority consensus to a given already computed
 consensus (from phylip consense)
*/
func TestMajorityConsensus(t *testing.T) {
	var treefile, treefile2 io.Closer
	var treereader, treereader2 *bufio.Reader
	var trees <-chan tree.Trees
	var err error
	var majority, expected_majority *tree.Tree

	/* File reader (plain text or gzip) */
	if treefile, treereader, err = utils.GetReader("data/bootstrap_trees.nw.gz"); err != nil {
		t.Error(err)
	}
	defer treefile.Close()
	trees = utils.ReadMultiTrees(treereader, utils.FORMAT_NEWICK)

	majority, err = tree.Consensus(trees, 0.5)
	if err != nil {
		t.Error(err)
	}

	edgeindex1 := tree.NewEdgeIndex(128, .75)
	for _, e := range majority.Edges() {
		edgeindex1.PutEdgeValue(e, 1, e.Length())
	}

	// Parsing single tree newick file
	if treefile2, treereader2, err = utils.GetReader("data/bootstrap_majority.nw.gz"); err != nil {
		t.Error(err)
	}
	defer treefile2.Close()
	expected_majority, err = newick.NewParser(treereader2).Parse()
	if err != nil {
		t.Error(err)
	}
	edgeindex2 := tree.NewEdgeIndex(128, .75)
	for _, e := range expected_majority.Edges() {
		edgeindex2.PutEdgeValue(e, 1, e.Length())
		_, ok := edgeindex1.Value(e)
		if !ok {
			t.Error("One edge of the majority consensus is present in the expected tree but absent from the consensus")
		}
	}

	for _, e := range majority.Edges() {
		_, ok := edgeindex2.Value(e)
		if !ok {
			t.Error("One edge of the majority consensus is present in the consensus but absent from the expected tree")
		}
	}
}

/*
 Function to test consensus tree generation
 It compares majority consensus to a given already computed
 consensus (from phylip consense)
*/
func TestStrictConsensus(t *testing.T) {
	var treefile, treefile2 io.Closer
	var treereader, treereader2 *bufio.Reader
	var trees <-chan tree.Trees
	var err error
	var strict, expected_strict *tree.Tree

	/* File reader (plain text or gzip) */
	if treefile, treereader, err = utils.GetReader("data/bootstrap_trees.nw.gz"); err != nil {
		t.Error(err)
	}
	defer treefile.Close()
	trees = utils.ReadMultiTrees(treereader, utils.FORMAT_NEWICK)

	strict, err = tree.Consensus(trees, 1)
	if err != nil {
		t.Error(err)
	}
	edgeindex1 := tree.NewEdgeIndex(128, .75)
	for _, e := range strict.Edges() {
		edgeindex1.PutEdgeValue(e, 1, e.Length())
	}

	// Parsing single tree newick file
	if treefile2, treereader2, err = utils.GetReader("data/bootstrap_strict.nw.gz"); err != nil {
		t.Error(err)
	}
	defer treefile2.Close()
	expected_strict, err = newick.NewParser(treereader2).Parse()
	if err != nil {
		t.Error(err)
	}
	edgeindex2 := tree.NewEdgeIndex(128, .75)
	for _, e := range expected_strict.Edges() {
		edgeindex2.PutEdgeValue(e, 1, e.Length())
		_, ok := edgeindex1.Value(e)
		if !ok {
			t.Error("One edge of the Strict consensus is present in the expected tree but absent from the consensus")
		}
	}

	for _, e := range strict.Edges() {
		_, ok := edgeindex2.Value(e)
		if !ok {
			t.Error("One edge of the Strict consensus is present in the consensus but absent from the expected tree")
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
		trees <- tree.Trees{randtree1, 1, nil}
		close(trees)
	}()
	consens, err := tree.Consensus(trees, 0.5)
	if err != nil {
		t.Error(err)
	}

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
		trees2 <- tree.Trees{randtree1, 1, nil}
		trees2 <- tree.Trees{randtree2, 2, nil}
		trees2 <- tree.Trees{randtree3, 3, nil}
		close(trees2)
	}()
	consens, err = tree.Consensus(trees2, 1)
	if err != nil {
		t.Error(err)
	}

	ntips, _ := consens.NbTips()
	if len(consens.Edges()) > ntips {
		t.Error("Strict Consensus of 3 random binary trees (1000 tips) should strongly probably be a star tree")
	}
}
