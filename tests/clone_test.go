package tests

import (
	"github.com/fredericlemoine/gotree/tree"
	"testing"
)

/*
Generates 100 random 1000 tip trees, clone them, and compare them to the original trees
*/
func TestCloneTree(t *testing.T) {
	for i := 0; i < 100; i++ {
		tr, err := tree.RandomYuleBinaryTree(1000, true)
		tr.ReinitIndexes()
		clone := tr.Clone()

		// Comparing tip names
		tips := tr.Tips()
		copyTips := clone.Tips()
		for i, _ := range tips {
			if tips[i].Name() != copyTips[i].Name() {
				t.Error("A tip is not found in the cloned tree")
			}
		}

		// Check wether the 2 trees have the same set of tip names
		if err = tr.CompareTipIndexes(clone); err != nil {
			t.Error(err)
		}

		// Comparing edges
		edges := tr.Edges()
		edges2 := clone.Edges()
		index := tree.NewEdgeIndex(int64(len(edges)*2), 0.75)
		for i, e := range edges {
			index.PutEdgeValue(e, i, e.Length())
		}
		for _, e2 := range edges2 {
			_, ok := index.Value(e2)
			if !ok {
				t.Error("An edge of the original tree is not found in the cloned tree")
			}
		}
	}
}
