package newick

import (
	"errors"

	"github.com/evolbioinfo/gotree/tree"
)

type nodeStackElt struct {
	n *tree.Node
	e *tree.Edge
}

type NodeStack struct {
	elt []nodeStackElt
}

/* Initialize a new Node Stack*/
func NewNodestack() (ns *NodeStack) {
	return &NodeStack{
		make([]nodeStackElt, 0, 10),
	}
}

/* Pushes a new node/edge pair to the Stack */
func (ns *NodeStack) Push(n *tree.Node, e *tree.Edge) {
	ns.elt = append(ns.elt, nodeStackElt{n, e})
}

/**
Pops and returns the head of the Stack. The calling function is
responsible for freeing the elt: free(elt).

Returns an error if the stack is empty
*/
func (ns *NodeStack) Pop() (n *tree.Node, e *tree.Edge, err error) {
	var last nodeStackElt
	if len(ns.elt) == 0 {
		err = errors.New("Cannot Pop an empty stack")
		return
	}
	last, ns.elt = ns.elt[len(ns.elt)-1], ns.elt[:len(ns.elt)-1]
	n, e = last.n, last.e
	last.n = nil
	last.e = nil
	return
}

/**
Returns the head of the Stack, and an error if the stack is empty
*/
func (ns *NodeStack) Head() (n *tree.Node, e *tree.Edge, err error) {
	if len(ns.elt) == 0 {
		err = errors.New("An empty stack has no head")
		return
	}
	head := ns.elt[len(ns.elt)-1]
	n, e = head.n, head.e
	return
}

/* Clears the whole stack and all its elements */
func (ns *NodeStack) Clear() {
	for _, el := range ns.elt {
		el.e = nil
		el.n = nil
	}
	ns.elt = ns.elt[:0]
}
