package acr

import (
	"fmt"
	"strings"
	"testing"

	"github.com/evolbioinfo/gotree/io/newick"
	"github.com/evolbioinfo/gotree/tree"
)

func TestACRDELTRANParsimony(t *testing.T) {
	// +--- t1                            A
	// |
	// t21 +--- t2                        A
	// |   |
	// |   |       +---- t3               B
	// +---|t20+---|t7
	//     |   |   |    +--- t4           B
	//     |   |   +----|t6
	//     |   |        +--- t5           A
	//     +---|t19
	//         |   +---- t8               B
	//         |   |
	//         +---|t18     +--- t9       B
	//             |    +---|t11
	//             |    |   +--- t10      A
	//             +----|t17
	//                  |       +--- t12  A
	//                  |   +---|t14
	//                  +---|t16+--- t13  A
	//                      |
	//                      +--- t15      A
	treeString := "(t1,(t2,((t3,(t4,t5)t6)t7,(t8,((t9,t10)t11,((t12,t13)t14,t15)t16)t17)t18)t19)t20)t21;"
	tipstates := map[string]string{
		"t1": "A", "t2": "A", "t3": "B", "t4": "B", "t5": "A", "t8": "B",
		"t9": "B", "t10": "A", "t12": "A", "t13": "A", "t15": "A",
	}
	states := make([]AncestralState, 21)
	upstates := make([]AncestralState, 21)

	alphabet := []string{"A", "B"}
	stateIndices := AncestralStateIndices(alphabet)
	tr, err := newick.NewParser(strings.NewReader(treeString)).Parse()
	if err != nil {
		t.Error(err)
	}
	ni := tree.NewAllNodeIndex(tr)

	for i, n := range tr.Nodes() {
		n.SetId(i)
		states[i] = make(AncestralState, 2)
		upstates[i] = make(AncestralState, 2)
	}

	// Test UPPASS
	err = parsimonyUPPASS(tr.Root(), nil, tipstates, states, stateIndices)
	if err != nil {
		t.Error(err)
	}
	testCheckStates(t, 2, ni, "t1", states, stateIndices, "A")
	testCheckStates(t, 2, ni, "t2", states, stateIndices, "A")
	testCheckStates(t, 2, ni, "t3", states, stateIndices, "B")
	testCheckStates(t, 2, ni, "t4", states, stateIndices, "B")
	testCheckStates(t, 2, ni, "t5", states, stateIndices, "A")
	testCheckStates(t, 2, ni, "t6", states, stateIndices, "B", "A")
	testCheckStates(t, 2, ni, "t7", states, stateIndices, "B")
	testCheckStates(t, 2, ni, "t8", states, stateIndices, "B")
	testCheckStates(t, 2, ni, "t9", states, stateIndices, "B")
	testCheckStates(t, 2, ni, "t10", states, stateIndices, "A")
	testCheckStates(t, 2, ni, "t11", states, stateIndices, "A", "B")
	testCheckStates(t, 2, ni, "t12", states, stateIndices, "A")
	testCheckStates(t, 2, ni, "t13", states, stateIndices, "A")
	testCheckStates(t, 2, ni, "t14", states, stateIndices, "A")
	testCheckStates(t, 2, ni, "t15", states, stateIndices, "A")
	testCheckStates(t, 2, ni, "t16", states, stateIndices, "A")
	testCheckStates(t, 2, ni, "t17", states, stateIndices, "A")
	testCheckStates(t, 2, ni, "t18", states, stateIndices, "A", "B")
	testCheckStates(t, 2, ni, "t19", states, stateIndices, "B")
	testCheckStates(t, 2, ni, "t20", states, stateIndices, "A", "B")
	testCheckStates(t, 2, ni, "t21", states, stateIndices, "A")

	// Test DOWNPASS
	parsimonyDOWNPASS(tr.Root(), nil, states, upstates, stateIndices, false)
	testCheckStates(t, 2, ni, "t1", states, stateIndices, "A")
	testCheckStates(t, 2, ni, "t2", states, stateIndices, "A")
	testCheckStates(t, 2, ni, "t3", states, stateIndices, "B")
	testCheckStates(t, 2, ni, "t4", states, stateIndices, "B")
	testCheckStates(t, 2, ni, "t5", states, stateIndices, "A")
	testCheckStates(t, 2, ni, "t6", states, stateIndices, "A", "B")
	testCheckStates(t, 2, ni, "t7", states, stateIndices, "A", "B")
	testCheckStates(t, 2, ni, "t8", states, stateIndices, "B")
	testCheckStates(t, 2, ni, "t9", states, stateIndices, "B")
	testCheckStates(t, 2, ni, "t10", states, stateIndices, "A")
	testCheckStates(t, 2, ni, "t11", states, stateIndices, "A", "B")
	testCheckStates(t, 2, ni, "t12", states, stateIndices, "A")
	testCheckStates(t, 2, ni, "t13", states, stateIndices, "A")
	testCheckStates(t, 2, ni, "t14", states, stateIndices, "A")
	testCheckStates(t, 2, ni, "t15", states, stateIndices, "A")
	testCheckStates(t, 2, ni, "t16", states, stateIndices, "A")
	testCheckStates(t, 2, ni, "t17", states, stateIndices, "A", "B")
	testCheckStates(t, 2, ni, "t18", states, stateIndices, "A", "B")
	testCheckStates(t, 2, ni, "t19", states, stateIndices, "A", "B")
	testCheckStates(t, 2, ni, "t20", states, stateIndices, "A")
	testCheckStates(t, 2, ni, "t21", states, stateIndices, "A")

	// Test DELTRAN
	parsimonyDELTRAN(tr.Root(), nil, states, stateIndices, false)
	testCheckStates(t, 2, ni, "t1", states, stateIndices, "A")
	testCheckStates(t, 2, ni, "t2", states, stateIndices, "A")
	testCheckStates(t, 2, ni, "t3", states, stateIndices, "B")
	testCheckStates(t, 2, ni, "t4", states, stateIndices, "B")
	testCheckStates(t, 2, ni, "t5", states, stateIndices, "A")
	testCheckStates(t, 2, ni, "t6", states, stateIndices, "A")
	testCheckStates(t, 2, ni, "t7", states, stateIndices, "A")
	testCheckStates(t, 2, ni, "t8", states, stateIndices, "B")
	testCheckStates(t, 2, ni, "t9", states, stateIndices, "B")
	testCheckStates(t, 2, ni, "t10", states, stateIndices, "A")
	testCheckStates(t, 2, ni, "t11", states, stateIndices, "A")
	testCheckStates(t, 2, ni, "t12", states, stateIndices, "A")
	testCheckStates(t, 2, ni, "t13", states, stateIndices, "A")
	testCheckStates(t, 2, ni, "t14", states, stateIndices, "A")
	testCheckStates(t, 2, ni, "t15", states, stateIndices, "A")
	testCheckStates(t, 2, ni, "t16", states, stateIndices, "A")
	testCheckStates(t, 2, ni, "t17", states, stateIndices, "A")
	testCheckStates(t, 2, ni, "t18", states, stateIndices, "A")
	testCheckStates(t, 2, ni, "t19", states, stateIndices, "A")
	testCheckStates(t, 2, ni, "t20", states, stateIndices, "A")
	testCheckStates(t, 2, ni, "t21", states, stateIndices, "A")

}

func TestACRACCTRANParsimony(t *testing.T) {
	// +--- t1                            A
	// |
	// t21 +--- t2                        A
	// |   |
	// |   |       +---- t3               B
	// +---|t20+---|t7
	//     |   |   |    +--- t4           B
	//     |   |   +----|t6
	//     |   |        +--- t5           A
	//     +---|t19
	//         |   +---- t8               B
	//         |   |
	//         +---|t18     +--- t9       B
	//             |    +---|t11
	//             |    |   +--- t10      A
	//             +----|t17
	//                  |       +--- t12  A
	//                  |   +---|t14
	//                  +---|t16+--- t13  A
	//                      |
	//                      +--- t15      A
	treeString := "(t1,(t2,((t3,(t4,t5)t6)t7,(t8,((t9,t10)t11,((t12,t13)t14,t15)t16)t17)t18)t19)t20)t21;"
	tipstates := map[string]string{
		"t1": "A", "t2": "A", "t3": "B", "t4": "B", "t5": "A", "t8": "B",
		"t9": "B", "t10": "A", "t12": "A", "t13": "A", "t15": "A",
	}
	states := make([]AncestralState, 21)
	alphabet := []string{"A", "B"}
	stateIndices := AncestralStateIndices(alphabet)
	tr, err := newick.NewParser(strings.NewReader(treeString)).Parse()
	if err != nil {
		t.Error(err)
	}
	ni := tree.NewAllNodeIndex(tr)

	for i, n := range tr.Nodes() {
		n.SetId(i)
		states[i] = make(AncestralState, 2)
	}

	// Test UPPASS
	err = parsimonyUPPASS(tr.Root(), nil, tipstates, states, stateIndices)
	if err != nil {
		t.Error(err)
	}
	testCheckStates(t, 2, ni, "t1", states, stateIndices, "A")
	testCheckStates(t, 2, ni, "t2", states, stateIndices, "A")
	testCheckStates(t, 2, ni, "t3", states, stateIndices, "B")
	testCheckStates(t, 2, ni, "t4", states, stateIndices, "B")
	testCheckStates(t, 2, ni, "t5", states, stateIndices, "A")
	testCheckStates(t, 2, ni, "t6", states, stateIndices, "B", "A")
	testCheckStates(t, 2, ni, "t7", states, stateIndices, "B")
	testCheckStates(t, 2, ni, "t8", states, stateIndices, "B")
	testCheckStates(t, 2, ni, "t9", states, stateIndices, "B")
	testCheckStates(t, 2, ni, "t10", states, stateIndices, "A")
	testCheckStates(t, 2, ni, "t11", states, stateIndices, "A", "B")
	testCheckStates(t, 2, ni, "t12", states, stateIndices, "A")
	testCheckStates(t, 2, ni, "t13", states, stateIndices, "A")
	testCheckStates(t, 2, ni, "t14", states, stateIndices, "A")
	testCheckStates(t, 2, ni, "t15", states, stateIndices, "A")
	testCheckStates(t, 2, ni, "t16", states, stateIndices, "A")
	testCheckStates(t, 2, ni, "t17", states, stateIndices, "A")
	testCheckStates(t, 2, ni, "t18", states, stateIndices, "A", "B")
	testCheckStates(t, 2, ni, "t19", states, stateIndices, "B")
	testCheckStates(t, 2, ni, "t20", states, stateIndices, "A", "B")
	testCheckStates(t, 2, ni, "t21", states, stateIndices, "A")

	// Test ACCTRAN
	parsimonyACCTRAN(tr.Root(), nil, states, stateIndices, false)
	testCheckStates(t, 2, ni, "t1", states, stateIndices, "A")
	testCheckStates(t, 2, ni, "t2", states, stateIndices, "A")
	testCheckStates(t, 2, ni, "t3", states, stateIndices, "B")
	testCheckStates(t, 2, ni, "t4", states, stateIndices, "B")
	testCheckStates(t, 2, ni, "t5", states, stateIndices, "A")
	testCheckStates(t, 2, ni, "t6", states, stateIndices, "B")
	testCheckStates(t, 2, ni, "t7", states, stateIndices, "B")
	testCheckStates(t, 2, ni, "t8", states, stateIndices, "B")
	testCheckStates(t, 2, ni, "t9", states, stateIndices, "B")
	testCheckStates(t, 2, ni, "t10", states, stateIndices, "A")
	testCheckStates(t, 2, ni, "t11", states, stateIndices, "A")
	testCheckStates(t, 2, ni, "t12", states, stateIndices, "A")
	testCheckStates(t, 2, ni, "t13", states, stateIndices, "A")
	testCheckStates(t, 2, ni, "t14", states, stateIndices, "A")
	testCheckStates(t, 2, ni, "t15", states, stateIndices, "A")
	testCheckStates(t, 2, ni, "t16", states, stateIndices, "A")
	testCheckStates(t, 2, ni, "t17", states, stateIndices, "A")
	testCheckStates(t, 2, ni, "t18", states, stateIndices, "B")
	testCheckStates(t, 2, ni, "t19", states, stateIndices, "B")
	testCheckStates(t, 2, ni, "t20", states, stateIndices, "A")
	testCheckStates(t, 2, ni, "t21", states, stateIndices, "A")
}

func TestALLDELTRAN(t *testing.T) {
	// +--- t1                            A
	// |
	// t21 +--- t2                        A
	// |   |
	// |   |       +---- t3               B
	// +---|t20+---|t7
	//     |   |   |    +--- t4           B
	//     |   |   +----|t6
	//     |   |        +--- t5           A
	//     +---|t19
	//         |   +---- t8               B
	//         |   |
	//         +---|t18     +--- t9       B
	//             |    +---|t11
	//             |    |   +--- t10      A
	//             +----|t17
	//                  |       +--- t12  A
	//                  |   +---|t14
	//                  +---|t16+--- t13  A
	//                      |
	//                      +--- t15      A
	treeString := "(t1,(t2,((t3,(t4,t5)t6)t7,(t8,((t9,t10)t11,((t12,t13)t14,t15)t16)t17)t18)t19)t20)t21;"
	tipstates := map[string]string{
		"t1": "A", "t2": "A", "t3": "B", "t4": "B", "t5": "A", "t8": "B",
		"t9": "B", "t10": "A", "t12": "A", "t13": "A", "t15": "A",
	}
	tr, err := newick.NewParser(strings.NewReader(treeString)).Parse()
	if err != nil {
		t.Error(err)
	}

	statemap, err := ParsimonyAcr(tr, tipstates, ALGO_DELTRAN, false)

	if err != nil {
		t.Error(err)
	}
	testCheckMap(t, "t6", statemap, "A")
	testCheckMap(t, "t7", statemap, "A")
	testCheckMap(t, "t11", statemap, "A")
	testCheckMap(t, "t14", statemap, "A")
	testCheckMap(t, "t16", statemap, "A")
	testCheckMap(t, "t17", statemap, "A")
	testCheckMap(t, "t18", statemap, "A")
	testCheckMap(t, "t19", statemap, "A")
	testCheckMap(t, "t20", statemap, "A")
	testCheckMap(t, "t21", statemap, "A")
}
func TestALLDOWNPASS(t *testing.T) {
	// +--- t1                            A
	// |
	// t21 +--- t2                        A
	// |   |
	// |   |       +---- t3               B
	// +---|t20+---|t7
	//     |   |   |    +--- t4           B
	//     |   |   +----|t6
	//     |   |        +--- t5           A
	//     +---|t19
	//         |   +---- t8               B
	//         |   |
	//         +---|t18     +--- t9       B
	//             |    +---|t11
	//             |    |   +--- t10      A
	//             +----|t17
	//                  |       +--- t12  A
	//                  |   +---|t14
	//                  +---|t16+--- t13  A
	//                      |
	//                      +--- t15      A
	treeString := "(t1,(t2,((t3,(t4,t5)t6)t7,(t8,((t9,t10)t11,((t12,t13)t14,t15)t16)t17)t18)t19)t20)t21;"
	tipstates := map[string]string{
		"t1": "A", "t2": "A", "t3": "B", "t4": "B", "t5": "A", "t8": "B",
		"t9": "B", "t10": "A", "t12": "A", "t13": "A", "t15": "A",
	}
	tr, err := newick.NewParser(strings.NewReader(treeString)).Parse()
	if err != nil {
		t.Error(err)
	}

	statemap, err := ParsimonyAcr(tr, tipstates, ALGO_DOWNPASS, false)
	if err != nil {
		t.Error(err)
	}
	testCheckMap(t, "t6", statemap, "A", "B")
	testCheckMap(t, "t7", statemap, "A", "B")
	testCheckMap(t, "t11", statemap, "A", "B")
	testCheckMap(t, "t14", statemap, "A")
	testCheckMap(t, "t16", statemap, "A")
	testCheckMap(t, "t17", statemap, "A", "B")
	testCheckMap(t, "t18", statemap, "A", "B")
	testCheckMap(t, "t19", statemap, "A", "B")
	testCheckMap(t, "t20", statemap, "A")
	testCheckMap(t, "t21", statemap, "A")
}

func TestALLACCTRAN(t *testing.T) {
	// +--- t1                            A
	// |
	// t21 +--- t2                        A
	// |   |
	// |   |       +---- t3               B
	// +---|t20+---|t7
	//     |   |   |    +--- t4           B
	//     |   |   +----|t6
	//     |   |        +--- t5           A
	//     +---|t19
	//         |   +---- t8               B
	//         |   |
	//         +---|t18     +--- t9       B
	//             |    +---|t11
	//             |    |   +--- t10      A
	//             +----|t17
	//                  |       +--- t12  A
	//                  |   +---|t14
	//                  +---|t16+--- t13  A
	//                      |
	//                      +--- t15      A
	treeString := "(t1,(t2,((t3,(t4,t5)t6)t7,(t8,((t9,t10)t11,((t12,t13)t14,t15)t16)t17)t18)t19)t20)t21;"
	tipstates := map[string]string{
		"t1": "A", "t2": "A", "t3": "B", "t4": "B", "t5": "A", "t8": "B",
		"t9": "B", "t10": "A", "t12": "A", "t13": "A", "t15": "A",
	}
	tr, err := newick.NewParser(strings.NewReader(treeString)).Parse()
	if err != nil {
		t.Error(err)
	}

	statemap, err := ParsimonyAcr(tr, tipstates, ALGO_ACCTRAN, false)
	if err != nil {
		t.Error(err)
	}
	testCheckMap(t, "t6", statemap, "B")
	testCheckMap(t, "t7", statemap, "B")
	testCheckMap(t, "t11", statemap, "A")
	testCheckMap(t, "t14", statemap, "A")
	testCheckMap(t, "t16", statemap, "A")
	testCheckMap(t, "t17", statemap, "A")
	testCheckMap(t, "t18", statemap, "B")
	testCheckMap(t, "t19", statemap, "B")
	testCheckMap(t, "t20", statemap, "A")
	testCheckMap(t, "t21", statemap, "A")
}

func testCheckStates(t *testing.T, nstates int, ni tree.NodeIndex, nodename string,
	states []AncestralState, stateIndices map[string]int, teststates ...string) {
	n, ok := ni.GetNode(nodename)
	if !ok {
		t.Error(fmt.Errorf("Node %s does not exist in the tree", nodename))
	}
	a := make(AncestralState, 2)
	for _, s := range teststates {
		idx, ok := stateIndices[s]
		if !ok {
			t.Error(fmt.Errorf("State %s does not exist in the index", s))
		}
		a[idx] = 1
	}
	// We compare the two ancestral states now
	for i, v := range a {
		if v != states[n.Id()][i] {
			t.Error(fmt.Errorf("Node %s should have states : %v but has states %v", nodename, a, states[n.Id()]))
		}
	}
}
func testCheckMap(t *testing.T, nodename string, statemap map[string]string, states ...string) {
	st, ok := statemap[nodename]
	if !ok {
		t.Error(fmt.Errorf("Node %s does not exist in the output state map", nodename))
	}
	if strings.Join(states, ",") != st {
		t.Error(fmt.Errorf("Node %s should have states : %s but has states %s", nodename, strings.Join(states, ","), st))
	}
}
