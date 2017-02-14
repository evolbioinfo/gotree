package support

import (
	"errors"
	"github.com/fredericlemoine/gotree/io"
	"github.com/fredericlemoine/gotree/io/utils"
	"github.com/fredericlemoine/gotree/tree"
	"runtime"
	"sync"
)

func Classical(reftreefile, boottreefile string, cpus int) (*tree.Tree, error) {
	var reftree *tree.Tree
	var err error
	var readerr error

	maxcpus := runtime.NumCPU()
	if cpus > maxcpus {
		cpus = maxcpus
	}

	if reftree, err = utils.ReadRefTree(reftreefile); err != nil {
		io.LogError(err)
		return nil, err
	}

	if boottreefile == "none" {
		er := errors.New("You must provide a file containing bootstrap trees")
		io.LogError(er)
		return nil, er
	}

	if cpus > maxcpus {
		cpus = maxcpus
	}

	var nbtrees int
	compareChannel := make(chan tree.Trees, 15)
	go func() {
		if nbtrees, readerr = utils.ReadCompTrees(boottreefile, compareChannel); err != nil {
			io.LogError(readerr)
		}
		close(compareChannel)
	}()

	edges := reftree.Edges()
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
	for cpu := 0; cpu < cpus; cpu++ {
		wg.Add(1)
		go func(cpu int) {
			for treeV := range compareChannel {
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
			edges[i].SetSupport(float64(count) / float64(nbtrees))
		}
	}
	if readerr != nil {
		return nil, readerr
	}
	return reftree, nil
}
