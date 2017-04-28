package tests

import (
	"fmt"
	"strings"
	"testing"

	"github.com/fredericlemoine/gotree/io/newick"
	"github.com/fredericlemoine/gotree/tree"
)

func TestClearLengths(t *testing.T) {
	treeString := "(Tip4:0.1,Tip0:0.1,(Tip3:0.1,(Tip2:0.2,Tip1:0.2)0.8:0.3)0.9:0.4);"
	tr, err := newick.NewParser(strings.NewReader(treeString)).Parse()
	if err != nil {
		t.Error(err)
	}
	tr.ClearLengths()
	if tr.Newick() != "(Tip4,Tip0,(Tip3,(Tip2,Tip1)0.8)0.9);" {
		t.Error(fmt.Sprintf("Tree after clear supports is not valid: %s", tr.Newick()))
	}
}

func TestClearSupports(t *testing.T) {
	treeString := "(Tip4:0.1,Tip0:0.1,(Tip3:0.1,(Tip2:0.2,Tip1:0.2)0.8:0.3)0.9:0.4);"
	tr, err := newick.NewParser(strings.NewReader(treeString)).Parse()
	if err != nil {
		t.Error(err)
	}
	tr.ClearSupports()
	if tr.Newick() != "(Tip4:0.1,Tip0:0.1,(Tip3:0.1,(Tip2:0.2,Tip1:0.2):0.3):0.4);" {
		t.Error(fmt.Sprintf("Tree after clear lengths is not valid: %s", tr.Newick()))
	}
}

func TestCollapseDepth(t *testing.T) {
	treeString := "(Tip4,Tip0,(Tip3,(Tip2,Tip1)));"
	tr, err := newick.NewParser(strings.NewReader(treeString)).Parse()
	if err != nil {
		t.Error(err)
	}
	tr.CollapseTopoDepth(2, 3)
	if tr.Newick() != "(Tip4,Tip0,Tip3,Tip2,Tip1);" {
		t.Error(fmt.Sprintf("Tree after collapse depth is not valid: %s", tr.Newick()))
	}
}

func TestCollapseLength(t *testing.T) {
	treeString := "(Tip4:0.1,Tip0:0.1,(Tip3:0.1,(Tip2:0.2,Tip1:0.2):0.001):0.4);"
	tr, err := newick.NewParser(strings.NewReader(treeString)).Parse()
	if err != nil {
		t.Error(err)
	}
	tr.CollapseShortBranches(0.01)
	if tr.Newick() != "(Tip4:0.1,Tip0:0.1,(Tip3:0.1,Tip2:0.2,Tip1:0.2):0.4);" {
		t.Error(fmt.Sprintf("Tree after collapse lengths is not valid: %s", tr.Newick()))
	}
}

func TestCollapseSupport(t *testing.T) {
	treeString := "(Tip4,Tip0,(Tip3,(Tip2,Tip1)0.2)0.9);"
	tr, err := newick.NewParser(strings.NewReader(treeString)).Parse()
	if err != nil {
		t.Error(err)
	}
	tr.CollapseLowSupport(0.7)
	if tr.Newick() != "(Tip4,Tip0,(Tip3,Tip2,Tip1)0.9);" {
		t.Error(fmt.Sprintf("Tree after collapse support is not valid: %s", tr.Newick()))
	}
}

func TestBipartitionTree(t *testing.T) {
	rightTips := []string{"T1", "T2", "T3", "T4"}
	leftTips := []string{"T5", "T6", "T7"}

	tr, err := tree.BipartitionTree(leftTips, rightTips)
	if err != nil {
		t.Error(err)
	}

	if len(tr.Tips()) != 7 {
		t.Error(fmt.Sprintf("Tree should have 7 tips but have %d", len(tr.Tips())))
	}

	if len(tr.Edges()) != 8 {
		t.Error(fmt.Sprintf("Tree should have 8 Edges but have %d", len(tr.Edges())))
	}
	nbInternal := 0
	nbExternal := 0
	var internal *tree.Edge
	for _, e := range tr.Edges() {
		if e.Right().Tip() {
			nbExternal++
		} else {
			nbInternal++
			internal = e
		}
	}

	if nbExternal != 7 {
		t.Error(fmt.Sprintf("Tree should have 7 external Edges but have %d", nbExternal))
	}

	if nbInternal != 1 {
		t.Error(fmt.Sprintf("Tree should have 1 internal Edge but have %d", nbInternal))
	}

	if internal.NumTips() != 4 {
		t.Error(fmt.Sprintf("Number of tips on the rightSide of the internal edge should be 4, but is %d", internal.NumTips()))
	}
}
