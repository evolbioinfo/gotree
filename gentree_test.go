package gotree

import (
	"github.com/fredericlemoine/gotree/tree"
	"testing"
)

var prevtree *tree.Tree

func benchmarkBinaryTreeGeneration(nbtips int, b *testing.B) {
	var t *tree.Tree
	for n := 0; n < b.N; n++ {
		var err error
		t, err = tree.RandomBinaryTree(nbtips)
		if err != nil {
			b.Error(err)
		}
	}
	prevtree = t
}

func BenchmarkBinaryTreeGeneration10(b *testing.B)     { benchmarkBinaryTreeGeneration(10, b) }
func BenchmarkBinaryTreeGeneration100(b *testing.B)    { benchmarkBinaryTreeGeneration(100, b) }
func BenchmarkBinaryTreeGeneration1000(b *testing.B)   { benchmarkBinaryTreeGeneration(1000, b) }
func BenchmarkBinaryTreeGeneration10000(b *testing.B)  { benchmarkBinaryTreeGeneration(10000, b) }
func BenchmarkBinaryTreeGeneration100000(b *testing.B) { benchmarkBinaryTreeGeneration(100000, b) }
func BenchmarkBinaryTreeGeneration200000(b *testing.B) { benchmarkBinaryTreeGeneration(200000, b) }
