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
// Else: we do not know, we do the full postorder traversal, and speciestoadd && speciestoremove are filled
// speciestoadd/speciestoremove: The list of species to add and remove, for each branch of the minedges slice. The size of the speciestoadd/speciestoremiove slices are the same as the minedges slice.
func MinTransferDist(refedge *tree.Edge, reftree, boottree *tree.Tree, ntips int, bootedges []*tree.Edge, absent bool) (dist int, minedges []*tree.Edge, speciestoadd, speciestoremove [][]*tree.Node) {
	ones := make([]int, len(bootedges))
	p, _ := refedge.TopoDepth()
	dist = p - 1
	minedges = make([]*tree.Edge, 0, 2)

	// If ref edge is a terminal edge
	if p == 1 {
		tip := refedge.Right()
		boottip, _ := boottree.TipNode(tip.Name())
		if boottip != nil {
			minedges = append(minedges, boottip.Edges()[0])
		}
		return
	}

	stop := false
	minTransferDistRecur(reftree, ntips, boottree.Root(), nil, nil, refedge, p, ones, &dist, &minedges, absent, &stop)

	speciestoadd = make([][]*tree.Node, len(minedges))
	speciestoremove = make([][]*tree.Node, len(minedges))
	for i, _ := range minedges {
		speciestoadd[i] = make([]*tree.Node, 0, 2)
		speciestoremove[i] = make([]*tree.Node, 0, 2)
	}

	if !absent {
		for i, minedge := range minedges {
			// computing species to move
			/////////////////////////////////////////
			n_subtree := minedge.NumTipsRight()
			ones_subtree := ones[minedge.Id()]
			zeros_subtree := n_subtree - ones_subtree
			ones_total := ntips - p
			zeros_total := p

			/* move ones into subtree, zeros outside subtree? */
			ops_ones_in_subtree := zeros_subtree + (ones_total - ones_subtree)
			/* move zeros into subtree, ones outside subtree? */
			ops_zeros_in_subtree := ones_subtree + (zeros_total - zeros_subtree)

			want_ones_outside := ops_zeros_in_subtree < ops_ones_in_subtree

			speciesToMoveRecursive(minedge, boottree.Root(), nil, nil, ones, want_ones_outside, &speciestoadd[i], &speciestoremove[i])
		}
	}
	return
}

func speciesToMoveRecursive(bootedge *tree.Edge, cur, prev *tree.Node, edge *tree.Edge, ones []int, want_ones_now bool, speciestoadd, speciestoremove *[]*tree.Node) {

	if edge == bootedge {
		want_ones_now = !want_ones_now
	}

	if cur.Tip() {
		if want_ones_now && ones[edge.Id()] == 0 {
			*speciestoadd = append(*speciestoadd, cur)
		}
		if !want_ones_now && ones[edge.Id()] == 1 {
			*speciestoremove = append(*speciestoremove, cur)
		}
	}

	if edge != nil {
		subtreesize := edge.NumTipsRight()
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

func minTransferDistRecur(refTree *tree.Tree, ntips int, cur, prev *tree.Node, curEdge, refEdge *tree.Edge, p int, ones []int, dist *int, minedges *[]*tree.Edge, absent bool, stop *bool) {
	if *stop {
		return
	}
	curOnes := 0
	if cur.Tip() {
		tipIndex := cur.TipIndex()
		light := refEdge.TipPresent(uint(tipIndex))
		if r := refEdge.NumTipsRight(); r > ntips/2 {
			light = !light
		}
		if !light {
			curOnes = 1
		}
	} else {
		for i, n := range cur.Neigh() {
			if n != prev {
				nextEdge := cur.Edges()[i]
				minTransferDistRecur(refTree, ntips, n, cur, nextEdge, refEdge, p, ones, dist, minedges, absent, stop)
				curOnes += ones[nextEdge.Id()]
				if *stop {
					return
				}
			}
		}
	}

	if curEdge != nil {
		ones[curEdge.Id()] = curOnes
		r := curEdge.NumTipsRight()
		zero := r - curOnes
		d := p - zero + curOnes
		if d > ntips/2 {
			d = ntips - d
		}
		// <= because even if d==p-1 (max dist)
		// we want to output a min dist edge
		if d < *dist {
			*minedges = make([]*tree.Edge, 0, 2)
		}
		if d <= *dist {
			*dist = d
			*minedges = append(*minedges, curEdge)
			if d == 1 && absent {
				(*stop) = true
			}
		}
	}
}

// computes the transfer dist for each edges of the ref tree
// outrawtree: if tree with average transfer distance (non normalized) must be computed
// if false: then output rawtree is null
func TBE(reftree *tree.Tree, boottrees <-chan tree.Trees, cpu int,
	outrawtree bool, computeavgtaxa, computeperbranchtaxa bool, distcutoff float64,
	logfile *os.File, sup *Supporter) (rawtree *tree.Tree, err error) {
	tips := reftree.Tips()

	//vals := make([]int, len(edges))
	// Number of branches that have a normalized similarity (1- (min_dist/(n-1)) to
	// bootstrap trees > 0.8
	//var nb_branches_close int

	var edges []*tree.Edge = reftree.Edges()
	var movedspeciestmp []float64
	var movedspecies []float64
	var sumNbClosestBranches []float64 // Average number (over all the bootstrap trees) of branches having the minimum distance to this ref branch
	var movedperbranch [][]float64
	var nbranchclose int = 0
	var nboot int = 0
	var mindepth int = int(math.Ceil(1.0/distcutoff + 1.0)) // For taxa move computation

	if sup == nil {
		sup = &Supporter{}
	}

	if computeavgtaxa {
		movedspecies = make([]float64, len(tips))
		movedspeciestmp = make([]float64, len(tips))
	}

	sumNbClosestBranches = make([]float64, len(edges))
	if computeperbranchtaxa {
		movedperbranch = make([][]float64, len(edges))
		for i := range edges {
			movedperbranch[i] = make([]float64, len(tips))
		}
	}

	for _, e := range edges {
		e.SetSupport(tree.NIL_SUPPORT)
		if !e.Right().Tip() {
			e.Right().SetName("")
		}
		if !e.Left().Tip() {
			e.Left().SetName("")
		}
	}

	for boot := range boottrees {
		if sup.Canceled() {
			break
		}
		if boot.Err != nil {
			io.LogError(boot.Err)
			err = boot.Err
			return
		} else {
			if err = boot.Tree.ReinitIndexes(); err != nil {
				io.LogError(err)
				return
			}
			if err = reftree.CompareTipIndexes(boot.Tree); err != nil {
				io.LogError(err)
			}
			nbranchclose = 0
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
			var mux sync.Mutex
			wg.Add(cpu)
			edgechan := make(chan *tree.Edge, cpu*10)

			go func() {
				for _, e := range edges {
					edgechan <- e
				}
				close(edgechan)
			}()

			for c := 0; c < cpu; c++ {
				go func() {
					for e := range edgechan {
						if p, _ := e.TopoDepth(); p > 1 {
							if _, ok := bootedgeindex.Value(e); ok {
								if p >= mindepth {
									mux.Lock()
									nbranchclose++
									mux.Unlock()
								}
								e.IncrementSupport(0.0)
								sumNbClosestBranches[e.Id()] += 1.0
								//} else if p == 2 {
								//} else if p == 2 {
								//e.IncrementSupport(1.0)
							} else {
								dist, minedges, sptoadd, sptoremove := MinTransferDist(e, reftree, boot.Tree, len(tips), bootedges, !(computeavgtaxa || computeperbranchtaxa))
								//dist, edge, sptoadd, sptoremove := MinTransferDist(e, reftree, boot.Tree, len(tips), bootedges)
								e.IncrementSupport(float64(dist))
								if computeavgtaxa || computeperbranchtaxa {
									UpdateTaxaMoveArrays(e, minedges, dist, p,
										movedspeciestmp, movedperbranch, sumNbClosestBranches,
										&nbranchclose, sptoadd, sptoremove, distcutoff, mindepth,
										computeavgtaxa, computeperbranchtaxa, &mux)
								}
							}
						}
					}
					wg.Done()
				}()
			}
			wg.Wait()
		}
		if computeavgtaxa || computeperbranchtaxa {
			for _, t := range tips {
				if nbranchclose > 0 {
					movedspecies[t.TipIndex()] += float64(movedspeciestmp[t.TipIndex()]) / float64(nbranchclose)
				}
				movedspeciestmp[t.TipIndex()] = 0
			}
		}
		nboot++
		boot.Tree.Delete()
		sup.IncrementProgress()
	}

	if outrawtree {
		rawtree = reftree.Clone()
		ReformatAvgDistance(rawtree, nboot)
	}
	NormalizeTransferDistancesByDepth(edges, nboot)
	// Write in log file
	if computeavgtaxa {
		fmt.Fprintf(logfile, "Taxon\ttIndex\n")
		for _, t := range tips {
			movedtaxaindex := movedspecies[t.TipIndex()] * 100.0 / float64(nboot)
			fmt.Fprintf(logfile, "%s\t%f\n", t.Name(), movedtaxaindex)
		}
	}

	if computeperbranchtaxa {
		fmt.Fprintf(logfile, "Edge\tLength\tP\tSupport\tAvgNbClosestBranches")
		for _, t := range tips {
			fmt.Fprintf(logfile, "\t%s", t.Name())
		}
		fmt.Fprintf(logfile, "\n")
		for _, e := range edges {
			if e.Right().Tip() {
				continue
			}
			d, _ := e.TopoDepth()
			fmt.Fprintf(logfile, "%d\t%s\t%d\t%s\t%f", e.Id(), e.LengthString(), d, e.SupportString(), sumNbClosestBranches[e.Id()]/float64(nboot))
			for _, t := range tips {
				fmt.Fprintf(logfile, "\t%f", float64(movedperbranch[e.Id()][t.TipIndex()])*1.0/float64(nboot))
			}
			fmt.Fprintf(logfile, "\n")
		}
	}

	return
}

// This function writes on the child node name the string: "branch_id|avg_dist|depth"
// and removes support information from each branch
func ReformatAvgDistance(t *tree.Tree, nboot int) {
	for i, e := range t.Edges() {
		if e.Support() != tree.NIL_SUPPORT {
			td, _ := e.TopoDepth()
			e.Right().SetName(fmt.Sprintf("%d|%.6f|%d", i, e.Support()/float64(nboot), td))
			e.SetSupport(tree.NIL_SUPPORT)
		}
	}
}

// This function takes all branch support values (that are considered as average
// transfer distances over bootstrap trees), normalizes them by the depth and
// convert them to similarity, i.e:
//
//	1-avg_dist/(depth-1)
func NormalizeTransferDistancesByDepth(edges []*tree.Edge, nboot int) {
	for _, e := range edges {
		if e.Support() != tree.NIL_SUPPORT {
			avgdist := e.Support() / float64(nboot)
			td, _ := e.TopoDepth()
			e.SetSupport(1.0 - avgdist/float64(td-1))
		}
	}
}

// Looking at number of times each taxon moves around low distance branches
// moved_species: array of size number of nodes in the tree
// species_to_add & species_to_remove: Array of species (tree nodes) that move
// distcutoff: if the bootstrap branch is too distant from the ref branch in terms of normalized transfer dist, then does not count
func UpdateTaxaMoveArrays(ref *tree.Edge, boot []*tree.Edge, dist, p int,
	moved_species []float64, moved_species_per_branch [][]float64, sumNbClosestBranches []float64,
	nb_branches_close *int, speciestoadd, speciestoremove [][]*tree.Node, distcutoff float64,
	mindepth int, computeavgtaxa, computeperbranchtaxa bool, mux *sync.Mutex) {

	norm := float64(dist) / (float64(p) - 1.0)
	nbminedges := float64(len(boot))

	mux.Lock()
	if computeavgtaxa && norm <= distcutoff && p >= mindepth {
		for _, spa := range speciestoadd {
			for _, t := range spa {
				moved_species[t.TipIndex()] += (1.0 / nbminedges)
			}
		}
		for _, spr := range speciestoremove {
			for _, t := range spr {
				moved_species[t.TipIndex()] += (1.0 / nbminedges)
			}
		}
		(*nb_branches_close)++
	}
	mux.Unlock()

	sumNbClosestBranches[ref.Id()] += nbminedges
	if computeperbranchtaxa {
		for _, spa := range speciestoadd {
			for _, t := range spa {
				moved_species_per_branch[ref.Id()][t.TipIndex()] += (1.0 / nbminedges)
			}
		}

		for _, spr := range speciestoremove {
			for _, t := range spr {
				moved_species_per_branch[ref.Id()][t.TipIndex()] += (1.0 / nbminedges)
			}
		}
	}
}
