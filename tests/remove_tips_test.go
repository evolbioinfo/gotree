package tests

import (
	"github.com/fredericlemoine/gotree/tree"
	"testing"
)

/*
Generates a 1000 tip random tree, then reroot it at each tip
and compare all bipartitions of the rerooted tree with the original tree
*/
func TestRemoveTips(t *testing.T) {
	tr, err := tree.RandomYuleBinaryTree(1000, true)
	t2 := tr.Clone()

	toremove := []string{"Tip1", "Tip2", "Tip3"}

	nodeindex, err2 := tree.NewNodeIndex(tr)
	if err2 != nil {
		t.Error(err2)
	}
	for _, tr := range toremove {
		_, ok := nodeindex.GetNode(tr)
		if !ok {
			t.Error("The tip " + tr + " does not exist in the tree")
		}
	}

	if err = tr.RemoveTips(false, toremove...); err != nil {
		t.Error(err)
	}

	if nodeindex, err2 = tree.NewNodeIndex(tr); err2 != nil {
		t.Error(err2)
	}
	for _, tr := range toremove {
		_, ok := nodeindex.GetNode(tr)
		if ok {
			t.Error("The tip " + tr + " should not exist anymore in the tree")
		}
	}

	if err = t2.RemoveTips(true, toremove...); err != nil {
		t.Error(err)
	}
	if nodeindex, err2 = tree.NewNodeIndex(t2); err2 != nil {
		t.Error(err2)
	}
	for _, tr := range toremove {
		_, ok := nodeindex.GetNode(tr)
		if !ok {
			t.Error("The tip " + tr + " should still exist in the tree")
		}
	}
	if len(t2.Tips()) != 3 {
		t.Error("There should be 3 tips left in the tree")
	}
}
