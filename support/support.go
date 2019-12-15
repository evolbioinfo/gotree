package support

type bootval struct {
	value   int
	edgeid  int
	randsup bool
}

type speciesmoved struct {
	taxid   uint
	nbtimes float64
}

func min(a int, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a int, b int) int {
	if a < b {
		return b
	}
	return a
}

func min_uint(a uint16, b uint16) uint16 {
	if a < b {
		return a
	}
	return b
}
