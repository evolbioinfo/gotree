package support

import (
	"container/list"
	"fmt"
	"os"
	"runtime"
	"sync"

	"github.com/fredericlemoine/gotree/io"
	"github.com/fredericlemoine/gotree/sort"
	"github.com/fredericlemoine/gotree/tree"
)

type bootval struct {
	value   int
	edgeid  int
	randsup bool
}

type speciesmoved struct {
	taxid   uint
	nbtimes float64
}

type Supporter interface {
	Init(maxdepth int, nbtips int)
	ExpectedRandValues(depth int) float64
	ComputeValue(refTree *tree.Tree, cpu int, edges []*tree.Edge,
		bootTreeChannel <-chan tree.Trees, valChan chan<- bootval, speciesChannel chan<- speciesmoved,
		taxPerBranchChannel chan<- []*list.List) error
	// Returns the number of bootstrap trees that have been computed
	Progress() int
	// Increments the number of trees processed
	NewBootTreeComputed()
	// Tells the supported to stop accepting new bootstrap trees from the bootTreeChannel
	// It will just finish the current computations
	Cancel()
	// Tells if hasbeen canceled or not
	Canceled() bool
	// Print moving taxa
	PrintMovingTaxa() bool
	PrintTaxPerBranches() bool
	PrintHighTaxPerBranches() bool

	// If true, the support is 1-avg_value/ExpectedRandValues()
	// If false, the support is avg
	NormalizeByExpected() bool
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

func ComputeSupport(reftree *tree.Tree, boottrees <-chan tree.Trees, logfile *os.File, cpus int, supporter Supporter) error {
	var deptherr error   // error reading comp file
	var computeerr error // error in support computation

	var maxcpus int = runtime.NumCPU() // max number of cpus
	var max_depth int                  // Maximum topo depth of all edges of ref tree
	var tips []*tree.Node              // Tip nodes of the ref tree
	var edges []*tree.Edge             // Edges of the reference tree
	var valuesBoot []int               // Sum of number of bootValues per edge over boot trees
	var speciesMovedCount []float64    // Number of times each species has been moved (for booster mainly)
	var taxPerBranches [][]int         // Number of times (bs tree) each species has been moved around each branch (for booster mainly)

	var wg sync.WaitGroup  // For waiting end of step computation
	var wg2 sync.WaitGroup // For waiting end of final counting

	var valuesChan chan bootval               // Channel of values computed for a given edge
	var speciesChannel chan speciesmoved      // Channel of number of times each species has been moved. Used only for booster support
	var taxPerBranchChannel chan []*list.List // Channel of Lists of moved species per reference branches compared to 1 bs tree

	valuesChan = make(chan bootval, 100)

	speciesChannel = make(chan speciesmoved, 100)
	taxPerBranchChannel = make(chan []*list.List, 100)

	if cpus > maxcpus {
		cpus = maxcpus
	}
	reftree.ReinitIndexes()
	tips = reftree.Tips()
	edges = reftree.Edges()
	if max_depth, deptherr = maxDepth(edges); deptherr != nil {
		return deptherr
	}

	valuesBoot = make([]int, len(edges))
	speciesMovedCount = make([]float64, len(tips))
	taxPerBranches = make([][]int, len(edges))

	//Initialize supporter
	supporter.Init(max_depth, len(tips))

	// Assign an id to every edge
	// And init taxPerBranches
	for i, e := range edges {
		e.SetId(i)
		taxPerBranches[i] = make([]int, len(tips))
	}

	// We compute value for each bootstrap tree
	for cpu := 0; cpu < cpus; cpu++ {
		wg.Add(1)
		curcpu := cpu
		go func() {
			if err := supporter.ComputeValue(reftree, curcpu, edges, boottrees, valuesChan, speciesChannel, taxPerBranchChannel); err != nil {
				io.LogError(err)
				computeerr = err
			}
			wg.Done()
		}()
	}

	// Wait for step computation to close output channels
	go func() {
		wg.Wait()
		// Read remaining trees from boottreeChannel if computations have been stoped before
		for _ = range boottrees {
		}
		close(valuesChan)
		close(speciesChannel)
		close(taxPerBranchChannel)
	}()

	// Now count Values from the output channels
	wg2.Add(3)
	go func() {
		for val := range valuesChan {
			valuesBoot[val.edgeid] += val.value
		}
		wg2.Done()
	}()

	// We gather all species moves numbers per bootstrap trees
	go func() {
		for val := range speciesChannel {
			speciesMovedCount[val.taxid] += val.nbtimes
		}
		wg2.Done()
	}()

	// We gather all taxperbranches lists coming from all bootstrap trees
	go func() {
		for taxlists := range taxPerBranchChannel {
			for i, l := range taxlists {
				for e := l.Front(); e != nil; e = e.Next() {
					taxPerBranches[i][e.Value.(uint)]++
				}
				l.Init()
			}
		}
		wg2.Done()
	}()

	wg2.Wait()

	if computeerr != nil {
		io.LogError(computeerr)
		return computeerr
	}

	names := reftree.SortedTips()

	if supporter.PrintMovingTaxa() {
		logfile.WriteString("Taxon : Instability\n")
		for i, n := range speciesMovedCount {
			logfile.WriteString(fmt.Sprintf("%s : %f\n", names[i], n*100.0/float64(supporter.Progress())))
		}
	}
	if supporter.PrintTaxPerBranches() {
		logfile.WriteString("Edge\tDepth\tAvgDist")
		for sp := 0; sp < len(tips); sp++ {
			logfile.WriteString("\t" + names[sp])
		}
		logfile.WriteString("\n")
	} else if supporter.PrintHighTaxPerBranches() {
		logfile.WriteString("Edge\tDepth\tAvgDist\tHighlyTransferedTaxa\n")
	}

	// Finally we compute support and write it in the tree
	for i, e := range edges {
		if !edges[i].Right().Tip() {
			d, err := e.TopoDepth()
			if err != nil {
				io.LogError(err)
				return err
			}
			support := float64(valuesBoot[i]) / float64(supporter.Progress())
			if supporter.PrintTaxPerBranches() {
				logfile.WriteString(fmt.Sprintf("%d\t%d\t%.6f", i, d, support))
				for sp := 0; sp < len(tips); sp++ {
					logfile.WriteString(fmt.Sprintf("\t%.6f", float64(taxPerBranches[i][sp])/float64(supporter.Progress())))
				}
				logfile.WriteString("\n")
			} else if supporter.PrintHighTaxPerBranches() {
				logfile.WriteString(fmt.Sprintf("%d\t%d\t%.6f", i, d, support))
				orderedtipsbytransferindex := sort.OrderInt(taxPerBranches[i], true)
				prevtindex := 0
				for ind, sp := range orderedtipsbytransferindex {
					if taxPerBranches[i][sp] > 0.0 && (ind <= int(support+0.5) || taxPerBranches[i][sp] == prevtindex) {
						if ind > 0 {
							logfile.WriteString(";")
						} else {
							logfile.WriteString("\t")
						}
						logfile.WriteString(fmt.Sprintf("%s(%.6f)", names[sp], float64(taxPerBranches[i][sp])/float64(supporter.Progress())))
						prevtindex = taxPerBranches[i][sp]
					} else {
						prevtindex = -1
					}
				}
				logfile.WriteString("\n")
			}
			if supporter.NormalizeByExpected() {
				support = float64(1) - support/supporter.ExpectedRandValues(d)
			}

			edges[i].SetSupport(support)
		}
	}

	return nil
}

func maxDepth(edges []*tree.Edge) (int, error) {
	max_depth := 0
	for _, e := range edges {
		var d int
		var err error
		if d, err = e.TopoDepth(); err != nil {
			io.LogError(err)
			return d, err
		}
		max_depth = max(d, max_depth)
	}
	return max_depth, nil
}
