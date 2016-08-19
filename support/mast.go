package support

import (
	"fmt"
	"github.com/fredericlemoine/gotree/io/utils"
	"github.com/fredericlemoine/gotree/tree"
	"sync"
)

type mastSupporter struct {
	expected_rand_val     []float64
	distribution_rand_val [][]float64
}

func (supporter *mastSupporter) ExpectedRandValues(depth int) float64 {
	return supporter.expected_rand_val[depth]
}

func (supporter *mastSupporter) ProbaDepthValue(d int, v int) float64 {
	return supporter.distribution_rand_val[d][v]
}

func (supporter *mastSupporter) Init(maxdepth int, nbtips int) {
	if supporter.expected_rand_val == nil {
		supporter.expected_rand_val = make([]float64, maxdepth+1)
		supporter.distribution_rand_val = make([][]float64, maxdepth+1)
		for i := 0; i <= maxdepth; i++ {
			supporter.distribution_rand_val[i] = make([]float64, nbtips+1)
			if i > 0 {
				supporter.expected_rand_val[i] = float64(i) - 1
				supporter.distribution_rand_val[i][i-1] = 1.0
			}
		}
	}
}

func update_all_i_c_post_order_ref_tree(refTree *tree.Tree, edges *[]*tree.Edge, bootTree *tree.Tree, bootEdges *[]*tree.Edge, i_matrix *[][]uint16, c_matrix *[][]uint16) {
	for i, child := range refTree.Root().Neigh() {
		update_i_c_post_order_ref_tree(refTree, edges, child, i, refTree.Root(), bootTree, bootEdges, i_matrix, c_matrix)
	}
}

// this function does the post-order traversal (recursive from the pseudoroot to the leaves, updating knowledge for the subtrees)
//   of the reference tree, examining only leaves (terminal edges) of the bootstrap tree.
//   It sends a probe from the orig node to the target node (nodes in ref_tree), calculating I_ij and C_ij
//   (see Brehelin, Gascuel, Martin 2008).
func update_i_c_post_order_ref_tree(refTree *tree.Tree, edges *[]*tree.Edge,
	current *tree.Node, curidx int, prev *tree.Node,
	bootTree *tree.Tree, bootEdges *[]*tree.Edge,
	i_matrix *[][]uint16, c_matrix *[][]uint16) {

	var e, be, e2 *tree.Edge
	var child *tree.Node
	var edge_id, edge_id2, be_id, k int

	e = prev.Edges()[curidx]
	edge_id = e.Id() /* all this is in ref_tree */

	if current.Tip() {
		for be_id, be = range *bootEdges { // for all the terminal edges of boot_tree
			if !be.Right().Tip() {
				continue
			}
			/* we only want to scan terminal edges of boot_tree, where the right son is a leaf */
			/* else we update all the I_ij and C_ij with i = edge_id */
			if current.Name() != be.Right().Name() {
				/* here the taxa are different */
				(*i_matrix)[edge_id][be_id] = 0
				(*c_matrix)[edge_id][be_id] = 1
			} else {
				/* same taxa here in T_ref and T_boot */
				(*i_matrix)[edge_id][be_id] = 1
				(*c_matrix)[edge_id][be_id] = 0
			}
		} /* end for on all edges of T_boot, for my_br being terminal */
	} else {
		/* now the case where my_br is not a terminal edge */
		/* first initialise (zero) the cells we are going to update */
		for be_id, be = range *bootEdges {
			// We initialize the i and c matrices for the edge edge_id with :
			// 	* 0 for i : because afterwards we do i[edge_id] = i[edge_id] || i[edge_id2]
			// 	* 1 for c : because afterwards we do c[edge_id] = c[edge_id] && c[edge_id2]
			if be.Right().Tip() {
				(*i_matrix)[edge_id][be_id] = 0
				(*c_matrix)[edge_id][be_id] = 1
			}
		}

		for k, child = range current.Neigh() {
			if child != prev {
				e2 = current.Edges()[k]
				edge_id2 = e2.Id()
				update_i_c_post_order_ref_tree(refTree, edges, child, k, current, bootTree, bootEdges, i_matrix, c_matrix)

				for be_id, be = range *bootEdges { /* for all the terminal edges of boot_tree */
					if !be.Right().Tip() {
						continue
					}

					// OR between two integers, result is 0 or 1 */
					if (*i_matrix)[edge_id][be_id] != 0 || (*i_matrix)[edge_id2][be_id] != 0 {
						(*i_matrix)[edge_id][be_id] = 1
					} else {
						(*i_matrix)[edge_id][be_id] = 0
					}

					// AND between two integers, result is 0 or 1
					if (*c_matrix)[edge_id][be_id] != 0 && (*c_matrix)[edge_id2][be_id] != 0 {
						(*c_matrix)[edge_id][be_id] = 1
					} else {
						(*c_matrix)[edge_id][be_id] = 0
					}
				}
			}
		}
	}
}

func update_all_i_c_post_order_boot_tree(refTree *tree.Tree, ntips uint, edges *[]*tree.Edge,
	bootTree *tree.Tree, bootEdges *[]*tree.Edge,
	i_matrix *[][]uint16, c_matrix *[][]uint16, hamming *[][]uint16, min_dist *[]uint16) {
	for i, child := range bootTree.Root().Neigh() {
		update_i_c_post_order_boot_tree(refTree, ntips, edges, bootTree, bootEdges, child, i, bootTree.Root(), i_matrix, c_matrix, hamming, min_dist)
	}

	/* and then some checks to make sure everything went ok */
	for _, e := range *edges {
		if (*min_dist)[e.Id()] < 0 {
			panic("Min dist should be >= 0")
		}
		if e.Right().Tip() && (*min_dist)[e.Id()] != 0 {
			panic("any terminal edge should have an exact match in any bootstrap tree")
		}
	}
}

// here we implement the second part of the Brehelin/Gascuel/Martin algorithm:
//    post-order traversal of the bootstrap tree, and numerical recurrence.
// in this function, orig and target are nodes of boot_tree (aka T_boot).
// min_dist is an array whose size is equal to the number of edges in T_ref.
//    It gives for each edge of T_ref its min distance to a split in T_boot.
func update_i_c_post_order_boot_tree(refTree *tree.Tree, ntips uint, edges *[]*tree.Edge,
	bootTree *tree.Tree, bootEdges *[]*tree.Edge,
	current *tree.Node, curindex int, prev *tree.Node,
	i_matrix *[][]uint16, c_matrix *[][]uint16, hamming *[][]uint16, min_dist *[]uint16) {
	var e, e2, e3 *tree.Edge
	var edge_id, edge_id2, edge_id3, j int
	var child *tree.Node

	e = prev.Edges()[curindex]
	edge_id = e.Id()

	if !current.Tip() {
		// because nothing to do in the case where target is a leaf: intersection and union already ok.
		// otherwise, keep on posttraversing in all other directions

		// first initialise (zero) the cells we are going to update
		for edge_id3 = 0; edge_id3 < len(*edges); edge_id3++ {
			(*i_matrix)[edge_id3][edge_id] = 0
			(*c_matrix)[edge_id3][edge_id] = 0
		}

		for j, child = range current.Neigh() {
			if child != prev {
				e2 = current.Edges()[j]
				edge_id2 = e2.Id()
				update_i_c_post_order_boot_tree(refTree, ntips, edges, bootTree, bootEdges, child, j, current,
					i_matrix, c_matrix, hamming, min_dist)
				for edge_id3 = 0; edge_id3 < len(*edges); edge_id3++ { /* for all the edges of ref_tree */
					(*i_matrix)[edge_id3][edge_id] += (*i_matrix)[edge_id3][edge_id2]
					(*c_matrix)[edge_id3][edge_id] += (*c_matrix)[edge_id3][edge_id2]
				}
			}
		}
	}

	for edge_id3, e3 = range *edges { // for all the edges of ref_tree
		// at this point we can calculate in all cases (internal branch or not) the Hamming distance at [i][edge_id],
		(*hamming)[edge_id3][edge_id] = // card of union minus card of intersection
			uint16(e3.NumTips()) + // #taxa in the cluster i of T_ref
				(*c_matrix)[edge_id3][edge_id] - // #taxa in cluster edge_id of T_boot BUT NOT in cluster i of T_ref
				(*i_matrix)[edge_id3][edge_id] // #taxa in the intersection of the two clusters

		/* NEW!! Let's immediately calculate the right distance, taking into account the fact that the true disance is min (dist, N-dist) */
		if (*hamming)[edge_id3][edge_id] > uint16(ntips)/2 { // floor value
			(*hamming)[edge_id3][edge_id] = uint16(ntips) - (*hamming)[edge_id3][edge_id]
		}

		/*   and update the min of all Hamming (MAST-like) distances hamming[i][j] over all j */
		(*min_dist)[edge_id3] = min_uint((*hamming)[edge_id3][edge_id], (*min_dist)[edge_id3])
	}
}

// Thread that takes bootstrap trees from the channel,
// computes the mastlike dist for each edges of the ref tree
// and send it to the result channel
func (supporter *mastSupporter) ComputeValue(refTree *tree.Tree, empiricalTrees []*tree.Tree, cpu int, empirical bool, edges []*tree.Edge, randEdges [][]*tree.Edge,
	wg *sync.WaitGroup, bootTreeChannel <-chan utils.Trees, valChan chan<- bootval, randvalChan chan<- bootval) {
	tips := refTree.Tips()
	var min_dist []uint16 = make([]uint16, len(edges))
	var i_matrix [][]uint16 = make([][]uint16, len(edges))
	var c_matrix [][]uint16 = make([][]uint16, len(edges))
	var hamming [][]uint16 = make([][]uint16, len(edges))

	vals := make([]int, len(edges))

	for treeV := range bootTreeChannel {
		fmt.Println(fmt.Sprintf("CPU : %d - Bootstrap tree %d", cpu, treeV.Id))
		bootEdges := treeV.Tree.Edges()
		for i, _ := range edges {
			min_dist[i] = uint16(len(tips))
			if len(bootEdges) > len(i_matrix[i]) {
				i_matrix[i] = make([]uint16, len(edges))
				c_matrix[i] = make([]uint16, len(edges))
				hamming[i] = make([]uint16, len(edges))
			}
		}

		for i, e := range bootEdges {
			e.SetId(i)
		}

		update_all_i_c_post_order_ref_tree(refTree, &edges, treeV.Tree, &bootEdges, &i_matrix, &c_matrix)
		update_all_i_c_post_order_boot_tree(refTree, uint(len(tips)), &edges, treeV.Tree, &bootEdges, &i_matrix, &c_matrix, &hamming, &min_dist)

		for i, e := range edges {
			if e.Right().Tip() {
				continue
			}
			vals[i] = int(min_dist[i])
			valChan <- bootval{
				vals[i],
				i,
				false,
			}
		}
		// We compute the empirical values
		if empirical {
			for j, et := range empiricalTrees {
				for i, _ := range edges {
					min_dist[i] = uint16(len(tips))
				}
				update_all_i_c_post_order_ref_tree(et, &randEdges[j], treeV.Tree, &bootEdges, &i_matrix, &c_matrix)
				update_all_i_c_post_order_boot_tree(et, uint(len(tips)), &randEdges[j], treeV.Tree, &bootEdges, &i_matrix, &c_matrix, &hamming, &min_dist)

				for i, _ := range randEdges[j] {
					val := int(min_dist[i])
					randvalChan <- bootval{
						val,
						i,
						vals[i] >= val,
					}
				}
			}
		}
	}
	wg.Done()
}

func MastLike(reftreefile, boottreefile string, empirical bool, cpus int) *tree.Tree {
	var supporter *mastSupporter = &mastSupporter{}
	return ComputeSupport(reftreefile, boottreefile, empirical, cpus, supporter)
}
