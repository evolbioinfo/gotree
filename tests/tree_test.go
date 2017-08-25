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
func TestClearComments(t *testing.T) {
	treeString := "(Tip4:0.1[c1],Tip0:0.1[c2],(Tip3:0.1[c3],(Tip2:0.2[c4],Tip1:0.2[c5])0.8:0.3[c6])0.9:0.4[c7])[c8];"
	tr, err := newick.NewParser(strings.NewReader(treeString)).Parse()
	if err != nil {
		t.Error(err)
	}
	tr.ClearComments()
	if tr.Newick() != "(Tip4:0.1,Tip0:0.1,(Tip3:0.1,(Tip2:0.2,Tip1:0.2)0.8:0.3)0.9:0.4);" {
		t.Error(fmt.Sprintf("Tree after clear comments is not valid: %s", tr.Newick()))
	}
}

func TestCollapseDepth(t *testing.T) {
	treeString := "(Tip4,Tip0,(Tip3,(Tip2,Tip1)));"
	tr, err := newick.NewParser(strings.NewReader(treeString)).Parse()
	tr.ReinitIndexes()
	if err != nil {
		t.Error(err)
	}
	if err = tr.CollapseTopoDepth(2, 3); err != nil {
		t.Error(err)
	}
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
	tr.ReinitIndexes()
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

// We merge two trees, and compare all bipartitions to the expected tree
func TestMerge(t *testing.T) {
	treeString := "(Tip0,(Tip3,(Tip2,Tip1)0.2)0.9);"
	treeString2 := "(Tip0_2,(Tip3_2,(Tip2_2,Tip1_2)0.2)0.9);"
	expected := "((Tip0,(Tip3,(Tip2,Tip1)0.2)0.9),(Tip0_2,(Tip3_2,(Tip2_2,Tip1_2)0.2)0.9));"
	tr, err := newick.NewParser(strings.NewReader(treeString)).Parse()
	if err != nil {
		t.Error(err)
	}
	tr2, err2 := newick.NewParser(strings.NewReader(treeString2)).Parse()
	if err2 != nil {
		t.Error(err2)
	}
	tr3, err3 := newick.NewParser(strings.NewReader(expected)).Parse()
	if err3 != nil {
		t.Error(err3)
	}
	tr.ReinitIndexes()
	tr2.ReinitIndexes()
	tr3.ReinitIndexes()

	compchan := make(chan tree.Trees)
	err4 := tr.Merge(tr2)
	if err4 != nil {
		t.Error(err4)
	}

	stats, err := tree.Compare(tr, compchan, false, true, 1)
	compchan <- tree.Trees{tr3, 0, nil}
	st := <-stats
	if st.Err != nil {
		t.Error(st.Err)
	}
	if !st.Sametree {
		t.Error(fmt.Sprintf("Merged tree %s does not correspond to the expected tree %s", tr3.Newick(), expected))
	}
}
