package tests

import (
	"github.com/fredericlemoine/gotree/io/utils"
	"github.com/fredericlemoine/gotree/tree"
	"sync"
	"testing"
)

// Benchmark for reading 1000 bootstrap trees for example
// With gunzip
func BenchmarkEdgeIndex(b *testing.B) {

	for n := 0; n < b.N; n++ {

		reftree, err := utils.ReadRefTree("data/benchmark_ref.nw.gz")
		if err != nil {
			b.Error(err.Error)
		}

		intrees := make(chan tree.Trees, 15)
		/* Read ref tree(s) */
		go func() {
			if _, err := utils.ReadCompTrees("data/benchmark_boot.nw.gz", intrees); err != nil {
				b.Error(err.Error)
			}
			close(intrees)
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

		for _, e := range reftree.Edges() {
			edgeindex.Value(e)
		}
	}
}
