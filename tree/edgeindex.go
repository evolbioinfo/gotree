package tree

import (
	"github.com/fredericlemoine/bitset"
)

type EdgeIndex struct {
	idx *bitset.BitSetIndex
}

// Initializes an Edge Count Index
func NewEdgeIndex(size int64, loadfactor float64) *EdgeIndex {
	return &EdgeIndex{
		bitset.NewBitSetIndex(size, loadfactor),
	}
}

// Returns the count for the given Edge
// If the edge is not present, returns 0 and false
// If the edge is present, returns the value and true
func (em *EdgeIndex) Value(e *Edge) (int, bool) {
	return em.idx.Value(e.Bitset())
}

// Increment edge count for an edge if it already exists in the map
// Otherwise adds it with count 1
func (em *EdgeIndex) AddEdgeCount(e *Edge) {
	em.idx.AddCount(e.Bitset())
}

// Adds the edge in the map, with given value
// If the edge already exists in the index
// The old value is erased
func (em *EdgeIndex) PutEdgeValue(e *Edge, value int) {
	em.idx.PutValue(e.Bitset(), value)
}
