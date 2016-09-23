package tests

import (
	"bufio"
	"github.com/fredericlemoine/gotree/io/newick"
	"github.com/fredericlemoine/gotree/tree"
	"os"
	"testing"
)

var prevtree *tree.Tree

func benchmarkTreeParse(nbtips int, b *testing.B) {
	var t *tree.Tree
	for n := 0; n < b.N; n++ {
		var err error
		t, err = tree.RandomBinaryTree(nbtips)
		if err != nil {
			b.Error(err)
		}
		f, err2 := os.Create("/tmp/tree")
		if err2 != nil {
			b.Error(err2)
		} else {
			f.WriteString(t.Newick() + "\n")
			f.Close()
		}
		fi, err3 := os.Open("/tmp/tree")
		if err3 != nil {
			b.Error(err3)
		}
		r := bufio.NewReader(fi)
		_, err4 := newick.NewParser(r).Parse()
		if err4 != nil {
			b.Error(err4)
		}
		if err5 := fi.Close(); err5 != nil {
			b.Error(err5)
		}
	}
	prevtree = t
}

func BenchmarkTreeParse10(b *testing.B)     { benchmarkTreeParse(10, b) }
func BenchmarkTreeParse100(b *testing.B)    { benchmarkTreeParse(100, b) }
func BenchmarkTreeParse1000(b *testing.B)   { benchmarkTreeParse(1000, b) }
func BenchmarkTreeParse10000(b *testing.B)  { benchmarkTreeParse(10000, b) }
func BenchmarkTreeParse100000(b *testing.B) { benchmarkTreeParse(100000, b) }
