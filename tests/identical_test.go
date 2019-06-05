package tests

import (
	"strings"
	"testing"

	"github.com/evolbioinfo/gotree/io/newick"
	"github.com/evolbioinfo/gotree/tree"
)

func TestAddIdenticalTip(t *testing.T) {
	var tr *tree.Tree
	var err error
	var treeString string = "(Tip4:0.1,Tip0:0.1,(Tip3:0.1,(Tip2:0.2,Tip1:0.2)0.8:0.3)0.9:0.4);"

	var expectednewick = "((Tip8:0,Tip4:0):0.1,Tip0:0.1,(Tip3:0.1,((Tip5:0,Tip2:0,Tip6:0,Tip7:0):0.2,Tip1:0.2)0.8:0.3)0.9:0.4);"
	if tr, err = newick.NewParser(strings.NewReader(treeString)).Parse(); err != nil {
		t.Error(err)
	}

	groups := [][]string{[]string{"Tip2", "Tip5", "Tip6", "Tip7"}, []string{"Tip4", "Tip8"}}
	tr.UpdateTipIndex()

	if err = tr.InsertIdenticalTips(groups); err != nil {
		t.Error(err)
	}
	if tr.Newick() != expectednewick {
		t.Errorf("Tree with identical tips is not as expected: \n%s vs.\n%s\n", tr.Newick(), expectednewick)
	}

	tr.UpdateTipIndex()
	if err = tr.ClearBitSets(); err != nil {
		t.Error(err)
	}
	if err = tr.UpdateBitSet(); err != nil {
		t.Error(err)
	}
	if !tr.CheckTree() {
		t.Error("The tree is maformed after adding identical tips")
	}
}
