package tree

import (
	"errors"
	"sort"
)

// Here are structures and functions to handle bags of tips
// Example of usage when we want to cut the some edges of the tree
// And get tips of connected components

// A bag of tips
type TipBag struct {
	tips map[string]*Node // The tips contained in the bag, by name
}

// Initialize a new TipBag
func NewTipBag() *TipBag {
	return &TipBag{
		tips: make(map[string]*Node),
	}
}

// Add a tip to the bag
//
// - If the same tip is already present: do nothing
// - If another tip of the tree with the same name is already present: returns an error.
//   it means that the tree contains several tips that have the same name
// - If t is nil or if t is not a tip: returns an error
func (tb *TipBag) AddTip(t *Node) error {
	if t == nil {
		return errors.New("Nil node given to TipBag.AddTip")
	}
	if !t.Tip() {
		return errors.New("Internal node given to TipBag.AddTip")
	}
	if n, ok := tb.tips[t.Name()]; !ok {
		tb.tips[t.Name()] = t
	} else {
		if n != t {
			return errors.New("TipBag.AddTip: TipBag already contains another tip of the tree having the same name: May be several tips have the same name?")
		}
	}
	return nil
}

// Removes all tips of the given TipBag
func (tb *TipBag) Clear() {
	tb.tips = make(map[string]*Node)
}

// List of tips in the bag
// Always returns tip in the same order (alphanumeric by tip name)
func (tb *TipBag) Tips() []*Node {
	v := make([]*Node, 0, len(tb.tips))
	names := make([]string, 0, len(tb.tips))
	for k, _ := range tb.tips {
		names = append(names, k)
	}
	sort.Strings(names)
	for _, n := range names {
		node, _ := tb.tips[n]
		v = append(v, node)
	}

	return v

}

// Size of the TipBag in number of contained tips
func (tb *TipBag) Size() int {
	return len(tb.tips)
}
