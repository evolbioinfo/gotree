/**
Package with simple math util functions such as min, max, and abs
*/
package mutils

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func Max(a, b int) int {
	if a < b {
		return b
	}
	return a
}

func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
