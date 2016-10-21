package tests

import (
	"fmt"
	"github.com/fredericlemoine/gotree/io/newick"
	"os"
	"strings"
	"testing"
)

var unroottree string = "((1:1,2:1):1,5:1,(3:1,4:1):1);"
var startree string = "(1:1,2:1,3:1,4:1,5:1);"

func TestLeastCommonAncestorUnrooted(t *testing.T) {
	tree, err := newick.NewParser(strings.NewReader(unroottree)).Parse()
	if err != nil {
		t.Error(err)
	}
	n, e, mono := tree.LeastCommonAncestorUnrooted("3", "4")
	if n == nil {
		t.Error("No common ancestor found")
	}
	if e == nil || len(e) == 0 {
		t.Error("No common ancestor Edges found")
	}
	if !mono {
		t.Error("The group should be monophyletic")
	}
	if len(e) != 2 {
		t.Error("Edge Length should be 2 and is %d", len(e))
	}

	for _, ed := range e {
		if ed.Right().Name() != "4" &&
			ed.Right().Name() != "3" &&
			ed.Left().Name() != "4" &&
			ed.Left().Name() != "3" {
			t.Error("Ancestor of wrong tips found")
		}
	}

	n, e, mono = tree.LeastCommonAncestorUnrooted("3", "5")
	if n == nil {
		t.Error("No common ancestor found")
	}
	if e == nil || len(e) == 0 {
		t.Error("No common ancestor Edges found")
	}
	if mono {
		t.Error("The group should not be monophyletic")
	}

	if len(e) != 2 {
		t.Error("Edge Length should be 2 and is %d", len(e))
	}

	star, err2 := newick.NewParser(strings.NewReader(startree)).Parse()
	if err2 != nil {
		t.Error(err2)
	}
	n, e, mono = star.LeastCommonAncestorUnrooted("3", "4", "5", "2")
	if n == nil {
		t.Error("No common ancestor found")
	}
	if e == nil || len(e) == 0 {
		t.Error("No common ancestor Edges found")
	}

	if len(e) != 4 {
		t.Error("Edge Length should be 4 and is %d", len(e))
	}
	if !mono {
		t.Error("The group should be monophyletic")
	}

}

func TestAddBipartition(t *testing.T) {
	star, err2 := newick.NewParser(strings.NewReader(startree)).Parse()
	if err2 != nil {
		t.Error(err2)
	}
	n, e, mono := star.LeastCommonAncestorUnrooted("3", "4", "5")
	if n == nil {
		t.Error("No common ancestor found")
	}
	if e == nil || len(e) == 0 {
		t.Error("No common ancestor Edges found")
	}
	if !mono {
		t.Error("The group should be monophyletic")
	}

	if len(e) != 3 {
		t.Error("Edge Length should be 3 and is %d", len(e))
	}
	star.AddBipartition(n, e, 0.001, 0.9)
	n, e, mono = star.LeastCommonAncestorUnrooted("4", "5")
	if n == nil {
		t.Error("No common ancestor found")
	}
	if e == nil || len(e) == 0 {
		t.Error("No common ancestor Edges found")
	}
	if !mono {
		t.Error("The group should be monophyletic")
	}

	if len(e) != 2 {
		t.Error("Edge Length should be 2 and is %d", len(e))
	}
	fmt.Fprintf(os.Stdout, "%s\n", star.Newick())
	star.AddBipartition(n, e, 0.001, 0.9)
	fmt.Fprintf(os.Stdout, "%s\n", star.Newick())
	n, e, mono = star.LeastCommonAncestorUnrooted("1", "2")
	if n == nil {
		t.Error("No common ancestor found")
	}
	if e == nil || len(e) == 0 {
		t.Error("No common ancestor Edges found")
	}
	if !mono {
		t.Error("The group should be monophyletic")
	}

	if len(e) != 2 {
		t.Error("Edge Length should be 2 and is %d", len(e))
	}
	fmt.Fprintf(os.Stdout, "%s\n", star.Newick())
}
