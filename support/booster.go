package support

import (
	"container/list"
	"errors"
	"fmt"
	"os"
	"sync"

	"github.com/evolbioinfo/gotree/io"
	"github.com/evolbioinfo/gotree/tree"
)

type BoosterSupporter struct {
	currentTree                    int
	mutex                          *sync.RWMutex
	stop                           bool
	silent                         bool
	computeMovedSpecies            bool
	computeTransferPerBranches     bool // All  transfered taxa per branches are computed
	computeHighTransferPerBranches bool // Only highest transfered taxa per branches are computed

	/* Cutoff for considering a branch : ok if norm distance to current
	bootstrap tree < cutoff (ex=0.05) => allows to compute a minimum depth also:
	norm_dist = distance / (p-1)
	=> If we want at least one species that moves (distance=1) at a given norm_dist cutoff, we need a depth p :
	p=(1/norm_dist)+1
	*/
	movedSpeciesCutoff float64
	// If false, then we do not normalize by the expected value: 1-avg/expected
	normalizeByExpected bool
}

// computeTransferPerBranches and computeHighTransferPerBranches are mutually expclusive
// computeTransferPerBranches has priority over computeHighTransferPerBranches
func NewBoosterSupporter(silent, computeMovedSpecies, computeTransferPerBranches, computeHighTransferPerBranches bool, movedSpeciesCutoff float64, normalizeByExpected bool) *BoosterSupporter {
	if movedSpeciesCutoff < 0 {
		movedSpeciesCutoff = 1.0
	}
	if movedSpeciesCutoff > 1 {
		movedSpeciesCutoff = 1.0
	}

	return &BoosterSupporter{
		currentTree:                    0,
		mutex:                          &sync.RWMutex{},
		stop:                           false,
		silent:                         silent,
		computeMovedSpecies:            computeMovedSpecies,
		movedSpeciesCutoff:             movedSpeciesCutoff,
		normalizeByExpected:            normalizeByExpected,
		computeTransferPerBranches:     computeTransferPerBranches,
		computeHighTransferPerBranches: computeHighTransferPerBranches,
	}
}

func (supporter *BoosterSupporter) ExpectedRandValues(depth int) float64 {
	return float64(depth - 1)
}

func (supporter *BoosterSupporter) NormalizeByExpected() bool {
	return supporter.normalizeByExpected
}

func (supporter *BoosterSupporter) NewBootTreeComputed() {
	supporter.mutex.Lock()
	supporter.currentTree++
	supporter.mutex.Unlock()
}

func (supporter *BoosterSupporter) Progress() int {
	supporter.mutex.RLock()
	defer supporter.mutex.RUnlock()
	return supporter.currentTree
}

func (supporter *BoosterSupporter) PrintMovingTaxa() bool {
	return supporter.computeMovedSpecies
}

func (supporter *BoosterSupporter) PrintTaxPerBranches() bool {
	return supporter.computeTransferPerBranches
}

func (supporter *BoosterSupporter) PrintHighTaxPerBranches() bool {
	return supporter.computeHighTransferPerBranches
}

func (supporter *BoosterSupporter) Cancel() {
	supporter.stop = true
}

func (supporter *BoosterSupporter) Canceled() bool {
	return supporter.stop
}

func (supporter *BoosterSupporter) Init(maxdepth int, nbtips int) {
	supporter.stop = false
	supporter.mutex = &sync.RWMutex{}
	supporter.currentTree = 0
}

func minTransferDist(refEdge *tree.Edge, refTree, bootTree *tree.Tree, ntips int, bootEdges []*tree.Edge) (dist int) {
	numBootEdges := len(bootEdges)
	ones := make([]int, numBootEdges)
	//fmt.Fprintf(os.Stderr, "r=%s\n", refEdge.DumpBitSet())
	p, _ := refEdge.TopoDepth()
	dist = p - 1
	minTransferDistRecur(refTree, ntips, bootTree.Root(), nil, nil, refEdge, p, ones, &dist)
	//fmt.Fprintf(os.Stderr, "Final d=%d\n", dist)
	return
}

func minTransferDistRecur(refTree *tree.Tree, ntips int, cur, prev *tree.Node, curEdge, refEdge *tree.Edge, p int, ones []int, dist *int) (stop bool) {
	if cur.Tip() {
		tipIndex, _ := refTree.TipIndex(cur.Name())
		light := refEdge.TipPresent(tipIndex)
		if r, _ := refEdge.NumTipsRight(); r > ntips/2 {
			light = !light
		}
		if light {
			ones[curEdge.Id()] = 0
		} else {
			ones[curEdge.Id()] = 1
		}
	} else {
		curOnes := 0
		for i, n := range cur.Neigh() {
			if n != prev {
				nextEdge := cur.Edges()[i]
				if minTransferDistRecur(refTree, ntips, n, cur, nextEdge, refEdge, p, ones, dist) {
					return true
				}
				curOnes += ones[nextEdge.Id()]
			}
		}
		if curEdge != nil {
			//fmt.Fprintf(os.Stderr, "b=%s\n", curEdge.DumpBitSet())
			ones[curEdge.Id()] = curOnes
			r, _ := curEdge.NumTipsRight()
			zero := r - curOnes
			d := p - zero + ones[curEdge.Id()]
			//fmt.Fprintf(os.Stderr, "d=%d\n", d)
			if d > ntips/2 {
				d = ntips - d
			}
			//fmt.Fprintf(os.Stderr, "d=%d\n", d)
			if d < *dist {
				*dist = d
				if d == 1 {
					return true
				}
			}
		}
		//fmt.Fprintf(os.Stderr, "%d\n", curOnes)
	}

	return false
}

// Thread that takes bootstrap trees from the channel,
// computes the transfer dist for each edges of the ref tree
// and send it to the result channel
// At the end, returns the number of treated trees
func (supporter *BoosterSupporter) ComputeValue(refTree *tree.Tree, cpu int, edges []*tree.Edge,
	bootTreeChannel <-chan tree.Trees, valChan chan<- bootval, speciesChannel chan<- speciesmoved,
	taxPerBranchChannel chan<- []*list.List) error {
	tips := refTree.Tips()
	//var movedSpecies []int = make([]int, len(tips))
	// List of moved species per reference branches
	// It is initialized at each new bootstrap tree
	// It is to the taxperbranch channel reciever
	// to clear all lists
	var taxaTransferedPerBranch []*list.List

	//vals := make([]int, len(edges))
	// Number of branches that have a normalized similarity (1- (min_dist/(n-1)) to
	// bootstrap trees > 0.8
	//var nb_branches_close int
	var err error

	for bootTree := range bootTreeChannel {
		if bootTree.Err != nil {
			io.LogError(bootTree.Err)
			err = bootTree.Err
		} else {
			bootTree.Tree.ReinitIndexes()
			err = refTree.CompareTipIndexes(bootTree.Tree)
			if err == nil {
				//nb_branches_close = 0
				if !supporter.silent {
					fmt.Fprintf(os.Stderr, "CPU : %02d - Bootstrap tree %d\r", cpu, bootTree.Id)
				}
				bootEdges := bootTree.Tree.Edges()
				bootEdgeIndex := tree.NewEdgeIndex(int64(len(bootEdges)*2), 0.75)
				taxaTransferedPerBranch = make([]*list.List, len(edges))

				for i, e := range bootEdges {
					e.SetId(i)
					if !e.Right().Tip() {
						e.Right().SetName("")
					}
					if !e.Left().Tip() {
						e.Left().SetName("")
					}
					bootEdgeIndex.PutEdgeValue(e, i, e.Length())
				}

				for i, e := range edges {
					if e.Right().Tip() {
						taxaTransferedPerBranch[i] = list.New()
						continue
					}

					if _, ok := bootEdgeIndex.Value(e); ok {
						valChan <- bootval{
							0,
							i,
							false,
						}
					} else if p, _ := e.TopoDepth(); p == 2 {
						valChan <- bootval{
							1,
							i,
							false,
						}
					} else {
						valChan <- bootval{
							minTransferDist(e, refTree, bootTree.Tree, len(tips), bootEdges),
							i,
							false,
						}
					}

					// vals[i] = int(min_dist[i])
					// if supporter.computeMovedSpecies || supporter.computeTransferPerBranches || supporter.computeHighTransferPerBranches {
					// 	td, err := e.TopoDepth()
					// 	if err != nil {
					// 		io.LogError(err)
					// 		return err
					// 	}
					// 	be := bootEdges[min_dist_edge[i]]
					// 	norm := float64(vals[i]) / (float64(td) - 1.0)
					// 	mindepth := int(math.Ceil(1.0/supporter.movedSpeciesCutoff + 1.0))
					// 	if sm, er := speciesToMove(e, be, vals[i]); er != nil {
					// 		io.LogError(er)
					// 		return er
					// 	} else {
					// 		if supporter.computeMovedSpecies && norm <= supporter.movedSpeciesCutoff && td >= mindepth {
					// 			for e := sm.Front(); e != nil; e = e.Next() {
					// 				movedSpecies[e.Value.(uint)]++
					// 			}
					// 			nb_branches_close++
					// 		}
					// 		if supporter.computeTransferPerBranches || supporter.computeHighTransferPerBranches {
					// 			// The list of taxons that move around branch i in that bootstrap tree
					// 			taxaTransferedPerBranch[i] = sm
					// 		} else {
					// 			sm.Init() // Clear List
					// 		}
					// 	}
					//}
				}

				// if supporter.computeMovedSpecies {
				// 	for sp, nb := range movedSpecies {
				// 		speciesChannel <- speciesmoved{
				// 			uint(sp),
				// 			float64(nb) / float64(nb_branches_close),
				// 		}
				// 		movedSpecies[sp] = 0
				// 	}
				// }
				// if supporter.computeTransferPerBranches || supporter.computeHighTransferPerBranches {
				// 	taxPerBranchChannel <- taxaTransferedPerBranch
				// }
				supporter.NewBootTreeComputed()
				if supporter.stop {
					break
				}
			}
		}
		bootTree.Tree.Delete()
	}
	return err
}

func Booster(reftree *tree.Tree, boottrees <-chan tree.Trees, logfile *os.File, silent, computeMovedSpecies, computeTransferPerBranches, computeHighTransferPerBranches bool, movedSpeciedCutoff float64, normalizedByExpected bool, cpus int) error {
	var supporter *BoosterSupporter = NewBoosterSupporter(silent, computeMovedSpecies, computeTransferPerBranches, computeHighTransferPerBranches, movedSpeciedCutoff, normalizedByExpected)
	return ComputeSupport(reftree, boottrees, logfile, cpus, supporter)
}

// Returns the list of species to move to go from one branch to the other
// Its length should correspond to given dist
// If not, exit with an error
func speciesToMove(e, be *tree.Edge, dist int) (*list.List, error) {
	var i uint
	diff := list.New()
	equ := list.New()

	for i = 0; i < e.Bitset().Len(); i++ {
		if e.Bitset().Test(i) != be.Bitset().Test(i) {
			diff.PushBack(i)
		} else {
			equ.PushBack(i)
		}
	}
	if diff.Len() < equ.Len() {
		if diff.Len() != dist {
			er := errors.New(fmt.Sprintf("Length of moved species array (%d) is not equal to the minimum distance found (%d)", diff.Len(), dist))
			io.LogError(er)
			return nil, er
		}
		equ.Init()
		return diff, nil
	}
	if equ.Len() != dist {
		er := errors.New(fmt.Sprintf("Length of moved species array (%d) is not equal to the minimum distance found (%d)", equ.Len(), dist))
		io.LogError(er)
		return nil, er
	}
	diff.Init()
	return equ, nil
}

// This function writes on the child node name the string: "branch_id|avg_dist|depth"
// and removes support information from each branch
func ReformatAvgDistance(t *tree.Tree) {
	for i, e := range t.Edges() {
		if e.Support() != tree.NIL_SUPPORT {
			td, _ := e.TopoDepth()
			e.Right().SetName(fmt.Sprintf("%d|%s|%d", i, e.SupportString(), td))
			e.SetSupport(tree.NIL_SUPPORT)
		}
	}
}

// This function takes all branch support values (that are considered as average
// transfer distances over bootstrap trees), normalizes them by the depth and
// convert them to similarity, i.e:
//     1-avg_dist/(depth-1)
func NormalizeTransferDistancesByDepth(t *tree.Tree) {
	for _, e := range t.Edges() {
		avgdist := e.Support()
		if avgdist != tree.NIL_SUPPORT {
			td, _ := e.TopoDepth()
			e.SetSupport(float64(1) - avgdist/float64(td-1))
		}
	}
}
