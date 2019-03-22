package tree

import (
	"errors"

	"github.com/fredericlemoine/bitset"
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

// The key for an edge is its bitset
type EdgeKey struct {
	key *bitset.BitSet
}

// HashCode for an Edge, computed from its bitset.
//
// Used for insertion in an HashMap
func (k *EdgeKey) HashCode() int64 {
	var hashCodeSet int64 = 1
	var hashCodeUnset int64 = 1
	var hashCodeAll int64 = 1
	nbset := 0
	nbunset := 0
	var bit uint
	for bit = 0; bit < k.key.Len(); bit++ {
		if k.key.Test(bit) {
			hashCodeSet = 31*hashCodeSet + int64(bit)
			nbset++
		} else {
			hashCodeUnset = 31*hashCodeUnset + int64(bit)
			nbunset++
		}
		hashCodeAll = 31*hashCodeAll + int64(bit)
	}
	// If the number of species on the left is the same
	// than the number of species on the right
	// We return the hashcode of the all species
	// Otherwise, we return the hashcode for the minimum
	// between left and right
	// Allows an edge to be kind of "unique"
	if nbset == nbunset {
		return hashCodeAll
	} else if nbset < nbunset {
		return hashCodeSet
	}
	return hashCodeUnset
}

// HashCode for an edge bitset.
//
// Used for insertion in an EdgeMap
func (k *EdgeKey) HashEquals(h hashmap.Hasher) bool {
	return k.key.EqualOrComplement(h.(*EdgeKey).key)
}

// Value stored in the HashMap
type EdgeIndexInfo struct {
	Count int     // Number of occurences of the branch
	Len   float64 // Mean length of branches occurences
}

// KeyValue Pair stored in the HashMap
type KeyValue struct {
	key *bitset.BitSet
	val *EdgeIndexInfo
}

// Initializes an Edge Count Index
func NewEdgeIndex(size int64, loadfactor float64) *EdgeIndex {
	return &EdgeIndex{
		hashmap.NewHashMap(size, loadfactor),
	}
}

// Returns the count for the given Edge
//	* If the edge is not present, returns 0 and false
//	* If the edge is present, returns the value and true
func (em *EdgeIndex) Value(e *Edge) (*EdgeIndexInfo, bool) {
	v, ok := em.hash.Value(&EdgeKey{e.Bitset()})
	if ok {
		return v.(*EdgeIndexInfo), ok
	} else {
		return nil, false
	}
}

// Increments edge count for an edge if it already exists in the map.
// If it does not exist, adds it with count 1
func (em *EdgeIndex) AddEdgeCount(e *Edge) error {
	if e.Bitset() == nil {
		io.LogError(errors.New("Bitset not initialized"))
		return errors.New("Bitset not initialized")
	}
	v, ok := em.hash.Value(&EdgeKey{e.Bitset()})
	if !ok {
		em.hash.PutValue(&EdgeKey{e.Bitset()}, &EdgeIndexInfo{1, e.Length()})
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
	em.hash.PutValue(&EdgeKey{e.Bitset()}, &EdgeIndexInfo{count, length})
	return nil
}

// Returns all the Bipartitions of the index (bitset) with their counts
// included in ]min,max]. If min==Max==1 : [1].
//
// Keys of the index
func (em *EdgeIndex) BitSets(minCount, maxCount int) []*KeyValue {
	keyvalues := em.hash.KeyValues()
	bitsets := make([]*KeyValue, 0, len(keyvalues))
	for _, kv := range keyvalues {
		b := kv.Key.(*EdgeKey).key
		v := (kv.Value).(*EdgeIndexInfo)
		if (v.Count > minCount && v.Count <= maxCount) || v.Count == maxCount {
			bitsets = append(bitsets, &KeyValue{b, v})
		}
	}
	return bitsets
}
