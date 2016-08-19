package tree

import (
	"errors"
	"github.com/fredericlemoine/bitset"
)

type Edge struct {
	left, right *Node   // Left and right nodes
	length      float64 // length of branch
	support     float64 // -1 if no support
	// a Bit at index i in the bitset corresponds to the position of the tip i
	//left:0/right:1 .
	// i is the index of the tip in the sorted tip name array
	bitset *bitset.BitSet // Bitset of length Number of taxa each
}

/* Edge functions */
/******************/

func (e *Edge) setLeft(left *Node) {
	e.left = left
}
func (e *Edge) setRight(right *Node) {
	e.right = right
}
func (e *Edge) SetLength(length float64) {
	e.length = length
}

func (e *Edge) SetSupport(support float64) {
	e.support = support
}

func (e *Edge) Length() float64 {
	return e.length
}

func (e *Edge) Support() float64 {
	return e.support
}

func (e *Edge) Right() *Node {
	return e.right
}

func (e *Edge) Left() *Node {
	return e.left
}

func (e *Edge) TopoDepth() (int, error) {
	if e.bitset == nil {
		return -1, errors.New("Cannot compute topodepth, Bitset is nil")
	}
	if e.bitset.None() {
		return -1, errors.New("Cannot compute topodepth, Bitset is 000...0")
	}
	count := int(e.bitset.Count())
	total := int(e.bitset.Len())
	return min(count, total-count), nil
}

func (e *Edge) DumpBitSet() string {
	if e.bitset == nil {
		return "nil"
	}
	return e.bitset.DumpAsBits()
}

func (e *Edge) SameBipartition(e2 *Edge) bool {
	return e.bitset.EqualOrComplement(e2.bitset)
}

// Tests wether the tip with index id in the bitset
// is Set or not
// The index corresponds to tree.Tipindex(tipname)
func (e *Edge) TipPresent(id uint) bool {
	return e.bitset.Test(id)
}

func (e *Edge) NumTips() uint {
	return e.bitset.Count()
}

// Return the given edge in the array of edges comparing bitsets fields
// Return nil if not found
func (e *Edge) FindEdge(edges []*Edge) (*Edge, error) {
	if e.bitset == nil {
		return nil, errors.New("BitSets has not been initialized with tree.clearBitSetsRecur(nil, nil, uint(len(tree.tipIndex)))")
	}
	if e.bitset.None() {
		return nil, errors.New("One edge has a bitset of 0...000 : May be BitSets have not been updated with tree.UpdateBitSet()?")
	}
	for _, e2 := range edges {
		if e2.bitset == nil {
			return nil, errors.New("BitSets has not been initialized with tree.clearBitSetsRecur(nil, nil, uint(len(tree.tipIndex)))")
		}

		if e.Right().Tip() != e2.Right().Tip() {
			continue
		}
		// If we take all the edges, or if both edges are not tips
		if e.bitset.EqualOrComplement(e2.bitset) {
			if e2.bitset.None() {
				return nil, errors.New("One edge has a bitset of 0...000 : May be BitSets have not been updated with tree.UpdateBitSet()?")
			}
			return e, nil
		}
	}
	return nil, nil
}
