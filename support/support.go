package support

import (
	"fmt"
	"github.com/fredericlemoine/gotree/io/utils"
	"github.com/fredericlemoine/gotree/tree"
	"runtime"
	"sync"
)

type bootval struct {
	value   int
	edgeid  int
	randsup bool
}

type Supporter interface {
	ExpectedRandValues(maxdepth int, nbtips int) []float64
	ComputeValue(refTree *tree.Tree, empiricalTrees []*tree.Tree, cpu int, empirical bool, edges []*tree.Edge, randEdges [][]*tree.Edge,
		wg *sync.WaitGroup, bootTreeChannel <-chan utils.Trees, valChan chan<- bootval, randvalChan chan<- bootval)
}

func ComputeSupport(reftreefile, boottreefile string, empirical bool, cpus int, supporter Supporter) *tree.Tree {
	var reftree *tree.Tree             // reference tree
	var err error                      // error output
	var nbEmpiricalTrees int = 10      // number of empirical trees to simulate
	var maxcpus int = runtime.NumCPU() // max number of cpus
	var randEdges [][]*tree.Edge       // Edges of shuffled trees
	var nbtrees int                    // Number of bootstrap trees
	var max_depth int                  // Maximum topo depth of all edges of ref tree
	var tips []*tree.Node              // Tip nodes of the ref tree
	var edges []*tree.Edge             // Edges of the reference tree
	var valuesBoot []int               // Sum of number of bootValues per edge over boot trees
	var valuesRand []int               // Sum of number of bootValues per random edges over boot trees
	var gtRandom []int                 // Number of times edges have steps that are >= rand steps
	var randTrees []*tree.Tree         // Empirical rand trees

	var wg sync.WaitGroup                // For waiting end of step computation
	var wg2 sync.WaitGroup               // For waiting end of final counting
	var valuesChan chan bootval          // Channel of values computed for a given edge
	var randValuesChan chan bootval      // Channel of values computed for a given shuffled edge
	var bootTreeChannel chan utils.Trees // Channel of bootstrap trees

	valuesChan = make(chan bootval, 1000)
	randValuesChan = make(chan bootval, 1000)
	bootTreeChannel = make(chan utils.Trees, 100)

	if cpus > maxcpus {
		cpus = maxcpus
	}

	if reftree, err = utils.ReadRefTree(reftreefile); err != nil {
		panic(err)
	}

	tips = reftree.Tips()
	edges = reftree.Edges()
	max_depth = maxDepth(edges)
	valuesBoot = make([]int, len(edges))
	gtRandom = make([]int, len(edges))
	valuesRand = make([]int, len(edges))

	// Precomputation of expected number of parsimony steps per depth
	expected_rand_val := supporter.ExpectedRandValues(max_depth, len(tips))

	// We generate nbEmpirical shuffled trees and store their edges
	randEdges = make([][]*tree.Edge, nbEmpiricalTrees)
	randTrees = make([]*tree.Tree, nbEmpiricalTrees)
	if empirical {
		for i := 0; i < nbEmpiricalTrees; i++ {
			var randT *tree.Tree
			if randT, err = utils.ReadRefTree(reftreefile); err != nil {
				panic(err)
			}
			randT.ShuffleTips()
			randEdges[i] = randT.Edges()
			randTrees[i] = randT
		}
	}

	// We read bootstrap trees and put them in the channel
	if boottreefile == "none" {
		panic("You must provide a file containing bootstrap trees")
	}
	go func() {
		if nbtrees, err = utils.ReadCompTrees(boottreefile, bootTreeChannel); err != nil {
			panic(err)
		}
	}()

	// We compute parsimony steps for all bootstrap trees
	for cpu := 0; cpu < cpus; cpu++ {
		wg.Add(1)
		go supporter.ComputeValue(reftree, randTrees, cpu, empirical, edges, randEdges, &wg, bootTreeChannel, valuesChan, randValuesChan)
	}

	// Wait for step computation to close output channels
	go func() {
		wg.Wait()
		close(valuesChan)
		close(randValuesChan)
	}()

	// Now count Values from the output channels
	wg2.Add(2)
	go func() {
		for val := range valuesChan {
			valuesBoot[val.edgeid] += val.value
			d, err := edges[val.edgeid].TopoDepth()
			if err != nil {
				panic(err)
			}
			// If theoretical we must count number >= here
			if !empirical {
				randval := expected_rand_val[d]
				if float64(val.value) >= randval {
					gtRandom[val.edgeid]++
				}
			}
		}
		wg2.Done()
	}()

	// If "empirical" we read the randStepsChan
	if empirical {
		go func() {
			for val := range randValuesChan {
				if val.randsup {
					gtRandom[val.edgeid]++
				}
				valuesRand[val.edgeid] += val.value
			}
			wg2.Done()
		}()
	} else {
		wg2.Done()
	}

	wg2.Wait()

	// Finally we compute pvalues and support and write it in the tree
	for i, e := range edges {
		if !edges[i].Right().Tip() {
			d, err := e.TopoDepth()
			if err != nil {
				panic(err)
			}
			avg_val := float64(valuesBoot[i]) / float64(nbtrees)
			var pval, avg_rand_val float64
			if empirical {
				avg_rand_val = float64(valuesRand[i]) / (float64(nbEmpiricalTrees) * float64(nbtrees))
				pval = float64(gtRandom[i]) * 1.0 / (float64(nbEmpiricalTrees) * float64(nbtrees))
			} else {
				avg_rand_val = expected_rand_val[d]
				pval = float64(gtRandom[i]) * 1.0 / float64(nbtrees)
			}
			support := float64(1) - avg_val/avg_rand_val

			edges[i].SetSupport(support)
			edges[i].Right().SetName(fmt.Sprintf("%.2f/%.4f", support, pval))
		}
	}

	return reftree
}
