package tree

import (
	"bytes"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
)

// Node structure
type Node struct {
	name    string   // Name of the node
	comment []string // Comment if any in the newick file
	neigh   []*Node  // neighbors array
	br      []*Edge  // Branches array (same order than neigh)
	depth   int      // Depth of the node
	id      int      // this field is used at discretion of the user to store information
}

// Uninitialized depth is coded as -1
const (
	NIL_DEPTH = -1
)

// Adds a child n to the node p, connected with edge e
func (p *Node) addChild(n *Node, e *Edge) {
	p.neigh = append(p.neigh, n)
	p.br = append(p.br, e)
}

// Sets the name of the node. No verification if another node
// has the same name
func (n *Node) SetName(name string) {
	n.name = name
}

// Adds a comment to the node. It will be coded by a list of []
// In the Newick format.
func (n *Node) AddComment(comment string) {
	n.comment = append(n.comment, comment)
}

// Returns the list of comments associated to the node.
func (n *Node) Comments() []string {
	return n.comment
}

// Returns the string of comma separated comments
// surounded by [].
func (n *Node) CommentsString() string {
	var buf bytes.Buffer
	buf.WriteRune('[')
	for i, c := range n.comment {
		if i > 0 {
			buf.WriteRune(',')
		}
		buf.WriteString(c)
	}
	buf.WriteRune(']')
	return buf.String()
}

// Removes all comments associated to the node
func (n *Node) ClearComments() {
	n.comment = n.comment[:0]
}

// Sets the depth of the node
func (n *Node) SetDepth(depth int) {
	n.depth = depth
}

// Returns the name of the node
func (n *Node) Name() string {
	return n.name
}

// Returns the Id of the node. Id==NIL_ID means that
// it has not been set yet.
func (n *Node) Id() int {
	return n.id
}

// Sets the id of the node
func (n *Node) SetId(id int) {
	n.id = id
}

// Returns the depth of the node. Returns an error if the depth has
// not been set yet.
func (n *Node) Depth() (int, error) {
	if n.depth == NIL_DEPTH {
		return n.depth, errors.New("Node depth has not been computed")
	}
	return n.depth, nil
}

// Number of neighbors of this node
func (n *Node) Nneigh() int {
	return len(n.neigh)
}

// List of neighbors of this node
func (n *Node) Neigh() []*Node {
	return n.neigh
}

// Is a tip or not?
func (n *Node) Tip() bool {
	return len(n.neigh) == 1
}

// List of edges going from this node
func (n *Node) Edges() []*Edge {
	return n.br
}

// deletes the given neighbor n2 of this node n
//
// If n2 is not a neighbor of n, then returns an error.
func (n *Node) delNeighbor(n2 *Node) (err error) {
	var i int
	if i, err = n.NodeIndex(n2); err != nil {
		return
	}
	n.br = append(n.br[0:i], n.br[i+1:]...)
	n.neigh = append(n.neigh[0:i], n.neigh[i+1:]...)
	return
}

// Retrieve the parent node
//
// If several parents: Error (should not happen). If no parent: Error (it is the root?)
//
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

// Retrieves the Edge going from node n to its parent node.
//
// If several parents: Error (should not happen). If no parent: Error (it is the root?)
//
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

// Returns the index of the given edge in the list of edges going from this node.
//
// If the edge is not connected to the node n, then returns an error.
func (n *Node) EdgeIndex(e *Edge) (int, error) {
	for i := 0; i < len(n.br); i++ {
		if n.br[i] == e {
			return i, nil
		}
	}
	return -1, errors.New("The Edge is not in the neighbors of node")
}

// Returns the index of the given node next in the list of neighbors of node n.
//
// If next is not a neighbor of n, then returns an error.
func (n *Node) NodeIndex(next *Node) (int, error) {
	for i := 0; i < len(n.neigh); i++ {
		if n.neigh[i] == next {
			return i, nil
		}
	}
	return -1, errors.New("The Node is not in the neighbors of node")
}

// Returns if a node "next" is connected to "n" in the tree
func (n *Node) IsConnected(next *Node) bool {
	for i := 0; i < len(n.neigh); i++ {
		if n.neigh[i] == next {
			return true
		}
	}
	return false
}

// Randomly rotates order of neighbor nodes and edges
// of a given node.
//
// Topology is not changed, just the order of the tree traversal
func (n *Node) RotateNeighbors() {
	for i, _ := range n.neigh {
		j := rand.Intn(i + 1)
		n.neigh[i], n.neigh[j] = n.neigh[j], n.neigh[i]
		n.br[i], n.br[j] = n.br[j], n.br[i]
	}
}

// Recursive function that outputs newick representation
// from the current node
func (n *Node) Newick(parent *Node, newick *bytes.Buffer) {
	if len(n.neigh) > 0 {
		if len(n.neigh) > 1 {
			newick.WriteString("(")
		}
		nbchild := 0
		for i, child := range n.neigh {
			if child != parent {
				if nbchild > 0 {
					newick.WriteString(",")
				}
				child.Newick(n, newick)
				if n.br[i].support != NIL_SUPPORT && child.Name() == "" {
					newick.WriteString(strconv.FormatFloat(n.br[i].support, 'f', -1, 64))
					if n.br[i].pvalue != NIL_PVALUE {
						newick.WriteString(fmt.Sprintf("/%s", strconv.FormatFloat(n.br[i].pvalue, 'f', -1, 64)))
					}
				}
				if len(child.comment) != 0 {
					for _, c := range child.comment {
						newick.WriteString("[")
						newick.WriteString(c)
						newick.WriteString("]")
					}
				}
				if n.br[i].length != NIL_LENGTH {
					newick.WriteString(":")
					newick.WriteString(strconv.FormatFloat(n.br[i].length, 'f', -1, 64))
				}
				if len(n.br[i].comment) != 0 {
					for _, c := range n.br[i].comment {
						newick.WriteString("[")
						newick.WriteString(c)
						newick.WriteString("]")
					}
				}
				nbchild++
			}
		}
		if len(n.neigh) > 1 {
			newick.WriteString(")")
		}
	}
	newick.WriteString(n.name)
}
