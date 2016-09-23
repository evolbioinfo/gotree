package tree

import (
	"github.com/fredericlemoine/bitset"
	"sync"
)

type EdgeIndex struct {
	mapArray   []Bucket
	capacity   int64
	loadfactor float64
	total      int
	sync.RWMutex
}

type Bucket []*KeyValue
type KeyValue struct {
	key   *bitset.BitSet
	value int
}

// HashCode for an edge.
// Used for insertion in an EdgeMap
func hashCode(b *bitset.BitSet) int64 {
	var hashCodeSet int64 = 1
	var hashCodeUnset int64 = 1
	var hashCodeAll int64 = 1
	nbset := 0
	nbunset := 0
	var bit uint
	for bit = 0; bit < b.Len(); bit++ {
		if b.Test(bit) {
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
// Used for insertion in an EdgeMap
func equals(b *bitset.BitSet, b2 *bitset.BitSet) bool {
	return b.EqualOrComplement(b2)
}

// Initializes an Edge Count Index
func NewEdgeIndex(size int64, loadfactor float64) *EdgeIndex {
	return &EdgeIndex{
		mapArray:   make([]Bucket, size),
		capacity:   size,
		loadfactor: loadfactor,
		total:      0,
	}
}

// Returns the count for the given Edge
// If the edge is not present, returns 0 and false
// If the edge is present, returns the value and true
func (em *EdgeIndex) Value(e *Edge) (int, bool) {
	index := indexFor(hashCode(e.Bitset()), em.capacity)
	em.RLock()
	defer em.RUnlock()

	if em.mapArray[index] != nil {
		for _, kv := range em.mapArray[index] {
			if equals(kv.key, e.Bitset()) {
				return kv.value, true
			}
		}
	}
	return 0, false
}

// Increment edge count for an edge if it already exists in the map
// Otherwise adds it with count 1
func (em *EdgeIndex) AddEdgeCount(e *Edge) {
	index := indexFor(hashCode(e.Bitset()), em.capacity)
	em.Lock()
	defer em.Unlock()

	if em.mapArray[index] == nil {
		em.mapArray[index] = make(Bucket, 1, 5)
		em.mapArray[index][0] = &KeyValue{e.Bitset(), 1}
		em.total++
	} else {
		for _, kv := range em.mapArray[index] {
			if equals(kv.key, e.Bitset()) {
				kv.value++
				return
			}
		}
		em.mapArray[index] = append(em.mapArray[index], &KeyValue{e.Bitset(), 1})
		em.total++
	}
	em.rehash()
}

// Adds the edge in the map, with given value
// If the edge already exists in the index
// The old value is erased
func (em *EdgeIndex) PutEdgeValue(e *Edge, value int) {
	index := indexFor(hashCode(e.Bitset()), em.capacity)
	em.Lock()
	defer em.Unlock()

	if em.mapArray[index] == nil {
		em.mapArray[index] = make(Bucket, 1, 3)
		em.mapArray[index][0] = &KeyValue{e.Bitset(), value}
		em.total++
	} else {
		for _, kv := range em.mapArray[index] {
			if equals(kv.key, e.Bitset()) {
				kv.value = value
				return
			}
		}
		em.mapArray[index] = append(em.mapArray[index], &KeyValue{e.Bitset(), value})
		em.total++
	}
	em.rehash()
}

// returns the index in the hash map, given a hashcode
func indexFor(hashcode int64, capacity int64) int64 {
	return hashcode & (capacity - 1)
}

func (em *EdgeIndex) rehash() {
	// We rehash everything with a new capacity
	if float64(em.total) >= float64(em.capacity)*em.loadfactor {
		newcapacity := em.capacity * 2
		newmap := make([]Bucket, newcapacity)
		for _, b := range em.mapArray {
			if b != nil {
				for _, kv := range b {
					index := indexFor(hashCode(kv.key), newcapacity)
					if newmap[index] == nil {
						newmap[index] = make(Bucket, 1, 5)
						newmap[index][0] = kv
					} else {
						newmap[index] = append(newmap[index], kv)
					}
				}
			}
		}
		em.capacity = newcapacity
		em.mapArray = newmap
	}
}
