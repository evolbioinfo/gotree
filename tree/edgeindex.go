package tree

import (
	"errors"

	"github.com/evolbioinfo/gotree/hashmap"
	"github.com/evolbioinfo/gotree/io"
)

// Structure for an EdgeIndex.
// It is basically an HashMap, storing as key
// the bitset of the edge, and as value, the number
// of occurences of this edge already stored and the
// average branch lengths.
type EdgeIndex struct {
	hash *hashmap.HashMap
}

// Value stored in the HashMap
type EdgeIndexInfo struct {
	Count int     // Number of occurences of the branch
	Len   float64 // Mean length of branches occurences
}

// KeyValue Pair stored in the HashMap
type KeyValue struct {
	key *Edge
	val *EdgeIndexInfo
}

// Initializes an Edge Count Index
func NewEdgeIndex(size uint64, loadfactor float64) *EdgeIndex {
	return &EdgeIndex{
		hashmap.NewHashMap(size, loadfactor),
	}
}

// Returns the count for the given Edge
//	* If the edge is not present, returns 0 and false
//	* If the edge is present, returns the value and true
func (em *EdgeIndex) Value(e *Edge) (*EdgeIndexInfo, bool) {
	v, ok := em.hash.Value(e)
	if ok {
		return v.(*EdgeIndexInfo), ok
	} else {
		return nil, false
	}
}

// Increments edge count for an edge if it already exists in the map.
// If it does not exist, adds it with count 1
//
// Also adds edge length
func (em *EdgeIndex) AddEdgeCount(e *Edge) error {
	if e.Bitset() == nil {
		io.LogError(errors.New("Bitset not initialized"))
		return errors.New("Bitset not initialized")
	}
	v, ok := em.hash.Value(e)
	if !ok {
		em.hash.PutValue(e, &EdgeIndexInfo{1, e.Length()})
	} else {
		v.(*EdgeIndexInfo).Count++
		v.(*EdgeIndexInfo).Len += e.Length()
	}
	return nil
}

// Adds the edge in the map, with given value.
// If the edge already exists in the index
// The old value is erased
func (em *EdgeIndex) PutEdgeValue(e *Edge, count int, length float64) error {
	if e.Bitset() == nil {
		io.LogError(errors.New("Bitset not initialized"))
		return errors.New("Bitset not initialized")
	}
	em.hash.PutValue(e, &EdgeIndexInfo{count, length})
	return nil
}

// Returns all the Bipartitions of the index (bitset) with their counts
// included in ]min,max]. If min==Max==1 : [1].
//
// Keys of the index
func (em *EdgeIndex) Edges(minCount, maxCount int) []*KeyValue {
	keyvalues := em.hash.KeyValues()
	bitsets := make([]*KeyValue, 0, len(keyvalues))
	for _, kv := range keyvalues {
		e := kv.Key.(*Edge)
		v := (kv.Value).(*EdgeIndexInfo)
		if (v.Count > minCount && v.Count <= maxCount) || v.Count == maxCount {
			bitsets = append(bitsets, &KeyValue{e, v})
		}
	}
	return bitsets
}
