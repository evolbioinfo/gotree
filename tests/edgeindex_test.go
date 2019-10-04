package tests

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"sync"
	"testing"

	"github.com/evolbioinfo/gotree/io/newick"
	"github.com/evolbioinfo/gotree/io/utils"
	"github.com/evolbioinfo/gotree/tree"
)

func TestEdgeIndex(t *testing.T) {
	var treestring2 string = "(Tip2:1.00000,Node0:1.0000,((Tip7:1.00000,((Tip9:1.00000,Tip6:1.0000):1.0000,(Tip5:1.00000,Tip3:1.0000):1.0000):1.00):1.00,(Tip4:1.00000,(Tip8:1.00000,Tip1:1.000):0.126):0.127):0.125);"

	tr, err2 := newick.NewParser(strings.NewReader(treestring2)).Parse()

	if err2 != nil {
		t.Error(err2)
	}
	edges := tr.Edges()
	edgeindex := tree.NewEdgeIndex(128, .75)

	for i := 1; i <= 10000; i++ {
		for _, e := range edges {
			edgeindex.AddEdgeCount(e)
			val, ok := edgeindex.Value(e)
			if val.Count != i || !ok {
				t.Error(fmt.Sprintf("Edge Count must be == %d", i))
			}
		}
	}
}

func TestEdgeIndex2(t *testing.T) {
	var treefile io.Closer
	var treereader *bufio.Reader
	var err error
	var trees <-chan tree.Trees

	/* File reader (plain text or gzip) */
	if treefile, treereader, err = utils.GetReader("data/twotrees.nw.gz"); err != nil {
		t.Error(err)
	}
	defer treefile.Close()
	trees = utils.ReadMultiTrees(treereader, utils.FORMAT_NEWICK)

	edgeindex := tree.NewEdgeIndex(128, .75)
	nbtrees := 1
	numLoops := 10
	for tr := range trees {
		edges := tr.Tree.Edges()
		for i := 1; i <= numLoops; i++ {
			for _, e := range edges {
				edgeindex.AddEdgeCount(e)
				val, ok := edgeindex.Value(e)
				if !ok {
					t.Error(fmt.Sprintf("Edge not found in the index"))
				} else if !e.Right().Tip() && val.Count != i {
					t.Errorf("Non tip edge count must be == %d (actually %d)", i, val.Count)
				} else if e.Right().Tip() && val.Count != (nbtrees-1)*numLoops+i {
					t.Errorf("Tip edge count must be == %d (actually %d)", (nbtrees-1)*numLoops+i, val.Count)
				}
			}
		}
		nbtrees++
	}
}

// Benchmark for reading 1000 bootstrap trees for example
// With gunzip
func benchmarkEdgeIndex(ntrees int, b *testing.B) {

	for n := 0; n < ntrees; n++ {
		var treefile io.Closer
		var treereader *bufio.Reader
		var intrees <-chan tree.Trees

		reftree, err := utils.ReadTree("data/bootstrap_trees.nw.gz", utils.FORMAT_NEWICK)
		if err != nil {
			b.Error(err.Error())
		}

		if treefile, treereader, err = utils.GetReader("data/bootstrap_trees.nw.gz"); err != nil {
			b.Error(err)
		}
		intrees = utils.ReadMultiTrees(treereader, utils.FORMAT_NEWICK)

		var wg sync.WaitGroup
		edgeindex := tree.NewEdgeIndex(24000, .75)
		for tr := range intrees {
			if tr.Err != nil {
				b.Error(tr.Err)
			}
			edges := tr.Tree.Edges()
			for _, e := range edges {
				edgeindex.AddEdgeCount(e)
			}
		}
		wg.Wait()

		for _, e := range reftree.Edges() {
			edgeindex.Value(e)
		}
		treefile.Close()
	}
}

func BenchmarkEdgeIndex10(b *testing.B)  { benchmarkEdgeIndex(10, b) }
func BenchmarkEdgeIndex100(b *testing.B) { benchmarkEdgeIndex(100, b) }
