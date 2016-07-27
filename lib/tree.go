/*
   Package gotree implements a simple
   library for handling phylogenetic trees in go
*/
package lib

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strconv"
)

type Tree struct {
	nodes []*Node // array of all the tree nodes
	edges []*Edge // array of all the tree edges
	root  *Node   // root node
}

type Node struct {
	name    string  // Name of the node
	comment string  // Comment if any in the newick file
	id      int     // Id of the node: attributed when parsing
	neigh   []*Node // neighbors array
	br      []*Edge // Branches array (same order than neigh)
	depth   int     // Depth of the node
}

type Edge struct {
	id          int     // id of the branch: attribute when parsing
	left, right *Node   // Left and right nodes
	length      float64 // length of branch
	support     float64 // -1 if no support
}

func NewNode() *Node {
	return &Node{
		name:    "",
		id:      0,
		comment: "",
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
		nodes: make([]*Node, 0, 10),
		edges: make([]*Edge, 0, 10),
		root:  nil,
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

func (n *Node) SetComment(comment string) {
	n.comment = comment
}

func (n *Node) SetId(id int) {
	n.id = id
}

func (n *Node) SetDepth(depth int) {
	n.depth = depth
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

/* Tree functions */
/******************/

func (t *Tree) SetRoot(r *Node) {
	t.root = r
}

func (t *Tree) Root() *Node {
	return t.root
}
func (t *Tree) String() string {
	return t.Newick()
}

func (t *Tree) Newick() string {
	return t.root.Newick(nil) + ";"
}

func (n *Node) NeighborEdgeIndex(e *Edge) (int, error) {
	for i := 0; i < len(n.br); i++ {
		if n.br[i] == e {
			return i, nil
		}
	}
	return -1, errors.New("The Edge is not in the neighbors of node")
}

// This function take the first node having 3 neighbors
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

// This function take a node and reroot the tree on that node
// The node must be one of the tree nodes, otherwize it returns an error
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

	err := t.ReorderEdges(n, nil)

	return err
}

// This function reorders the edges of a tree
// in order to always have left-edge-right
// with left node being parent of right node
// with respect to the given root node
// Important even for unrooted trees
func (t *Tree) ReorderEdges(n *Node, prev *Node) error {
	for _, next := range n.br {
		if next.right != prev && next.left != prev {
			if next.right == n {
				next.right, next.left = next.left, next.right
			}
			t.ReorderEdges(next.right, n)
		}
	}
	return nil
}

// This function graft the Node n at the middle of the Edge e
// It divides the branch lenght by 2
func (t *Tree) GraftTipOnEdge(n *Node, e *Edge) error {
	newnode := NewNode()
	newedge := NewEdge()
	lnode := e.left
	rnode := e.right

	// index of edge in neighbors of l
	e_l_ind, err := lnode.NeighborEdgeIndex(e)
	if err != nil {
		return err
	}
	// index of edge in neighbors of r
	e_r_ind, err2 := rnode.NeighborEdgeIndex(e)
	if err2 != nil {
		return err2
	}

	newedge.id = len(t.edges)
	t.edges = append(t.edges, newedge)
	newedge.SetLength(1.0)
	newnode.id = len(t.nodes)
	t.nodes = append(t.nodes, newnode)

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

	newedge2 := NewEdge()
	newedge2.id = len(t.edges)
	t.edges = append(t.edges, newedge2)
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
		n := NewNode()
		n.SetName("Tip" + strconv.Itoa(i))
		n.id = len(t.nodes)
		t.nodes = append(t.nodes, n)
		switch len(t.edges) {
		case 0:
			n2 := NewNode()
			e := NewEdge()
			e.SetLength(1.0)
			n2.SetName("Node" + strconv.Itoa(i-1))
			t.SetRoot(n2)
			n2.AddChild(n, e)
			n.AddChild(n2, e)
			e.id = len(t.edges)
			t.edges = append(t.edges, e)
			n2.id = len(t.nodes)
			t.nodes = append(t.nodes, n2)
			e.left = n2
			e.right = n
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
		for i, child := range n.neigh {
			if child != parent {
				if i > 0 {
					newick += ","
				}
				newick += child.Newick(n)
				if n.br[i].support != -1 {
					newick += strconv.FormatFloat(n.br[i].support, 'f', 5, 64)
				}
				if n.br[i].length != -1 {
					newick += ":" + strconv.FormatFloat(n.br[i].length, 'f', 5, 64)
				}
			}
		}
		if len(n.neigh) > 1 {
			newick += ")"
		}
	}
	newick += n.name

	return newick
}
