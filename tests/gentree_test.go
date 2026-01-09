package tests

import (
	"math/rand"
	"testing"

	"github.com/evolbioinfo/gotree/tree"
)

var prevtree2 *tree.Tree

func benchmarkBinaryTreeGeneration(nbtips int, b *testing.B) {
	var t *tree.Tree
	for n := 0; n < b.N; n++ {
		var err error
		t, err = tree.RandomUniformBinaryTree(nbtips, false, rand.New(rand.NewSource(10)))
		if err != nil {
			b.Error(err)
		}
	}
	prevtree2 = t
}

func BenchmarkBinaryTreeGeneration10(b *testing.B)     { benchmarkBinaryTreeGeneration(10, b) }
func BenchmarkBinaryTreeGeneration100(b *testing.B)    { benchmarkBinaryTreeGeneration(100, b) }
func BenchmarkBinaryTreeGeneration1000(b *testing.B)   { benchmarkBinaryTreeGeneration(1000, b) }
func BenchmarkBinaryTreeGeneration10000(b *testing.B)  { benchmarkBinaryTreeGeneration(10000, b) }
func BenchmarkBinaryTreeGeneration100000(b *testing.B) { benchmarkBinaryTreeGeneration(100000, b) }
func BenchmarkBinaryTreeGeneration200000(b *testing.B) { benchmarkBinaryTreeGeneration(200000, b) }
