package tests

import (
	"fmt"
	"github.com/fredericlemoine/gotree/io/newick"
	"github.com/fredericlemoine/gotree/tree"
	"os"
	"strings"
	"testing"
)

var unroottree string = "((1:1,2:1):1,5:1,(3:1,4:1):1);"
var startree string = "(1:1,2:1,3:1,4:1,5:1);"
var longestpathtree string = "((1:5,2:3):1,5:1,(3:5,4:10):1);"
var midpointtree string = "((1:1,2:1):4,3:10,(4:1,(5:10,6:1):2):4);"

func TestLeastCommonAncestorUnrooted(t *testing.T) {
	tree, err := newick.NewParser(strings.NewReader(unroottree)).Parse()
	if err != nil {
		t.Error(err)
	}

	n, e, mono, err := tree.LeastCommonAncestorUnrooted(nil, "3", "4")
	if err != nil {
		t.Error(err)
	}
	if n == nil {
		t.Errorf("No common ancestor found")
	}
	if e == nil || len(e) == 0 {
		t.Errorf("No common ancestor Edges found")
	}
	if !mono {
		t.Errorf("The group should be monophyletic")
	}
	if len(e) != 2 {
		t.Errorf("Edge Length should be 2 and is %d", len(e))
	}

	for _, ed := range e {
		if ed.Right().Name() != "4" &&
			ed.Right().Name() != "3" &&
			ed.Left().Name() != "4" &&
			ed.Left().Name() != "3" {
			t.Errorf("Ancestor of wrong tips found")
		}
	}

	n, e, mono, err = tree.LeastCommonAncestorUnrooted(nil, "3", "5")
	if err != nil {
		t.Error(err)
	}

	if n == nil {
		t.Errorf("No common ancestor found")
	}
	if e == nil || len(e) == 0 {
		t.Errorf("No common ancestor Edges found")
	}
	if mono {
		t.Errorf("The group should not be monophyletic")
	}

	if len(e) != 2 {
		t.Errorf("Edge Length should be 2 and is %d", len(e))
	}

	star, err2 := newick.NewParser(strings.NewReader(startree)).Parse()
	if err2 != nil {
		t.Error(err2)
	}
	n, e, mono, err = star.LeastCommonAncestorUnrooted(nil, "3", "4", "5", "2")
	if err != nil {
		t.Error(err)
	}

	if n == nil {
		t.Errorf("No common ancestor found")
	}
	if e == nil || len(e) == 0 {
		t.Errorf("No common ancestor Edges found")
	}

	if len(e) != 4 {
		t.Errorf("Edge Length should be 4 and is %d", len(e))
	}
	if !mono {
		t.Errorf("The group should be monophyletic")
	}

}

func TestAddBipartition(t *testing.T) {
	star, err2 := newick.NewParser(strings.NewReader(startree)).Parse()
	if err2 != nil {
		t.Error(err2)
	}
	n, e, mono, err := star.LeastCommonAncestorUnrooted(nil, "3", "4", "5")
	if err != nil {
		t.Error(err)
	}

	if n == nil {
		t.Errorf("No common ancestor found")
	}
	if e == nil || len(e) == 0 {
		t.Errorf("No common ancestor Edges found")
	}
	if !mono {
		t.Errorf("The group should be monophyletic")
	}

	if len(e) != 3 {
		t.Errorf("Edge Length should be 3 and is %d", len(e))
	}
	star.AddBipartition(n, e, 0.001, 0.9)
	n, e, mono, err = star.LeastCommonAncestorUnrooted(nil, "4", "5")
	if err != nil {
		t.Error(err)
	}

	if n == nil {
		t.Errorf("No common ancestor found")
	}
	if e == nil || len(e) == 0 {
		t.Errorf("No common ancestor Edges found")
	}
	if !mono {
		t.Errorf("The group should be monophyletic")
	}

	if len(e) != 2 {
		t.Errorf("Edge Length should be 2 and is %d", len(e))
	}
	fmt.Fprintf(os.Stdout, "%s\n", star.Newick())
	star.AddBipartition(n, e, 0.001, 0.9)
	fmt.Fprintf(os.Stdout, "%s\n", star.Newick())
	n, e, mono, err = star.LeastCommonAncestorUnrooted(nil, "1", "2")
	if err != nil {
		t.Error(err)
	}

	if n == nil {
		t.Errorf("No common ancestor found")
	}
	if e == nil || len(e) == 0 {
		t.Errorf("No common ancestor Edges found")
	}
	if !mono {
		t.Errorf("The group should be monophyletic")
	}

	if len(e) != 2 {
		t.Errorf("Edge Length should be 2 and is %d", len(e))
	}
	fmt.Fprintf(os.Stdout, "%s\n", star.Newick())
}

func TestMaxLengthPath(t *testing.T) {

	tr, err := newick.NewParser(strings.NewReader(longestpathtree)).Parse()
	if err != nil {
		t.Error(err)
	}

	tipstr := []string{"1", "2", "3", "4", "5"}
	expmaxlen := []float64{17, 15, 15, 17, 12}
	expmaxpath := []int{4, 4, 2, 4, 3}
	expmaxtip := []string{"4", "4", "4", "1", "4"}

	nodeindex, err2 := tree.NewNodeIndex(tr)
	if err2 != nil {
		t.Error(err2)
	}

	for i, name := range tipstr {
		tip, ok := nodeindex.GetNode(name)
		if !ok {
			if err != nil {
				t.Errorf("Tip %s not found in the tree", name)
			}
		}
		e, l, err2 := tree.MaxLengthPath(tip, nil)
		if err2 != nil {
			t.Error(err2)

		}
		if l != expmaxlen[i] {
			t.Errorf("Maximum length from Tip %s should be %f and is %f", name, expmaxlen[i], l)
		}

		if len(e) != expmaxpath[i] {
			t.Errorf("Nb edges of the maximum length path from Tip %s should be %d and is %d", name, expmaxpath[i], len(e))
		}

		if e[0].Right().Name() != expmaxtip[i] {
			t.Errorf("Maximum length tip from tip %s should be %s and is %s", name, expmaxtip[i], e[0].Right().Name())
		}
	}
	tr.RerootMidPoint()
}

func TestRerootMidPoint(t *testing.T) {

	tr, err := newick.NewParser(strings.NewReader(longestpathtree)).Parse()
	if err != nil {
		t.Error(err)
	}
	tr.RerootMidPoint()

	for _, e := range tr.Root().Edges() {
		if e.Right().Name() == "4" {
			if e.Length() != 8.5 {
				t.Errorf("Length of the edge from root to 4 should be 8.5")
			}
		} else if e.Length() != 1.5 {
			t.Errorf("Length of the edge from root to internal node should be 1.5")
		}
	}
}

func TestRerootMidPoint2(t *testing.T) {

	tr, err := newick.NewParser(strings.NewReader(midpointtree)).Parse()
	if err != nil {
		t.Error(err)
	}
	tr.RerootMidPoint()

	fmt.Println(tr.Newick())
	for _, e := range tr.Root().Edges() {
		if e.Length() != 3 && e.Length() != 1 {
			t.Errorf("Length at root should be 1 or 3 but is %f", e.Length())
		}
	}
}
