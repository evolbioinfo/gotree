package tree

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/fredericlemoine/gotree/io"
)

type Node struct {
	name    string   // Name of the node
	comment []string // Comment if any in the newick file
	neigh   []*Node  // neighbors array
	br      []*Edge  // Branches array (same order than neigh)
	depth   int      // Depth of the node
	id      int      // this field is used at discretion of the user to store information
}

// Adds a child n to the node p, connected with edge e
func (p *Node) addChild(n *Node, e *Edge) {
	p.neigh = append(p.neigh, n)
	p.br = append(p.br, e)

}

func (n *Node) SetName(name string) {
	n.name = name
}

func (n *Node) AddComment(comment string) {
	n.comment = append(n.comment, comment)
}

func (n *Node) SetDepth(depth int) {
	n.depth = depth
}

func (n *Node) Name() string {
	return n.name
}

func (n *Node) Id() int {
	if n.id == -1 {
		io.ExitWithMessage(errors.New("Id has not been set"))
	}
	return n.id
}

func (n *Node) SetId(id int) {
	n.id = id
}

func (n *Node) Depth() (int, error) {
	if n.depth == -1 {
		return n.depth, errors.New("Node depth has not been computed")
	}
	return n.depth, nil
}

// Number of neighbors of this node
func (n *Node) Nneigh() int {
	return len(n.neigh)
}

// Neighbors of this node
func (n *Node) Neigh() []*Node {
	return n.neigh
}

// Is a tip or not?
func (n *Node) Tip() bool {
	return len(n.neigh) == 1
}

func (n *Node) Edges() []*Edge {
	return n.br
}

func (n *Node) delNeighbor(n2 *Node) error {
	i, err := n.NodeIndex(n2)
	if err != nil {
		return err
	}
	n.br = append(n.br[0:i], n.br[i+1:]...)
	n.neigh = append(n.neigh[0:i], n.neigh[i+1:]...)
	return nil
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
				if n.br[i].support != -1 {
					newick.WriteString(fmt.Sprintf("%f", n.br[i].support))
					if n.br[i].pvalue != -1 {
						newick.WriteString(fmt.Sprintf("/%f", n.br[i].support))
					}
				}
				if len(child.comment) != 0 {
					for _, c := range child.comment {
						newick.WriteString("[")
						newick.WriteString(c)
						newick.WriteString("]")
					}
				}
				if n.br[i].length != -1 {
					newick.WriteString(":")
					newick.WriteString(fmt.Sprintf("%f", n.br[i].length))
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
