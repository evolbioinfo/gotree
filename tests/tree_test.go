package tests

import (
	"fmt"
	"strings"
	"testing"

	"github.com/evolbioinfo/gotree/io/newick"
	"github.com/evolbioinfo/gotree/tree"
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

func TestRoundLengths(t *testing.T) {
	treeString := "(Tip4:0.000099999999,Tip0:0.000099999999,(Tip3:0.00006666666666666,(Tip2:0.00019999999,Tip1:0.00019999999)0.8:0.000266666666666)0.9:0.0003999999999);"
	tr, err := newick.NewParser(strings.NewReader(treeString)).Parse()
	if err != nil {
		t.Error(err)
	}
	tr.RoundLengths(4)
	if tr.Newick() != "(Tip4:0.0001,Tip0:0.0001,(Tip3:0.0001,(Tip2:0.0002,Tip1:0.0002)0.8:0.0003)0.9:0.0004);" {
		t.Error(fmt.Sprintf("Tree after round supports is not valid: %s", tr.Newick()))
	}
}

func TestRoundSupports(t *testing.T) {
	treeString := "(Tip4:0.1,Tip0:0.1,(Tip3:0.1,(Tip2:0.2,Tip1:0.2)0.00007888888888:0.3)0.0008666666666:0.4);"
	tr, err := newick.NewParser(strings.NewReader(treeString)).Parse()
	if err != nil {
		t.Error(err)
	}
	tr.RoundSupports(4)
	if tr.Newick() != "(Tip4:0.1,Tip0:0.1,(Tip3:0.1,(Tip2:0.2,Tip1:0.2)0.0001:0.3)0.0009:0.4);" {
		t.Error(fmt.Sprintf("Tree after round lengths is not valid: %s", tr.Newick()))
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
	if err != nil {
		t.Error(err)
	}
	tr.ReinitIndexes()
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

	if n, _ := internal.NumTipsRight(); n != 4 {
		t.Error(fmt.Sprintf("Number of tips on the rightSide of the internal edge should be 4, but is %d", n))
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
	compchan <- tree.Trees{Tree: tr3, Id: 0, Err: nil}
	st := <-stats
	if st.Err != nil {
		t.Error(st.Err)
	}
	if !st.Sametree {
		t.Error(fmt.Sprintf("Merged tree %s does not correspond to the expected tree %s", tr3.Newick(), expected))
	}
}

// Test counting the Number of cherries
func TestNbCherries(t *testing.T) {
	treeString := "(1,2,((3,4),(5,6)));"
	expected := 3
	tr, err := newick.NewParser(strings.NewReader(treeString)).Parse()
	if err != nil {
		t.Error(err)
	}
	found := tr.NbCherries()
	if found != expected {
		t.Error(fmt.Sprintf("%d cherries are found instead of %d", found, expected))
	}
}

// Test counting the Colless Index
func TestColless(t *testing.T) {
	treeString := "(1,2,((3,4),(5,6)));"
	expected := 2
	tr, err := newick.NewParser(strings.NewReader(treeString)).Parse()
	if err != nil {
		t.Error(err)
	}
	cidx := tr.CollessIndex()
	if cidx != expected {
		t.Error(fmt.Sprintf("%d cherries are found instead of %d", cidx, expected))
	}
}

// Test counting the Colless Index
func TestColless2(t *testing.T) {
	treeString := "(((1,2),3),((4,5),6));"
	expected := 2
	tr, err := newick.NewParser(strings.NewReader(treeString)).Parse()
	if err != nil {
		t.Error(err)
	}
	cidx := tr.CollessIndex()
	if cidx != expected {
		t.Error(fmt.Sprintf("%d cherries are found instead of %d", cidx, expected))
	}
}

// Test counting the Colless Index
func TestColless3(t *testing.T) {
	treeString := "(((1,2),(3,4)),((5,6),(7,8)));"
	expected := 0
	tr, err := newick.NewParser(strings.NewReader(treeString)).Parse()
	if err != nil {
		t.Error(err)
	}
	cidx := tr.CollessIndex()
	if cidx != expected {
		t.Error(fmt.Sprintf("%d cherries are found instead of %d", cidx, expected))
	}
}

// Test counting the Colless Index
func TestColless4(t *testing.T) {
	treeString := "((1,2),(3,4),((5,6),(7,8)));"
	expected := 0
	tr, err := newick.NewParser(strings.NewReader(treeString)).Parse()
	if err != nil {
		t.Error(err)
	}
	cidx := tr.CollessIndex()
	if cidx != expected {
		t.Error(fmt.Sprintf("%d cherries are found instead of %d", cidx, expected))
	}
}

// Test counting the Sackin Index on balanced rooted and unrooted trees
func TestSackin1(t *testing.T) {
	// rooted, depth 7: 128 tips : Sackin should be 7*128=896
	tr, err := tree.RandomBalancedBinaryTree(7, true)
	if err != nil {
		t.Error(err)
	}
	expected := 7 * 128
	sidx := tr.SackinIndex()
	if sidx != expected {
		t.Error(fmt.Sprintf("Sackin index found : %d instead of %d expected", sidx, expected))
	}

	// unrooted, depth 7: 128 tips : Sackin should be 7*128=896
	tr, err = tree.RandomBalancedBinaryTree(7, false)
	if err != nil {
		t.Error(err)
	}
	expected = 7 * 128
	sidx = tr.SackinIndex()
	if sidx != expected {
		t.Error(fmt.Sprintf("Sackin index found : %d instead of %d expected", sidx, expected))
	}

	// rooted caterpillar, 128 tips : Sackin should be sum(1:127)+127=8255
	tr, err = tree.RandomCaterpillarBinaryTree(128, true)
	if err != nil {
		t.Error(err)
	}
	expected = 8255
	sidx = tr.SackinIndex()
	if sidx != expected {
		t.Error(fmt.Sprintf("Sackin index found : %d instead of %d expected", sidx, expected))
	}

	// unrooted caterpillar, 128 tips : Sackin should be (sum(2:64)+64)*2=4286
	tr, err = tree.RandomCaterpillarBinaryTree(128, false)
	if err != nil {
		t.Error(err)
	}
	expected = 4286
	sidx = tr.SackinIndex()
	if sidx != expected {
		t.Error(fmt.Sprintf("Sackin index found : %d instead of %d expected", sidx, expected))
	}

	treeString := "(((1,2),(3,4)),5);"
	expected = 13
	tr, err = newick.NewParser(strings.NewReader(treeString)).Parse()
	if err != nil {
		t.Error(err)
	}
	sidx = tr.SackinIndex()
	if sidx != expected {
		t.Error(fmt.Sprintf("Sackin index found : %d instead of %d expected", sidx, expected))
	}

	treeString = "((1,2),(3,4),5);"
	expected = 12
	tr, err = newick.NewParser(strings.NewReader(treeString)).Parse()
	if err != nil {
		t.Error(err)
	}
	sidx = tr.SackinIndex()
	if sidx != expected {
		t.Error(fmt.Sprintf("Sackin index found : %d instead of %d expected", sidx, expected))
	}
}

// Test counting the Sackin Index on balanced rooted and unrooted trees
func TestDelete(t *testing.T) {
	tr, err := tree.RandomYuleBinaryTree(701, true)
	if err != nil {
		t.Error(err)
	}
	tr.Delete()
}

// Generates a 1000 tip random tree, then rotate randomly its node neighbors
//
// Topology must be the same afterwards
func TestRotateRandom(t *testing.T) {
	tr, err := tree.RandomYuleBinaryTree(1000, true)
	clone := tr.Clone()
	if err != nil {
		t.Error(err)
	}
	edges := tr.Edges()
	index := tree.NewEdgeIndex(uint64(len(edges)*2), 0.75)
	for i, e := range edges {
		index.PutEdgeValue(e, i, e.Length())
	}

	clone.RotateInternalNodes()
	edges2 := clone.Edges()
	// Check wether the 2 trees have the same set of tip names
	if err = tr.CompareTipIndexes(clone); err != nil {
		t.Error(err)
	}

	for _, e2 := range edges2 {
		_, ok := index.Value(e2)
		if !ok {
			t.Error("An edge of the original tree is not found in the rotated tree")
		}
	}
}

// Generates a 1000 tip random tree, then rotate its node neighbors to sort
// them by number of tips
//
// Topology must be the same afterwards
func TestRotateSort(t *testing.T) {
	tr, err := tree.RandomYuleBinaryTree(1000, true)
	clone := tr.Clone()
	if err != nil {
		t.Error(err)
	}
	edges := tr.Edges()
	index := tree.NewEdgeIndex(uint64(len(edges)*2), 0.75)
	for i, e := range edges {
		index.PutEdgeValue(e, i, e.Length())
	}

	clone.SortNeighborsByTips()
	edges2 := clone.Edges()
	// Check wether the 2 trees have the same set of tip names
	if err = tr.CompareTipIndexes(clone); err != nil {
		t.Error(err)
	}

	for _, e2 := range edges2 {
		_, ok := index.Value(e2)
		if !ok {
			t.Error("An edge of the original tree is not found in the sorted tree")
		}
	}
}

// Generates a 1000 tip random tree, then rotate its node neighbors to sort
// them by number of tips
//
// Topology must be the same afterwards
func TestGenerateAllTopologie(t *testing.T) {
	rooted := []int{1, 1, 3, 15, 105, 945, 10395}
	unrooted := []int{1, 1, 1, 3, 15, 105, 945}
	var topo []*tree.Tree
	var err error
	for i := 3; i <= 7; i++ {
		//rooted topologies
		if topo, err = tree.AllTopologies(i, true); err != nil {
			t.Error(err)
		}
		if len(topo) != rooted[i-1] {
			t.Error(fmt.Sprintf("Wrong number of rooted topologies for %d tips, is %d, must be %d", i, len(topo), rooted[i-1]))
		}
		//unrooted topologies
		if topo, err = tree.AllTopologies(i, false); err != nil {
			t.Error(err)
		}
		if len(topo) != unrooted[i-1] {
			t.Error(fmt.Sprintf("Wrong number of unrooted topologies for %d tips, is %d, must be %d", i, len(topo), rooted[i-1]))
		}
	}
}

func TestPostOrder(t *testing.T) {
	var tr *tree.Tree
	var err error
	var treeString string = "(Tip4:0.1,Tip0:0.1,(Tip3:0.1,(Tip2:0.2,Tip1:0.2)0.8:0.3)0.9:0.4);"
	var expected_tiporder []string = []string{"0.1", "Tip4", "0.1", "Tip0", "0.1", "Tip3", "0.2", "Tip2", "0.2", "Tip1", "0.3", "f", "0.4", "f", "f"}
	var result_tiporder []string = make([]string, 0, len(expected_tiporder))

	if tr, err = newick.NewParser(strings.NewReader(treeString)).Parse(); err != nil {
		t.Error(err)
	}

	tr.PostOrder(func(cur *tree.Node, prev *tree.Node, e *tree.Edge) bool {
		if e != nil {
			result_tiporder = append(result_tiporder, fmt.Sprintf("%v", e.Length()))
		}
		if cur.Tip() {
			result_tiporder = append(result_tiporder, cur.Name())
		} else {
			result_tiporder = append(result_tiporder, "f")
		}
		return true
	})

	if len(result_tiporder) != len(expected_tiporder) {
		t.Errorf("Resulting postorder is not equal to Expected postorder %v vs. %v", result_tiporder, expected_tiporder)
	}

	for i, n := range result_tiporder {
		if n != expected_tiporder[i] {
			t.Errorf("Resulting postorder is not equal to Expected postorder %v vs. %v", result_tiporder, expected_tiporder)
		}
	}
}

func TestPostOrder2(t *testing.T) {
	var tr *tree.Tree
	var err error
	var treeString string = "(Tip4:0.1,Tip0:0.1,(Tip3:0.1,(Tip2:0.2,Tip1:0.2)0.8:0.3)0.9:0.4);"
	var expected_tiporder []string = []string{"0.1", "Tip4", "0.1", "Tip0", "0.1", "Tip3"}
	var result_tiporder []string = make([]string, 0, len(expected_tiporder))

	if tr, err = newick.NewParser(strings.NewReader(treeString)).Parse(); err != nil {
		t.Error(err)
	}

	tr.PostOrder(func(cur *tree.Node, prev *tree.Node, e *tree.Edge) bool {
		if e != nil {
			result_tiporder = append(result_tiporder, fmt.Sprintf("%v", e.Length()))
		}
		if cur.Tip() {
			result_tiporder = append(result_tiporder, cur.Name())
		} else {
			result_tiporder = append(result_tiporder, "f")
		}
		if cur.Name() == "Tip3" {
			return false
		}
		return true
	})

	if len(result_tiporder) != len(expected_tiporder) {
		t.Errorf("Resulting postorder is not equal to Expected postorder %v vs. %v", result_tiporder, expected_tiporder)
	}

	for i, n := range result_tiporder {
		if n != expected_tiporder[i] {
			t.Errorf("Resulting postorder is not equal to Expected postorder %v vs. %v", result_tiporder, expected_tiporder)
		}
	}
}

func TestPreOrder(t *testing.T) {
	var tr *tree.Tree
	var err error
	var treeString string = "(Tip4:0.1,Tip0:0.1,(Tip3:0.1,(Tip2:0.2,Tip1:0.2)0.8:0.3)0.9:0.4);"
	var expected_tiporder []string = []string{"f", "0.1", "Tip4", "0.1", "Tip0", "0.4", "f", "0.1", "Tip3", "0.3", "f", "0.2", "Tip2", "0.2", "Tip1"}
	var result_tiporder []string = make([]string, 0, len(expected_tiporder))

	if tr, err = newick.NewParser(strings.NewReader(treeString)).Parse(); err != nil {
		t.Error(err)
	}

	tr.PreOrder(func(cur *tree.Node, prev *tree.Node, e *tree.Edge) bool {
		if e != nil {
			result_tiporder = append(result_tiporder, fmt.Sprintf("%v", e.Length()))
		}
		if cur.Tip() {
			result_tiporder = append(result_tiporder, cur.Name())
		} else {
			result_tiporder = append(result_tiporder, "f")
		}
		return true
	})

	if len(result_tiporder) != len(expected_tiporder) {
		t.Errorf("Resulting preorder is not equal to Expected preorder %v vs. %v", result_tiporder, expected_tiporder)
	}

	for i, n := range result_tiporder {
		if n != expected_tiporder[i] {
			t.Errorf("Resulting preorder is not equal to Expected preorder %v vs. %v", result_tiporder, expected_tiporder)
		}
	}
}

func TestPreOrder2(t *testing.T) {
	var tr *tree.Tree
	var err error
	var treeString string = "(Tip4:0.1,Tip0:0.1,(Tip3:0.1,(Tip2:0.2,Tip1:0.2)0.8:0.3)0.9:0.4);"
	var expected_tiporder []string = []string{"f", "0.1", "Tip4", "0.1", "Tip0", "0.4", "f", "0.1", "Tip3"}
	var result_tiporder []string = make([]string, 0, len(expected_tiporder))

	if tr, err = newick.NewParser(strings.NewReader(treeString)).Parse(); err != nil {
		t.Error(err)
	}

	tr.PreOrder(func(cur *tree.Node, prev *tree.Node, e *tree.Edge) bool {
		if e != nil {
			result_tiporder = append(result_tiporder, fmt.Sprintf("%v", e.Length()))
		}
		if cur.Tip() {
			result_tiporder = append(result_tiporder, cur.Name())
		} else {
			result_tiporder = append(result_tiporder, "f")
		}
		if cur.Name() == "Tip3" {
			return false
		}
		return true
	})

	if len(result_tiporder) != len(expected_tiporder) {
		t.Errorf("Resulting preorder is not equal to Expected preorder %v vs. %v", result_tiporder, expected_tiporder)
	}

	for i, n := range result_tiporder {
		if n != expected_tiporder[i] {
			t.Errorf("Resulting preorder is not equal to Expected preorder %v vs. %v", result_tiporder, expected_tiporder)
		}
	}
}
