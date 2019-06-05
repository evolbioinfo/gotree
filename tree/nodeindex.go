package tree

import "errors"

// Structure of a node index.
// Basically a hashmap that stores as keys node names,
// and values node of the trees.
type NodeIndex interface {
	GetNode(name string) (*Node, bool)
	AddNode(n *Node)
}

type nodeIndex struct {
	index map[string]*Node
}

// Only Tips
// Computes a node index for a given tree.
func NewNodeIndex(t *Tree) (*nodeIndex, error) {

	nodeindex := &nodeIndex{
		index: make(map[string]*Node, 0),
	}

	nodes := t.Nodes()

	for _, n := range nodes {
		// tip
		if n.Name() != "" {
			if _, ok := nodeindex.index[n.Name()]; ok {
				return nil, errors.New("NewNodeIndex error: Tree contains several node with the same name: " + n.Name())
			}
			nodeindex.index[n.Name()] = n
		}
	}

	return nodeindex, nil
}

// Tips + internal nodes
func NewAllNodeIndex(t *Tree) *nodeIndex {
	nodeindex := &nodeIndex{
		index: make(map[string]*Node, 0),
	}

	nodes := t.Nodes()

	for _, n := range nodes {
		nodeindex.index[n.Name()] = n
	}

	return nodeindex
}

// returns the node associated to the name in argument
// it may correspond to a tip node or an internal node
// with a name
func (ni *nodeIndex) GetNode(name string) (*Node, bool) {
	n, ok := ni.index[name]
	return n, ok
}

// Adds the given node to the index
// If it already exists, then replaces it
func (ni *nodeIndex) AddNode(n *Node) {
	ni.index[n.Name()] = n
}
