package tree

import (
	"errors"
	"math"
	"sort"
)

// if UpdateTipIndex has been called before ok
// otherwise returns an error
func (t *Tree) NbTips() (int, error) {
	if len(t.tipIndex) == 0 {
		return 0, errors.New("No tips in the index, tip name index is not initialized")
	}

	return len(t.tipIndex), nil

}

func (t *Tree) SumBranchLengths() float64 {
	sumlen := 0.0
	for _, e := range t.Edges() {
		if e.Length() == NIL_LENGTH {
			return math.NaN()
		}
		sumlen += e.Length()
	}
	return sumlen
}

func (t *Tree) MeanBranchLength() float64 {
	mean := 0.0
	edges := t.Edges()
	for _, e := range edges {
		if e.Length() == NIL_LENGTH {
			return math.NaN()
		}
		mean += e.Length()
	}
	return mean / float64(len(edges))
}

func (t *Tree) MeanSupport() float64 {
	mean := 0.0
	edges := t.Edges()
	i := 0
	for _, e := range edges {
		if !e.Right().Tip() {
			if e.Support() == NIL_SUPPORT {
				return math.NaN()
			}
			mean += e.Support()
			i++
		}
	}

	return mean / float64(i)
}

func (t *Tree) MedianSupport() float64 {
	edges := t.Edges()
	tips := t.Tips()
	supports := make([]float64, len(edges)-len(tips))
	if len(supports) == 0 {
		return math.NaN()
	}
	i := 0
	for _, e := range edges {
		if !e.Right().Tip() {
			if e.Support() == NIL_SUPPORT {
				return math.NaN()
			}
			supports[i] = e.Support()
			i++
		}
	}
	sort.Float64s(supports)

	middle := len(supports) / 2
	result := supports[middle]
	if len(supports)%2 == 0 {
		result = (result + supports[middle-1]) / 2
	}
	return result
}

func (t *Tree) NbCherries() (nbcherries int) {
	nbcherries = 0
	for _, n := range t.Nodes() {
		nbtips := 0
		nbchilds := 0
		for _, c := range n.Neigh() {
			if c.Tip() {
				nbtips++
			}
			nbchilds++
		}
		if nbtips == 2 && nbchilds == 3 {
			nbcherries++
		}
	}
	return
}
