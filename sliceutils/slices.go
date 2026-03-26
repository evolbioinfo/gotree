package sliceutils

import "sort"

// sortIndicesByValues returns indices sorted by their corresponding values in ascending order
func SortIndicesByValues(values []float64) []int {
	indices := make([]int, len(values))
	for i := range indices {
		indices[i] = i
	}
	sort.Slice(indices, func(i, j int) bool {
		return values[indices[i]] < values[indices[j]]
	})
	return indices
}

// This function returns the jacquard similarity between two slices of ints, considered as sets (i.e. it is the size of the intersection divided by the size of the union)
func Jacquard(slice1, slice2 []int) float64 {
	intersection := len(Intersect(slice1, slice2))
	union := len(Union(slice1, slice2))
	if union == 0 {
		return 1.0
	}
	return float64(intersection) / float64(union)
}

// This function returns the intersection of two slices of ints, considered as sets (i.e. it is the slice of the values that are present in both slices)
func Intersect(slice1, slice2 []int) []int {
	set := make(map[int]bool)
	for _, v := range slice1 {
		set[v] = true
	}
	var intersection []int
	for _, v := range slice2 {
		if set[v] {
			intersection = append(intersection, v)
		}
	}
	return intersection
}

// This function returns the union of two slices of ints, considered as sets (i.e. it is the slice of the values that are present in at least one of the slices)
func Union(slice1, slice2 []int) []int {
	set := make(map[int]bool)
	for _, v := range slice1 {
		set[v] = true
	}
	for _, v := range slice2 {
		set[v] = true
	}
	var union []int
	for v := range set {
		union = append(union, v)
	}
	return union
}

// This function returns the union of two slices of ints, considered as sets (i.e. it is the slice of the values that are present in at least one of the slices)
func UnionStr(slice1, slice2 []string) []string {
	set := make(map[string]bool)
	for _, v := range slice1 {
		set[v] = true
	}
	for _, v := range slice2 {
		set[v] = true
	}
	var union []string
	for v := range set {
		union = append(union, v)
	}
	return union
}
