package tree

import (
	"errors"
	"github.com/evolbioinfo/gotree/hashmap"
	"github.com/evolbioinfo/gotree/io"
)

const (
	QUARTET_EQUALS = iota
	QUARTET_CONFLICT
	QUARTET_DIFF
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
	      left-----right
	        /       \
	       b1       b3
	*/
	left  [][]uint
	right [][]uint
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

/**
Structure representing a "quartets"
*/
type Quartet struct {
	/** Indexes of the taxa in the node index
	  t1       t3
	   \       /
	     -----
	   /       \
	  t2       t4
	*/
	T1, T2, T3, T4 uint
}

/** Quartets are hashed according to their
unorder taxon index, i.e:
(1,2)(3,4) == (1,3)(4,8)
=> because we do not want conflicting quartets
 in the map
*/
func (q *Quartet) HashCode() uint64 {
	i1, i2, i3, i4 := int(q.T1), int(q.T2), int(q.T3), int(q.T4)
	// We sort the tax id before computing the hashcode
	if i2 < i1 {
		i1, i2 = i2, i1
	}
	if i3 < i4 {
		i3, i4 = i4, i3
	}
	if i3 < i1 {
		i1, i3 = i3, i1
	}
	if i4 < i2 {
		i2, i4 = i4, i2
	}
	if i3 < i2 {
		i3, i2 = i2, i3
	}
	var hashCode uint64 = 1
	hashCode = 31*(31*(31*(31+uint64(i1))+uint64(i2))+uint64(i3)) + uint64(i4)
	return hashCode
}

/**
Equals returns true if quartets are equals or
conflicting => for hashing
*/
func (q *Quartet) HashEquals(h hashmap.Hasher) bool {
	q2 := h.(*Quartet)
	return q.Compare(q2) != QUARTET_DIFF
}

/**
Iterate over all the quartets of the tree, edge by edge
(t1,t2)(t3,t4)
specific: If true gives the specific quartets
	       b0       b2
	        \       /
	      left-----right
	        /       \
	       b1       b3
Else gives all the quartets
            b0-|\       /|-b2
	       | >-----< |
            b1-|/       \|-b3
*/
func (t *Tree) Quartets(specific bool, it func(q *Quartet)) {
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
		if specific {
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
		} else {
			qs.left = append(qs.left, left[e.Right().Id()])
			qs.right = append(qs.right, right[e.Right().Id()])
		}

		// if specific {
		// 	for b1 := 0; b1 < len(qs.left); b1++ {
		// 		for b2 := b1 + 1; b2 < len(qs.left); b2++ {
		// 			for b3 := 0; b3 < len(qs.right); b3++ {
		// 				for b4 := b3 + 1; b4 < len(qs.right); b4++ {
		// 					fmt.Printf("%d\n", len(qs.left[b1])*len(qs.left[b2])*len(qs.right[b3])*len(qs.right[b4]))
		// 				}
		// 			}
		// 		}
		// 	}
		// } else {
		// 	if len(qs.left) != 1 || len(qs.right) != 1 {
		// 		io.ExitWithMessage(errors.New("A non specific quartetset should have only one set at left and one set at right"))
		// 	}
		// 	fmt.Printf("%f\n", float64(len(qs.left[0])*(len(qs.left[0])-1)*len(qs.right[0])*(len(qs.right[0])-1))/4.0)
		// }
		qs.iterate(specific, it)
	}
}

// Function that enumerates all quartets defined by a quartetset
// (t1,t2)(t3,t4)
func (qs *QuartetSet) iterate(specific bool, it func(q *Quartet)) {
	if specific {
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
										it(&Quartet{qs.left[b1][tb1], qs.left[b2][tb2], qs.right[b3][tb3], qs.right[b4][tb4]})
									}
								}
							}
						}
					}
				}
			}
		}
	} else {
		if len(qs.left) != 1 || len(qs.right) != 1 {
			io.ExitWithMessage(errors.New("A non specific quartetset should have only one set at left and one set at right"))
		}
		// First Taxon of set 1
		for tb1 := 0; tb1 < len(qs.left[0]); tb1++ {
			// Second Taxon of set 1
			for tb2 := tb1 + 1; tb2 < len(qs.left[0]); tb2++ {
				// First Taxon of set 2
				for tb3 := 0; tb3 < len(qs.right[0]); tb3++ {
					// Second Taxon of set 2
					for tb4 := tb3 + 1; tb4 < len(qs.right[0]); tb4++ {
						it(&Quartet{qs.left[0][tb1], qs.left[0][tb2], qs.right[0][tb3], qs.right[0][tb4]})
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

/* Compares the first quartet
(q1,q2)(q3,q4)
With the second quartet
(q5,q6)(q7,q8)

Returns:
- QUARTET_EQUALS if they have the same taxa
  and  the same topology
- QUARTET_CONFLICT if they have the same taxa
  and different topology
- QUARTET_DIFF if they have different taxa
*/
func (q *Quartet) Compare(q2 *Quartet) int {
	//q1 (1,2)(3,4)
	//q2 (1,2)(3,4)
	// Quartets equals
	if ((q.T1 == q2.T1 && q.T2 == q2.T2) || (q.T1 == q2.T2 && q.T2 == q2.T1)) &&
		((q.T3 == q2.T3 && q.T4 == q2.T4) || (q.T3 == q2.T4 && q.T4 == q2.T3)) {
		return QUARTET_EQUALS
	}
	if ((q.T1 == q2.T3 && q.T2 == q2.T4) || (q.T1 == q2.T4 && q.T2 == q2.T3)) &&
		((q.T3 == q2.T1 && q.T4 == q2.T2) || (q.T3 == q2.T2 && q.T4 == q2.T1)) {
		return QUARTET_EQUALS
	}
	// Quartets conflict
	//(q1,q3)(q2,q4)
	//(q5,q6)(q7,q8)
	//or
	//(q1,q4)(q2,q3)
	//(q5,q6)(q7,q8)
	if ((q.T3 == q2.T1 && q.T2 == q2.T2) || (q.T3 == q2.T2 && q.T2 == q2.T1)) &&
		((q.T1 == q2.T3 && q.T4 == q2.T4) || (q.T1 == q2.T4 && q.T4 == q2.T3)) {
		return QUARTET_CONFLICT
	}
	if ((q.T3 == q2.T3 && q.T2 == q2.T4) || (q.T3 == q2.T4 && q.T2 == q2.T3)) &&
		((q.T1 == q2.T1 && q.T4 == q2.T2) || (q.T1 == q2.T2 && q.T4 == q2.T1)) {
		return QUARTET_CONFLICT
	}
	if ((q.T4 == q2.T1 && q.T2 == q2.T2) || (q.T4 == q2.T2 && q.T2 == q2.T1)) &&
		((q.T3 == q2.T3 && q.T1 == q2.T4) || (q.T3 == q2.T4 && q.T1 == q2.T3)) {
		return QUARTET_CONFLICT
	}
	if ((q.T4 == q2.T3 && q.T2 == q2.T4) || (q.T4 == q2.T4 && q.T2 == q2.T3)) &&
		((q.T3 == q2.T1 && q.T1 == q2.T2) || (q.T3 == q2.T2 && q.T1 == q2.T1)) {
		return QUARTET_CONFLICT
	}

	return QUARTET_DIFF
}

func (t *Tree) IndexQuartets(specific bool) *hashmap.HashMap {
	index := hashmap.NewHashMap(12800000, .75)
	n := 0

	t.Quartets(specific, func(q *Quartet) {
		n++
		index.PutValue(q, q)
	})
	return index
}
