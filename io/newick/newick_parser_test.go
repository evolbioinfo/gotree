package newick_test

import (
	"bufio"
	"fmt"
	"github.com/fredericlemoine/gotree/io/newick"
	"github.com/fredericlemoine/gotree/tree"
	"os"
	"strings"
	"testing"
)

// Ensure the parser can parse strings into Statement ASTs.
func TestParser_ParseTree(t *testing.T) {
	goodtrees := [...]string{
		"(Tip2:1.00000,(Tip 7:1.00000,((Tip9[NODE COMMENT]:1.00000,((Tip5:1.00000,((Tip8:1.00000,Tip6:0.50000):0.50000,Tip4:0.25000):0.25000):0.50000,Tip3:0.50000):0.25000):0.25000,Node0:0.25000):0.12500):0.12500,Tip1:0.50000);",
		"(Tip2:1.00000,(Tip 7:1.00000,((Tip9:1.00000[NODE COMMENT],((Tip5:1.00000,((Tip8:1.00000,Tip6:0.50000):0.50000,Tip4:0.25000):0.25000):0.50000,Tip3:0.50000):0.25000):0.25000,Node0:0.25000):0.12500):0.12500,Tip1:0.50000);",
		"(Tip2:1.00000,(Tip 7:1.00000,((Tip9:1.00000,((Tip5:1.00000,((Tip8:1.00000,Tip6:0.50000)[NODE COMMENT]:0.50000,Tip4:0.25000):0.25000):0.50000,Tip3:0.50000):0.25000):0.25000,Node0:0.25000):0.12500):0.12500,Tip1:0.50000);",
		"(Tip2:1.00000,(Tip 7:1.00000,((Tip9:1.00000,((Tip5:1.00000,((Tip8:1.00000,Tip6:0.50000):0.50000[NODE COMMENT],Tip4:0.25000):0.25000):0.50000,Tip3:0.50000):0.25000):0.25000,Node0:0.25000):0.12500):0.12500,Tip1:0.50000);",
		"(Tip2:1.00000,(Tip 7[NODE COMMENT]:1.00000,((Tip9:1.00000,((Tip5:1.00000,((Tip8:1.00000,Tip6:0.50000):0.50000,Tip4:0.25000):0.25000):0.50000,Tip3:0.50000):0.25000):0.25000,Node0:0.25000):0.12500):0.12500,Tip1:0.50000);",
		"(Tip2:1.00000,(Tip 7:1.00000,((Tip9[NODE COMMENT]:1.00000,((Tip5:1.00000,((Tip8:1.00000,Tip6:0.50000):0.50000,Tip4:0.25000):0.25000):0.50000,Tip3:0.50000):0.25000):0.25000,Node0:0.25000):0.12500):0.12500,Tip1:0.50000);",
		"(Tip2:1.00000,(Tip 7:1.00000,((Tip9[NODE COMMENT]:1.00000,((Tip5:1.00000,((Tip8:1.00000,Tip6:0.50000)Hello:0.50000,Tip4:0.25000):0.25000):0.50000,Tip3:0.50000):0.25000):0.25000,Node0:0.25000):0.12500):0.12500,Tip1:0.50000);",
		"(Tip2:1.00000,(Tip 7:1.00000,((Tip9[NODE COMMENT]:1.00000,((Tip5:1.00000,((Tip8:1.00000,Tip6:0.50000)0.999:0.50000,Tip4:0.25000):0.25000):0.50000,Tip3:0.50000):0.25000):0.25000,Node0:0.25000):0.12500):0.12500,Tip1:0.50000);",
		"(Tip2:1.00000,(Tip 7:1.00000,((Tip9[NODE COMMENT][HELLO]:1.00000,((Tip5:1.00000,((Tip8:1.00000,Tip6:0.50000):0.50000,Tip4:0.25000)0.8:0.25000):0.50000,Tip3:0.50000):0.25000):0.25000,Node0:0.25000):0.12500):0.12500,Tip1:0.50000)qsdqsd;",
		"(Tip2:1.00000,(Tip 7:1.00000,((Tip9[NODE COMMENT][HELLO]:1.00000,((Tip5:1.00000,((Tip8:1.00000,Tip6:0.50000):0.50000,Tip4:0.25000)0.8:0.25000):0.50000,Tip3:0.50000):0.25000):0.25000,Node0:0.25000):0.12500):0.12500,Tip1:0.50000):456;",
	}
	badtrees := [...]string{
		"(Tip2:1.00000,(Tip 7:1.00000,((Tip9:1.00000,((Tip5:1.00000,((Tip8:1.00000,Tip6:0.50000):0.50000,Tip4:0.25000):0.25000):0.50000,Tip3:0.50000):0.25000):0.25000,Node0:0.25000):0.12500):0.12500,Tip1:0.50000));",
		"(Tip2:1.00000,(Tip 7:1.00000,((Tip9:1.00000[NODE COMMENT],((Tip5:1.00000,((Tip8:1.00000,Tip6:0.50000):0.50000,Tip4:0.25000):0.25000):0.50000,Tip3:0.50000):0.25000):0.25000,Node0:0.25000):0.12500):0.12500,Tip1:Hi);",
		"(Tip2:1.00000),(Tip 7:1.00000,((Tip9:1.00000,((Tip5:1.00000,((Tip8:1.00000,Tip6:0.50000)[NODE COMMENT]:0.50000,Tip4:0.25000):0.25000):0.50000,Tip3:0.50000):0.25000):0.25000,Node0:0.25000):0.12500):0.12500,Tip1:0.50000);",
		"(Tip2:1.00000,(Tip 7:1.00000,((Tip9:1.00000,((Tip5:1.00000,((Tip8:1.00000,Tip6:0.50000):0.50000[NODE COMMENT],Tip4:0.25000))))):0.25000):0.50000,Tip3:0.50000):0.25000:0.25000,Node0:0.25000):0.12500):0.12500,Tip1:0.50000);",
		"(Tip2:1.00000,(Tip 7:1.00000,((Tip9:1.00000,((Tip5:1.00000,((Tip8:1.00000,Tip6:0.50000):0.50000[NODE COMMENT],Tip4:0.25000)))))))):0.25000):0.50000,Tip3:0.50000):0.25000):0.25000,Node0:0.25000):0.12500):0.12500,Tip1:0.50000);",
		"(Tip2:1.00000,[Hello](Tip 7[NODE COMMENT]:1.00000,((Tip9:1.00000,((Tip5:1.00000,((Tip8:1.00000,Tip6:0.50000):0.50000,Tip4:0.25000):0.25000):0.50000,Tip3:0.50000):0.25000):0.25000,Node0:0.25000):0.12500):0.12500,Tip1:0.50000);",
		"(Tip2:1.00000,([Hello]Tip 7:1.00000,((Tip9[NODE COMMENT]:1.00000,((Tip5:1.00000,((Tip8:1.00000,Tip6:0.50000):0.50000,Tip4:0.25000):0.25000):0.50000,Tip3:0.50000):0.25000):0.25000,Node0:0.25000):0.12500):0.12500,Tip1:0.50000);",
		"(:1000,Tip2:1.00000,(Tip 7:1.00000,((Tip9[NODE COMMENT]:1.00000,((Tip5:1.00000,((Tip8:1.00000,Tip6:0.50000)Hello:0.50000,Tip4:0.25000):0.25000):0.50000,Tip3:0.50000):0.25000):0.25000,Node0:0.25000):0.12500):0.12500,Tip1:0.50000);",
		"(Tip2:1.00000,((((Tip 7:1.00000,((Tip9[NODE COMMENT]:1.00000,((Tip5:1.00000,((Tip8:1.00000,Tip6:0.50000)0.999:0.50000,Tip4:0.25000):0.25000):0.50000,Tip3:0.50000):0.25000):0.25000,Node0:0.25000):0.12500):0.12500,Tip1:0.50000);",
		"(Tip2:1.00000,((((Tip 7:1.00000,((Tip9[NODE COMMENT]:1.00000,((Tip5:1.00000,((Tip8:1.00000,Tip6:0.50000):0.50000,Tip4:0.25000)0.8:0.25000):0.50000,Tip3:0.50000):0.25000):0.25000,Node0:0.25000):0.12500):0.12500,Tip1:0.50000);",
		"(Tip2:1.00000,(Tip 7:1.00000,((Tip9[NODE COMMENT][HELLO]:1.00000:110,((Tip5:1.00000,((Tip8:1.00000,Tip6:0.50000):0.50000,Tip4:0.25000)0.8:0.25000):0.50000,Tip3:0.50000):0.25000):0.25000,Node0:0.25000):0.12500):0.12500,Tip1:0.50000);",
		"(Tip2:1.00000,(Tip 7:1.00000,((Tip9[NODE COMMENT]:1.00000:2.00000,((Tip5:1.00000,((Tip8:1.00000,Tip6:0.50000):0.50000,Tip4:0.25000):0.25000):0.50000,Tip3:0.50000):0.25000):0.25000,Node0:0.25000):0.12500):0.12500,Tip1:0.50000);",
	}

	for i, intree := range goodtrees {
		tree, err := newick.NewParser(strings.NewReader(intree)).Parse()
		if err != nil {
			t.Errorf("Tree %d ERROR: %s\n", i, err.Error())
		} else {
			fmt.Println(tree.Newick())
		}
	}

	for i, intree := range badtrees {
		_, err := newick.NewParser(strings.NewReader(intree)).Parse()
		if err == nil {
			t.Errorf("There should be an error with tree %d: %s", i, intree)
		} else {
			fmt.Printf("Tree %d OK: %s\n", i, err.Error())
		}
	}
}

var prevtree *tree.Tree

func benchmarkTreeParse(nbtips int, b *testing.B) {
	var t *tree.Tree
	for n := 0; n < b.N; n++ {
		var err error
		t, err = tree.RandomUniformBinaryTree(nbtips, false)
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
