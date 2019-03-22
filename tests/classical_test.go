package tests

import (
	"bufio"
	"fmt"
	"io"
	"testing"

	"github.com/evolbioinfo/gotree/io/newick"
	"github.com/evolbioinfo/gotree/io/utils"
	"github.com/evolbioinfo/gotree/support"
	"github.com/evolbioinfo/gotree/tree"
)

func TestClassicalSupport_1(t *testing.T) {
	var reftree *tree.Tree
	var treefile, reftreefile io.Closer
	var treereader, reftreereader *bufio.Reader
	var err error
	var trees <-chan tree.Trees

	// Parsing multi tree newick (compared trees
	if treefile, treereader, err = utils.GetReader("data/rand_tree_boot.nw.gz"); err != nil {
		t.Error(err)
	}
	defer treefile.Close()
	trees = utils.ReadMultiTrees(treereader, utils.FORMAT_NEWICK)

	// Parsing single tree newick file
	// Parsing multi tree newick (compared trees
	if reftreefile, reftreereader, err = utils.GetReader("data/rand_tree.nw.gz"); err != nil {
		t.Error(err)
	}
	defer reftreefile.Close()
	reftree, err = newick.NewParser(reftreereader).Parse()
	if err != nil {
		t.Error(err)
	}

	// Computing fbp
	err = support.Classical(reftree, trees, 1)
	if err != nil {
		t.Error(err)
	}

	for _, e := range reftree.Edges() {
		if !e.Right().Tip() && e.Support() != 0 {
			t.Error(fmt.Sprintf("Non Tip support should be 0.00 and is %.2f", e.Support()))
		} else if e.Right().Tip() && e.Support() != -1 {
			t.Error(fmt.Sprintf("Tip support should be -1.00 and is %.2f", e.Support()))
		}
	}
}

func TestClassicalSupport_2(t *testing.T) {
	var reftree *tree.Tree
	var treefile, reftreefile io.Closer
	var treereader, reftreereader *bufio.Reader
	var err error
	var trees <-chan tree.Trees

	// Parsing multi tree newick (compared trees
	if treefile, treereader, err = utils.GetReader("data/rand_tree_same.nw.gz"); err != nil {
		t.Error(err)
	}
	defer treefile.Close()
	trees = utils.ReadMultiTrees(treereader, utils.FORMAT_NEWICK)

	// Parsing single tree newick file
	if reftreefile, reftreereader, err = utils.GetReader("data/rand_tree.nw.gz"); err != nil {
		t.Error(err)
	}
	defer reftreefile.Close()

	reftree, err = newick.NewParser(reftreereader).Parse()
	if err != nil {
		t.Error(err)
	}

	// Computing fbp
	err = support.Classical(reftree, trees, 1)
	if err != nil {
		t.Error(err)
	}

	for _, e := range reftree.Edges() {
		if !e.Right().Tip() && e.Support() != 1.00 {
			t.Error(fmt.Sprintf("Non Tip support should be 1.00 and is %.2f", e.Support()))
		} else if e.Right().Tip() && e.Support() != -1 {
			t.Error(fmt.Sprintf("Tip support should be -1.00 and is %.2f", e.Support()))
		}
	}
}

func TestClassicalSupport_3(t *testing.T) {
	var reftree *tree.Tree
	var treefile, reftreefile io.Closer
	var treereader, reftreereader *bufio.Reader
	var err error
	var trees <-chan tree.Trees

	// Parsing multi tree newick (compared trees
	if treefile, treereader, err = utils.GetReader("data/rand_tree_half_same.nw.gz"); err != nil {
		t.Error(err)
	}
	defer treefile.Close()
	trees = utils.ReadMultiTrees(treereader, utils.FORMAT_NEWICK)

	// Parsing single tree newick file
	if reftreefile, reftreereader, err = utils.GetReader("data/rand_tree.nw.gz"); err != nil {
		t.Error(err)
	}
	defer reftreefile.Close()
	reftree, err = newick.NewParser(reftreereader).Parse()
	if err != nil {
		t.Error(err)
	}

	// Computing fbp
	err = support.Classical(reftree, trees, 1)
	if err != nil {
		t.Error(err)
	}

	for _, e := range reftree.Edges() {
		if !e.Right().Tip() && e.Support() != 0.50 {
			t.Error(fmt.Sprintf("Non Tip support should be 0.50 and is %.2f", e.Support()))
		} else if e.Right().Tip() && e.Support() != -1 {
			t.Error(fmt.Sprintf("Tip support should be -1.00 and is %.2f", e.Support()))
		}
	}
}

func TestClassicalSupport_4(t *testing.T) {
	var reftree *tree.Tree
	var treefile, reftreefile io.Closer
	var treereader, reftreereader *bufio.Reader
	var err error
	var trees <-chan tree.Trees

	// Parsing multi tree newick (compared trees
	if treefile, treereader, err = utils.GetReader("data/rand_tree_quarter_same.nw.gz"); err != nil {
		t.Error(err)
	}
	defer treefile.Close()
	trees = utils.ReadMultiTrees(treereader, utils.FORMAT_NEWICK)

	// Parsing single tree newick file
	if reftreefile, reftreereader, err = utils.GetReader("data/rand_tree.nw.gz"); err != nil {
		t.Error(err)
	}
	defer reftreefile.Close()
	reftree, err = newick.NewParser(reftreereader).Parse()
	if err != nil {
		t.Error(err)
	}

	// Computing fbp
	err = support.Classical(reftree, trees, 1)
	if err != nil {
		t.Error(err)
	}

	for _, e := range reftree.Edges() {
		if !e.Right().Tip() && e.Support() != 0.25 {
			t.Error(fmt.Sprintf("Non Tip support should be 0.25 and is %.2f", e.Support()))
		} else if e.Right().Tip() && e.Support() != -1 {
			t.Error(fmt.Sprintf("Tip support should be -1.00 and is %.2f", e.Support()))
		}
	}
}
