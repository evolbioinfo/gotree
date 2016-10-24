package tree

import (
	"github.com/fredericlemoine/bitset"
)

type EdgeIndex struct {
	idx *bitset.BitSetIndex
}

type EdgeIndexInfo struct {
	Count int     // Number of occurences of the branch
	Len   float64 // Mean length of branches occurences
}

type KeyValue struct {
	key *bitset.BitSet
	val *EdgeIndexInfo
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
func (em *EdgeIndex) Value(e *Edge) (*EdgeIndexInfo, bool) {
	v, ok := em.idx.Value(e.Bitset())
	if ok {
		return v.(*EdgeIndexInfo), ok
	} else {
		return nil, false
	}
}

// Increment edge count for an edge if it already exists in the map
// Otherwise adds it with count 1
func (em *EdgeIndex) AddEdgeCount(e *Edge) {
	v, ok := em.idx.Value(e.Bitset())
	if !ok {
		em.idx.PutValue(e.Bitset(), &EdgeIndexInfo{1, e.Length()})
	} else {
		v.(*EdgeIndexInfo).Count++
		v.(*EdgeIndexInfo).Len += e.Length()
	}
	//em.idx.AddCount(e.Bitset())
}

// Adds the edge in the map, with given value
// If the edge already exists in the index
// The old value is erased
func (em *EdgeIndex) PutEdgeValue(e *Edge, count int, length float64) {
	em.idx.PutValue(e.Bitset(), &EdgeIndexInfo{count, length})
}

// Returns all the Bipartitions of the index (bitset) with their counts
// That have a count included in ]min,max]. If min==Max==1 : [1]
// Keys of the index
func (em *EdgeIndex) BitSets(minCount, maxCount int) []*KeyValue {
	keyvalues := em.idx.KeyValues()
	bitsets := make([]*KeyValue, 0, len(keyvalues))
	for _, kv := range keyvalues {
		b := kv.Key
		v := (kv.Value).(*EdgeIndexInfo)
		if (v.Count > minCount && v.Count <= maxCount) || v.Count == maxCount {
			bitsets = append(bitsets, &KeyValue{b, v})
		}
	}
	return bitsets
}
