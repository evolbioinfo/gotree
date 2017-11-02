package sort

import (
	"fmt"
	"testing"
)

// Tests the function to get neighboring edges of a given edges
// Testing the edges of distance 1
func TestSortIntBy1(t *testing.T) {
	keys := []int{1, 2, 3, 4, 5}
	values := []int{10, 20, 30, 80, 50}
	expkeys := []int{1, 2, 3, 5, 4}
	expvalues := []int{10, 20, 30, 80, 50}

	SortIntBy(keys, values, false)

	for i, v := range values {
		if v != expvalues[i] {
			t.Error(fmt.Sprintf("Values should not be reordered after SortIntBy function (%d: %d instead of %d)", i, v, expvalues[i]))
		}
	}

	for i, v := range keys {
		if v != expkeys[i] {
			t.Error(fmt.Sprintf("Keys are not reordered appropriately after SortIntBy function (%d: %d instead of %d)", i, v, expkeys[i]))
		}
	}
}

func TestSortIntBy2(t *testing.T) {
	keys := []int{1, 2, 3, 4, 5}
	values := []int{10, 20, 30, 80, 50}
	expkeys := []int{4, 5, 3, 2, 1}
	expvalues := []int{10, 20, 30, 80, 50}

	SortIntBy(keys, values, true)

	for i, v := range values {
		if v != expvalues[i] {
			t.Error(fmt.Sprintf("Values should not be reordered after SortIntBy function (%d: %d instead of %d)", i, v, expvalues[i]))
		}
	}

	for i, v := range keys {
		if v != expkeys[i] {
			t.Error(fmt.Sprintf("Keys are not reordered appropriately after SortIntBy function (%d: %d instead of %d)", i, v, expkeys[i]))
		}
	}
}

func TestSortIntBy3(t *testing.T) {
	keys := []int{5, 4, 3, 2, 1}
	values := []int{10, 20, 30, 80, 50}
	expkeys := []int{5, 4, 3, 1, 2}
	expvalues := []int{10, 20, 30, 80, 50}

	SortIntBy(keys, values, false)

	for i, v := range values {
		if v != expvalues[i] {
			t.Error(fmt.Sprintf("Values should not be reordered after SortIntBy function (%d: %d instead of %d)", i, v, expvalues[i]))
		}
	}

	for i, v := range keys {
		if v != expkeys[i] {
			t.Error(fmt.Sprintf("Keys are not reordered appropriately after SortIntBy function (%d: %d instead of %d)", i, v, expkeys[i]))
		}
	}

}

func TestSortIntBy4(t *testing.T) {
	keys := []int{5, 4, 3, 2, 1}
	values := []int{10, 20, 30, 80, 50}
	expkeys := []int{2, 1, 3, 4, 5}
	expvalues := []int{10, 20, 30, 80, 50}

	SortIntBy(keys, values, true)

	for i, v := range values {
		if v != expvalues[i] {
			t.Error(fmt.Sprintf("Values should not be reordered after SortIntBy function (%d: %d instead of %d)", i, v, expvalues[i]))
		}
	}

	for i, v := range keys {
		if v != expkeys[i] {
			t.Error(fmt.Sprintf("Keys are not reordered appropriately after SortIntBy function (%d: %d instead of %d)", i, v, expkeys[i]))
		}
	}

}

func TestOrderInt1(t *testing.T) {
	values := []int{10, 20, 30, 80, 50}
	expvalues := []int{10, 20, 30, 80, 50}
	expIndices := []int{0, 1, 2, 4, 3}

	indices := OrderInt(values, false)

	for i, v := range values {
		if v != expvalues[i] {
			t.Error(fmt.Sprintf("Values should not be reordered after OrderInt function (%d: %d instead of %d)", i, v, expvalues[i]))
		}
	}

	for i, v := range indices {
		if v != expIndices[i] {
			t.Error(fmt.Sprintf("Output Indices are not ordered appropriately after OrderInt function (%d: %d instead of %d)", i, v, expIndices[i]))
		}
	}

	indices = OrderInt(values, true)

	for i, v := range values {
		if v != expvalues[i] {
			t.Error(fmt.Sprintf("Values should not be reordered after OrderInt function (%d: %d instead of %d)", i, v, expvalues[i]))
		}
	}

	for i, v := range indices {
		if v != expIndices[len(expIndices)-i-1] {
			t.Error(fmt.Sprintf("Output Indices are not ordered appropriately after OrderInt function (%d: %d instead of %d)", i, v, expIndices[i]))
		}
	}

}

func TestOrderInt2(t *testing.T) {
	values := []int{10, 9, 8, 7, 6, 5, 4, 3, 2, 1, 0}
	expvalues := []int{10, 9, 8, 7, 6, 5, 4, 3, 2, 1, 0}
	expIndices := []int{10, 9, 8, 7, 6, 5, 4, 3, 2, 1, 0}

	indices := OrderInt(values, false)

	for i, v := range values {
		if v != expvalues[i] {
			t.Error(fmt.Sprintf("Values should not be reordered after OrderInt function (%d: %d instead of %d)", i, v, expvalues[i]))
		}
	}

	for i, v := range indices {
		if v != expIndices[i] {
			t.Error(fmt.Sprintf("Output Indices are not ordered appropriately after OrderInt function (%d: %d instead of %d)", i, v, expIndices[i]))
		}
	}

	indices = OrderInt(values, true)

	for i, v := range values {
		if v != expvalues[i] {
			t.Error(fmt.Sprintf("Values should not be reordered after OrderInt function (%d: %d instead of %d)", i, v, expvalues[i]))
		}
	}

	for i, v := range indices {
		if v != expIndices[len(expIndices)-i-1] {
			t.Error(fmt.Sprintf("Output Indices are not ordered appropriately after OrderInt function (%d: %d instead of %d)", i, v, expIndices[i]))
		}
	}
}

func TestOrderInt3(t *testing.T) {
	values := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	expvalues := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	expIndices := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	indices := OrderInt(values, false)

	for i, v := range values {
		if v != expvalues[i] {
			t.Error(fmt.Sprintf("Values should not be reordered after OrderInt function (%d: %d instead of %d)", i, v, expvalues[i]))
		}
	}

	for i, v := range indices {
		if v != expIndices[i] {
			t.Error(fmt.Sprintf("Output Indices are not ordered appropriately after OrderInt function (%d: %d instead of %d)", i, v, expIndices[i]))
		}
	}

	indices = OrderInt(values, true)

	for i, v := range values {
		if v != expvalues[i] {
			t.Error(fmt.Sprintf("Values should not be reordered after OrderInt function (%d: %d instead of %d)", i, v, expvalues[i]))
		}
	}

	for i, v := range indices {
		if v != expIndices[len(expIndices)-i-1] {
			t.Error(fmt.Sprintf("Output Indices are not ordered appropriately after reversed OrderInt function (%d: %d instead of %d)", i, v, expIndices[i]))
		}
	}

}
