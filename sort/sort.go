/*
This package provides a few functions to sort int arrays.
	* Sort by values of another array
	* Return permutation allowing to sort input array (like R order() function)
Inspired from https://stackoverflow.com/questions/42707252/
*/
package sort

import (
	"sort"
)

type by struct {
	Indices []int
	Values  []int
}

func (b by) Len() int           { return len(b.Values) }
func (b by) Less(i, j int) bool { return b.Values[i] < b.Values[j] }
func (b by) Swap(i, j int) {
	b.Indices[i], b.Indices[j] = b.Indices[j], b.Indices[i]
	b.Values[i], b.Values[j] = b.Values[j], b.Values[i]
}

// Sorts "toSort" array according to "byValues" Array.
//	* "toSort" array is modified after the function.
//	* "byValues" array is not modified.
//
// The sort may be in decreasing order.
func SortIntBy(toSort []int, byValues []int, decreasing bool) {
	valuescopy := make([]int, len(byValues))
	copy(valuescopy, byValues)
	if decreasing {
		sort.Sort(sort.Reverse(by{Indices: toSort, Values: valuescopy}))
	} else {
		sort.Sort(by{Indices: toSort, Values: valuescopy})
	}
}

// Returns reordered indices (permutation) by values (like order() R function).
// It may be in decreasing order.
//
// Example:
//	input := []int{10,20,30,80,50}
//	output := OrderInt(input) // {0,1,2,4,3}
//	// To iterate over the sorted array :
// 	for _,ord := range output {
// 	    input[ord]
// 	}
func OrderInt(values []int, decreasing bool) []int {
	// init initial order indices
	indices := make([]int, len(values))
	for i, _ := range indices {
		indices[i] = i
	}
	// Sort indices by values
	SortIntBy(indices, values, decreasing)
	return indices
}
