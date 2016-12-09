/*
This package defines a generic hashmap
This hashmap can store any Hasher object as key
and any object as value.
Similarly to Java HashMap, Hasher objects must define two functions:
- HashCode and
- HashEquals
*/
package hashmap

import (
	"sync"
)

type HashMap struct {
	mapArray   []Bucket
	capacity   int64
	loadfactor float64
	total      int
	sync.RWMutex
}

type Bucket []*KeyValue
type KeyValue struct {
	Key   Hasher
	Value interface{}
}

type Hasher interface {
	HashCode() int64
	HashEquals(k Hasher) bool
}

// Initializes a HashMap
func NewHashMap(size int64, loadfactor float64) *HashMap {
	return &HashMap{
		mapArray:   make([]Bucket, size),
		capacity:   size,
		loadfactor: loadfactor,
		total:      0,
	}
}

// Returns the count for the given Edge
// If the edge is not present, returns 0 and false
// If the edge is present, returns the value and true
func (em *HashMap) Value(h Hasher) (interface{}, bool) {
	index := indexFor(h.HashCode(), em.capacity)
	em.RLock()
	defer em.RUnlock()

	if em.mapArray[index] != nil {
		for _, kv := range em.mapArray[index] {
			if h.HashEquals(kv.Key) {
				return kv.Value, true
			}
		}
	}
	return nil, false
}

func (em *HashMap) PutValue(h Hasher, value interface{}) {
	index := indexFor(h.HashCode(), em.capacity)
	em.Lock()
	defer em.Unlock()

	if em.mapArray[index] == nil {
		em.mapArray[index] = make(Bucket, 1, 3)
		em.mapArray[index][0] = &KeyValue{h, value}
		em.total++
	} else {
		for _, kv := range em.mapArray[index] {
			if h.HashEquals(kv.Key) {
				kv.Value = value
				return
			}
		}
		em.mapArray[index] = append(em.mapArray[index], &KeyValue{h, value})
		em.total++
	}
	em.rehash()
}

// returns the index in the hash map, given a hashcode
func indexFor(hashcode int64, capacity int64) int64 {
	return hashcode & (capacity - 1)
}

// Reconstructs the HashMap if the capacity is almost attained (loadfactor)
func (em *HashMap) rehash() {
	// We rehash everything with a new capacity
	if float64(em.total) >= float64(em.capacity)*em.loadfactor {
		newcapacity := em.capacity * 2
		newmap := make([]Bucket, newcapacity)
		for _, b := range em.mapArray {
			if b != nil {
				for _, kv := range b {
					index := indexFor(kv.Key.HashCode(), newcapacity)
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

/* Returns all keys of the index */
func (em *HashMap) Keys() []Hasher {
	keys := make([]Hasher, em.total)
	total := 0
	for _, b := range em.mapArray {
		if b != nil {
			for _, kv := range b {
				keys[total] = kv.Key
				total++
			}
		}
	}
	return keys
}

func (em *HashMap) KeyValues() []*KeyValue {
	keyvalues := make([]*KeyValue, em.total)
	total := 0
	for _, b := range em.mapArray {
		if b != nil {
			for _, kv := range b {
				keyvalues[total] = kv
				total++
			}
		}
	}
	return keyvalues
}
