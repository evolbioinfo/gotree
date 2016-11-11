package tree

import (
	"fmt"
	"github.com/fredericlemoine/goalign/io"
)

/**
Structure representing a "quartets set"
Actually X sets of taxa defined by a bipartition (Maybe multifurcated)
The enumeration of all the quartets is done by the
"Quartets" function
*/
type QuartetSet struct {
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
	left  [][]uint
	right [][]uint
}

/**
Iterate over all the quartets defined by the bipartition
(t1,t2)(t3,t4)
*/

func (t *Tree) Quartets(it func(t1, t2, t3, t4 uint)) {
	// We initialize the nodes Id of the tree
	nodes := t.Nodes()
	nnodes := len(nodes)
	for i, n := range nodes {
		n.SetId(i)
	}
	// And nodes in all the left and right side
	// of the edges
	right := make([][]uint, nnodes)
	left := make([][]uint, nnodes)

	for i := 0; i < nnodes; i++ {
		right[i] = make([]uint, 0, 4)
		left[i] = make([]uint, 0, 4)
	}

	postOrderQuartetSet(t, t.Root(), nil, right)
	preOrderQuartetSet(t, t.Root(), nil, left, right)

	// We use the information from left and rights arrays
	// To fill quartetsets for each edge
	for _, e := range t.Edges() {
		// If not possible to define a quartet we do nothing
		if e.Left().Nneigh() < 3 || e.Right().Nneigh() < 3 {
			continue
		}
		qs := NewQuartetSet()
		for i, n := range e.Left().Neigh() {
			if n != e.Right() {
				br := e.Left().Edges()[i]
				// if outgoing edge from e.Left()
				if br.Left() == e.Left() {
					qs.left = append(qs.left, right[n.Id()])
				} else {
					// Ingoing edge from e.Left()
					qs.left = append(qs.left, left[br.Right().Id()])
				}
			}
		}

		for _, n := range e.Right().Neigh() {
			if n != e.Left() {
				qs.right = append(qs.right, right[n.Id()])
			}
		}
		// for b1 := 0; b1 < len(qs.left); b1++ {
		// 	for b2 := b1 + 1; b2 < len(qs.left); b2++ {
		// 		for b3 := 0; b3 < len(qs.right); b3++ {
		// 			for b4 := b3 + 1; b4 < len(qs.right); b4++ {
		// 				fmt.Println(len(qs.left[b1]) * len(qs.left[b2]) * len(qs.right[b3]) * len(qs.right[b4]))
		// 			}
		// 		}
		// 	}
		// }
		qs.iterate(it)
	}
}

// Function that enumerates all quartets defined by a quartetset
// (t1,t2)(t3,t4)
func (qs *QuartetSet) iterate(it func(t1, t2, t3, t4 uint)) {
	// Foreach pairs of branches [b1,b2] on the left
	for b1 := 0; b1 < len(qs.left); b1++ {
		for b2 := b1 + 1; b2 < len(qs.left); b2++ {
			// Foreach pairs of branches [b3,b4] on the right
			for b3 := 0; b3 < len(qs.right); b3++ {
				for b4 := b3 + 1; b4 < len(qs.right); b4++ {
					// All the quartets
					// Taxa of branch 1
					for tb1 := 0; tb1 < len(qs.left[b1]); tb1++ {
						// Taxa of branch 2
						for tb2 := 0; tb2 < len(qs.left[b2]); tb2++ {
							// Taxa of branch 3
							for tb3 := 0; tb3 < len(qs.right[b3]); tb3++ {
								// Taxa of branch 4
								for tb4 := 0; tb4 < len(qs.right[b4]); tb4++ {
									it(qs.left[b1][tb1], qs.left[b2][tb2], qs.right[b3][tb3], qs.right[b4][tb4])
								}
							}
						}
					}
				}
			}
		}
	}
}

/*
 Compute information from all taxa at right side of every edges : postorder traversal
*/
func postOrderQuartetSet(t *Tree, n *Node, prev *Node, right [][]uint) []uint {
	output := make([]uint, 0, 4)
	if n.Tip() {
		taxindex, err := t.TipIndex(n.Name())
		if err != nil {
			io.ExitWithMessage(err)
		}
		output = append(output, taxindex)
	} else {
		for _, next := range n.Neigh() {
			if next != prev {
				output = append(output, postOrderQuartetSet(t, next, n, right)...)
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
func preOrderQuartetSet(t *Tree, n *Node, prev *Node, left [][]uint, right [][]uint) {
	for _, next := range n.Neigh() {
		if next != prev {
			for _, next2 := range n.Neigh() {
				// We append right of childs of n other than next
				// to left of next
				if next2 != prev {
					if next2 != next {
						left[next.Id()] = append(left[next.Id()], right[next2.Id()]...)
					}
				} else {
					// We append left of n to left of next
					left[next.Id()] = append(left[next.Id()], left[n.Id()]...)
				}
			}
			preOrderQuartetSet(t, next, n, left, right)
		}
	}
}

/**
Initializes a new empty quartet with 4 sets of taxa
*/
func NewQuartetSet() *QuartetSet {
	return &QuartetSet{
		left:  make([][]uint, 0, 4),
		right: make([][]uint, 0, 4),
	}
}
