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
func FBP(reftree *tree.Tree, boottrees <-chan tree.Trees, cpus int, sup *Supporter) error {
	var err error

	if err = reftree.ReinitIndexes(); err != nil {
		return err
	}

	if sup == nil {
		sup = &Supporter{}
	}

	maxcpus := runtime.NumCPU()
	if cpus > maxcpus {
		cpus = maxcpus
	}
	edges := reftree.Edges()
	var ntrees int32 = 0
	foundEdges := make(chan int, 100)
	foundBoot := make([]int, len(edges))
	for _, e := range edges {
		if !e.Right().Tip() {
			e.Right().SetName("")
		}
		if !e.Left().Tip() {
			e.Left().SetName("")
		}
	}
	var wg sync.WaitGroup
	for cpu := 0; cpu < cpus; cpu++ {
		wg.Add(1)
		go func(cpu int) {
			var inerr error
			for treeV := range boottrees {
				edgeIndex := tree.NewEdgeIndex(uint64(len(edges)*2), 0.75)
				if sup.Canceled() {
					break
				}
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
					for i, e2 := range edges2 {
						if !e2.Right().Tip() {
							if inerr = edgeIndex.PutEdgeValue(e2, i, e2.Length()); inerr != nil {
								err = inerr
								return
							}
						}
					}
					for i, e := range edges {
						_, ok := edgeIndex.Value(e)
						if ok {
							foundEdges <- i
						}
					}
				}
				sup.IncrementProgress()
			}
			wg.Done()
		}(cpu)
	}

	go func() {
		wg.Wait()
		close(foundEdges)
	}()

	for edgeI := range foundEdges {
		foundBoot[edgeI]++
	}

	for i, count := range foundBoot {
		if !edges[i].Right().Tip() {
			//fmt.Printf("%d: %d/%d\n", i, count, ntrees)
			edges[i].SetSupport(float64(count) / float64(ntrees))
		}
	}
	return err
}
