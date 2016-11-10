package tree

import (
	"errors"
	"github.com/fredericlemoine/goalign/io"
)

/**
Structure representing a "quadruplet set"
Actually X sets of taxa defined by a bipartition (Maybe multifurcated)
The enumeration of all the quadruplets is done by the
"Quadruplets" function
*/
type QuadrupletSet struct {
	/** Indexes of the taxa in the node index
	First dimension : The branch index: b0,b1 or b2,b3
	Second dimension : The tax indexes in the branche
	       b0       b2
	        \       /
	         \     /
	      left-----right
	         /     \
	        /       \
	       b1       b3
	*/
	left  [][]int
	right [][]int
}

type Quadruplet struct {

	/**
	structure of a single quadruplet
	either (t1,t2)(t3,t4)
	*/
	t1, t2, t3, t4 int
}

/* Iterate over all the quadruplets defined by the bipartition */
func (t *Tree) Quadruplets(it func(quad []int)) {

	// We initialize the nodes Id of the tree
	nodes := t.Nodes()
	nnodes := length(nodes)
	for i, n := range nodes {
		n.SetId(i)
	}
	// And nodes in all the left and right side
	// of the edges
	right := make([][]int, nnodes)
	left := make([][]int, nnodes)

	for i := 0; i < nnodes; i++ {
		right[i] = make([]int)
		left[i] = make([]int)
	}

	postOrderQuadrupletSet(t, t.Root(), nil, right)
	preOrderQuadrupletSet(t, t.Root().nil, left, right)
	qs := NewQuadrupletSet()
	// We use the information from left and rights arrays
	// To fill quadrupletsets for each edge
	for i, e := range t.Edges() {
		// If not possible to define a quadruplet we do nothing
		if e.Left().Nneigh() < 3 || e.Right().Nneigh() < 3 {
			next
		}
		for _, n := range e.Left().Neigh() {
			if n != e.Right() {
				qs.left = append(qs.left, left[n.Id()])
			}
		}
		for _, n := range e.Right().Neigh() {
			if n != e.Left() {
				qs.right = append(qs.right, right[n.Id()])
			}
		}
		qs.iterate(it)
	}
}

// Function that enumerates all quadruplets defined by a
// quadrupletset
func (qs *QuadrupletSet) iterate(it func(quad []int)) {

}

/*
 Compute information from all taxa at right side of every edges : postorder traversal
*/
func postOrderQuadrupletSet(t *Tree, n *Node, prev *Node, right [][]int) []uint {
	output := make([]uint)
	if n.Tip() {
		taxindex, err := t.TipIndex(n.Name())
		if err != nil {
			io.ExitWithMessage(err)
		}
		output = append(output, taxindex)
	} else {
		for next := range n.Neigh() {
			if next != prev {
				output = append(output, postOrderQuadrupletSet(t, next, n, right))
			}
		}
	}

	right[n.Id()] = append(right[n.Id()], output...)
	return output
}

/*
 Use information from all taxa at right side of every edges to know the taxa at left
side of every edges : preorder traversal
*/
func preOrderQuadrupletSet(t *Tree, n *Node, prev *Node, left [][]int, right [][]int) {
	for next := range n.Neigh() {
		if next != prev {
			for next2 := range n.Neigh() {
				// We append right of childs of n other than next
				// to left of next
				if next2 != prev && next2 != next {
					left[next.Id()] = append(left[next.Id()], right[next2.Id()]...)
				}
				// We append left of n to left of next
				left[next.Id()] = append(left[next.Id()], left[n.Id()])
			}
		}
		preOrderQuadrupletSet(t, next, n, left, right)
	}
}

/**
Initializes a new empty quadruplet with 4 sets of taxa
*/
func NewQuadrupletSet() *QuadrupletSet {
	return &Quadruplet{
		left:  make([][]int),
		right: make([][]int),
	}
}
