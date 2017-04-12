package tests

import (
	"fmt"
	"strings"
	"testing"

	"github.com/fredericlemoine/gotree/io/newick"
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
