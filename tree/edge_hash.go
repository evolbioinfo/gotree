package tree

import (
	"fmt"
	"os"

	"github.com/evolbioinfo/gotree/hashmap"
)

func (t *Tree) ComputeEdgeHashes(cur, prev *Node, e *Edge) {
	if cur == nil {
		cur = t.Root()
	}
	t.computeEdgeHashesRightRecur(cur, nil, nil)
	t.computeEdgeHashesLeftRecur(cur, nil, nil)
}

func (t *Tree) computeEdgeHashesRightRecur(cur, prev *Node, e *Edge) {
	if e != nil {
		e.ntaxright = 0
		e.hashcoderight = 0
	}
	if cur.Tip() {
		tipIndex, _ := t.TipIndex(cur.Name())
		e.hashcoderight = 31 + 31*int64(tipIndex)
		e.ntaxright++
	} else {
		for i, n := range cur.Neigh() {
			if n != prev {
				nextEdge := cur.Edges()[i]
				t.computeEdgeHashesRightRecur(n, cur, nextEdge)
				if e != nil {
					e.hashcoderight += nextEdge.hashcoderight
					e.ntaxright += nextEdge.ntaxright
				}
			}
		}
	}
}

func (t *Tree) computeEdgeHashesLeftRecur(cur, prev *Node, e *Edge) {
	if e != nil {
		e.ntaxleft = 0
		e.hashcodeleft = 0
		// We traverse other edges than e connected to prev
		for i, n := range prev.Neigh() {
			if n != cur {
				prevE := prev.Edges()[i]
				// Descending edge
				if n == prevE.Right() {
					e.hashcodeleft += prevE.hashcoderight
					e.ntaxleft += prevE.ntaxright
				} else {
					// Ascending edge
					if n == prevE.Left() {
						e.hashcodeleft += prevE.hashcodeleft
						e.ntaxleft += prevE.ntaxleft
					} else {
						fmt.Fprintf(os.Stderr, "Error: The edge is not oriented as it should be")
					}
				}
			}
		}
	}
	for i, n := range cur.Neigh() {
		if n != prev {
			nextEdge := cur.Edges()[i]
			t.computeEdgeHashesLeftRecur(n, cur, nextEdge)
		}
	}
}

// HashCode for an Edge, computed from its bitset.
//
// Used for insertion in an HashMap
// If the bitsets are not initialized, then returns 0
func (e *Edge) HashCode() int64 {
	var hashcode int64 = 0
	if e.ntaxleft == e.ntaxright {
		hashcode = 31 * (e.hashcodeleft + e.hashcoderight)
	} else if e.ntaxleft < e.ntaxright {
		hashcode = 31 * e.hashcodeleft
	} else {
		hashcode = 31 * e.hashcoderight
	}
	return hashcode
}

// HashCode for an edge bitset.
//
// Used for insertion in an EdgeMap
func (e *Edge) HashEquals(h hashmap.Hasher) bool {
	return e.bitset.EqualOrComplement(h.(*Edge).bitset)
}
