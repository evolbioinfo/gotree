package support

import (
	"fmt"
	"math"
	"os"
	"sync"

	"github.com/evolbioinfo/gotree/io"
	"github.com/evolbioinfo/gotree/tree"
)

// This function computes the min transfer distance between the refedge and the bootstrap tree.
// If "absent" is true, then we know that the ref branch is not present in the bootstrap tree (it is faster to compute then), and we stop if dist == 1
// Else: we do not know, then the branch is tested until we fund a dist == 0
func MinTransferDist(refedge *tree.Edge, reftree, boottree *tree.Tree, ntips int, bootedges []*tree.Edge, absent bool) (dist int, minedge *tree.Edge, speciestoadd, speciestoremove []string) {
	numbootedges := len(bootedges)
	ones := make([]int, numbootedges)
	//fmt.Fprintf(os.Stderr, "r=%s\n", refEdge.DumpBitSet())
	p, _ := refedge.TopoDepth()
	dist = p - 1

	stop := false
	minTransferDistRecur(reftree, ntips, boottree.Root(), nil, nil, refedge, p, ones, &dist, &minedge, absent, &stop)

	distcutoff := 0.7
	norm := float64(dist) * 1.0 / (float64(p) - 1.0)
	mindepth := int(math.Ceil(1.0/distcutoff + 1.0))
	//fmt.Printf("Dist : %d, p : %d\n", dist, p)
	if p > mindepth && norm >= distcutoff {
		// computing species to move
		/////////////////////////////////////////
		//fmt.Printf("p= %d, d=%d\n", p, dist)
		n_subtree, _ := minedge.NumTipsRight()
		ones_subtree := ones[minedge.Id()]
		zeros_subtree := n_subtree - ones_subtree
		ones_total := ntips - p
		zeros_total := p

		/* move ones into subtree, zeros outside subtree? */
		ops_ones_in_subtree := zeros_subtree + (ones_total - ones_subtree)
		/* move zeros into subtree, ones outside subtree? */
		ops_zeros_in_subtree := ones_subtree + (zeros_total - zeros_subtree)

		want_ones_outside := ops_zeros_in_subtree < ops_ones_in_subtree

		speciestoadd = make([]string, 0, 10)
		speciestoremove = make([]string, 0, 10)
		speciesToMoveRecursive(minedge, boottree.Root(), nil, nil, ones, want_ones_outside, &speciestoadd, &speciestoremove)
	}
	return
}

func speciesToMoveRecursive(bootedge *tree.Edge, cur, prev *tree.Node, edge *tree.Edge, ones []int, want_ones_now bool, speciestoadd, speciestoremove *[]string) {
	if cur.Tip() {
		if want_ones_now && ones[edge.Id()] == 0 {
			*speciestoadd = append(*speciestoadd, cur.Name())
		}
		if !want_ones_now && ones[edge.Id()] == 1 {
			*speciestoadd = append(*speciestoadd, cur.Name())
		}
	}
	if edge == bootedge {
		want_ones_now = !want_ones_now
	}

	if edge != nil {
		subtreesize, _ := edge.NumTipsRight()
		if (want_ones_now && ones[edge.Id()] == subtreesize) || (!want_ones_now && ones[edge.Id()] == 0) {
			return
		}
	}

	for i, c := range cur.Neigh() {
		if c != prev {
			speciesToMoveRecursive(bootedge, c, cur, cur.Edges()[i], ones, want_ones_now, speciestoadd, speciestoremove)
		}
	}
}

func minTransferDistRecur(refTree *tree.Tree, ntips int, cur, prev *tree.Node, curEdge, refEdge *tree.Edge, p int, ones []int, dist *int, minedge **tree.Edge, absent bool, stop *bool) {
	if *stop {
		return
	}
	curOnes := 0
	if cur.Tip() {
		tipIndex, _ := refTree.TipIndex(cur.Name())
		light := refEdge.TipPresent(tipIndex)
		if r, _ := refEdge.NumTipsRight(); r > ntips/2 {
			light = !light
		}
		if !light {
			curOnes = 1
		}
	} else {
		for i, n := range cur.Neigh() {
			if n != prev {
				nextEdge := cur.Edges()[i]
				minTransferDistRecur(refTree, ntips, n, cur, nextEdge, refEdge, p, ones, dist, minedge, absent, stop)
				curOnes += ones[nextEdge.Id()]
				if *stop {
					return
				}
			}
		}
	}

	if curEdge != nil {
		ones[curEdge.Id()] = curOnes
		r, _ := curEdge.NumTipsRight()
		zero := r - curOnes
		d := p - zero + curOnes
		//fmt.Fprintf(os.Stderr, "d=%d, ones=%d, r=%d, zero=%d, ntips=%d\n", d, curOnes, r, zero, ntips)
		if d > ntips/2 {
			d = ntips - d
		}
		// <= because even if d==p-1 (max dist)
		// we want to output a min dist edge
		if d <= *dist {
			//fmt.Fprintf(os.Stderr, "d=%d, dist=%d, p=%d\n", d, *dist, p)
			*dist = d
			*minedge = curEdge
			if (d == 1 && absent) || d == 0 {
				(*stop) = true
			}
		}
	}
}

// computes the transfer dist for each edges of the ref tree
// outrawtree: if tree with average transfer distance (non normalized) must be computed
// if false: then output rawtree is null
func Booster(reftree *tree.Tree, boottrees <-chan tree.Trees, cpu int, outrawtree bool) (rawtree *tree.Tree, err error) {
	tips := reftree.Tips()

	//vals := make([]int, len(edges))
	// Number of branches that have a normalized similarity (1- (min_dist/(n-1)) to
	// bootstrap trees > 0.8
	//var nb_branches_close int

	var edges []*tree.Edge = reftree.Edges()

	var nboot int = 0

	for i, e := range edges {
		e.SetId(i)
	}

	for boot := range boottrees {
		if boot.Err != nil {
			io.LogError(boot.Err)
			err = boot.Err
		} else {
			boot.Tree.ReinitIndexes()
			err = reftree.CompareTipIndexes(boot.Tree)
			if err == nil {
				//nb_branches_close = 0
				fmt.Fprintf(os.Stderr, "CPU : %02d - Bootstrap tree %d\r", cpu, boot.Id)
				bootedges := boot.Tree.Edges()
				bootedgeindex := tree.NewEdgeIndex(uint64(len(bootedges)*2), 0.75)

				for i, e := range bootedges {
					e.SetId(i)
					if !e.Right().Tip() {
						e.Right().SetName("")
					}
					if !e.Left().Tip() {
						e.Left().SetName("")
					}
					bootedgeindex.PutEdgeValue(e, i, e.Length())
				}

				var wg sync.WaitGroup
				for c := 0; c < cpu; c++ {
					wg.Add(1)
					go func() {
						for _, e := range edges {
							if p, _ := e.TopoDepth(); p > 1 {
								if _, ok := bootedgeindex.Value(e); ok {
									e.IncrementSupport(0.0)
								} else if p == 2 {
									e.IncrementSupport(1.0)
								} else {
									dist, _, _, _ := MinTransferDist(e, reftree, boot.Tree, len(tips), bootedges, true)
									//dist, edge, sptoadd, sptoremove := MinTransferDist(e, reftree, boot.Tree, len(tips), bootedges)
									e.IncrementSupport(float64(dist))
								}
							}
						}
						wg.Done()
					}()
				}
				wg.Wait()
			}
		}

		nboot++
		boot.Tree.Delete()
	}

	if outrawtree {
		rawtree = reftree.Clone()
		ReformatAvgDistance(rawtree)
	}
	NormalizeTransferDistancesByDepth(edges, nboot)

	return
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
func NormalizeTransferDistancesByDepth(edges []*tree.Edge, nboot int) {
	for _, e := range edges {
		if e.Support() != tree.NIL_SUPPORT {
			avgdist := e.Support() / float64(nboot)
			td, _ := e.TopoDepth()
			e.SetSupport(1.0 - avgdist/float64(td-1))
		}
	}
}
