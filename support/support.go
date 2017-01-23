package support

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"sort"
	"sync"

	"github.com/fredericlemoine/gotree/io"
	"github.com/fredericlemoine/gotree/io/utils"
	"github.com/fredericlemoine/gotree/tree"
)

type bootval struct {
	value   int
	edgeid  int
	randsup bool
}

type speciesmoved struct {
	taxid   uint
	nbtimes int
}

type Supporter interface {
	Init(maxdepth int, nbtips int)
	ExpectedRandValues(depth int) float64
	ProbaDepthValue(d int, v int) float64
	ComputeValue(refTree *tree.Tree, empiricalTrees []*tree.Tree, cpu int, empirical bool, edges []*tree.Edge, randEdges [][]*tree.Edge,
		wg *sync.WaitGroup, bootTreeChannel <-chan tree.Trees, valChan chan<- bootval, randvalChan chan<- bootval, speciesChannel chan<- speciesmoved)
}

func min(a int, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a int, b int) int {
	if a < b {
		return b
	}
	return a
}

func min_uint(a uint16, b uint16) uint16 {
	if a < b {
		return a
	}
	return b
}
func ComputeSupport(reftreefile, boottreefile string, logfile *os.File, empirical bool, cpus int, supporter Supporter) *tree.Tree {
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
	var gtRandom []float64             // Number of times edges have steps that are >= rand steps
	var randTrees []*tree.Tree         // Empirical rand trees
	var speciesMovedCount []int        // Number of times each species has been moved (for mast mainly)

	var wg sync.WaitGroup  // For waiting end of step computation
	var wg2 sync.WaitGroup // For waiting end of final counting

	var valuesChan chan bootval          // Channel of values computed for a given edge
	var randValuesChan chan bootval      // Channel of values computed for a given shuffled edge
	var bootTreeChannel chan tree.Trees  // Channel of bootstrap trees
	var speciesChannel chan speciesmoved // Channel of number of times each species has been moved. Used only for mast support
	valuesChan = make(chan bootval, 100)
	randValuesChan = make(chan bootval, 100)
	bootTreeChannel = make(chan tree.Trees, 15)
	speciesChannel = make(chan speciesmoved, 100)

	if cpus > maxcpus {
		cpus = maxcpus
	}

	if reftree, err = utils.ReadRefTree(reftreefile); err != nil {
		io.ExitWithMessage(err)
	}

	tips = reftree.Tips()
	edges = reftree.Edges()
	max_depth = maxDepth(edges)
	valuesBoot = make([]int, len(edges))
	gtRandom = make([]float64, len(edges))
	valuesRand = make([]int, len(edges))
	speciesMovedCount = make([]int, len(tips))

	// Precomputation of expected number of parsimony steps per depth
	supporter.Init(max_depth, len(tips))

	// We generate nbEmpirical shuffled trees and store their edges
	randEdges = make([][]*tree.Edge, nbEmpiricalTrees)
	randTrees = make([]*tree.Tree, nbEmpiricalTrees)
	if empirical {
		for i := 0; i < nbEmpiricalTrees; i++ {
			var randT *tree.Tree
			if randT, err = utils.ReadRefTree(reftreefile); err != nil {
				io.ExitWithMessage(err)
			}
			randT.ShuffleTips()
			randEdges[i] = randT.Edges()
			randTrees[i] = randT
			for j, e := range randEdges[i] {
				e.SetId(j)
			}
		}
	}
	for i, e := range edges {
		e.SetId(i)
	}

	// We read bootstrap trees and put them in the channel
	if boottreefile == "none" {
		io.ExitWithMessage(errors.New("You must provide a file containing bootstrap trees"))
	}
	go func() {
		if nbtrees, err = utils.ReadCompTrees(boottreefile, bootTreeChannel); err != nil {
			io.ExitWithMessage(err)
		}
	}()

	// We compute parsimony steps for all bootstrap trees
	for cpu := 0; cpu < cpus; cpu++ {
		wg.Add(1)
		go supporter.ComputeValue(reftree, randTrees, cpu, empirical, edges, randEdges, &wg, bootTreeChannel, valuesChan, randValuesChan, speciesChannel)
	}

	// Wait for step computation to close output channels
	go func() {
		wg.Wait()
		close(valuesChan)
		close(randValuesChan)
		close(speciesChannel)
	}()

	// Now count Values from the output channels
	wg2.Add(3)
	go func() {
		for val := range valuesChan {
			valuesBoot[val.edgeid] += val.value
			d, err := edges[val.edgeid].TopoDepth()
			if err != nil {
				io.ExitWithMessage(err)
			}
			// If theoretical we must count number >= here
			if !empirical {
				for v := 0; v <= val.value; v++ {
					gtRandom[val.edgeid] += supporter.ProbaDepthValue(d, v)
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

	go func() {
		for val := range speciesChannel {
			speciesMovedCount[val.taxid] += val.nbtimes
		}
		wg2.Done()
	}()

	wg2.Wait()

	names := reftree.AllTipNames()
	sort.Strings(names)

	logfile.WriteString("Moving taxa\n")
	for i, n := range speciesMovedCount {
		logfile.WriteString(fmt.Sprintf("%s : %d\n", names[i], n))
	}

	// Finally we compute pvalues and support and write it in the tree
	for i, e := range edges {
		if !edges[i].Right().Tip() {
			d, err := e.TopoDepth()
			if err != nil {
				io.ExitWithMessage(err)
			}
			avg_val := float64(valuesBoot[i]) / float64(nbtrees)
			var pval, avg_rand_val float64
			if empirical {
				avg_rand_val = float64(valuesRand[i]) / (float64(nbEmpiricalTrees) * float64(nbtrees))
				pval = gtRandom[i] * 1.0 / (float64(nbEmpiricalTrees) * float64(nbtrees))
			} else {
				avg_rand_val = supporter.ExpectedRandValues(d)
				pval = gtRandom[i] * 1.0 / float64(nbtrees)
			}
			support := float64(1) - avg_val/avg_rand_val

			edges[i].SetSupport(support)
			edges[i].SetPValue(pval)
		}
	}

	return reftree
}
