package support

import (
	"bufio"
	"errors"
	"os"
	"runtime"
	"sync"
	"sync/atomic"

	"github.com/fredericlemoine/gotree/io"
	"github.com/fredericlemoine/gotree/io/utils"
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

func ClassicalFile(reftreefile, boottreefile string, cpus int) (*tree.Tree, error) {
	var reftree *tree.Tree
	var err error
	var bootfile *os.File
	var bootreader *bufio.Reader
	var compareChannel <-chan tree.Trees

	maxcpus := runtime.NumCPU()
	if cpus > maxcpus {
		cpus = maxcpus
	}

	if reftree, err = utils.ReadTree(reftreefile); err != nil {
		io.LogError(err)
		return nil, err
	}

	if boottreefile == "none" {
		er := errors.New("You must provide a file containing bootstrap trees")
		io.LogError(er)
		return nil, er
	}

	if bootfile, bootreader, err = utils.GetReader(boottreefile); err != nil {
		return nil, err
	}
	defer bootfile.Close()

	compareChannel = utils.ReadMultiTrees(bootreader)

	if err = Classical(reftree, compareChannel, cpus); err != nil {
		return nil, err
	}
	return reftree, nil
}
