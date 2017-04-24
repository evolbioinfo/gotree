package support

import (
	"runtime"
	"sync"
	"sync/atomic"

	"github.com/fredericlemoine/gotree/tree"
)

/*
Computes bootstrap supports of reftree branches, given trees in boottrees channel
*/
func Classical(reftree *tree.Tree, boottrees <-chan tree.Trees, cpus int) error {
	maxcpus := runtime.NumCPU()
	if cpus > maxcpus {
		cpus = maxcpus
	}
	edges := reftree.Edges()
	var ntrees int32 = 0
	foundEdges := make(chan int, 100)
	foundBoot := make([]int, len(edges))
	edgeIndex := tree.NewEdgeIndex(int64(len(edges)*2), 0.75)
	for i, e := range edges {
		if !e.Right().Tip() {
			e.Right().SetName("")
		}
		if !e.Left().Tip() {
			e.Left().SetName("")
		}
		edgeIndex.PutEdgeValue(e, i, e.Length())
	}
	var wg sync.WaitGroup
	var err error
	for cpu := 0; cpu < cpus; cpu++ {
		wg.Add(1)
		go func(cpu int) {
			for treeV := range boottrees {
				if treeV.Err != nil {
					err = treeV.Err
				} else {
					err = reftree.CompareTipIndexes(treeV.Tree)
					if err == nil {
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

type ClassicalSupporter struct {
	currentTree int
	mutex       *sync.RWMutex
	stop        bool
	silent      bool
}

func NewClassicalSupporter(silent bool) *ClassicalSupporter {
	return &ClassicalSupporter{
		currentTree: 0,
		mutex:       &sync.RWMutex{},
		stop:        false,
		silent:      silent,
	}
}

func (supporter *ClassicalSupporter) NormalizeByExpected() bool {
	return false
}

func (supporter *ClassicalSupporter) ExpectedRandValues(depth int) float64 {
	return 0
}

func (supporter *ClassicalSupporter) NewBootTreeComputed() {
	supporter.mutex.Lock()
	supporter.currentTree++
	supporter.mutex.Unlock()
}

func (supporter *ClassicalSupporter) Progress() int {
	supporter.mutex.RLock()
	defer supporter.mutex.RUnlock()
	return supporter.currentTree
}
func (supporter *ClassicalSupporter) PrintMovingTaxa() bool {
	return false
}

func (supporter *ClassicalSupporter) Cancel() {
	supporter.stop = true
}
func (supporter *ClassicalSupporter) Canceled() bool {
	return supporter.stop
}

func (supporter *ClassicalSupporter) Init(maxdepth int, nbtips int) {
	supporter.stop = false
	supporter.mutex = &sync.RWMutex{}
	supporter.currentTree = 0
}

// Thread that takes bootstrap trees from the channel,
// computes the transfer dist for each edges of the ref tree
// and send it to the result channel
func (supporter *ClassicalSupporter) ComputeValue(refTree *tree.Tree, cpu int, edges []*tree.Edge,
	bootTreeChannel <-chan tree.Trees, valChan chan<- bootval, speciesChannel chan<- speciesmoved) error {

	edgeIndex := tree.NewEdgeIndex(int64(len(edges)*2), 0.75)
	for i, e := range edges {
		if !e.Right().Tip() {
			e.Right().SetName("")
		}
		if !e.Left().Tip() {
			e.Left().SetName("")
		}
		edgeIndex.PutEdgeValue(e, i, e.Length())
	}
	var err error

	for treeV := range bootTreeChannel {
		if treeV.Err != nil {
			err = treeV.Err
		} else {
			err = refTree.CompareTipIndexes(treeV.Tree)
			if err == nil {
				edges2 := treeV.Tree.Edges()
				for _, e2 := range edges2 {
					if !e2.Right().Tip() {
						val, ok := edgeIndex.Value(e2)
						if ok {
							valChan <- bootval{
								1,
								val.Count,
								false,
							}
						}
					}
				}
			}
		}
		supporter.NewBootTreeComputed()
		if supporter.stop {
			break
		}
	}
	return err
}
