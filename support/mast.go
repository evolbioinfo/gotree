package support

import (
	"github.com/fredericlemoine/gotree/io/utils"
	"github.com/fredericlemoine/gotree/tree"
	"sync"
)

type mastSupporter struct {
	expected_rand_steps []float64
}

func (supporter *mastSupporter) ExpectedRandValues(maxdepth int, nbtips int) []float64 {
	if supporter.expected_rand_steps == nil {
		supporter.expected_rand_steps = make([]float64, maxdepth+1)
		for i := 0; i <= maxdepth; i++ {
			supporter.expected_rand_steps[i] = float64(i) - 1
		}
	}
	return supporter.expected_rand_steps
}

func update_all_i_c_post_order_ref_tree(refTree *tree.Tree, edges map[*tree.Edge]uint, bootTree *tree.Tree, bootEdges map[*tree.Edge]uint, i_matrix [][]uint, c_matrix [][]uint) {
	for _, child := range refTree.Root().Neigh() {
		update_i_c_post_order_ref_tree(refTree, edges, child, refTree.Root(), bootTree, bootEdges, i_matrix, c_matrix)
	}
}

// this function does the post-order traversal (recursive from the pseudoroot to the leaves, updating knowledge for the subtrees)
//   of the reference tree, examining only leaves (terminal edges) of the bootstrap tree.
//   It sends a probe from the orig node to the target node (nodes in ref_tree), calculating I_ij and C_ij
//   (see Brehelin, Gascuel, Martin 2008).
func update_i_c_post_order_ref_tree(refTree *tree.Tree, edges map[*tree.Edge]uint, current *tree.Node, prev *tree.Node, bootTree *tree.Tree, bootEdges map[*tree.Edge]uint, i_matrix [][]uint, c_matrix [][]uint) {

	idx, err := prev.NodeIndex(current)
	if err != nil {
		panic(err)
	}
	e := prev.Edges()[idx]
	edge_id, ok := edges[e] /* all this is in ref_tree */
	if !ok {
		panic("Edge has no id")
	}

	if current.Tip() {
		for be, j := range bootEdges { // for all the terminal edges of boot_tree
			if !be.Right().Tip() {
				continue
			}
			/* we only want to scan terminal edges of boot_tree, where the right son is a leaf */
			/* else we update all the I_ij and C_ij with i = edge_id */
			if current.Name() != be.Right().Name() {
				/* here the taxa are different */
				i_matrix[edge_id][j] = 0
				c_matrix[edge_id][j] = 1
			} else {
				/* same taxa here in T_ref and T_boot */
				i_matrix[edge_id][j] = 1
				c_matrix[edge_id][j] = 0
			}
		} /* end for on all edges of T_boot, for my_br being terminal */
	} else {
		/* now the case where my_br is not a terminal edge */
		/* first initialise (zero) the cells we are going to update */
		for be, j := range bootEdges {
			// We initialize the i and c matrices for the edge edge_id with :
			// 	* 0 for i : because afterwards we do i[edge_id] = i[edge_id] || i[edge_id2]
			// 	* 1 for c : because afterwards we do c[edge_id] = c[edge_id] && c[edge_id2]
			if be.Right().Tip() {
				i_matrix[edge_id][j] = 0
				c_matrix[edge_id][j] = 1
			}
		}

		for k, child := range current.Neigh() {
			if child != prev {
				e2 := current.Edges()[k]
				edge_id2, ok := edges[e2]
				if !ok {
					panic("Edge has no id")
				}
				update_i_c_post_order_ref_tree(refTree, edges, child, current, bootTree, bootEdges, i_matrix, c_matrix)

				//edge_id2 = current->br[dir]->id;
				for be, j := range bootEdges { /* for all the terminal edges of boot_tree */
					if !be.Right().Tip() {
						continue
					}

					if i_matrix[edge_id][j] != 0 || i_matrix[edge_id2][j] != 0 {
						i_matrix[edge_id][j] = 1
					} else {
						i_matrix[edge_id][j] = 0
					}
					//above is an OR between two integers, result is 0 or 1 */

					if c_matrix[edge_id][j] != 0 && c_matrix[edge_id2][j] != 0 {
						c_matrix[edge_id][j] = 1
					} else {
						c_matrix[edge_id][j] = 0
					}
					// above is an AND between two integers, result is 0 or 1
				} /* end for j */
			}
		} /* end for on all edges of T_boot, for my_br being internal */

	} /* ending the case where my_br is an internal edge */

} /* end update_i_c_post_order_ref_tree */

func update_all_i_c_post_order_boot_tree(refTree *tree.Tree, ntips uint, edges map[*tree.Edge]uint, bootTree *tree.Tree, bootEdges map[*tree.Edge]uint, i_matrix [][]uint, c_matrix [][]uint, hamming [][]uint, min_dist []uint) {
	for _, child := range bootTree.Root().Neigh() {
		update_i_c_post_order_boot_tree(refTree, ntips, edges, bootTree, bootEdges, child, bootTree.Root(), i_matrix, c_matrix, hamming, min_dist)
	}

	/* and then some checks to make sure everything went ok */
	for e, i := range edges {
		if min_dist[i] < 0 {
			panic("Min dist should be >= 0")
		}
		if e.Right().Tip() && min_dist[i] != 0 {
			panic("any terminal edge should have an exact match in any bootstrap tree")
		}
	}
}

func update_i_c_post_order_boot_tree(refTree *tree.Tree, ntips uint, edges map[*tree.Edge]uint, bootTree *tree.Tree, bootEdges map[*tree.Edge]uint, current *tree.Node, prev *tree.Node, i_matrix [][]uint, c_matrix [][]uint,
	hamming [][]uint, min_dist []uint) {
	// here we implement the second part of the Brehelin/Gascuel/Martin algorithm:
	//    post-order traversal of the bootstrap tree, and numerical recurrence.
	// in this function, orig and target are nodes of boot_tree (aka T_boot).
	// min_dist is an array whose size is equal to the number of edges in T_ref.
	//    It gives for each edge of T_ref its min distance to a split in T_boot.
	idx, err := prev.NodeIndex(current)
	if err != nil {
		panic(err)
	}
	e := prev.Edges()[idx]
	edge_id, ok := bootEdges[e] /* all this is in ref_tree */
	if !ok {
		panic("Edge has no id")
	}

	if !current.Tip() {
		/* because nothing to do in the case where target is a leaf: intersection and union already ok. */
		/* otherwise, keep on posttraversing in all other directions */

		/* first initialise (zero) the cells we are going to update */
		for _, i := range edges {
			i_matrix[i][edge_id] = 0
			c_matrix[i][edge_id] = 0
		}

		for j, child := range current.Neigh() {
			if child != prev {
				e2 := current.Edges()[j]
				edge_id2, ok := bootEdges[e2]
				if !ok {
					panic("Edge has no id")
				}
				update_i_c_post_order_boot_tree(refTree, ntips, edges, bootTree, bootEdges, child, current,
					i_matrix, c_matrix, hamming, min_dist)
				for _, i := range edges { /* for all the edges of ref_tree */
					i_matrix[i][edge_id] += i_matrix[i][edge_id2]
					c_matrix[i][edge_id] += c_matrix[i][edge_id2]
				}
			}
		}
	}

	for e, i := range edges { // for all the edges of ref_tree
		// at this point we can calculate in all cases (internal branch or not) the Hamming distance at [i][edge_id],
		hamming[i][edge_id] = // card of union minus card of intersection
			e.NumTips() + // #taxa in the cluster i of T_ref
				c_matrix[i][edge_id] - // #taxa in cluster edge_id of T_boot BUT NOT in cluster i of T_ref
				i_matrix[i][edge_id] // #taxa in the intersection of the two clusters

		/* NEW!! Let's immediately calculate the right ditance, taking into account the fact that the true disance is min (dist, N-dist) */
		if hamming[i][edge_id] > ntips/2 { // floor value
			hamming[i][edge_id] = ntips - hamming[i][edge_id]
		}

		/*   and update the min of all Hamming (MAST-like) distances hamming[i][j] over all j */
		if hamming[i][edge_id] < min_dist[i] {
			min_dist[i] = hamming[i][edge_id]
		}
	}
}

// Thread that takes bootstrap trees from the channel,
// computes the mastlike dist for each edges of the ref tree
// and send it to the result channel
func (supporter *mastSupporter) ComputeValue(refTree *tree.Tree, empiricalTrees []*tree.Tree, cpu int, empirical bool, edges []*tree.Edge, randEdges [][]*tree.Edge,
	wg *sync.WaitGroup, bootTreeChannel <-chan utils.Trees, valChan chan<- bootval, randvalChan chan<- bootval) {
	func(cpu int) {
		tips := refTree.Tips()
		var min_dist []uint = make([]uint, len(edges))
		var i_matrix [][]uint = make([][]uint, len(edges))
		var c_matrix [][]uint = make([][]uint, len(edges))
		var hamming [][]uint = make([][]uint, len(edges))

		// Map to assign ids to edges of reftree, empirical trees and boottrees
		var edgemap map[*tree.Edge]uint = make(map[*tree.Edge]uint)
		for i, e := range edges {
			edgemap[e] = uint(i)
		}

		var randedgemap []map[*tree.Edge]uint = make([]map[*tree.Edge]uint, len(randEdges))
		for j, _ := range randEdges {
			randedgemap[j] = make(map[*tree.Edge]uint)
			for i, e := range randEdges[j] {
				randedgemap[j][e] = uint(i)
			}
		}

		for treeV := range bootTreeChannel {
			bootEdges := treeV.Tree.Edges()
			var bootEdgesmap map[*tree.Edge]uint = make(map[*tree.Edge]uint)
			for i, _ := range edges {
				min_dist[i] = uint(len(tips))
				i_matrix[i] = make([]uint, len(bootEdges))
				c_matrix[i] = make([]uint, len(bootEdges))
				hamming[i] = make([]uint, len(bootEdges))
			}

			for i, e := range bootEdges {
				bootEdgesmap[e] = uint(i)
			}

			update_all_i_c_post_order_ref_tree(refTree, edgemap, treeV.Tree, bootEdgesmap, i_matrix, c_matrix)
			update_all_i_c_post_order_boot_tree(refTree, uint(len(tips)), edgemap, treeV.Tree, bootEdgesmap, i_matrix, c_matrix, hamming, min_dist)

			vals := make([]int, len(edges))
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
						min_dist[i] = uint(len(tips))
					}
					update_all_i_c_post_order_ref_tree(et, randedgemap[j], treeV.Tree, bootEdgesmap, i_matrix, c_matrix)
					update_all_i_c_post_order_boot_tree(et, uint(len(tips)), randedgemap[j], treeV.Tree, bootEdgesmap, i_matrix, c_matrix, hamming, min_dist)

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
	}(cpu)
}

func MastLike(reftreefile, boottreefile string, empirical bool, cpus int) *tree.Tree {
	var supporter *mastSupporter = &mastSupporter{}
	return ComputeSupport(reftreefile, boottreefile, empirical, cpus, supporter)
}
