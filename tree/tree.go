/*
   Package gotree implements a simple
   library for handling phylogenetic trees in go
*/
package tree

import (
	"bytes"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"regexp"
	"sort"
	"strconv"

	"github.com/evolbioinfo/gotree/io"
	"github.com/evolbioinfo/gotree/mutils"
	"github.com/fredericlemoine/bitset"
)

// Tree structure having a root and a tip index, that maps tip names to their index
type Tree struct {
	root     *Node           // root node: If the tree is unrooted the root node should have 3 children
	tipIndex map[string]uint // Map between tip name and bitset index
}

// Type for channel of trees
type Trees struct {
	Tree *Tree
	Id   int
	Err  error
}

// Initialize a new empty Tree
func NewTree() *Tree {
	return &Tree{
		root:     nil,
		tipIndex: make(map[string]uint, 0),
	}
}

// Initialize a new empty Node
func (t *Tree) NewNode() *Node {
	return &Node{
		name:    "",
		comment: make([]string, 0),
		neigh:   make([]*Node, 0, 3),
		br:      make([]*Edge, 0, 3),
		depth:   NIL_DEPTH,
		id:      NIL_ID,
	}
}

// Set to nil the node and all its branches
func (t *Tree) delNode(n *Node) {
	for i, _ := range n.neigh {
		n.neigh[i] = nil
	}
	n.neigh = nil

	for i, e := range n.br {
		e.left = nil
		e.right = nil
		e.bitset = nil
		n.br[i] = nil
	}
	n.br = nil
}

// Set to nil the node but do not touch
// branches
func (t *Tree) unconnectNode(n *Node) {
	for i, _ := range n.neigh {
		n.neigh[i] = nil
	}
	n.neigh = nil

	for i, _ := range n.br {
		n.br[i] = nil
	}
	n.br = nil
}

// Initialize a new empty Edge
func (t *Tree) NewEdge() *Edge {
	return &Edge{
		length:  NIL_LENGTH,
		support: NIL_SUPPORT,
		id:      NIL_ID,
		pvalue:  NIL_PVALUE,
	}
}

/* Tree functions */
/******************/

// Set a root for the tree. This does not check that the
// node is part of the tree. It may be useful to call
//	t.ReinitIndexes()
// After setting a new root, to update branch bitsets.
func (t *Tree) SetRoot(r *Node) {
	t.root = r
}

// Returns the current root of the tree
func (t *Tree) Root() *Node {
	return t.root
}

// Returns true if the tree is rooted (i.e. root node
// has 2 neighbors), and false otherwise.
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

// Recursive function to list all edges of the tree
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

// Returns all internal edges of the tree (do it recursively)
func (t *Tree) InternalEdges() []*Edge {
	edges := make([]*Edge, 0, 2000)
	for _, e := range t.Root().br {
		if !e.Right().Tip() {
			edges = append(edges, e)
			t.internalEdgesRecur(e, &edges)
		}
	}
	return edges
}

// recursive function that lists all internal edges of the tree
func (t *Tree) internalEdgesRecur(edge *Edge, edges *[]*Edge) {
	if len(edge.right.neigh) > 1 {
		for _, child := range edge.right.br {
			if child.left == edge.right && !child.Right().Tip() {
				*edges = append((*edges), child)
				t.edgesRecur(child, edges)
			}
		}
	}
}

// Returns all the external edges (tip) of the tree (do it recursively)
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

// recursive function that lists all external edges (tips) of the tree
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

// recursive function that lists all nodes of the tree
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

// recursive function that lists all tips of the tree
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

// Returns the list of nodes having a name matching the given regexp
// May return an error if the regexp is malformed.
// In this case, returns an empty (not nil) slice of nodes.
func (t *Tree) SelectNodes(re string) ([]*Node, error) {
	nodes := make([]*Node, 0)
	if r, err := regexp.Compile(re); err == nil {
		for _, n := range t.Nodes() {
			if r.MatchString(n.Name()) {
				nodes = append(nodes, n)
			}
		}
		return nodes, err
	} else {
		return nodes, err
	}
}

// Removes a set of tips from the tree, given their names
//
// if revert is true, then keeps only tips with the given names
//
// Removed tips
func (t *Tree) RemoveTips(revert bool, names ...string) error {
	namemap := make(map[string]bool)

	for _, name := range names {
		namemap[name] = true
	}

	for _, tip := range t.Tips() {
		if len(tip.neigh) != 1 {
			return errors.New("The node named " + tip.Name() + " is not a tip")
		}

		_, ok := namemap[tip.Name()]
		if (!revert && ok) || (revert && !ok) {
			if err := t.removeTip(tip); err != nil {
				return err
			}
		}
	}
	t.ReinitInternalIndexes()
	return nil
}

// Removes one tip from the tree. The internal node may be removed, example:
//	          t1
//	         /
//	 n0--n2--
//	         \
//	          t2
// If we remove t2, then n2 must be removed. In that case, we remove n2 and
// connect n0 to t1 with a branch having :
//	* length=length(n0--n2)+length(n2--t1)
//	* support=max(support(n0--n2),support(n2--t1))
func (t *Tree) removeTip(tip *Node) error {
	if len(tip.neigh) != 1 {
		return errors.New("Cannot remove node, it is not a tip")
	}
	tip.neigh = nil
	internal := tip.br[0].left
	if err := internal.delNeighbor(tip); err != nil {
		return err
	}
	t.delNode(tip)

	// Then 3 solutions :
	// 1 - Internal node is now terminal : it means it was the root of a rooted tree : we delete it and new root is its child
	// 2 - Internal node is now a bifurcation : we do not want to keep it thus we will delete it and connect the two neighbors
	// 3 - Internal node still has a degree > 2 : We do not do anything => the node should remain

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
		// 4) n1--->internal<---n2 : Error : A node cannot have 2 parents
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

		if length1 != NIL_LENGTH || length2 != NIL_LENGTH {
			e.SetLength(math.Max(0, length1) + math.Max(0, length2))
		}

		// We attribute a support to the new branch only if it is not a tip branch
		if (sup1 != NIL_SUPPORT || sup2 != NIL_SUPPORT) && len(n1.neigh) > 1 && len(n2.neigh) > 1 {
			e.SetSupport(math.Max(sup1, sup2))
		}

		t.delNode(internal)
		return nil
	}
	//Case 3 : Nothing
	return nil
	//return errors.New("Unknown problem: The internal node remaining after removing the tip has a unexpected number of neighbors")
}

// Returns a newick string representation of this tree
// It calls t.Newick()
func (t *Tree) String() string {
	return t.Newick()
}

// Returns a newick string representation of this tree
func (t *Tree) Newick() string {
	var buffer bytes.Buffer
	t.root.Newick(nil, &buffer)
	if len(t.root.comment) != 0 {
		for _, c := range t.root.comment {
			buffer.WriteString("[")
			buffer.WriteString(c)
			buffer.WriteString("]")
		}
	}
	buffer.WriteString(";")
	return buffer.String()
}

// returns a Nexus string representation of this tree
func (t *Tree) Nexus() string {
	newick := t.Newick()
	var buffer bytes.Buffer
	buffer.WriteString("#NEXUS\n")
	buffer.WriteString("BEGIN TAXA;\n")
	tips := t.Tips()
	buffer.WriteString(" DIMENSIONS NTAX=")
	buffer.WriteString(strconv.Itoa(len(tips)))
	buffer.WriteString(";\n")
	buffer.WriteString(" TAXLABELS")
	for _, tip := range tips {
		buffer.WriteString(" " + tip.Name())
	}
	buffer.WriteString(";\n")
	buffer.WriteString("END;\n")
	buffer.WriteString("BEGIN TREES;\n")
	buffer.WriteString("  TREE tree1 = ")
	buffer.WriteString(newick)
	buffer.WriteString("\n")
	buffer.WriteString("END;\n")
	return buffer.String()
}

// Updates the tipindex which maps tip names to
// their index in the bitsets.
//
// Bitset indexes correspond to the position
// of the tip in the alphabetically ordered tip
// name list
func (t *Tree) UpdateTipIndex() (err error) {
	names := t.SortedTips()
	for k := range t.tipIndex {
		delete(t.tipIndex, k)
	}
	for i, n := range names {
		if _, ok := t.tipIndex[n]; ok {
			err = errors.New("Cannot create a tip index when several tips have the same name")
			return
		}
		t.tipIndex[n] = uint(i)
	}
	return
}

/* Tips, sorted by their order in the bitsets*/
func (t *Tree) SortedTips() []string {
	names := t.AllTipNames()
	sort.Strings(names)
	return names
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

// Connects the two nodes in argument by an edge that is returned.
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
// It then recomputes indices
func (t *Tree) RerootFirst() error {
	for _, n := range t.Nodes() {
		if len(n.neigh) == 3 {
			err := t.Reroot(n)
			return err
		}
	}
	return errors.New("No nodes with 3 neighors have been found for rerooting")
}

// Clears all bitsets associated to all edges
func (t *Tree) ClearBitSets() error {
	length := uint(len(t.tipIndex))
	if length == 0 {
		return errors.New("No tips in the index, tip name index is not initialized")
	}
	t.clearBitSetsRecur(nil, nil, length)
	return nil
}

// This Function initializes or reinitializes
// memory consuming structures:
//	* bitset indexes
//	* Tipindex
//	* And computes node depths
func (t *Tree) ReinitIndexes() (err error) {
	if err = t.UpdateTipIndex(); err != nil {
		return
	}

	t.ClearBitSets()
	t.UpdateBitSet()
	t.ComputeEdgeHashes(nil, nil, nil)
	t.ComputeDepths()
	return
}

// This Function initializes or reinitializes
// memory consuming structures:
//	* bitset indexes
//	* Tipindex
//	* And computes node depths
func (t *Tree) ReinitInternalIndexes() {
	t.ClearBitSets()
	t.UpdateBitSet()
	t.ComputeEdgeHashes(nil, nil, nil)
	t.ComputeDepths()
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
			e.hashcodeleft = 0
			e.hashcoderight = 0
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

// This function compares 2 trees and returns the number of edges in common
// If the trees have different sets of tip names, returns an error.
//
// It assumes that functions
//	tree.UpdateTipIndex()
//	tree.ClearBitSets()
//	tree.UpdateBitSet()
// Have been called before, otherwise will output an error
//
// If tipedges is false: does not take into account tip edges
func (t *Tree) CommonEdges(t2 *Tree, tipEdges bool) (tree1 int, common int, err error) {

	if err = t.CompareTipIndexes(t2); err != nil {
		return
	}

	edges1 := t.Edges()
	edges2 := t2.Edges()

	tree1, common, err = CommonEdges(edges1, edges2, tipEdges)

	return tree1, common, nil
}

// This function compares 2 trees and returns the number of edges in common.
//
// It does not check if the trees have different sets of tip names,
// but just compare the bitsets. If called on two trees with the same
// number of tips with different names, it will give meaningless
// results.
//
// It assumes that functions
// 	tree.UpdateTipIndex()
//	tree.ClearBitSets()
//	tree.UpdateBitSet()
// Have been called before, otherwise will output an error
//
// If tipedges is false: does not take into account tip edges
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
//
// If the tipindexes have the same size (!=0) and have the
// same set of tip names, then returns nil, otherwise returns an error
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

	// for k := range t2.tipIndex {
	// 	_, ok := t.tipIndex[k]
	// 	if !ok {
	// 		return errors.New("Trees do not have the same tip names")
	// 	}
	// }
	return nil
}

// This function takes a node and reroots the tree on that node.
//
// It reorients edges left-edge-right : see ReorderEdges()
//
// The node must be part of the tree, otherwise it returns an error
func (t *Tree) Reroot(n *Node) error {
	if n.Nneigh() < 2 {
		return errors.New("Cannot reroot on a tip node")
	}
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
	err := t.ReorderEdges(n, nil, nil)
	t.ReinitInternalIndexes()
	return err
}

// This function takes a node and reroots the tree on that node.
//
// It reorients edges left-edge-right : see ReorderEdges()
//
// It does not check wether the node is part of tree or not
//
// Returns an error If the given node is a tip or if it can
// not reorder edges
func (t *Tree) reroot_nocheck(n *Node) error {
	if n.Nneigh() < 2 {
		return errors.New("Cannot reroot on a tip node")
	}
	t.root = n
	err := t.ReorderEdges(n, nil, nil)
	return err
}

// This function reorders the edges of a tree in order to always have
// left-edge-right with left node being parent of right node with respect
// to the given root node.
//
// Important even for unrooted trees. Useful mainly after a reroot.
//
// It updates "reversed" edge slice, edges that have been reversed
func (t *Tree) ReorderEdges(n *Node, prev *Node, reversed *[]*Edge) error {
	for _, next := range n.br {
		if next.right != prev && next.left != prev {
			if next.right == n {
				next.right, next.left = next.left, next.right
				if reversed != nil {
					(*reversed) = append((*reversed), next)
				}
			}
			t.ReorderEdges(next.right, n, reversed)
		}
	}
	return nil
}

// This function grafts the Tip n at the middle of the Edge e.
//
// Example:
//	* Before
//		*--e--*
//	* After
//		*--e1--newnode--e2--*
//		          |
//		          n
//
// To do so, it divides the branch lenght by 2,and returns the 2 new
// edges and the new internal node.
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

// Computes detphs of all nodes. Depth of internal node n is defined as
// the length of the path from n to the closest tip. Depth of tip nodes
// is 0.
//
// Depth is then accessible by n.Depth() for any node n.
func (t *Tree) ComputeDepths() {
	if t.Rooted() {
		t.computeDepthRecurRooted(t.Root(), nil)
	} else {
		t.computeDepthUnRooted()
	}
}

// Recursive function to compute depths for a rooted tree
func (t *Tree) computeDepthRecurRooted(n *Node, prev *Node) int {
	if n.Tip() {
		n.depth = 0
		return n.depth
	} else {
		mindepth := NIL_DEPTH
		for _, next := range n.neigh {
			if next != prev {
				depth := t.computeDepthRecurRooted(next, n)
				if mindepth == NIL_DEPTH || depth < mindepth {
					mindepth = depth
				}
			}
		}
		n.depth = mindepth + 1
		return n.depth
	}
}

// Recursive function to compute depths for an unrooted tree
func (t *Tree) computeDepthUnRooted() {
	nodes := t.Tips()
	currentlevel := 0
	nbchanged := 1
	for nbchanged != 0 {
		nbchanged = 0
		nextnodes := make([]*Node, 0, 2000)
		for _, n := range nodes {
			if n.depth == NIL_DEPTH {
				n.depth = currentlevel
				nbchanged++
			}
		}
		for _, n := range nodes {
			for _, next := range n.neigh {
				if next.depth == NIL_DEPTH {
					nextnodes = append(nextnodes, next)
				}
			}
		}
		nodes = nextnodes
		currentlevel++
	}
}

// This function shuffles the tips of the tree
// and recompute tipindex and bitsets.
//
// The tree have the same topology, but tip names
// are reassigned randomly.
func (t *Tree) ShuffleTips() {
	tips := t.Tips()
	names := t.AllTipNames()
	permutation := rand.Perm(len(names))

	for i, p := range permutation {
		tips[i].SetName(names[p])
	}

	t.ReinitIndexes()
}

// Randomly rotates neighbors of all internal nodes
//
// It does not change the topology, but just the way
// the tree is traversed.
func (t *Tree) RotateInternalNodes() {
	for _, n := range t.Nodes() {
		n.RotateNeighbors()
	}
}

// Sort neighbors of all nodes by their number of tips
//
// May give better results for drawing.
func (t *Tree) SortNeighborsByTips() {
	t.sortNeighbors(t.Root(), nil)
}

// Sorts the neighbors of current node after having recursively
// sorted the neigbors of its children
func (t *Tree) sortNeighbors(cur, prev *Node) int {
	// Number of tips at the end of each neighbor of cur
	neighbors := make([]struct {
		ntips int
		neigh *Node
		br    *Edge
	}, len(cur.Neigh()))
	total := 0
	for i, c := range cur.Neigh() {
		neighbors[i].neigh = c
		neighbors[i].br = cur.Edges()[i]
		if c != prev {
			neighbors[i].ntips = t.sortNeighbors(c, cur)
			total += neighbors[i].ntips
		} else {
			neighbors[i].ntips = 0
		}
	}
	// we sort neighbor slice
	sort.SliceStable(neighbors, func(i, j int) bool { return neighbors[i].ntips < neighbors[j].ntips })
	//
	for i, _ := range cur.Neigh() {
		cur.neigh[i] = neighbors[i].neigh
		cur.br[i] = neighbors[i].br
	}

	if cur.Tip() {
		return 1
	} else {
		return total
	}
}

// Collapses (removes) the branches having
// length <= length threshold
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

// Collapses (removes) the branches having
// support < support threshold && support != NIL_SUPPORT (exists)
func (t *Tree) CollapseLowSupport(support float64) {
	edges := t.Edges()
	lowsupportbranches := make([]*Edge, 0, 1000)
	for _, e := range edges {
		if e.Support() != NIL_SUPPORT && e.Support() < support {
			lowsupportbranches = append(lowsupportbranches, e)
		}
	}
	t.RemoveEdges(lowsupportbranches...)
}

// Collapses (removes) the branches having their depth d
// (# taxa on the lightest side of the bipartition) such that
// mindepththreshold<=d<=maxdepththreshold
func (t *Tree) CollapseTopoDepth(mindepthThreshold, maxdepthThreshold int) error {
	edges := t.Edges()
	depthbranches := make([]*Edge, 0, 1000)
	for _, e := range edges {
		if d, err := e.TopoDepth(); err != nil {
			return err
		} else {
			if d >= mindepthThreshold && d <= maxdepthThreshold {
				depthbranches = append(depthbranches, e)
			}
		}
	}
	t.RemoveEdges(depthbranches...)
	return nil
}

// Resolves multifurcating nodes (>3 neighbors).
//
// If any node has more than 3 neighbors, then
// neighbors are resolved randomly by adding 0 length
// branches until 3 neighbors are remaining.
//
// This function does not update bitsets on edges.
//
// If needed, the calling function must do it with:
//	err := t.ClearBitSets()
//	if err != nil {
//		return err
//	}
//	t.UpdateBitSet()
func (t *Tree) Resolve() {
	root := t.Root()

	t.resolveRecur(root, nil)
	t.ReinitInternalIndexes()
}

// Recursive function to resolve
// multifurcating nodes (see Resolve).
//
// Post order: We first resolve neighbors,
// and then resolve the current node.
//
// This function does not update bitsets on edges:
// The calling function must do it with:
//	err := t.ClearBitSets()
//	if err != nil {
//		return err
//	}
//	t.UpdateBitSet()
func (t *Tree) resolveRecur(current, previous *Node) {
	// Resolve neighbors
	for _, n := range current.Neigh() {
		if n != previous {
			t.resolveRecur(n, current)
		}
	}
	// Resolve the current node if needed
	if len(current.Neigh()) > 3 {
		// Neighbors to group : All neighbors except the "parent"
		// And random order in the array
		l := len(current.Neigh())
		if previous != nil {
			l--
		}
		togroup := make([]*Edge, l)
		perm := rand.Perm(l)
		nb := 0
		for i, n := range current.Neigh() {
			if n != previous {
				togroup[perm[nb]] = current.Edges()[i]
				nb++
			}
		}
		// Now we take neighbors 2 by 2 in reverse order of togroup
		for len(current.Neigh()) > 3 {
			// And add a new node that will connect the 2 neighbors
			n2 := t.NewNode()
			// Take the two last edges of perm
			for i := 0; i < 2; i++ {
				// And an edge between current and the new node
				e := togroup[len(togroup)-1]
				// Remove last element of togroup
				togroup = togroup[:len(togroup)-1]
				boot := e.Support()
				len := e.Length()
				pv := e.PValue()
				other := e.Right()
				other.delNeighbor(current)
				current.delNeighbor(other)
				etmp := t.ConnectNodes(n2, other)
				etmp.SetLength(len)
				etmp.SetSupport(boot)
				etmp.SetPValue(pv)
			}
			// Connect new node to current node
			e := t.ConnectNodes(current, n2)
			e.SetLength(0.0)
			e.SetSupport(NIL_SUPPORT)
			e.SetPValue(NIL_PVALUE)
			// Update togroup removing two last nodes and adding the new node at the end
			togroup = append(togroup, e)
		}
	}
}

// Removes Edges for which the left node has a unique child:
//
// Example:
//	           t1           t1
//	           /	       /
//	 n0--n1--n2   => n0--n2
//	           \	       \
//	            t2          t2
// Will remove edge n1-n2 and keep node n2 informations (name, etc.)
// It adds n1-n2 length to n0-n1 (if any) and keeps n0-n1 support (if any)
// Useful for cleaning ncbi taxonomy for example.
func (t *Tree) RemoveSingleNodes() {
	root := t.Root()

	t.removeSingleNodesRecur(root, nil, nil)
	t.ReinitInternalIndexes()
}

// Removes recursively Edges for which the left node has a unique child.
//
// Post order: We first remove in subtrees, and then look at the
// current node.
//
// This function does not update bitsets on edges, the calling function
// must do it with:
//	err := t.ClearBitSets()
//	if err != nil {
//		return err
//	}
//	t.UpdateBitSet()
func (t *Tree) removeSingleNodesRecur(current, previous *Node, e *Edge) error {
	// Resolve neighbors
	// Temporary slice of node neighbors (the real neighbor slice will be updated
	// during traversal)
	tmpnodes := make([]*Node, len(current.Neigh()))
	tmpedges := make([]*Edge, len(current.Edges()))
	copy(tmpnodes, current.Neigh())
	copy(tmpedges, current.Edges())
	for i, n := range tmpnodes {
		if n != previous {
			t.removeSingleNodesRecur(n, current, tmpedges[i])
		}
	}
	tmpnodes = nil
	tmpedges = nil
	// Remove the current node if needed connect descendant node to parent
	if len(current.Neigh()) == 2 && current != t.Root() {
		// Remove the edge from left and right node
		length := e.Length()
		support := e.Support()
		current.delNeighbor(previous)
		previous.delNeighbor(current)
		e.left = nil
		e.right = nil
		// Connect the edges of children if current to parent node (previous)
		for _, child := range current.Neigh() {
			if child != previous {
				idx, err := child.NodeIndex(current)
				if err != nil {
					return err
				}
				child.neigh[idx] = previous
				if child.br[idx].left == current {
					child.br[idx].left = previous
					child.br[idx].support = math.Max(child.br[idx].support, support)
				} else {
					return errors.New("Problem in edge orientation")
				}
				previous.addChild(child, child.br[idx])
				if child.br[idx].Length() != NIL_LENGTH && length != NIL_LENGTH {
					child.br[idx].SetLength(child.br[idx].Length() + length)
				}
			}
		}
		t.unconnectNode(current)
	}
	return nil
}

// Clears support (set to NIL_SUPPORT) of all branches of the tree
func (t *Tree) ClearSupports() {
	edges := t.Edges()
	for _, e := range edges {
		e.SetSupport(NIL_SUPPORT)
		e.SetPValue(NIL_PVALUE)
	}
}

// Clears pvalues associated with supports (set to NIL_PVALUE) of all branches of the tree
func (t *Tree) ClearPvalues() {
	edges := t.Edges()
	for _, e := range edges {
		e.SetPValue(NIL_PVALUE)
	}
}

// Clears length (set to NIL_LENGTH) of all branches of the tree
func (t *Tree) ClearLengths() {
	edges := t.Edges()
	for _, e := range edges {
		e.SetLength(NIL_LENGTH)
	}
}

// Clears comments associated to all nodes, tips and edges of the tree
func (t *Tree) ClearComments() {
	t.ClearNodeComments()
	t.ClearEdgeComments()
}

func (t *Tree) ClearNodeComments() {
	nodes := t.Nodes()
	for _, n := range nodes {
		n.ClearComments()
	}
}

func (t *Tree) ClearEdgeComments() {
	edges := t.Edges()
	for _, e := range edges {
		e.ClearComments()
	}
}

// Removes the given branches from the tree if they are not
// tip edges and if they do not connect to the root of a rooted tree.
//
// Merges the 2 nodes and creates multifurcations.
//
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
		t.unconnectNode(e.Right())
	}
	t.ReinitInternalIndexes()
}

// Unroots a rooted tree by removing the bifurcating root, and
// rerooting to one of the non tip direct children of the previous root.
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

	if e1.Length() != NIL_LENGTH || e2.Length() != NIL_LENGTH {
		e3.SetLength(math.Max(0, e1.Length()) + math.Max(0, e2.Length()))
	}
	if !n1.Tip() && !n2.Tip() && (e1.Support() != NIL_SUPPORT || e2.Support() != NIL_SUPPORT) {
		e3.SetSupport(math.Max(math.Max(0, e1.Support()), math.Max(0, e2.Support())))
	}
	t.delNode(root)

	t.ReinitIndexes()
}

// Annotates internal branches of a tree with given data using the
// given list of names:
//	* each element is a slice of names whose first element is the new name to set
//      * following names are internal node name or tip names defining an internal node (LCA)
//
// For a given element different possibilities:
// 1) If < 2 names: an error is returned
// 1) If 2 names [n1,n2]: we search for n2 in the tree (tip or internal node)
//    and rename it as n1
// 2) If > 2 names [n1,n2,...,ni]: We find the LCA of every tips whose name is in [n2,...,ni]
//    and rename it as n1
//
// If comment is true, then we do not change the name, but the comment of the given nodes.
//
// The output tree won't have bootstrap support at the given branches anymore.
//
// It considers the tree as rooted (even if multifurcation at root).
func (t *Tree) Annotate(names [][]string, comment bool) error {
	nodeindex, err := NewNodeIndex(t)
	if err != nil {
		return err
	}

	for _, line := range names {
		// Only one node (internal or tip): Then sets the name to this specific node is found on the tree
		if len(line) < 2 {
			return errors.New("Error in tree annotation: Wrongly formatted annotation slice")
		} else if len(line) == 2 {
			if node, found := nodeindex.GetNode(line[1]); found {
				if comment {
					node.AddComment(line[0])
				} else {
					node.SetName(line[0])
				}
			}
		} else {
			n, _, _, err := t.LeastCommonAncestorRooted(nodeindex, line[1:]...)
			if err != nil {
				return err
			}
			if comment {
				n.AddComment(line[0])
			} else {
				n.SetName(line[0])
			}
		}
	}
	return nil
}

// This function renames nodes of the tree based on the map in argument
// If a name in the map does not exist in the tree, then returns an error
// If a node/tip in the tree does not have a name in the map: OK
// After rename, tip index is updated, as well as bitsets of the edges
func (t *Tree) Rename(namemap map[string]string) error {
	nodeindex, err := NewNodeIndex(t)
	if err != nil {
		return err
	}
	for name, newname := range namemap {
		node, ok := nodeindex.GetNode(name)
		if ok {
			node.SetName(newname)
		}
	}
	// After we update bitsets if any, and node indexes
	//t.ReinitIndexes()

	if err = t.UpdateTipIndex(); err != nil {
		return err
	}
	// err = t.ClearBitSets()
	// if err != nil {
	// 	return err
	// }
	// t.UpdateBitSet()
	return nil
}

// Renames automatically nodes/tips of the Tree, with names of the form T000001, N00001, etc.
//
// internals: Renames internal nodes
// tips: Renames tips
// length: length of the generated names
// curid: index that is incremented, to keep track of the number of generated names
// namemap: map that associates old names to new names - allows to keep track of renames for previous trees
func (t *Tree) RenameAuto(internals, tips bool, length int, curid *int, namemap map[string]string) error {
	for i, n := range t.Nodes() {
		if (tips && n.Tip()) || (internals && !n.Tip()) {
			prefix := 'T'
			if !n.Tip() {
				prefix = 'N'
				if n.Name() == "" {
					n.SetName(fmt.Sprintf("%d", i))
				}
			}
			newname, ok := namemap[n.Name()]
			if !ok {
				newname = fmt.Sprintf(fmt.Sprintf("%c%%0%dd", prefix, (length-1)), *curid)
				if len(newname) != length {
					return (fmt.Errorf("Id length %d does not allow to generate as much ids: %d (%s)",
						length, curid, newname))
				}
				namemap[n.Name()] = newname
				*curid++
			}
			n.SetName(newname)
		}
	}
	// After we update bitsets if any, and node indexes
	//t.ReinitIndexes()
	if err := t.UpdateTipIndex(); err != nil {
		return err
	}
	// err := t.ClearBitSets()
	// if err != nil {
	// 	return err
	// }
	// t.UpdateBitSet()
	return nil
}

// Renames nodes/tips of the Tree, given the regexp and a replaement string
//
// internals: Renames internal nodes
// tips: Renames tips
// regexp: the regexp to match node names
// replace: Replacement string
// namemap: map that associates old names to new names
func (t *Tree) RenameRegexp(internals, tips bool, regex, replace string, namemap map[string]string) error {
	r, err := regexp.Compile(regex)
	if err != nil {
		return err
	}
	for _, n := range t.Nodes() {
		if (tips && n.Tip()) || (internals && !n.Tip()) {
			newname := r.ReplaceAllString(n.Name(), replace)
			namemap[n.Name()] = newname
			n.SetName(newname)
		}
	}
	// After we update bitsets if any, and node indexes
	//t.ReinitIndexes()
	if err = t.UpdateTipIndex(); err != nil {
		return err
	}
	// err = t.ClearBitSets()
	// if err != nil {
	// 	return err
	// }
	// t.UpdateBitSet()
	return nil
}

// Clone the given node, copy attributes of the given
// node into a new node
func (t *Tree) CopyNode(n *Node) *Node {
	out := t.NewNode()
	out.name = n.name
	out.depth = n.depth
	out.id = n.id
	out.comment = make([]string, len(n.comment))
	for i, c := range n.comment {
		out.comment[i] = c
	}
	return out
}

func (t *Tree) Delete() {
	t.tipIndex = nil
	for _, n := range t.Nodes() {
		t.delNode(n)
	}
}

// Copy attributes of the given edge to the other given edge:
//	* Length
//	* Support
//	* id
//	* bitset (if not nil)
func (t *Tree) CopyEdge(e *Edge, copy *Edge) {
	copy.length = e.length
	copy.support = e.support
	copy.pvalue = e.pvalue
	copy.id = e.id
	if e.bitset != nil {
		copy.bitset = e.bitset.Clone()
	}
	copy.ntaxleft = e.ntaxleft
	copy.ntaxright = e.ntaxright
	copy.hashcodeleft = e.hashcodeleft
	copy.hashcoderight = e.hashcoderight
}

// Clone the input tree
func (t *Tree) Clone() *Tree {
	copy := NewTree()
	root := t.CopyNode(t.Root())
	copy.SetRoot(root)
	for _, e := range t.Root().br {
		t.copyTreeRecur(copy, root, t.Root(), e)
	}
	if t.tipIndex != nil {
		copy.UpdateTipIndex()
	}
	return (copy)
}

// Recursive function to clone the tree
func (t *Tree) copyTreeRecur(copytree *Tree, copynode, node *Node, edge *Edge) {
	child := edge.Right()
	copychild := t.CopyNode(child)
	copyedge := copytree.ConnectNodes(copynode, copychild)
	t.CopyEdge(edge, copyedge)
	for _, e := range child.br {
		if e != edge {
			t.copyTreeRecur(copytree, copychild, child, e)
		}
	}
}

// Assumes that the tree is rooted.
//
// Otherwise, will consider the pseudo root
// defined by the initial newick file
func (t *Tree) SubTree(n *Node) *Tree {
	subtree := NewTree()
	root := t.CopyNode(n)
	subtree.SetRoot(root)
	for _, e := range n.br {
		if e.Left() == n {
			t.copyTreeRecur(subtree, root, n, e)
		}
	}
	subtree.ReinitIndexes()
	return (subtree)
}

// Merges Two rooted trees t and t2 in t by adding a new root node with two children
// Corresponding to the roots of the 2 trees.
//
// If one of the tree is not rooted, returns an error.
//
// Tip set must be different between the two trees, otherwise returns an error.
//
// it is advised not to use t2 after the merge, since it may conflict with t.
//
// Edges connecting the new root with old roots have length of 1.0, but can be modified
// afterwards.
func (t *Tree) Merge(t2 *Tree) error {
	if !t.Rooted() || !t2.Rooted() {
		return errors.New("One of the two tree (or both) is not rooted")
	}

	//Comparing tip names
	if len(t.tipIndex) == 0 || len(t2.tipIndex) == 0 {
		return errors.New("No tips in the index, tip name index is not initialized")
	}
	for k := range t.tipIndex {
		_, ok := t2.tipIndex[k]
		if ok {
			return errors.New("Trees should not have common tip names")
		}
	}

	// Now we add a new root
	newroot := t.NewNode()
	t.ConnectNodes(newroot, t.Root())
	t.ConnectNodes(newroot, t2.Root())
	t.SetRoot(newroot)

	if err := t.ReinitIndexes(); err != nil {
		return err
	}

	return nil
}

// Returns the deepest edge of the tree (considered unrooted)
// in terms of number of tips on the light side of it.
//
// It does not use bitsets, thus they may be uninitialized.
func (t *Tree) DeepestEdge() (maxedge *Edge) {
	// We choose the deepest edge
	for i, e := range t.Edges() {
		e.SetId(i)
	}
	numtips := len(t.Tips())
	maxedge, _, _, _, _ = t.deepestEdgeRecur(t.Root(), nil, nil, numtips)
	return
}

func (t *Tree) deepestEdgeRecur(node, prev *Node, edge *Edge, numtips int) (maxedge *Edge, lefttips, righttips, maxdepth, curtips int) {
	if node.Tip() {
		return edge, 1, curtips - 1, 1, 1
	}
	curtips = 0
	maxdepth = 0
	lefttips = 0  // Nb tips on the left of current max depth branch
	righttips = 0 // Nb tips on the right of current max depth branch
	for i, c := range node.Neigh() {
		if c != prev {
			e, l, r, d, t := t.deepestEdgeRecur(c, node, node.Edges()[i], numtips)
			if d > maxdepth {
				maxdepth = d
				maxedge = e
				lefttips = l
				righttips = r
			}
			curtips += t
		}
	}
	if mutils.Min(numtips-curtips, curtips) > maxdepth {
		maxdepth = mutils.Min(numtips-curtips, curtips)
		maxedge = edge
		lefttips = numtips - curtips
		righttips = curtips
	}
	return
}

// Returns the deepest node of the tree (considered unrooted).
//
// We define the deepest node as the heavy side of the deepest
// edge (See tree.DeepestEdge())
//
// It does not use bitsets, thus they may be uninitialized.
func (t *Tree) DeepestNode() (maxnode *Node) {
	// We choose the deepest edge
	for i, e := range t.Edges() {
		e.SetId(i)
	}
	var maxedge *Edge
	var lefttips, righttips int
	numtips := len(t.Tips())
	maxedge, lefttips, righttips, _, _ = t.deepestEdgeRecur(t.Root(), nil, nil, numtips)
	if lefttips > righttips {
		maxnode = maxedge.Left()
	} else {
		maxnode = maxedge.Right()
	}
	return
}

// Scale branch supports by a given factor. Precision 10^-6.
//
// Does not do anything for branches with a support of NIL_SUPPORT
func (t *Tree) ScaleSupports(factor float64) {
	for _, e := range t.Edges() {
		if s := e.Support(); s != NIL_SUPPORT {
			e.SetSupport(float64(int(1000000*(s*factor))) / 1000000)
		}
	}
}

// Rounds branch supports by a given precision. Precision is defined as 1/(power of 10)
// Example: 6, means 10^-6
//
// Does not do anything for branches with a support of NIL_SUPPORT or if precision is <=0
// If precision > 15, then 15 is taken
func (t *Tree) RoundSupports(precision int) {
	if precision > 15 {
		precision = 15
	}
	for _, e := range t.Edges() {
		if s := e.Support(); s != NIL_SUPPORT {
			e.SetSupport(math.Round(math.Pow10(precision)*s) / math.Pow10(precision))
		}
	}
}

// Scale branch lengths by a given factor.
//
// Does not do anything for branches with a length of NIL_LENGTH
func (t *Tree) ScaleLengths(factor float64) {
	for _, e := range t.Edges() {
		if s := e.Length(); s != NIL_LENGTH {
			e.SetLength(e.Length() * factor)
		}
	}
}

// Rounds branch lengths by a given precision. Precision is defined as 1/(power of 10)
// Example: 6, means 10^-6.
//
// Does not do anything for branches with a length of NIL_LENGTH or if precision is <=0.
// If precision > 15, then 15 is taken
func (t *Tree) RoundLengths(precision int) {
	if precision > 15 {
		precision = 15
	}
	for _, e := range t.Edges() {
		if s := e.Length(); s != NIL_LENGTH {
			e.SetLength(math.Round(math.Pow10(precision)*s) / math.Pow10(precision))
		}
	}
}

// Get all connected components (only the tips) of the tree
// when edges with length less than maxLen (included) are removed
// If a connected components have no tips, then it is not taken into account.
//
// The edges are not actually removed from the input tree.
func (t *Tree) CutEdgesMaxLength(maxlen float64) (bags []*TipBag, err error) {
	edges := t.Edges()
	visited := make([]bool, len(edges))
	bags = make([]*TipBag, 0, 10)

	// Associate an index to each edges
	for i, e := range edges {
		e.SetId(i)
		visited[i] = false
	}

	// Visit all the edges if they are not already visited
	for _, e := range edges {
		if !visited[e.Id()] {
			visited[e.Id()] = true
			if e.Length() < maxlen {
				tipBag := NewTipBag()
				// We go both directions of the edge
				if err = t.cutEdgesMaxLengthRecur(tipBag, e.Left(), e.Right(), maxlen, visited); err != nil {
					return
				}
				if err = t.cutEdgesMaxLengthRecur(tipBag, e.Right(), e.Left(), maxlen, visited); err != nil {
					return
				}
				if tipBag.Size() > 0 {
					bags = append(bags, tipBag)
				}
			} else {
				// We consider that edge as removed
				// but if it connects tips, we keep the tips in a new bag
				if e.Left().Tip() {
					tipBag := NewTipBag()
					tipBag.AddTip(e.Left())
					bags = append(bags, tipBag)
				}
				if e.Right().Tip() {
					tipBag := NewTipBag()
					tipBag.AddTip(e.Right())
					bags = append(bags, tipBag)
				}
			}
		}
	}
	return
}

// recursive function that lists all tips of the current connected
// component defined by edges with length < maxlen
func (t *Tree) cutEdgesMaxLengthRecur(tipBag *TipBag, cur *Node, prev *Node, maxlen float64, visited []bool) error {
	if cur == nil {
		return errors.New("Nil node given to Tree.cutEdgesMaxLengthRecur")
	}
	if cur.Tip() {
		if err := tipBag.AddTip(cur); err != nil {
			return err
		}
	}
	for i, n := range cur.neigh {
		b := cur.br[i]
		if n != prev && b.Length() < maxlen {
			visited[b.Id()] = true
			if err := t.cutEdgesMaxLengthRecur(tipBag, n, cur, maxlen, visited); err != nil {
				return err
			}
		}
	}
	return nil
}

// This function adds Tips to the tree, adding 0 length branches.
// To do so, it takes all identical tipnames of the given slice
// And add the new tip names next to the existing ones, by adding
// 0 length branches.
// Each identical group must contain exactly 1 already present tip
// otherwise returns an error
// If a new tip is identical to several already present tips, then
// returns and error.
func (t *Tree) InsertIdenticalTips(identicalgroups [][]string) (err error) {
	var e bool
	var nodeindex NodeIndex
	var old, newtip *Node
	var name string

	if nodeindex, err = NewNodeIndex(t); err != nil {
		return
	}
	for _, group := range identicalgroups {
		// Search for the already present tip
		newtips := make([]string, 0)
		oldtip := ""
		for _, name = range group {
			if e, err = t.ExistsTip(name); err != nil {
				return
			}
			if e && oldtip == "" {
				oldtip = name
			} else if e && oldtip != "" {
				err = fmt.Errorf("Several already existing tips are present in an identical group")
				return
			} else {
				newtips = append(newtips, name)
			}
		}
		if oldtip == "" {
			err = fmt.Errorf("No existing tip is present in the given identical group: %s", name)
			return
		}
		// Now we add the new tips
		if len(newtips) > 0 {
			if old, e = nodeindex.GetNode(oldtip); e {
				for _, tip := range newtips {
					if newtip, err = t.InsertIdenticalTip(old, tip); err != nil {
						return
					}
					nodeindex.AddNode(newtip)
				}
			}
		}
	}
	if err = t.ReinitIndexes(); err != nil {
		return err
	}
	return
}

// This function adds a new tip next to the given tip node.
// This adds two 0 length branches. The given node reference is still a tip after this function.
// The new node is the internal node.
//
// Before:
//   l
// *----*tip1
// After:
//       *tip1
//  l   /.0
// ----*
//      \.0
//       *newTipName
//
// If l==0.0 then, after:
//   *tip1
//  /.0
// *
//  \.0
//   *newTipName
//
// Warning: This function may add polytomies if l==0.0.
//
// It updates the tipindex temporarily but if needed in downstream analysis t.ReinitIndexes()
// must be called.
// If initialized, Bitsets and Depths are not up2date after this function. They should be
// updated with t.ReinitIndexes() if needed.
func (t *Tree) InsertIdenticalTip(n *Node, newTipName string) (newtipnode *Node, err error) {
	var e bool

	if !n.Tip() {
		err = fmt.Errorf("The node to add the new tip to is not a tip")
		return
	}
	if e, err = t.ExistsTip(newTipName); err != nil {
		return
	}
	if e {
		err = fmt.Errorf("The tip to add to %s is already present in the tree", n.Name())
		return
	}

	var parentedge *Edge
	var parentnode *Node
	if parentedge, err = n.ParentEdge(); err != nil {
		return
	}
	if parentnode, err = n.Parent(); err != nil {
		return
	}
	if parentedge.Length() == 0.0 {
		newtipnode = t.NewNode()
		newtipnode.SetName(newTipName)
		e1 := t.ConnectNodes(parentnode, newtipnode)
		e1.SetLength(0.0)
	} else {
		var e_ind, n_ind int
		// index of parent edge in neighbors of parent node
		if e_ind, err = parentnode.EdgeIndex(parentedge); err != nil {
			return
		}
		// index of parent node in neighbors of node
		if n_ind, err = n.NodeIndex(parentnode); err != nil {
			return
		}

		newtipnode = t.NewNode()
		newtipnode.SetName(newTipName)
		newinternal := t.NewNode()
		newedge := t.ConnectNodes(newinternal, newtipnode)
		newedge.SetLength(0.0)

		newedge1 := t.NewEdge()
		newedge1.SetLength(0.0)
		parentedge.setRight(newinternal)
		newinternal.addChild(parentnode, parentedge)
		newinternal.addChild(n, newedge1)
		parentnode.neigh[e_ind] = newinternal

		newedge1.setRight(n)
		newedge1.setLeft(newinternal)
		n.neigh[n_ind] = newinternal
		n.br[n_ind] = newedge1
	}
	t.tipIndex[newTipName] = uint(len(t.tipIndex))
	return
}

// Check that tree is well structured, with the
// right child, parent, edges, and node pointers
func (t *Tree) CheckTree() bool {
	return t.checkTreeRecur(nil, nil)
}

// Check that tree is well structured, with the
// right child, parent, edges, and node pointers
func (t *Tree) checkTreeRecur(cur *Node, prev *Node) (ok bool) {
	ok = true
	if cur == nil {
		cur = t.Root()
	}

	for i, n := range cur.neigh {
		if n != prev {
			if cur.br[i].left != cur {
				return false
			}
			if cur.br[i].right != n {
				return false
			}
			if ok = ok && t.checkTreeRecur(n, cur); !ok {
				return
			}
		} else if n == prev {
			if prev == nil {
				return false
			}
			if cur.br[i].right != cur {
				return false
			}
			if cur.br[i].left != prev {
				return false
			}
		}
	}
	return
}
