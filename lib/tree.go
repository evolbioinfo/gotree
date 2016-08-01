/*
   Package gotree implements a simple
   library for handling phylogenetic trees in go
*/
package lib

import (
	"errors"
	"fmt"
	"github.com/willf/bitset"
	"math/rand"
	"os"
	"sort"
	"strconv"
)

type Tree struct {
	nodes    []*Node         // array of all the tree nodes
	edges    []*Edge         // array of all the tree edges
	root     *Node           // root node
	tipIndex map[string]uint // Map between tip name and bitset index
}

type Node struct {
	name    string   // Name of the node
	comment []string // Comment if any in the newick file
	id      int      // Id of the node: attributed when parsing
	neigh   []*Node  // neighbors array
	br      []*Edge  // Branches array (same order than neigh)
	depth   int      // Depth of the node
}

type Edge struct {
	id          int     // id of the branch: attribute when parsing
	left, right *Node   // Left and right nodes
	length      float64 // length of branch
	support     float64 // -1 if no support
	// a Bit at index i in the bitset corresponds to the position of the tip i
	//left:0/right:1 .
	// i is the index of the tip in the sorted tip name array
	bitset *bitset.BitSet // Bitset of length Number of taxa each
}

func NewNode() *Node {
	return &Node{
		name:    "",
		id:      0,
		comment: make([]string, 0),
		neigh:   make([]*Node, 0, 3),
		br:      make([]*Edge, 0, 3),
		depth:   0,
	}
}

func NewEdge() *Edge {
	return &Edge{
		id:      0,
		length:  -1.0,
		support: -1.0,
	}
}

func NewTree() *Tree {
	return &Tree{
		nodes:    make([]*Node, 0, 10),
		edges:    make([]*Edge, 0, 10),
		root:     nil,
		tipIndex: make(map[string]uint, 0),
	}
}

/* Node functions */
/******************/

func (p *Node) AddChild(n *Node, e *Edge) {
	p.neigh = append(p.neigh, n)
	p.br = append(p.br, e)

}

func (n *Node) SetName(name string) {
	n.name = name
}

func (n *Node) AddComment(comment string) {
	n.comment = append(n.comment, comment)
}

func (n *Node) SetId(id int) {
	n.id = id
}

func (n *Node) SetDepth(depth int) {
	n.depth = depth
}

func (n *Node) Name() string {
	return n.name
}

// Retrieve the parent node
// If several parents: Error
// Parent is defined as the node n2 connected to n
// by an edge e with e.left == n2 and e.right == n
func (n *Node) Parent() (*Node, error) {
	var n2 *Node
	for _, e := range n.br {
		if e.right == n {
			if n2 != nil {
				return nil, errors.New("The node has more than one parent")
			}
			n2 = e.left
		}
	}
	if n2 == nil {
		return nil, errors.New("The node has no parent : May be the root?")
	}
	return n2, nil
}

// Retrieve the Edge going to Parent node
// If several parents: Error
// Parent is defined as the node n2 connected to n
// by an edge e with e.left == n2 and e.right == n
func (n *Node) ParentEdge() (*Edge, error) {
	var e2 *Edge
	for _, e := range n.br {
		if e.right == n {
			if e2 != nil {
				return nil, errors.New("The node has more than one parent")
			}
			e2 = e
		}
	}
	if e2 == nil {
		return nil, errors.New("The node has no parent : May be the root?")
	}
	return e2, nil
}

/* Edge functions */
/******************/

func (e *Edge) SetId(id int) {
	e.id = id
}
func (e *Edge) SetLeft(left *Node) {
	e.left = left
}
func (e *Edge) SetRight(right *Node) {
	e.right = right
}
func (e *Edge) SetLength(length float64) {
	e.length = length
}

func (e *Edge) SetSupport(support float64) {
	e.support = support
}

func (e *Edge) Length() float64 {
	return e.length
}

func (e *Edge) DumpBitSet() string {
	if e.bitset == nil {
		return "nil"
	}
	return e.bitset.DumpAsBits()
}

/* Tree functions */
/******************/

func (t *Tree) SetRoot(r *Node) {
	t.root = r
}

func (t *Tree) Root() *Node {
	return t.root
}

func (t *Tree) Edges() []*Edge {
	return t.edges
}

func (t *Tree) String() string {
	return t.Newick()
}

func (t *Tree) Newick() string {
	return t.root.Newick(nil) + ";"
}

func (t *Tree) UpdateTipIndex() {
	names := t.AllTipNames(nil, nil)
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

// Returns all the tip name in the tree
// Starts with n==nil (root)
func (t *Tree) AllTipNames(n *Node, parent *Node) []string {
	names := make([]string, 0)
	if n == nil {
		n = t.Root()
	}
	// is a tip
	if len(n.neigh) == 1 {
		names = append(names, n.name)

	} else {
		for _, child := range n.neigh {
			if child != parent {
				for _, name := range t.AllTipNames(child, n) {
					names = append(names, name)
				}
			}
		}
	}
	return names
}

func (n *Node) EdgeIndex(e *Edge) (int, error) {
	for i := 0; i < len(n.br); i++ {
		if n.br[i] == e {
			return i, nil
		}
	}
	return -1, errors.New("The Edge is not in the neighbors of node")
}

func (n *Node) NodeIndex(next *Node) (int, error) {
	for i := 0; i < len(n.neigh); i++ {
		if n.neigh[i] == next {
			return i, nil
		}
	}
	return -1, errors.New("The Node is not in the neighbors of node")
}

func (t *Tree) AddNewNode() *Node {
	newnode := NewNode()
	newnode.id = len(t.nodes)
	t.nodes = append(t.nodes, newnode)
	return newnode
}

func (t *Tree) AddNewEdge() *Edge {
	newedge := NewEdge()
	newedge.id = len(t.edges)
	t.edges = append(t.edges, newedge)
	return newedge
}

func (t *Tree) ConnectNodes(parent *Node, child *Node) *Edge {
	newedge := t.AddNewEdge()
	newedge.SetLeft(parent)
	newedge.SetRight(child)
	parent.AddChild(child, newedge)
	child.AddChild(parent, newedge)
	return newedge
}

// This function takes the first node having 3 neighbors
// and reroot the tree on this node
func (t *Tree) RerootFirst() error {
	for _, n := range t.nodes {
		if len(n.neigh) == 3 {
			err := t.Reroot(n)
			return err
		}
	}
	return errors.New("No nodes with 3 neighors have been found for rerooting")
}

// Recursively update bitsets of edges from the Node n
// If node == nil then it starts from the root
func (t *Tree) ClearBitSetsRecur(n *Node, parent *Node, ntip uint) {
	if n == nil {
		n = t.Root()
	}

	for i, child := range n.neigh {
		e := n.br[i]
		e.bitset.ClearAll()
		e.bitset = bitset.New(ntip)
		if child != parent {
			t.ClearBitSetsRecur(child, n, ntip)
		}
	}
}

// Updates bitsets of all edges in the tree
// Assumes that the hashmap tip name : index is
// initialized with UpdateTipIndex function
func (t *Tree) UpdateBitSet() error {
	for _, e := range t.Root().br {
		err := t.FillRightBitSet(e, make([]*Edge, 0))
		if err != nil {
			return err
		}
	}
	return nil
}

// This function compares 2 trees and output
// the number of edges in common
// If the trees have different sets of tip names, returns an error
// It assumes that functions
// 	tree.UpdateTipIndex()
//	tree.clearBitSetsRecur(nil, nil, uint(len(tree.tipIndex)))
//	tree.UpdateBitSet()
// Have been called before, otherwise will output an error
func (t *Tree) CommonEdges(t2 *Tree) (int, error) {
	common := 0

	err := t.CompareTipIndexes(t2)

	if err != nil {
		return 0, err
	}

	for _, e := range t.edges {
		for _, e2 := range t2.edges {
			if e.bitset == nil || e2.bitset == nil {
				return 0, errors.New("BitSets has not been initialized with tree.clearBitSetsRecur(nil, nil, uint(len(tree.tipIndex)))")
			}
			if !e.bitset.Any() || !e2.bitset.Any() {
				return 0, errors.New("One edge has a bitset of 0...000 : May be BitSets have not been updated with tree.UpdateBitSet()?")
			}
			if e.bitset.Equal(e2.bitset) ||
				e.bitset.Complement().Equal(e2.bitset) {
				common++
				break
			}
		}
	}
	return common, nil
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

// Recursively clears and sets the bitsets of the descending edges
//
func (t *Tree) FillRightBitSet(currentEdge *Edge, rightEdges []*Edge) error {
	if currentEdge.bitset == nil {
		return errors.New("BitSets has not been initialized with tree.clearBitSetsRecur(nil, nil, uint(len(tree.tipIndex)))")
	}
	currentEdge.bitset = currentEdge.bitset.ClearAll()
	rightEdges = append(rightEdges, currentEdge)
	// If we are at a tip edge
	// We set at 1 the bits of the tip in
	// the bitsets of all rightEdges
	if len(currentEdge.right.neigh) == 1 {
		i, err := t.tipIndexNode(currentEdge.right)
		if err != nil {
			return err
		}
		for _, e := range rightEdges {
			e.bitset = e.bitset.Set(i)
		}
	} else {
		// Else
		for _, e2 := range currentEdge.right.br {
			if e2.left == currentEdge.right {
				t.FillRightBitSet(e2, rightEdges)
			}
		}
	}
	return nil
}

// This function takes a node and reroot the tree on that node
// It reorients edges left-edge-right : see ReorderEdges
// The node must be one of the tree nodes, otherwise it returns an error
func (t *Tree) Reroot(n *Node) error {
	intree := false
	for _, n2 := range t.nodes {
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
func (t *Tree) GraftTipOnEdge(n *Node, e *Edge) error {
	newnode := t.AddNewNode()
	newedge := t.AddNewEdge()
	lnode := e.left
	rnode := e.right

	// index of edge in neighbors of l
	e_l_ind, err := lnode.EdgeIndex(e)
	if err != nil {
		return err
	}
	// index of edge in neighbors of r
	e_r_ind, err2 := rnode.EdgeIndex(e)
	if err2 != nil {
		return err2
	}

	newedge.SetLength(1.0)
	newedge.SetLeft(newnode)
	newedge.SetRight(n)
	newnode.AddChild(n, newedge)
	n.AddChild(newnode, newedge)
	e.SetRight(newnode)
	newnode.AddChild(lnode, e)
	lnode.neigh[e_l_ind] = newnode

	if lnode.br[e_l_ind] != e {
		return errors.New("The Edge is not at the same index")
	}

	newedge2 := t.AddNewEdge()
	newedge2.SetLength(e.length / 2)
	e.SetLength(e.length / 2)
	newedge2.SetLeft(newnode)
	newedge2.SetRight(rnode)
	newnode.AddChild(rnode, newedge2)
	if rnode.br[e_r_ind] != e {
		return errors.New("The Edge is not at the same index")
	}
	rnode.neigh[e_r_ind] = newnode
	rnode.br[e_r_ind] = newedge2

	return nil
}

//Creates a Random Binary tree
//nbtips : Number of tips of the random binary tree to create
func RandomBinaryTree(nbtips int) (*Tree, error) {
	t := NewTree()
	if nbtips < 2 {
		return nil, errors.New("Cannot create a random binary tree with less than 2 tips")
	}
	for i := 1; i < nbtips; i++ {
		n := t.AddNewNode()
		n.SetName("Tip" + strconv.Itoa(i))
		switch len(t.edges) {
		case 0:
			n2 := t.AddNewNode()
			n2.SetName("Node" + strconv.Itoa(i-1))
			e := t.ConnectNodes(n2, n)
			e.SetLength(1.0)
			t.SetRoot(n2)
		default:
			// Where to insert the new tip
			i_edge := rand.Intn(len(t.edges))
			e := t.edges[i_edge]
			err := t.GraftTipOnEdge(n, e)
			if err != nil {
				fmt.Fprintf(os.Stderr, err.Error()+"\n")
			}
		}
	}
	err := t.RerootFirst()
	t.UpdateTipIndex()
	t.ClearBitSetsRecur(nil, nil, uint(len(t.tipIndex)))
	t.UpdateBitSet()
	return t, err
}

// Recursive function that outputs newick representation
// from the current node
func (n *Node) Newick(parent *Node) string {
	newick := ""
	if len(n.neigh) > 0 {
		if len(n.neigh) > 1 {
			newick += "("
		}
		nbchild := 0
		for i, child := range n.neigh {
			if child != parent {
				if nbchild > 0 {
					newick += ","
				}
				newick += child.Newick(n)
				if n.br[i].support != -1 {
					newick += strconv.FormatFloat(n.br[i].support, 'f', 5, 64)
				}
				if len(child.comment) != 0 {
					for _, c := range child.comment {
						newick += "[" + c + "]"
					}
				}
				if n.br[i].length != -1 {
					newick += ":" + strconv.FormatFloat(n.br[i].length, 'f', 5, 64)
				}
				nbchild++
			}
		}
		if len(n.neigh) > 1 {
			newick += ")"
		}
	}
	newick += n.name

	return newick
}
