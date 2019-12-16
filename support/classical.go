package support

import (
	"runtime"
	"sync"
	"sync/atomic"

	"github.com/evolbioinfo/gotree/tree"
)

/*
Computes bootstrap supports of reftree branches, given trees in boottrees channel
*/
func Classical(reftree *tree.Tree, boottrees <-chan tree.Trees, cpus int) error {
	var err error

	if err = reftree.ReinitIndexes(); err != nil {
		return err
	}

	maxcpus := runtime.NumCPU()
	if cpus > maxcpus {
		cpus = maxcpus
	}
	edges := reftree.Edges()
	var ntrees int32 = 0
	foundEdges := make(chan int, 100)
	foundBoot := make([]int, len(edges))
	edgeIndex := tree.NewEdgeIndex(uint64(len(edges)*2), 0.75)
	for i, e := range edges {
		if !e.Right().Tip() {
			e.Right().SetName("")
		}
		if !e.Left().Tip() {
			e.Left().SetName("")
		}
		if err = edgeIndex.PutEdgeValue(e, i, e.Length()); err != nil {
			return err
		}
	}
	var wg sync.WaitGroup
	for cpu := 0; cpu < cpus; cpu++ {
		wg.Add(1)
		go func(cpu int) {
			var inerr error
			for treeV := range boottrees {
				if treeV.Err != nil {
					err = treeV.Err
					return
				} else {
					if inerr = treeV.Tree.ReinitIndexes(); err != nil {
						err = inerr
						return
					}
					if inerr = reftree.CompareTipIndexes(treeV.Tree); err != nil {
						err = inerr
						return
					}
					atomic.AddInt32(&ntrees, 1)
					edges2 := treeV.Tree.Edges()
					for _, e2 := range edges2 {
						if !e2.Right().Tip() {
							val, ok := edgeIndex.Value(e2)
							if ok {
								foundEdges <- val.Count
							}
						}
					}
				}
			}
			wg.Done()
		}(cpu)
	}

	go func() {
		wg.Wait()
		close(foundEdges)
	}()

	for edge_i := range foundEdges {
		foundBoot[edge_i]++
	}

	for i, count := range foundBoot {
		if !edges[i].Right().Tip() {
			edges[i].SetSupport(float64(count) / float64(ntrees))
		}
	}
	return err
}
