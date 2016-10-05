/*
   Package gotree implements a simple
   library for handling phylogenetic trees in go
*/
package tree

import (
	"bytes"
	"errors"
	"github.com/fredericlemoine/bitset"
	"github.com/fredericlemoine/gotree/io"
	"math"
	"math/rand"
	"sort"
	"strconv"
)

type Tree struct {
	root     *Node           // root node: If the tree is unrooted the root node should have 3 children
	tipIndex map[string]uint // Map between tip name and bitset index
}

func NewTree() *Tree {
	return &Tree{
		root:     nil,
		tipIndex: make(map[string]uint, 0),
	}
}

func (t *Tree) NewNode() *Node {
	return &Node{
		name:    "",
		comment: make([]string, 0),
		neigh:   make([]*Node, 0, 3),
		br:      make([]*Edge, 0, 3),
		depth:   -1,
	}
}

// put at nil the node and all its branches
func (t *Tree) delNode(n *Node) {
	for i, _ := range n.neigh {
		n.neigh[i] = nil
	}
	n.neigh = nil

	for i, e := range n.br {
		e.left = nil
		e.right = nil
		n.br[i] = nil
	}
	n.br = nil
}

func (t *Tree) NewEdge() *Edge {
	return &Edge{
		length:  -1.0,
		support: -1.0,
		id:      -1,
	}
}

/* Tree functions */
/******************/

func (t *Tree) SetRoot(r *Node) {
	t.root = r
}

func (t *Tree) Root() *Node {
	return t.root
}

func (t *Tree) Rooted() bool {
	return t.root.Nneigh() == 2
}

// Returns all the edges of the tree (do it recursively)
func (t *Tree) Edges() []*Edge {
	edges := make([]*Edge, 0, 2000)
	for _, e := range t.Root().br {
		edges = append(edges, e)
		t.edgesRecur(e, &edges)
	}
	return edges
}

func (t *Tree) edgesRecur(edge *Edge, edges *[]*Edge) {
	if len(edge.right.neigh) > 1 {
		for _, child := range edge.right.br {
			if child.left == edge.right {
				*edges = append((*edges), child)
				t.edgesRecur(child, edges)
			}
		}
	}
}

// Returns all the tip edges of the tree (do it recursively)
func (t *Tree) TipEdges() []*Edge {
	edges := make([]*Edge, 0, 2000)
	for _, e := range t.Root().br {
		if e.Right().Tip() {
			edges = append(edges, e)
		}
		t.tipEdgesRecur(e, &edges)
	}
	return edges
}

func (t *Tree) tipEdgesRecur(edge *Edge, edges *[]*Edge) {
	if len(edge.right.neigh) > 1 {
		for _, child := range edge.right.br {
			if child.left == edge.right {
				if child.Right().Tip() {
					*edges = append((*edges), child)
				}
				t.tipEdgesRecur(child, edges)
			}
		}
	}
}

// Returns all the nodes of the tree (do it recursively)
func (t *Tree) Nodes() []*Node {
	nodes := make([]*Node, 0, 2000)
	t.nodesRecur(&nodes, nil, nil)
	return nodes
}

func (t *Tree) nodesRecur(nodes *[]*Node, cur *Node, prev *Node) {
	if cur == nil {
		cur = t.Root()
	}
	*nodes = append((*nodes), cur)
	for _, n := range cur.neigh {
		if n != prev {
			t.nodesRecur(nodes, n, cur)
		}
	}
}

// Returns all the tips of the tree (do it recursively)
func (t *Tree) Tips() []*Node {
	tips := make([]*Node, 0, 2000)
	t.tipsRecur(&tips, nil, nil)
	return tips
}

func (t *Tree) tipsRecur(tips *[]*Node, cur *Node, prev *Node) {
	if cur == nil {
		cur = t.Root()
	}
	if cur.Tip() {
		*tips = append((*tips), cur)
	}
	for _, n := range cur.neigh {
		if n != prev {
			t.tipsRecur(tips, n, cur)
		}
	}
}

// Removes a set of tips from the tree, from tip names
func (t *Tree) RemoveTips(names ...string) error {
	nodeindex := NewNodeIndex(t)

	for _, name := range names {
		n, ok := nodeindex.GetNode(name)
		if !ok {
			return errors.New("No tip named " + name + " in the Tree")
		}
		if len(n.neigh) != 1 {
			return errors.New("The node named " + name + " is not a tip")
		}
		if err := t.removeTip(n); err != nil {
			return err
		}
	}

	return nil
}

// Remove one tip from the tree
func (t *Tree) removeTip(tip *Node) error {
	if len(tip.neigh) != 1 {
		return errors.New("Cannot remove node, it is not a tip")
	}
	tip.neigh = nil
	internal := tip.br[0].left
	if err := internal.delNeighbor(tip); err != nil {
		return err
	}
	tip.neigh = nil
	tip.br[0].left = nil
	tip.br[0].right = nil
	tip.br = nil

	// Then 2 solutions :
	// 1 - Internal node is now terminal : it means it was the root of a rooted tree : we delete it and new root is its child
	// 2 - Internal node is now a bifurcation : we do not want to keep it thus we will delete it and connect the two neighbors
	// Case 1
	if len(internal.neigh) == 1 {
		if t.Root() != internal {
			return errors.New("After tip removal, this node should not have degre 1 without being the root")
		}
		t.root = internal.neigh[0]
		if err := t.root.delNeighbor(internal); err != nil {
			return err
		}
		t.delNode(internal)
		return nil
	}

	// Case 2: We remove the node
	if len(internal.neigh) == 2 {
		n1, n2 := internal.neigh[0], internal.neigh[1]
		b1, b2 := internal.br[0], internal.br[1]
		length1, length2 := b1.Length(), b2.Length()
		sup1, sup2 := b1.Support(), b2.Support()
		var e *Edge
		// Direction : true if n1-->internal
		dir1 := b1.left == n1
		// Direction : true if internal-->n2
		dir2 := b2.right == n2
		if err := n1.delNeighbor(internal); err != nil {
			return err
		}
		if err := n2.delNeighbor(internal); err != nil {
			return err
		}

		// Now we have two options to connect n1 and n2: (n1 parent of n2) or (n2 parent of n1)
		// This direction depends on the directions of the previous edges:
		// 1) n1--->internal--->n2 : t.ConnectNodes(n1,n2)
		// 2) n1<---internal<---n2 : t.ConnectNodes(n2,n1)
		// 3) n1<---internal--->n2 : internal is the root of an unrooted tree so:
		//        1 - we connect the two nodes from n1 to n2 if n1 is not a tip or n2 to n1 otherwise
		//        2 - we choose a new root (n1 if n1->n2, n2 otherwise)
		// 4) n1--->internal<---n2 : Error
		if dir1 && dir2 { // 1)
			e = t.ConnectNodes(n1, n2)
		} else if !dir1 && !dir2 { // 2)
			e = t.ConnectNodes(n2, n1)
		} else if !dir1 && dir2 { // 3
			if t.Root() != internal {
				return errors.New("The tree root is not the internal node, but it should be")
			}
			if len(n1.neigh) > 1 { // Not a tip
				e = t.ConnectNodes(n1, n2)
				t.SetRoot(n1)
			} else if len(n1.neigh) == 1 {
				return errors.New("The neighbor n1 should not have only one neighbor")
			} else if len(n2.neigh) > 1 { // Not a tip
				e = t.ConnectNodes(n2, n1)
				t.SetRoot(n2)
			} else if len(n2.neigh) == 1 {
				return errors.New("The neighbor n2 should not have only one neighbor")
			} else {
				return errors.New("The tree after tip removal is only made of two tips")
			}
		} else {
			return errors.New("Branches of internal node are not oriented as they should be")
		}

		if length1 != -1 || length2 != -1 {
			e.SetLength(math.Max(0, length1) + math.Max(0, length2))
		}

		// We attribute a support to the new branch only if it is not a tip branch
		if (sup1 != -1 || sup2 != -1) && len(n1.neigh) > 1 && len(n2.neigh) > 1 {
			e.SetSupport(math.Max(sup1, sup2))
		}

		t.delNode(internal)
		return nil
	}
	return errors.New("Unknown problem: The internal node remaining after removing the tip has a unexpected number of neighbors")
}

func (t *Tree) String() string {
	return t.Newick()
}

func (t *Tree) Newick() string {
	var buffer bytes.Buffer
	t.root.Newick(nil, &buffer)
	buffer.WriteString(";")
	return buffer.String()
}

// Updates the tipindex which maps tip names
// To index in the bitsets
// Bitset indexes correspond to the position
// of the tip in the alphabetically ordered tip
// name list
func (t *Tree) UpdateTipIndex() {
	names := t.AllTipNames()
	sort.Strings(names)
	for k := range t.tipIndex {
		delete(t.tipIndex, k)
	}
	for i, n := range names {
		t.tipIndex[n] = uint(i)
	}
}

// if UpdateTipIndex has been called before ok
// otherwise returns an error
func (t *Tree) NbTips() (int, error) {
	if len(t.tipIndex) == 0 {
		return 0, errors.New("No tips in the index, tip name index is not initialized")
	}

	return len(t.tipIndex), nil

}

// Returns the bitset index of the tree in the Tree
// Returns an error if the node is not a tip
func (t *Tree) tipIndexNode(n *Node) (uint, error) {
	if len(n.neigh) != 1 {
		return 0, errors.New("Cannot get bitset index of a non tip node")
	}
	return t.TipIndex(n.name)
}

// Return the tip index if the tip with given name exists in the tree
// May return an error if tip index has not been initialized
// With UpdateTipIndex or if the tip does not exist
func (t *Tree) TipIndex(name string) (uint, error) {
	if len(t.tipIndex) == 0 {
		return 0, errors.New("No tips in the index, tip name index is not initialized")
	}
	v, ok := t.tipIndex[name]
	if !ok {
		return 0, errors.New("No tip named " + name + " in the index")
	}
	return v, nil
}

// Return true if the tip with given name exists in the tree
// May return an error if tip index has not been initialized
// With UpdateTipIndex
func (t *Tree) ExistsTip(name string) (bool, error) {
	if len(t.tipIndex) == 0 {
		return false, errors.New("No tips in the index, tip name index is not initialized")
	}
	_, ok := t.tipIndex[name]
	return ok, nil
}

// Returns all the tip name in the tree
// Starts with n==nil (root)
func (t *Tree) AllTipNames() []string {
	names := make([]string, 0, 1000)
	t.allTipNamesRecur(&names, nil, nil)
	return names
}

// Returns all the tip name in the tree
// Starts with n==nil (root)
// It is an internal recursive function
func (t *Tree) allTipNamesRecur(names *[]string, n *Node, parent *Node) {
	if n == nil {
		n = t.Root()
	}
	// is a tip
	if len(n.neigh) == 1 {
		*names = append(*names, n.name)
	} else {
		for _, child := range n.neigh {
			if child != parent {
				t.allTipNamesRecur(names, child, n)
			}
		}
	}
}

func (t *Tree) ConnectNodes(parent *Node, child *Node) *Edge {
	newedge := t.NewEdge()
	newedge.setLeft(parent)
	newedge.setRight(child)
	parent.addChild(child, newedge)
	child.addChild(parent, newedge)
	return newedge
}

// This function takes the first node having 3 neighbors
// and reroot the tree on this node
func (t *Tree) RerootFirst() error {
	for _, n := range t.Nodes() {
		if len(n.neigh) == 3 {
			err := t.Reroot(n)
			return err
		}
	}
	return errors.New("No nodes with 3 neighors have been found for rerooting")
}

func (t *Tree) ClearBitSets() error {
	length := uint(len(t.tipIndex))
	if length == 0 {
		return errors.New("No tips in the index, tip name index is not initialized")
	}
	t.clearBitSetsRecur(nil, nil, length)
	return nil
}

// Recursively update bitsets of edges from the Node n
// If node == nil then it starts from the root
func (t *Tree) clearBitSetsRecur(n *Node, parent *Node, ntip uint) {
	if n == nil {
		n = t.Root()
	}

	for i, child := range n.neigh {
		if child != parent {
			e := n.br[i]
			e.bitset = nil
			e.bitset = bitset.New(ntip)
			t.clearBitSetsRecur(child, n, ntip)
		}
	}
}

// Updates bitsets of all edges in the tree
// Assumes that the hashmap tip name : index is
// initialized with UpdateTipIndex function
func (t *Tree) UpdateBitSet() error {
	rightedges := make([]*Edge, 0, 2000)
	for _, e := range t.Root().br {
		rightedges = rightedges[:0]
		rightedges = append(rightedges, e)
		err := t.fillRightBitSet(e, &rightedges)
		if err != nil {
			return err
		}
	}
	return nil
}

// Recursively clears and sets the bitsets of the descending edges
//
func (t *Tree) fillRightBitSet(currentEdge *Edge, rightEdges *[]*Edge) error {
	if currentEdge.bitset == nil {
		return errors.New("BitSets has not been initialized with tree.clearBitSetsRecur(nil, nil, uint(len(tree.tipIndex)))")
	}
	currentEdge.bitset.ClearAll()
	// If we are at a tip edge
	// We set at 1 the bits of the tip in
	// the bitsets of all rightEdges
	if len(currentEdge.right.neigh) == 1 {
		i, err := t.tipIndexNode(currentEdge.right)
		if err != nil {
			return err
		}
		for _, e := range *rightEdges {
			e.bitset.Set(i)
		}
	} else {
		// Else
		for _, e2 := range currentEdge.right.br {
			if e2.left == currentEdge.right {
				*rightEdges = append(*rightEdges, e2)
				err := t.fillRightBitSet(e2, rightEdges)
				if err != nil {
					return err
				}
				*rightEdges = (*rightEdges)[:len(*rightEdges)-1]
			}
		}
	}
	return nil
}

// This function compares 2 trees and output
// the number of edges in common
// If the trees have different sets of tip names, returns an error
// It assumes that functions
// 	tree.UpdateTipIndex()
//	tree.ClearBitSets()
//	tree.UpdateBitSet()
// If tipedges is false: does not take into account tip edges
// Have been called before, otherwise will output an error
func (t *Tree) CommonEdges(t2 *Tree, tipEdges bool) (tree1 int, common int, err error) {

	err = t.CompareTipIndexes(t2)

	if err != nil {
		return 0, 0, err
	}

	edges1 := t.Edges()
	edges2 := t2.Edges()

	tree1, common, err = CommonEdges(edges1, edges2, tipEdges)

	return tree1, common, nil
}

// This function compares 2 trees and output
// the number of edges in common
// It does not check if the trees have different sets of tip names, but just compare the bitsets
// If applied on two tree with the same number of tips with different names, it will give results anyway
// It assumes that functions
// 	tree.UpdateTipIndex()
//	tree.ClearBitSets()
//	tree.UpdateBitSet()
// If tipedges is false: does not take into account tip edges
// Have been called before, otherwise will output an error
func CommonEdges(edges1 []*Edge, edges2 []*Edge, tipEdges bool) (tree1 int, common int, err error) {
	var e, e2 *Edge
	for _, e = range edges1 {
		if tipEdges || !e.right.Tip() {
			tree1++
			if e2, err = e.FindEdge(edges2); err != nil {
				return -1, -1, err
			}
			if e2 != nil {
				common++
			}
		}
	}
	tree1 = tree1 - common
	return tree1, common, nil
}

// This function compares the tip name indexes of 2 trees
// If the tipindexes have the same size (!=0) and have the same set of tip names,
// The returns nil, otherwise returns an error
func (t *Tree) CompareTipIndexes(t2 *Tree) error {
	if len(t.tipIndex) == 0 ||
		len(t2.tipIndex) == 0 ||
		len(t.tipIndex) != len(t2.tipIndex) {
		return errors.New("Tip name index is not initialized or trees do not have the same number of tips")
	}

	for k := range t.tipIndex {
		_, ok := t2.tipIndex[k]
		if !ok {
			return errors.New("Trees do not have the same tip names")
		}
	}

	for k := range t2.tipIndex {
		_, ok := t.tipIndex[k]
		if !ok {
			return errors.New("Trees do not have the same tip names")
		}
	}
	return nil
}

// This function takes a node and reroot the tree on that node
// It reorients edges left-edge-right : see ReorderEdges
// The node must be one of the tree nodes, otherwise it returns an error
func (t *Tree) Reroot(n *Node) error {
	intree := false
	for _, n2 := range t.Nodes() {
		if n2 == n {
			intree = true
		}
	}
	if !intree {
		return errors.New("The node is not part of the tree")
	}
	t.root = n
	err := t.reorderEdges(n, nil)
	return err
}

// This function reorders the edges of a tree
// in order to always have left-edge-right
// with left node being parent of right node
// with respect to the given root node
// Important even for unrooted trees
// Useful mainly after a reroot
func (t *Tree) reorderEdges(n *Node, prev *Node) error {
	for _, next := range n.br {
		if next.right != prev && next.left != prev {
			if next.right == n {
				next.right, next.left = next.left, next.right
			}
			t.reorderEdges(next.right, n)
		}
	}
	return nil
}

// This function graft the Node n at the middle of the Edge e
// It divides the branch lenght by 2
// It returns the added edges and the added nodes
func (t *Tree) GraftTipOnEdge(n *Node, e *Edge) (*Edge, *Edge, *Node, error) {
	newnode := t.NewNode()
	newedge := t.NewEdge()

	lnode := e.left
	rnode := e.right

	// index of edge in neighbors of l
	e_l_ind, err := lnode.EdgeIndex(e)
	if err != nil {
		return nil, nil, nil, err
	}
	// index of edge in neighbors of r
	e_r_ind, err2 := rnode.EdgeIndex(e)
	if err2 != nil {
		return nil, nil, nil, err2
	}

	newedge.SetLength(1.0)
	newedge.setLeft(newnode)
	newedge.setRight(n)
	newnode.addChild(n, newedge)
	n.addChild(newnode, newedge)
	e.setRight(newnode)
	newnode.addChild(lnode, e)
	lnode.neigh[e_l_ind] = newnode

	if lnode.br[e_l_ind] != e {
		return nil, nil, nil, errors.New("The Edge is not at the same index")
	}

	newedge2 := t.NewEdge()
	newedge2.SetLength(e.length / 2)
	e.SetLength(e.length / 2)
	newedge2.setLeft(newnode)
	newedge2.setRight(rnode)
	newnode.addChild(rnode, newedge2)
	if rnode.br[e_r_ind] != e {
		return nil, nil, nil, errors.New("The Edge is not at the same index")
	}
	rnode.neigh[e_r_ind] = newnode
	rnode.br[e_r_ind] = newedge2
	return newedge, newedge2, newnode, nil
}

//Creates a Random Binary tree
//nbtips : Number of tips of the random binary tree to create
func RandomBinaryTree(nbtips int) (*Tree, error) {
	t := NewTree()
	if nbtips < 2 {
		return nil, errors.New("Cannot create a random binary tree with less than 2 tips")
	}
	edges := make([]*Edge, 0, 2000)
	for i := 1; i < nbtips; i++ {
		n := t.NewNode()
		n.SetName("Tip" + strconv.Itoa(i))
		switch len(edges) {
		case 0:
			n2 := t.NewNode()
			n2.SetName("Node" + strconv.Itoa(i-1))
			e := t.ConnectNodes(n2, n)
			edges = append(edges, e)
			e.SetLength(1.0)
			t.SetRoot(n2)
		default:
			// Where to insert the new tip
			i_edge := rand.Intn(len(edges))
			e := edges[i_edge]
			newedge, newedge2, _, err := t.GraftTipOnEdge(n, e)

			edges = append(edges, newedge)
			edges = append(edges, newedge2)

			if err != nil {
				return nil, err
			}
		}
	}
	err := t.RerootFirst()
	t.UpdateTipIndex()
	t.ClearBitSets()
	t.UpdateBitSet()
	t.ComputeDepths()
	return t, err
}
func (t *Tree) ComputeDepths() {
	if t.Rooted() {
		t.computeDepthRecurRooted(t.Root(), nil)
	} else {
		t.computeDepthUnRooted()
	}
}

func (t *Tree) computeDepthRecurRooted(n *Node, prev *Node) int {
	if n.Tip() {
		n.depth = 0
		return n.depth
	} else {
		mindepth := -1
		for _, next := range n.neigh {
			if next != prev {
				depth := t.computeDepthRecurRooted(next, n)
				if mindepth == -1 || depth < mindepth {
					mindepth = depth
				}
			}
		}
		n.depth = mindepth + 1
		return n.depth
	}
}

func (t *Tree) computeDepthUnRooted() {
	nodes := t.Tips()
	currentlevel := 0
	nbchanged := 1
	for nbchanged != 0 {
		nbchanged = 0
		nextnodes := make([]*Node, 0, 2000)
		for _, n := range nodes {
			if n.depth == -1 {
				n.depth = currentlevel
				nbchanged++
			}
		}
		for _, n := range nodes {
			for _, next := range n.neigh {
				if next.depth == -1 {
					nextnodes = append(nextnodes, next)
				}
			}
		}
		nodes = nextnodes
		currentlevel++
	}
}

// This function shuffles the tips of the tree
// and recompute tipindex and bitsets
func (t *Tree) ShuffleTips() {
	tips := t.Tips()
	names := t.AllTipNames()
	permutation := rand.Perm(len(names))

	for i, p := range permutation {
		tips[i].SetName(names[p])
	}

	t.UpdateTipIndex()
	t.ClearBitSets()
	t.UpdateBitSet()
}

func (t *Tree) CollapseShortBranches(length float64) {
	edges := t.Edges()
	shortbranches := make([]*Edge, 0, 1000)
	for _, e := range edges {
		if e.Length() <= length {
			shortbranches = append(shortbranches, e)
		}
	}
	t.RemoveEdges(shortbranches...)
}

// Removes branches from the tree if they are not tip edges
// And if they do not connects the root of a rooted tree
// Merges the 2 nodes and creates multifurcations
// At the end, bitsets should not need to be updated
func (t *Tree) RemoveEdges(edges ...*Edge) {
	for _, e := range edges {
		// Tip node
		if e.Right().Tip() {
			continue
		}
		// Root node
		if e.Right().Nneigh() == 2 || e.Left().Nneigh() == 2 {
			continue
		}
		// Remove the edge from left and right node
		e.Left().delNeighbor(e.Right())
		e.Right().delNeighbor(e.Left())

		// Move the edges on right node to left node
		for _, child := range e.Right().Neigh() {
			if child != e.Left() {
				idx, err := child.NodeIndex(e.Right())
				if err != nil {
					io.ExitWithMessage(err)
				}
				child.neigh[idx] = e.Left()
				if child.br[idx].left == e.Right() {
					child.br[idx].left = e.Left()
				} else {
					io.ExitWithMessage(errors.New("Problem in edge orientation"))
				}
				e.Left().addChild(child, child.br[idx])
			}
		}
	}
}

func (t *Tree) SumBranchLengths() float64 {
	sumlen := 0.0
	for _, e := range t.Edges() {
		sumlen += e.Length()
	}
	return sumlen
}

func (t *Tree) UnRoot() {
	if !t.Rooted() {
		return
	}

	root := t.Root()
	n1 := root.Neigh()[0]
	n2 := root.Neigh()[1]

	n1tip := n1.Tip()

	e1 := root.br[0]
	e2 := root.br[1]

	n1.delNeighbor(t.Root())
	n2.delNeighbor(t.Root())

	var e3 *Edge

	if n1tip {
		e3 = t.ConnectNodes(n2, n1)
		t.SetRoot(n2)
	} else {
		e3 = t.ConnectNodes(n1, n2)
		t.SetRoot(n1)
	}

	if e1.Length() != -1 || e2.Length() != -1 {
		e3.SetLength(math.Max(0, e1.Length()) + math.Max(0, e2.Length()))
	}
	if !n1.Tip() && !n2.Tip() && (e1.Support() != -1 || e2.Support() != -1) {
		e3.SetSupport(math.Max(0, e1.Support()) + math.Max(0, e2.Support()))
	}
	t.delNode(root)
}

// This function renames nodes of the tree based on the map in argument
// If a name in the map does not exist in the tree, then returns an error
// If a node/tip in the tree does not have a name in the map: OK
// After rename, tip index is updated, as well as bitsets of the edges
func (t *Tree) Rename(namemap map[string]string) error {
	nodeindex := NewNodeIndex(t)
	for name, newname := range namemap {
		node, ok := nodeindex.GetNode(name)
		if ok {
			node.SetName(newname)
		}
	}
	// After we update bitsets if any, and node indexes
	t.UpdateTipIndex()
	err := t.ClearBitSets()
	if err != nil {
		return err
	}
	t.UpdateBitSet()
	return nil
}

func (t *Tree) MeanBrLength() float64 {
	mean := 0.0
	edges := t.Edges()
	for _, e := range edges {
		mean += e.Length()
	}
	return mean / float64(len(edges))
}

func (t *Tree) MeanSupport() float64 {
	mean := 0.0
	edges := t.Edges()
	i := 0
	for _, e := range edges {
		if !e.Right().Tip() {
			mean += e.Support()
			i++
		}
	}
	return mean / float64(i)
}

func (t *Tree) MedianSupport() float64 {
	edges := t.Edges()
	tips := t.Tips()
	supports := make([]float64, len(edges)-len(tips))

	i := 0
	for _, e := range edges {
		if !e.Right().Tip() {
			supports[i] = e.Support()
			i++
		}
	}
	sort.Float64s(supports)

	middle := len(supports) / 2
	result := supports[middle]
	if len(supports)%2 == 0 {
		result = (result + supports[middle-1]) / 2
	}
	return result
}
