package tree

import (
	"errors"
	"github.com/fredericlemoine/gostats"
	"github.com/fredericlemoine/gotree/io"
	"math/rand"
	"strconv"
)

//Creates a Random uniform Binary tree by successively adding
//new tips to a random edge of the tree.
//nbtips : Number of tips of the random binary tree to create
//rooted: if true, generates a rooted tree
//branch length: follow an exponential distribution with param lambda=1/0.1
func RandomUniformBinaryTree(nbtips int, rooted bool) (*Tree, error) {
	t := NewTree()
	if nbtips < 2 {
		return nil, errors.New("Cannot create an unrooted random binary tree with less than 2 tips")
	}
	if nbtips < 3 && rooted {
		return nil, errors.New("Cannot create a rooted random binary tree with less than 3 tips")
	}

	edges := make([]*Edge, 0, 2000)
	for i := 1; i < nbtips; i++ {
		n := t.NewNode()
		n.SetName("Tip" + strconv.Itoa(i))
		switch len(edges) {
		case 0:
			n2 := t.NewNode()
			n2.SetName("Tip" + strconv.Itoa(i-1))
			e := t.ConnectNodes(n2, n)
			edges = append(edges, e)
			e.SetLength(gostats.Exp(1 / 0.1))
			if rooted {
				n2.SetName("")
				n3 := t.NewNode()
				n3.SetName("Tip" + strconv.Itoa(i-1))
				e2 := t.ConnectNodes(n2, n3)
				edges = append(edges, e2)
				e2.SetLength(gostats.Exp(1 / 0.1))
			}
			t.SetRoot(n2)
		default:
			// Where to insert the new tip
			i_edge := rand.Intn(len(edges))
			e := edges[i_edge]
			e.SetLength(gostats.Exp(1 / 0.1))
			newedge, newedge2, _, err := t.GraftTipOnEdge(n, e)
			newedge.SetLength(gostats.Exp(1 / 0.1))
			newedge2.SetLength(gostats.Exp(1 / 0.1))

			edges = append(edges, newedge)
			edges = append(edges, newedge2)

			if err != nil {
				return nil, err
			}
		}
	}
	var err error = nil
	if !rooted {
		err = t.RerootFirst()
	}
	t.UpdateTipIndex()
	t.ClearBitSets()
	t.UpdateBitSet()
	t.ComputeDepths()
	return t, err
}

//Creates a Random Balanced Binary tree. Does it recursively
// until the given depth is attained. Root is at depth 0.
// So a depth "d" will generate a tree with 2^(d) tips.
//depth : Depth of the balanced binary tree
//rooted: if true, generates a rooted tree
//branch length: follow an exponential distribution with param lambda=1/0.1
func RandomBalancedBinaryTree(depth int, rooted bool) (*Tree, error) {
	t := NewTree()
	if depth < 1 {
		return nil, errors.New("Cannot create an random binary tree of depth < 1")
	}

	curdepth := 0
	root := t.NewNode()
	t.SetRoot(root)
	id := 0
	randomBalancedBinaryTreeRecur(t, root, curdepth+1, depth, &id)
	if !rooted {
		t.UnRoot()
	}
	t.UpdateTipIndex()
	t.ClearBitSets()
	t.UpdateBitSet()
	t.ComputeDepths()
	return t, nil
}

// Recursive function called by RandomBalancedBinaryTree function
func randomBalancedBinaryTreeRecur(t *Tree, node *Node, curdepth int, targetdepth int, id *int) {
	child1 := t.NewNode()
	child2 := t.NewNode()

	e1 := t.ConnectNodes(node, child1)
	e2 := t.ConnectNodes(node, child2)
	e1.SetLength(gostats.Exp(1 / 0.1))
	e2.SetLength(gostats.Exp(1 / 0.1))

	if curdepth < targetdepth {
		randomBalancedBinaryTreeRecur(t, child1, curdepth+1, targetdepth, id)
		randomBalancedBinaryTreeRecur(t, child2, curdepth+1, targetdepth, id)
	} else {
		child1.SetName("Tip" + strconv.Itoa(*id))
		(*id)++
		child2.SetName("Tip" + strconv.Itoa(*id))
		(*id)++
	}
}

//Creates a Random Yule tree by successively adding new tips
//to random terminal edges of the tree.
//nbtips : Number of tips of the random binary tree to create
//rooted: if true, generates a rooted tree
//branch lengths: follow an exponential distribution with param lambda=1/0.1
func RandomYuleBinaryTree(nbtips int, rooted bool) (*Tree, error) {
	t := NewTree()
	if nbtips < 2 {
		return nil, errors.New("Cannot create an unrooted random binary tree with less than 2 tips")
	}
	if nbtips < 3 && rooted {
		return nil, errors.New("Cannot create a rooted random binary tree with less than 3 tips")
	}

	edges := make([]*Edge, 0, 2000)
	tips := make([]*Node, 0, 2000)
	for i := 1; i < nbtips; i++ {
		n := t.NewNode()
		n.SetName("Tip" + strconv.Itoa(i))
		switch len(edges) {
		case 0:
			n2 := t.NewNode()
			n2.SetName("Tip" + strconv.Itoa(i-1))
			e := t.ConnectNodes(n2, n)
			edges = append(edges, e)
			e.SetLength(gostats.Exp(1 / 0.1))
			if !rooted {
				tips = append(tips, n2)
			} else {
				n2.SetName("")
				n3 := t.NewNode()
				n3.SetName("Tip" + strconv.Itoa(i-1))
				e2 := t.ConnectNodes(n2, n3)
				edges = append(edges, e2)
				e2.SetLength(gostats.Exp(1 / 0.1))
				tips = append(tips, n3)
			}
			t.SetRoot(n2)
		default:
			// Where to insert the new tip
			i_tip := rand.Intn(len(tips))
			ntemp := tips[i_tip]
			e := ntemp.br[0]
			e.SetLength(gostats.Exp(1 / 0.1))
			newedge, newedge2, _, err := t.GraftTipOnEdge(n, e)
			newedge.SetLength(gostats.Exp(1 / 0.1))
			newedge2.SetLength(gostats.Exp(1 / 0.1))
			edges = append(edges, newedge)
			edges = append(edges, newedge2)
			if err != nil {
				return nil, err
			}
		}
		tips = append(tips, n)
	}
	var err error = nil

	if !rooted {
		err = t.RerootFirst()
	}
	t.UpdateTipIndex()
	t.ClearBitSets()
	t.UpdateBitSet()
	t.ComputeDepths()
	return t, err
}

//Creates a Random Caterpilar tree by adding new tips to the last
//added terminal edge of the tree.
//nbtips : Number of tips of the random binary tree to create
//rooted: if true, generates a rooted tree
//branch length: follow an exponential distribution with param lambda=1/0.1
func RandomCaterpilarBinaryTree(nbtips int, rooted bool) (*Tree, error) {
	t := NewTree()
	if nbtips < 2 {
		return nil, errors.New("Cannot create an unrooted random binary tree with less than 2 tips")
	}
	if nbtips < 3 && rooted {
		return nil, errors.New("Cannot create a rooted random binary tree with less than 3 tips")
	}

	var lasttip *Node = nil

	for i := 1; i < nbtips; i++ {
		n := t.NewNode()
		n.SetName("Tip" + strconv.Itoa(i))
		switch i {
		case 1:
			n2 := t.NewNode()
			n2.SetName("Tip" + strconv.Itoa(i-1))
			e := t.ConnectNodes(n2, n)
			e.SetLength(gostats.Exp(1 / 0.1))
			if !rooted {
				lasttip = n2
			} else {
				n2.SetName("")
				n3 := t.NewNode()
				n3.SetName("Tip" + strconv.Itoa(i-1))
				e2 := t.ConnectNodes(n2, n3)
				e2.SetLength(gostats.Exp(1 / 0.1))
				lasttip = n3
			}
			t.SetRoot(n2)
		default:
			e := lasttip.br[0]
			e.SetLength(gostats.Exp(1 / 0.1))
			newedge, newedge2, _, err := t.GraftTipOnEdge(n, e)
			newedge.SetLength(gostats.Exp(1 / 0.1))
			newedge2.SetLength(gostats.Exp(1 / 0.1))
			if err != nil {
				return nil, err
			}
		}
		lasttip = n
	}
	var err error = nil

	if !rooted {
		err = t.RerootFirst()
	}
	t.UpdateTipIndex()
	t.ClearBitSets()
	t.UpdateBitSet()
	t.ComputeDepths()
	return t, err
}

// Creates a Star tree with nbtips tips.
// Since there is only one possible labelled tree, no need
// of randomness.
// nbtips : Number of tips of the star tree.
// Branch lengths are all set to 1.0
func StarTree(nbtips int) (*Tree, error) {
	t := NewTree()
	if nbtips < 2 {
		return nil, errors.New("Cannot create a star tree with less than 2 tips")
	}

	// Central node of the star (root)
	n := t.NewNode()
	t.SetRoot(n)
	//n.SetName("N" + strconv.Itoa(i))
	for i := 0; i < nbtips; i++ {
		n2 := t.NewNode()
		n2.SetName("Tip" + strconv.Itoa(i))
		e := t.ConnectNodes(n, n2)
		e.SetLength(1.0)
	}
	//err := t.RerootFirst()
	t.UpdateTipIndex()
	t.ClearBitSets()
	t.UpdateBitSet()
	t.ComputeDepths()
	return t, nil
}

// Creates a star tree using tipnames in argument
// Since there is only one possible labelled tree, no need
// of randomness.
// Branch lengths are all set to 1.0
func StarTreeFromName(names ...string) (*Tree, error) {
	if t, err := StarTree(len(names)); err != nil {
		return nil, err
	} else {
		for i, t := range t.Tips() {
			t.SetName(names[i])
		}
		t.UpdateTipIndex()
		t.ClearBitSets()
		t.UpdateBitSet()
		t.ComputeDepths()
		return t, nil
	}
}

// Creates a Star tree with the same tips as the tree in argument
// Lengths of branches in the star trees are the same as lengths of
// terminal edges of the input tree
func StarTreeFromTree(t *Tree) (*Tree, error) {
	edges := t.TipEdges()
	if star, err := StarTree(len(edges)); err != nil {
		return nil, err
	} else {
		for i, te := range star.TipEdges() {
			te.Right().SetName(edges[i].Right().Name())
			te.SetLength(edges[i].Length())
		}
		star.UpdateTipIndex()
		star.ClearBitSets()
		star.UpdateBitSet()
		star.ComputeDepths()
		return star, nil
	}
}

/**
Builds a tree whose only internal edge is the given edge e
The two internal nodes are multifurcated
\     /
-*---*-
/     \
alltips is the slice containing all tip names of the tree
if nil, it will be recomputed
*/
func EdgeTree(t *Tree, e *Edge, alltips []string) *Tree {
	edgeTree := NewTree()
	n := edgeTree.NewNode()
	n2 := edgeTree.NewNode()
	et := edgeTree.ConnectNodes(n2, n)
	et.SetLength(1.0)
	edgeTree.SetRoot(n2)
	if alltips == nil {
		alltips = t.AllTipNames()
	}
	// We add tips on the left or on the right of the first edge
	for _, name := range alltips {
		if idx, err := t.TipIndex(name); err != nil {
			io.ExitWithMessage(err)
		} else {
			ntmp := edgeTree.NewNode()
			ntmp.SetName(name)
			// Right
			if e.Bitset().Test(idx) {
				etmp := edgeTree.ConnectNodes(n, ntmp)
				etmp.SetLength(1.0)
			} else {
				// Left
				etmp := edgeTree.ConnectNodes(n2, ntmp)
				etmp.SetLength(1.0)
			}
		}
	}
	return edgeTree
}

/**
Builds a single edge tree, given left taxa and right taxa
\     /
-*---*-
/     \
*/
func BipartitionTree(leftTips []string, rightTips []string) (*Tree, error) {

	if len(leftTips) <= 1 || len(rightTips) <= 1 {
		return nil, errors.New("Left and Right tip sets must have length > 1")
	}

	tipMap := make(map[string]bool)
	for _, t := range leftTips {
		tipMap[t] = true
	}
	for _, t := range rightTips {
		if _, ok := tipMap[t]; ok {
			return nil, errors.New("One or more tips are common between left set and right set")
		}
	}

	bipartitionTree := NewTree()
	n := bipartitionTree.NewNode()
	n2 := bipartitionTree.NewNode()
	et := bipartitionTree.ConnectNodes(n2, n)
	et.SetLength(1.0)
	bipartitionTree.SetRoot(n2)
	// We add left tips on the left of the first edge
	for _, name := range leftTips {
		ntmp := bipartitionTree.NewNode()
		ntmp.SetName(name)
		// Left
		etmp := bipartitionTree.ConnectNodes(n2, ntmp)
		etmp.SetLength(1.0)
	}
	// We add right tips on the right of the first edge
	for _, name := range rightTips {
		ntmp := bipartitionTree.NewNode()
		ntmp.SetName(name)
		// Right
		etmp := bipartitionTree.ConnectNodes(n, ntmp)
		etmp.SetLength(1.0)
	}
	bipartitionTree.UpdateTipIndex()
	bipartitionTree.ClearBitSets()
	bipartitionTree.UpdateBitSet()
	bipartitionTree.ComputeDepths()

	return bipartitionTree, nil
}
