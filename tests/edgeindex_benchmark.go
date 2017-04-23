package tests

import (
	"bufio"
	"os"
	"sync"
	"testing"

	"github.com/fredericlemoine/gotree/io/utils"
	"github.com/fredericlemoine/gotree/tree"
)

// Benchmark for reading 1000 bootstrap trees for example
// With gunzip
func BenchmarkEdgeIndex(b *testing.B) {

	for n := 0; n < b.N; n++ {
		var treefile *os.File
		var treereader *bufio.Reader
		var intrees <-chan tree.Trees

		reftree, err := utils.ReadTree("data/benchmark_ref.nw.gz")
		if err != nil {
			b.Error(err.Error)
		}

		if treefile, treereader, err = utils.GetReader("data/benchmark_boot.nw.gz"); err != nil {
			b.Error(err)
		}
		defer treefile.Close()
		intrees = utils.ReadMultiTrees(treereader)

		var wg sync.WaitGroup
		edgeindex := tree.NewEdgeIndex(24000, .75)
		for cpu := 0; cpu < 4; cpu++ {
			wg.Add(1)
			go func() {
				for tr := range intrees {
					if tr.Err != nil {
						b.Error(tr.Err)
					}
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
