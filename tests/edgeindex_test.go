package tests

import (
	"fmt"
	"github.com/fredericlemoine/gotree/io/newick"
	"github.com/fredericlemoine/gotree/io/utils"
	"github.com/fredericlemoine/gotree/tree"
	"os"
	"strings"
	"sync"
	"testing"
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
			if val != i || !ok {
				t.Error(fmt.Sprintf("Edge Count must be == %d", i))
			}
		}
	}
}

func TestEdgeIndex2(t *testing.T) {

	intrees := make(chan utils.Trees, 15)
	/* Read ref tree(s) */
	go func() {
		if _, err := utils.ReadCompTrees("data/twotrees.nw", intrees); err != nil {
			t.Error(err.Error)
		}
	}()

	edgeindex := tree.NewEdgeIndex(128, .75)
	nbtrees := 1
	numLoops := 10
	for tr := range intrees {
		edges := tr.Tree.Edges()
		for i := 1; i <= numLoops; i++ {
			for _, e := range edges {
				edgeindex.AddEdgeCount(e)
				val, ok := edgeindex.Value(e)
				if !ok {
					t.Error(fmt.Sprintf("Edge not found in the index"))
				} else if !e.Right().Tip() && val != i {
					t.Error(fmt.Sprintf("Non tip edge count must be == %d (actually %d)", i, val))
				} else if e.Right().Tip() && val != (nbtrees-1)*numLoops+i {
					t.Error(fmt.Sprintf("Tip edge count must be == %d (actually %d)", (nbtrees-1)*numLoops+i, val))
				}
			}
		}
		nbtrees++
	}
}

// Benchmark for reading 1000 bootstrap trees for example
// With gunzip
func BenchmarkEdgeIndex(b *testing.B) {

	for n := 0; n < b.N; n++ {

		reftree, err := utils.ReadRefTree("data/ref_16_14588.nw.gz")
		if err != nil {
			b.Error(err.Error)
		}

		intrees := make(chan utils.Trees, 15)
		/* Read ref tree(s) */
		go func() {
			if _, err := utils.ReadCompTrees("data/boot_16_14588.nw.gz", intrees); err != nil {
				b.Error(err.Error)
			}
		}()
		var wg sync.WaitGroup
		edgeindex := tree.NewEdgeIndex(24000, .75)
		for cpu := 0; cpu < 4; cpu++ {
			wg.Add(1)
			go func() {
				for tr := range intrees {
					edges := tr.Tree.Edges()
					for _, e := range edges {
						edgeindex.AddEdgeCount(e)
					}
				}
				wg.Done()
			}()
		}
		wg.Wait()

		for i, e := range reftree.Edges() {
			val, _ := edgeindex.Value(e)
			fmt.Fprintf(os.Stderr, "Edge %d, Support: %d\n", i, val)
		}
	}
}
