package tree

import (
	"errors"
	"fmt"
	"log"
	"sort"
	"sync"

	//"os"

	"github.com/evolbioinfo/gotree/io"
)

// Given a set of tip names, this function
// returns the node that is the common ancestor of them
// and the edges that connects this node to the subtree.
//
// It considers the tree as unrooted
//
//	       e2---1
//	 ----a|
//	|      e1---2
//	|     ---3
//	 ----|
//	|     ---4
//	|     ---5
//	 ----|
//	      ---6
//
// # LeastCommonAncestorUnrooted(1,2) returns a,e1,e2,true
//
// The returned boolean value telling if the group is monophyletic
// (i.e. contains all tips descending from LCA).
func (t *Tree) LeastCommonAncestorUnrooted(nodeindex *nodeIndex, tips ...string) (*Node, []*Edge, bool, error) {
	var err error
	if nodeindex == nil {
		nodeindex, err = NewNodeIndex(t)
		if err != nil {
			return nil, nil, false, err
		}
	}
	tipindex := make(map[string]*Node, 0)
	for _, name := range tips {
		node, found := nodeindex.GetNode(name)
		if found && node.Tip() {
			tipindex[name] = node
		} else {
			io.LogWarning(fmt.Errorf("tip not found in the tree : %s", name))
		}
	}
	if len(tipindex) == 0 {
		return nil, nil, false, errors.New("none of the given tips are present in the tree")
	}

	// We search a tip that is not in the input tips
	// It will serve as a temporary root for the tree
	var temproot *Node = nil
	for _, othertip := range t.Tips() {
		_, found := tipindex[othertip.Name()]
		if !found {
			temproot = othertip
			break
		}
	}

	// If temproot == nil : Means that the input tips consist of all the tips of the tree
	if temproot == nil {
		return nil, nil, false, errors.New("all tips of the tree given : Nothing to do")
	}
	// otherwise we take the only child of the tip as first root
	ancestor, goodedges, _, diff, _, err2 := t.LeastCommonAncestorRecur(temproot.neigh[0], nil, tipindex)
	if err != nil {
		return nil, nil, false, err2
	}

	return ancestor, goodedges, diff == 0, nil
}

// Given a set of tip names, this function returns
// the node that is the common ancestor of them
// and the edges that connect this LCA node to the subtree.
//
// It considers the tree as Rooted.
//
// The returned boolean value tell if the group is monophyletic or not
// (i.e. contains all tips descending from LCA).
func (t *Tree) LeastCommonAncestorRooted(nodeindex *nodeIndex, tips ...string) (*Node, []*Edge, bool, error) {
	var err error
	if nodeindex == nil {
		if nodeindex, err = NewNodeIndex(t); err != nil {
			return nil, nil, false, err
		}
	}
	tipindex := make(map[string]*Node, 0)
	for _, name := range tips {
		node, found := nodeindex.GetNode(name)
		if found {
			tipindex[name] = node
		} else {
			io.LogWarning(fmt.Errorf("tip not found in the tree : %s", name))
		}
	}
	if len(tipindex) == 0 {
		return nil, nil, false, errors.New("none of the given tips are present in the tree")
	}

	// We search a tip that is not in the input tips
	// It will serve as a temporary root for the tree
	var temproot *Node = t.Root()

	// We take the only child of the tip as first root
	ancestor, goodedges, _, diff, _, err2 := t.LeastCommonAncestorRecur(temproot, nil, tipindex)
	if err2 != nil {
		return nil, nil, false, err2
	}

	return ancestor, goodedges, diff == 0, nil
}

// recursive function for getting the least common ancestor.
func (t *Tree) LeastCommonAncestorRecur(current *Node, prev *Node, tipIndex map[string]*Node) (*Node, []*Edge, int, int, bool, error) {
	common := 0
	edges := make([]*Edge, 0, 3)
	different := 0
	allFound := false

	// If current is a tip
	if current.Tip() {
		//fmt.Println(current.Name())
		_, found := tipIndex[current.Name()]
		if found {
			common++
			if idx, e := current.NodeIndex(prev); e == nil {
				edges = append(edges, current.br[idx])
			} else {
				return nil, nil, -1, -1, false, e
			}
		} else {
			different = 1
		}
	}

	// If current is not a tip
	tmpdiff := 0
	for i, n := range current.neigh {
		if n != prev {
			node, succedges, com, diff, found, err := t.LeastCommonAncestorRecur(n, current, tipIndex)
			if err != nil {
				return nil, nil, -1, -1, false, err
			}
			if found {
				//fmt.Println("int found - diff:", diff)
				return node, succedges, com, diff, found, nil
			} else if com > 0 {
				edges = append(edges, current.br[i])
				common += com
				different += diff
			} else {
				tmpdiff += diff
			}
		}
	}
	//fmt.Println("tmpdiff: ", tmpdiff)
	allFound = common == len(tipIndex)
	if allFound {
		//fmt.Println("found - diff:", different)
		return current, edges, common, different, allFound, nil
	} else {
		different += tmpdiff
		//fmt.Println("diff:", different)
		return nil, nil, common, different, allFound, nil
	}
}

// This function adds a branch/bipartition between the given node n
// and the given edges. To do so, it creates a new node between n and the edges,
// and connects it with a new edge.
//
// Imagine a star tree with central node n,
//
//	     1
//	     |
//	     |
//	6----n-----2
//	    /|\
//	   / | \
//	 e5 e4  e3
//
// if we call AddBipartition(n,{e3,e4,e5}), at the end we have:
//
//	     1
//	     |
//	     |
//	6----n-----2
//	     |
//	     |
//	     n2
//	    /|\
//	   / | \
//	 e5 e4  e3
//
// Useful for building consensus tree.
//
// If the edges are not initially directly connected to n, then returns an error. If ony one edge is given,
// returns an error (no need to add a new edge).
func (t *Tree) AddBipartition(n *Node, edges []*Edge, length, support float64) (*Edge, error) {
	n2 := t.NewNode()
	// Number of edges in direction n->e->other
	nbout := 0
	// Number of edges in direction n<-e<-other
	nbin := 0
	if len(edges) <= 1 || len(edges) >= len(n.br)-1 {
		return nil, errors.New("we cannot add the bipartition, it already exists")
	}
	for _, e := range edges {
		// We check if the edges are connected to the node
		// Else it exits with an error
		if e.Left() != n && e.Right() != n {
			return nil, errors.New("Edges need to be connected to the node to add a bipartition")
		}
		// Direction : true if n->e->other..., false if n<-e<-other
		// According to left / right
		dir := e.Left() == n
		var other *Node
		boot := e.Support()
		len := e.Length()
		pv := e.PValue()
		var etmp *Edge
		if dir {
			nbout++
			other = e.Right()
			other.delNeighbor(n)
			n.delNeighbor(other)
			etmp = t.ConnectNodes(n2, other)
		} else {
			nbin++
			other = e.Left()
			other.delNeighbor(n)
			n.delNeighbor(other)
			etmp = t.ConnectNodes(other, n2)
		}
		etmp.SetLength(len)
		etmp.SetSupport(boot)
		etmp.SetPValue(pv)
	}

	var e *Edge
	if nbin == 0 {
		e = t.ConnectNodes(n, n2)
	} else {
		e = t.ConnectNodes(n2, n)
	}
	e.SetLength(length)
	e.SetSupport(support)
	return e, nil
}

// Builds the consensus of trees given in the input channel.
//   - If the cutoff is 0.5 : The majority rule consensus is computed;
//   - If tht cutoff is 1   : The strict consensus is computed
//
// In the output consensus tree:
//  1. Branch supports are computed as the proportion of trees in which the bipartitions are present
//  2. Branch lengths are computed as the average length of the same branch over all the trees where it is present
//
// There can be errors if:
//   - The cutoff <0.5 or >1
//   - The tip names are different in the different trees
//   - Incompatible bipartition are generated to build the consensus (It should not happen since cutoff should be >=0.5)
func Consensus(trees <-chan Trees, cutoff float64) (*Tree, error) {
	if cutoff < 0.5 || cutoff > 1 {
		return nil, errors.New("min frequency for bipartition must be >=0.5 and <=1")
	}
	nbtrees := 0
	edgeindex := NewEdgeIndex(128, .75)
	var nodeindex *nodeIndex
	var startree *Tree = nil
	nbtips := 0
	var alltips []string
	var err error
	// We fill the edge index with all the bipartition and their count
	for curtree := range trees {
		if curtree.Err != nil {
			/* We empty the channel if needed */
			for range trees {
			}
			return nil, curtree.Err
		}

		if err = curtree.Tree.ReinitIndexes(); err != nil {
			return nil, err
		}

		// If the star tree is not initialized, we create it with the tips of the first tree
		if startree == nil {
			alltips = curtree.Tree.AllTipNames()
			if startree, err = StarTreeFromTree(curtree.Tree); err != nil {
				return nil, err
			}
			if err = startree.UpdateTipIndex(); err != nil {
				return nil, err
			}

			nbtips = len(alltips)
			// We first build the node index
			if nodeindex, err = NewNodeIndex(startree); err != nil {
				return nil, err
			}
		} else {
			// Compare tip names between star tree and current tree
			// Error if different sets (use already computed indexes)
			names := curtree.Tree.AllTipNames()
			if len(names) != nbtips {
				return nil, errors.New("Trees do not have the same set of tips")
			}
			for _, name := range names {
				if ok, err3 := startree.ExistsTip(name); err3 != nil {
					return nil, err3
				} else if !ok {
					return nil, errors.New("Trees do not have the same set of tips")
				}
			}
		}
		// We add the edge into the index
		for _, e := range curtree.Tree.Edges() {
			edgeindex.AddEdgeCount(e)
		}
		nbtrees++
	}

	// We take the bipartitions that are present in more than cutoff trees and less
	// than or equal the number of trees
	// And we add it to the startree
	for _, bs := range edgeindex.Edges(int(cutoff*float64(nbtrees)), nbtrees) {
		names := make([]string, 0, bs.key.Bitset().Count())
		for _, n := range alltips {
			if idx, err := startree.TipIndex(n); err != nil {
				return nil, err
			} else {
				if bs.key.Bitset().Test(uint(idx)) {
					names = append(names, n)
				}
			}
		}

		// Names of the tips in one side of the bipartition
		if len(names) < 2 {
			if len(names) == 1 {
				if t, ok := nodeindex.GetNode(names[0]); !ok || !t.Tip() {
					return nil, fmt.Errorf("this taxon name does not exist in the consensus: %s", names[0])
				} else {
					t.br[0].SetLength(float64(bs.val.Len) / float64(bs.val.Count))
				}
			} else {
				return nil, errors.New("this bipartition has a side with no taxa")
			}
		} else {
			node, edges, monophyletic, err := startree.LeastCommonAncestorUnrooted(nodeindex, names...)
			if err != nil {
				return nil, err
			}
			if node == nil {
				return nil, errors.New("consensus error: No common ancestor found for biparition")
			}
			if edges == nil || len(edges) == 0 {
				return nil, errors.New("consensus error: No common ancestor Edges found")
			}
			if !monophyletic {
				return nil, errors.New("the group should be monophyletic")
			}
			// We add the bipartition with a support value corresponding to the percentage of
			// trees in which it appears
			// TODO: Average branch length : Need to change the data structure
			startree.AddBipartition(node, edges, float64(bs.val.Len)/float64(bs.val.Count), float64(bs.val.Count)/float64(nbtrees))
		}
	}

	if err = startree.ReinitIndexes(); err != nil {
		return nil, err
	}
	return startree, nil
}

// CollapseClade takes a list of tips, and collapses the last common ancestor of these tips
// i.e. it removes all the descendants of that node and considers it as a tip.
// It returns the collapsed clade as a new Tree
//
// The tree is considered rooted
func (t *Tree) CollapseClade(strict bool, name string, tips ...string) (clade *Tree, err error) {
	var n, p, tip *Node
	var e *Edge
	var idx int
	var monophyletic bool

	if n, _, monophyletic, err = t.LeastCommonAncestorRooted(nil, tips...); err != nil {
		return
	}
	if !monophyletic {
		if strict {
			err = errors.New("the given outgroup is not monophyletic, cannot reroot")
			return
		}
		log.Println("warning: The given outgroup is not Monophyloetic, and may result in inappropriate rerooting")
	}

	if p, err = n.Parent(); err != nil {
		return
	}

	if e, err = n.ParentEdge(); err != nil {
		return
	}

	if idx, err = p.NodeIndex(n); err != nil {
		return
	}
	tip = t.NewNode()
	tip.SetName(name)

	p.neigh[idx] = tip
	e.setRight(tip)
	tip.addChild(p, e)
	t.UpdateTipIndex()

	// We build a new tree from the clade
	if idx, err = n.NodeIndex(p); err != nil {
		return
	}
	// We remove p from neighbors of n
	n.delNeighbor(p)
	clade = NewTree()
	clade.SetRoot(n)
	clade.UpdateTipIndex()

	return
}

// This function first unroots the input tree and reroots it using the outgroup in argument.
//
// If the outgroup is not monophyletic and strict is false, it will take all the descendant
// of the LCA and print a warning. If strict is true, it returns an error.
//
// An error is returned if the LCA is multifurcated, and several branches are possible
// for the placement of the root.
//
// If the outgroup includes a tip that is not present in the tree,
// this tip will not be taken into account for the rerooting.
//
// If removeoutgroup is true, then the outgrouped is removed from the rerooted tree.
func (t *Tree) RerootOutGroup(removeoutgroup, strict bool, tips ...string) error {
	t.UnRoot()

	n, edges, monophyletic, err := t.LeastCommonAncestorUnrooted(nil, tips...)
	if err != nil {
		return err
	}
	if !monophyletic {
		if strict {
			return errors.New("the given outgroup is not monophyletic, cannot reroot")
		}
		log.Println("Warning! The given outgroup is not Monophyloetic, and may result in inappropriate rerooting")
	}
	var rootedge *Edge

	if len(n.br) == 1 {
		rootedge = n.br[0]
	} else {
		if len(n.br)-len(edges) != 1 {
			return errors.New("Reroot error: Several possible branches for root placement (multifurcated node)")
		}
		// We search the unique branch connecting "n" and that is not part of the outgroup
		// If there were several branches, there should have been an error above
		for _, e := range n.br {
			found := false
			for _, e2 := range edges {
				if e == e2 {
					found = true
					break
				}
			}
			// That branch (e) is not part of the ougroup
			// => OK
			if !found {
				rootedge = e
				break
			}
		}
	}

	var root *Node
	// Here we will remove outgroup nodes
	if removeoutgroup {
		// We get the new root as the node on the
		// other side of the edge
		root = rootedge.Left()
		if root == n {
			root = rootedge.Right()
		}
		root.delNeighbor(n)
		n.delNeighbor(root)
		// We retrieve all nodes of the
		// subtree we want to remove
		nodes := make([]*Node, 0, len(tips))
		t.nodesRecur(&nodes, n, root)
		for _, no := range nodes {
			t.delNode(no)
		}
	} else {
		// Else we add a new root at the middle of the edge
		// connecting the outgroup to the other subtree
		root = t.NewNode()
		length := rootedge.Length()
		support := rootedge.Support()

		lnode := rootedge.Left()
		rnode := rootedge.Right()
		lnode.delNeighbor(rnode)
		rnode.delNeighbor(lnode)

		ne := t.ConnectNodes(root, lnode)
		ne2 := t.ConnectNodes(root, rnode)

		if length > 0 {
			ne.SetLength(length / 2.0)
			ne2.SetLength(length / 2.0)
			ne.SetSupport(support)
			ne2.SetSupport(support)
		}
	}
	if err = t.reroot_nocheck(root); err != nil {
		return err
	}
	if removeoutgroup {
		if err = t.UpdateTipIndex(); err != nil {
			return err
		}
	}
	t.ReinitInternalIndexes()
	return nil
}

// This function reroots the tree at the midpoint position.
// To do so, it first gets the 2 farthest tips of the tree,
// and takes the middle of the path between these tips as the
// new root position.
func (t *Tree) RerootMidPoint() error {
	// We first unroot the tree
	t.UnRoot()

	// All tips of the tree
	tips := t.Tips()
	// Maximum length path
	var potentialedges []*Edge
	// Length of the path
	curlength := 0.0

	// We take the max length path of all the tips
	for _, t := range tips {
		edges, length, err := MaxLengthPath(t, nil)
		if err != nil {
			return err
		}
		if length > curlength {
			curlength = length
			potentialedges = edges
		}
	}
	// Path potentialedges starts from tip 1:
	// potentialedges[0].Right()
	// And ends at tip 2:
	// potentialedges[len(potentialedges)-1].Right()

	// Find the right edge in the path to place the root
	i := 0
	len := 0.0
	// We need to orient the edge we find.
	// To know from which node the cut will be done.
	// Necessary because orientation changes during the path
	// when we cross the root node.
	var node1, node2 *Node
	for float64(len) < curlength/2.0 {
		// First tip
		if i == 0 {
			node1 = potentialedges[i].Right()
			node2 = potentialedges[i].Left()
		} else {
			if potentialedges[i].Right() == node2 {
				// We did not cross the root node, and we go up
				node1 = potentialedges[i].Right()
				node2 = potentialedges[i].Left()
			} else if potentialedges[i].Left() == node2 {
				// We already crossed the root node and we now go done
				node1 = potentialedges[i].Left()
				node2 = potentialedges[i].Right()
			}
		}
		len += potentialedges[i].Length()
		i++
	}

	// Where I cut the current edge
	// The cut is done at "cut" distance from node1 on the edge.
	cut := len - curlength/2.0

	newroot := t.NewNode()
	l := potentialedges[i-1].Length()
	b := potentialedges[i-1].Support()
	node1.delNeighbor(node2)
	node2.delNeighbor(node1)
	e := t.ConnectNodes(newroot, node1)
	e2 := t.ConnectNodes(newroot, node2)

	e.SetLength(l - cut)
	e2.SetLength(cut)
	e.SetSupport(b)
	e2.SetSupport(b)

	t.Reroot(newroot)
	t.ReinitInternalIndexes()
	return nil
}

// Computes the path of maximum length between the given node
// and any other node.
//
// It takes as argument the node from which we want to get the farthest
// tip (longest possible path).
//
// It returns the path (slice of edges), and the sum of branch lengths
// of this path.
//
// Returns an error if a branch has no length
func MaxLengthPath(cur *Node, prev *Node) ([]*Edge, float64, error) {
	var potentialedges []*Edge
	curlength := 0.0
	for i, child := range cur.neigh {
		if child != prev {
			e := cur.br[i]
			if e.Length() == NIL_LENGTH {
				return nil, -1, errors.New("some branches have no length")
			}
			edges, l, err := MaxLengthPath(child, cur)
			if err != nil {
				return nil, -1, err
			}
			if l+e.Length() > curlength {
				curlength = l + e.Length()
				potentialedges = append(edges, e)
			}
		}
	}
	return potentialedges, curlength, nil
}

const (
	DISTANCE_METRIC_BRLEN = iota // The distance of each edge corresponds to length (patristic distance).
	DISTANCE_METRIC_BOOTS        // The distance of each edge corresponds to its bootstrap support.
	DISTANCE_METRIC_NONE         // Each edge will count a distance of 1 (topological distance).
)

// ToDistanceMatrix computes and returns distances
// between all pairs of tips in the tree.
// metric can be :
//   - DISTANCE_METRIC_BRLEN : The distance of each edge corresponds to length (patristic distance).
//   - DISTANCE_METRIC_BOOTS : The distance of each edge corresponds to its bootstrap support.
//   - DISTANCE_METRIC_NONE : Each edge will count a distance of 1 (topological distance).
//   - All other values will be considered as DISTANCE_METRIC_BRLEN
//
// Computes patristic distance matrix. Outputs the distance matrix and the list of tips in the
// same order as the lines and columns of the matrix
func (t *Tree) ToDistanceMatrix(metric int) ([][]float64, []*Node) {
	// All tips of the tree
	tips := t.Tips()
	// Sort by name
	sort.Slice(tips, func(i, j int) bool {
		return tips[i].Name() < tips[j].Name()
	})

	// Init distance Matrix
	var matrix [][]float64 = make([][]float64, len(tips))
	for i := range tips {
		matrix[i] = make([]float64, len(tips))
		tips[i].SetId(i)
	}

	for i, t := range tips {
		pathLengths(t, nil, matrix[i], 0, metric)
	}
	return matrix, tips
}

func pathLengths(cur *Node, prev *Node, lengths []float64, curlength float64, metric int) {
	if cur.Tip() && prev != nil {
		lengths[cur.Id()] = curlength
	} else {
		for i, child := range cur.neigh {
			if child != prev {
				e := cur.br[i]
				l := 1.0
				switch metric {
				case DISTANCE_METRIC_BOOTS:
					l = e.Support()
					if l == NIL_SUPPORT {
						l = 1.0
					}
				case DISTANCE_METRIC_NONE:
					l = 1.0
				default:
					l = e.Length()
					if l == NIL_LENGTH {
						l = 0.0
					}
				}
				pathLengths(child, cur, lengths, curlength+l, metric)
			}
		}
	}
}

// Type for channel of tree stats
type BipartitionStats struct {
	Id       int   // Identifier of the tree analyzed
	Tree1    int   // Number of bipartitions specific to the first tree
	Tree2    int   // Number of bipartitions specific to the second tree
	Common   int   // Number of common bipartitions specific to the second tree
	Sametree bool  // True if the trees are identical
	Err      error // Wether an error occured or not in the computation
}

// This function compares bipartitions of a reference tree with a set of trees given in the input channel.
//
// If tips is true, then comparison includes external branches. If comparetreeidentical is true, does not compute
// the exact number of common and specific branches, but just put sametree=true or sametree=false in the stat channel.
//
// This function returns almost immediately because computation is done in several go routines in background.
// However it returns a Channel that will contain bipartition statistics computed so far. This channel is closed at the end of the computations,
// so on the calling functin, you can iterate over this channel in order to wait for the end of computations.
//
// It First Initializes bitsets of the reference tree
func Compare(refTree *Tree, compTrees <-chan Trees, tips, comparetreeidentical bool, cpus int) (<-chan BipartitionStats, error) {
	var edges []*Edge
	var err error

	stats := make(chan BipartitionStats)

	if refTree == nil {
		return nil, errors.New("Tree 1 in comparison is null")
	}
	if err = refTree.ReinitIndexes(); err != nil {
		return nil, err
	}
	edges = refTree.Edges()
	index := NewEdgeIndex(uint64(len(edges)*2), 0.75)
	total := 0
	for i, e := range edges {
		index.PutEdgeValue(e, i, e.Length())
		if tips || !e.Right().Tip() {
			total++
		}
	}

	var wg sync.WaitGroup
	for cpu := 0; cpu < cpus; cpu++ {
		wg.Add(1)
		go func(cpu int) {
			for treeV := range compTrees {
				total2 := 0
				common := 0
				var inerr error
				inerr = treeV.Err
				// Check wether the 2 trees have the same set of tip names
				// Else an error is included in the stats
				sametree := false
				if inerr == nil {
					if inerr = treeV.Tree.ReinitIndexes(); inerr == nil {
						edges2 := treeV.Tree.Edges()
						if inerr = refTree.CompareTipIndexes(treeV.Tree); err == nil {
							sametree = true
							for _, e2 := range edges2 {
								ok := true
								if tips || !e2.Right().Tip() {
									total2++
								}
								if !e2.Right().Tip() {
									_, ok = index.Value(e2)
								}
								if !ok {
									sametree = false
									if comparetreeidentical {
										break
									}
								}
								if ok && (tips || !e2.Right().Tip()) {

									common++
								}
							}
						}
					}
				}
				stats <- BipartitionStats{
					treeV.Id,
					total - common,
					total2 - common,
					common,
					sametree,
					inerr,
				}
			}
			wg.Done()
		}(cpu)
	}

	go func() {
		wg.Wait()
		close(stats)
	}()

	return stats, nil
}

type WeightedBipartitionStats struct {
	Id       int       // Identifier of the tree analyzed
	Tree1    []float64 // Lengths of bipartitions specific to the first tree
	Tree2    []float64 // Lengths of bipartitions specific to the second tree
	Common   []float64 // Difference of lengths for common bipartitions
	Sametree bool      // True if the trees are identical
	Err      error     // Wether an error occured or not in the computation
}

// This function compares bipartitions of a reference tree with a set of trees given in the input channel, taking into account branch lengths.
//
// If tips is true, then comparison includes external branches. If comparetreeidentical is true, does not compute
// the exact number of common and specific branches, but just put sametree=true or sametree=false in the stat channel.
//
// This function returns almost immediately because computation is done in several go routines in background.
// However it returns a Channel that will contain bipartition statistics computed so far. This channel is closed at the end of the computations,
// so on the calling functin, you can iterate over this channel in order to wait for the end of computations.
//
// Since this function builds an edge index for the reference tree and the compared trees it will use about twice as much memory
// as `Compare`, so if you do not need the branch length differences it will be more efficient to use `Compare` than `CompareWeighted`
func CompareWeighted(refTree *Tree, compTrees <-chan Trees, tips, comparetreeidentical bool, cpus int) (<-chan WeightedBipartitionStats, error) {
	var refEdges []*Edge
	var compEdges []*Edge

	var err error

	stats := make(chan WeightedBipartitionStats)

	if refTree == nil {
		return nil, errors.New("Tree 1 in comparison is null")
	}
	if err = refTree.ReinitIndexes(); err != nil {
		return nil, err
	}

	// Edge index of reference tree
	refEdges = refTree.Edges()
	refIndex := NewEdgeIndex(uint64(len(refEdges)*2), 0.75)
	for i, e := range refEdges {
		refIndex.PutEdgeValue(e, i, e.Length())
	}

	var wg sync.WaitGroup
	for cpu := 0; cpu < cpus; cpu++ {
		wg.Add(1)
		go func(cpu int) {
			for treeV := range compTrees {
				var inerr error
				inerr = treeV.Err

				var Common []float64
				var Ref []float64
				var Comp []float64

				// Check wether the 2 trees have the same set of tip names
				// Else an error is included in the stats
				sametree := false
				if inerr == nil {
					if inerr = treeV.Tree.ReinitIndexes(); inerr == nil {

						// Edge index of compared tree
						compEdges = treeV.Tree.Edges()
						compIndex := NewEdgeIndex(uint64(len(compEdges)*2), 0.75)
						for i, e := range compEdges {
							compIndex.PutEdgeValue(e, i, e.Length())
						}

						// The trees have the same tips, we can compare them
						if inerr = refTree.CompareTipIndexes(treeV.Tree); err == nil {
							sametree = true

							// Check compared edges against reference index
							for _, compEdge := range compEdges {
								if tips || !compEdge.Right().Tip() {
									refEdge, ok := refIndex.Value(compEdge)
									if ok { // Common edge
										refLen := refEdge.Len
										compLen := compEdge.Length()
										if refLen != compLen {
											sametree = false
											if comparetreeidentical {
												break
											}
										}

										Common = append(Common, refLen-compLen)
									} else { // Unique to compared tree
										sametree = false
										if comparetreeidentical {
											break
										}

										Comp = append(Comp, compEdge.Length())
									}
								}
							}

							// Check reference edges against compared index
							for _, refEdge := range refEdges {
								if tips || !refEdge.Right().Tip() {
									_, ok := compIndex.Value(refEdge)
									if !ok { // Unique to reference tree
										sametree = false
										if comparetreeidentical {
											break
										}

										Ref = append(Ref, refEdge.Length())
									}
								}
							}
						}
					}
				}

				stats <- WeightedBipartitionStats{
					treeV.Id,
					Ref,
					Comp,
					Common,
					sametree,
					inerr,
				}
			}
			wg.Done()
		}(cpu)
	}

	go func() {
		wg.Wait()
		close(stats)
	}()

	return stats, nil
}
