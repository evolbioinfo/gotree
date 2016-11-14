package hashmap

import (
	"fmt"
	"testing"
)

type PairKey struct {
	i1, i2 int
}

type FloatValue struct {
	val float64
}

// Does not depend on the order
func (k *PairKey) HashEquals(h2 Hasher) bool {
	k2 := h2.(*PairKey)
	return (k.i1 == k2.i1 && k.i2 == k2.i2) ||
		(k.i1 == k2.i2 && k.i2 == k2.i1)
}

// Does not depend on the order
func (k *PairKey) HashCode() int64 {
	var hashCode int64 = 1
	if k.i1 < k.i2 {
		hashCode = 31*hashCode + int64(k.i1)
		hashCode = 31*hashCode + int64(k.i2)
	} else {
		hashCode = 31*hashCode + int64(k.i2)
		hashCode = 31*hashCode + int64(k.i1)
	}
	return hashCode
}

func TestHashMap(t *testing.T) {
	index := NewHashMap(128, .75)
	for i := 0; i < 20000; i++ {
		k := &PairKey{i, i + 1}
		val, ok := index.Value(k)
		if ok {
			t.Error(fmt.Sprintf("Key should not already be present in the map : %f", val.(*FloatValue).val))
		}
		v := &FloatValue{float64(i+i*2) / 2}
		index.PutValue(k, v)
	}
	for i := 0; i < 20000; i++ {
		// Original Key
		k2 := &PairKey{i, i + 1}
		val, ok := index.Value(k2)
		expected := float64(i+i*2) / 2
		if val.(*FloatValue).val != expected || !ok {
			t.Error(fmt.Sprintf("Value must be == %f and is %f", expected, val.(*FloatValue).val))
		}
	}
	for i := 0; i < 20000; i++ {
		// Reversed Key
		k2 := &PairKey{i + 1, i}
		val, ok := index.Value(k2)
		expected := float64(i+i*2) / 2
		if val.(*FloatValue).val != expected || !ok {
			t.Error(fmt.Sprintf("Value must be == %f and is %f", expected, val.(*FloatValue).val))
		}
	}
}
