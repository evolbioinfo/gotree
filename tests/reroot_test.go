package tests

import (
	"fmt"
	"github.com/evolbioinfo/gotree/tree"
	"testing"
)

/*
Generates a 1000 tip random tree, then reroot it at each tip
and compare all bipartitions of the rerooted tree with the original tree
*/
func TestRootOutgroup(t *testing.T) {
	tr, err := tree.RandomYuleBinaryTree(1000, true)
	clone := tr.Clone()
	if err != nil {
		t.Error(err)
	}
	edges := tr.Edges()
	index := tree.NewEdgeIndex(uint64(len(edges)*2), 0.75)
	for i, e := range edges {
		index.PutEdgeValue(e, i, e.Length())
	}
	tips := tr.Tips()

	for _, tip := range tips {
		err = clone.RerootOutGroup(false, true, tip.Name())
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

/*
Generates a 1000 tip ROOTED random tree, then reroot it at each tip
and compare all bipartitions of the rerooted tree with the original tree
*/
func TestReRootOutgroupRemove(t *testing.T) {
	tr, err := tree.RandomYuleBinaryTree(1000, true)
	edges := tr.Edges()
	tips := tr.Tips()
	nodes := tr.Nodes()

	for _, tip := range tips {
		clone := tr.Clone()
		if err != nil {
			t.Error(err)
		}

		err = clone.RerootOutGroup(true, true, tip.Name())

		if !clone.Rooted() {
			t.Error("Output tree should be rooted")
		}
		if len(clone.Edges()) != len(edges)-2 {
			t.Error(fmt.Sprintf("Rerooted tree should have %d edges and has %d", len(edges)-1, len(clone.Edges())))
		}
		if len(clone.Nodes()) != len(nodes)-2 {
			t.Error(fmt.Sprintf("Rerooted tree should have %d nodes and has %d", len(nodes)-1, len(clone.Nodes())))
		}

		for _, t2 := range clone.Tips() {
			if t2.Name() == tip.Name() {
				t.Error(fmt.Sprintf("Outgroup Tip %s should not be present in the rerooted tree ", t2.Name()))
			}
		}
	}
}

/*
Generates a 1000 tip UNROOTED random tree, then reroot it at each tip
and compare all bipartitions of the rerooted tree with the original tree
*/
func TestReRootOutgroupRemoveUnRooted(t *testing.T) {
	tr, err := tree.RandomYuleBinaryTree(1000, false)
	edges := tr.Edges()
	tips := tr.Tips()
	nodes := tr.Nodes()

	for _, tip := range tips {
		clone := tr.Clone()
		if err != nil {
			t.Error(err)
		}

		err = clone.RerootOutGroup(true, true, tip.Name())

		if !clone.Rooted() {
			t.Error("Output tree should be rooted")
		}
		if len(clone.Edges()) != len(edges)-1 {
			t.Error(fmt.Sprintf("Rerooted tree should have %d edges and has %d", len(edges)-1, len(clone.Edges())))
		}
		if len(clone.Nodes()) != len(nodes)-1 {
			t.Error(fmt.Sprintf("Rerooted tree should have %d nodes and has %d", len(nodes)-1, len(clone.Nodes())))
		}
		for _, t2 := range clone.Tips() {
			if t2.Name() == tip.Name() {
				t.Error(fmt.Sprintf("Outgroup Tip %s should not be present in the rerooted tree ", t2.Name()))
			}
		}
	}
}
