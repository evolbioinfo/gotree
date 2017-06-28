package tests

import (
	"github.com/fredericlemoine/gotree/tree"
	"testing"
)

/*
Generates a 1000 tip random tree, then reroot it at each tip
and compare all bipartitions of the rerooted tree with the original tree
*/
func TestRootOutgroup(t *testing.T) {
	tr, err := tree.RandomYuleBinaryTree(1000, true)
	tr.ReinitIndexes()
	clone := tr.Clone()
	if err != nil {
		t.Error(err)
	}
	edges := tr.Edges()
	index := tree.NewEdgeIndex(int64(len(edges)*2), 0.75)
	for i, e := range edges {
		index.PutEdgeValue(e, i, e.Length())
	}
	tips := tr.Tips()

	for _, tip := range tips {
		err = clone.RerootOutGroup(tip.Name())
		found := false
		for _, n := range clone.Root().Neigh() {
			if n.Tip() && n.Name() == tip.Name() {
				found = true
			}
		}
		if !found {
			t.Error("Outgroup (tip) not found in the children of the root on the rerooted tree")
		}
		edges2 := clone.Edges()
		// Check wether the 2 trees have the same set of tip names
		if err = tr.CompareTipIndexes(clone); err != nil {
			t.Error(err)
		}

		for _, e2 := range edges2 {
			_, ok := index.Value(e2)
			if !ok {
				t.Error("An edge of the original tree is not found in the rerooted tree")
			}
		}
	}
}
