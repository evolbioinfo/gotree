package tree

import (
	"fmt"
)

// Particulat rearrangement (NNI, SPR, etc.)
type Rearrangement interface {
	Apply() error
	Undo() error
}

// rearranger provides functions to modify input tree in order to search tree space
type Rearranger interface {
	// List all rearragements
	Rearrange(t *Tree, f func(r Rearrangement))
}

type NNIRearranger struct {
}

func (nnir *NNIRearranger) Rearrange(t *Tree, f func(r Rearrangement)) {
	for _, e := range t.Edges() {
		if e.Left().Nneigh() == 3 && e.Right().Nneigh() == 3 {
			f(newNNI(t, e.Left(), e.Right(), false))
			f(newNNI(t, e.Left(), e.Right(), true))
		}
	}
}

// Applies only to binary trees
type nni struct {
	t *Tree
	//  n1_1          n2_1
	//      \        /
	//       n1-----n2
	//	/        \
	//  n1_2          n2_2
	n1, n2     *Node
	n1_1, n1_2 *Node
	n2_1, n2_2 *Node
	// if false: after applying the NNI
	// n1_1 is grouped with n2_1 and
	// n1_2 is grouped with n2_2
	// if true: after applying the NNI:
	// n1_1 is grouped with n2_2 and
	// n1_2 is grouped with n2_1
	cross bool
	// If the NNI has already been applied
	applied bool
}

//  n1_1          n2_1
//      \        /
//       n1-----n2
//	/        \
//  n1_2          n2_2
// cross:
// if false: after applying the NNI
// n1_1 is grouped with n2_1 and
// n1_2 is grouped with n2_2
// if true: after applying the NNI:
// n1_1 is grouped with n2_2 and
// n1_2 is grouped with n2_1
func newNNI(t *Tree, n1, n2 *Node, cross bool) (n *nni) {
	var n1_1, n1_2, n2_1, n2_2 *Node

	// We search n1_1 and n1_2 in n1 neighbors
	n2index, _ := n1.NodeIndex(n2)
	n1_1 = n1.Neigh()[(n2index+1)%3]
	n1_2 = n1.Neigh()[(n2index+2)%3]

	// We search n2_1 and n2_2 in n2 neighbors
	n1index, _ := n2.NodeIndex(n1)
	n2_1 = n2.Neigh()[(n1index+1)%3]
	n2_2 = n2.Neigh()[(n1index+2)%3]

	n = &nni{t, n1, n2, n1_1, n1_2, n2_1, n2_2, cross, false}
	return
}

func (n *nni) Apply() (err error) {
	if n.applied {
		return
	}
	var n12index, n1index int
	var n1n2index int
	var n22index, n2index int
	var e1, e2 *Edge
	var n22node *Node

	if n1n2index, err = n.n1.NodeIndex(n.n2); err != nil {
		err = fmt.Errorf("Cannot create NNI with unconnected nodes n1 n2")
	}

	// We first get n12 index for n1 node
	if n12index, err = n.n1.NodeIndex(n.n1_2); err != nil {
		err = fmt.Errorf("Cannot apply NNI with unconnected nodes n1 n1_2")
		return
	}
	// then we get n1 index for  n12 node
	if n1index, err = n.n1_2.NodeIndex(n.n1); err != nil {
		err = fmt.Errorf("Cannot apply NNI with unconnected nodes n1_2 n1")
		return
	}

	// Same for n2 and n12 or n21 nodes
	if n.cross {
		n22node = n.n2_1
	} else {
		n22node = n.n2_2
	}
	if n22index, err = n.n2.NodeIndex(n22node); err != nil {
		err = fmt.Errorf("Cannot apply NNI with unconnected nodes n1 n2_1")
		return
	}
	if n2index, err = n22node.NodeIndex(n.n2); err != nil {
		err = fmt.Errorf("Cannot apply NNI with unconnected nodes n2_1 n2")
		return
	}

	e1 = n.n1.Edges()[n12index]
	e2 = n.n2.Edges()[n22index]

	// The root is somwhere in the
	// clade on the n1_2 side
	if e1.Right() == n.n1 {
		// Reorient n1-n2 edge
		n.n1.Edges()[n1n2index].Inverse()
	}

	n.n1.Edges()[n12index] = e2
	n.n2.Edges()[n22index] = e1

	n.n1.Neigh()[n12index] = n22node
	n22node.Neigh()[n2index] = n.n1

	n.n2.Neigh()[n22index] = n.n1_2
	n.n1_2.Neigh()[n1index] = n.n2

	if e1.Left() == n.n1 {
		e1.setLeft(n.n2)
	} else {
		e1.setRight(n.n2)
	}
	if e2.Left() == n.n2 {
		e2.setLeft(n.n1)
	} else {
		e2.setRight(n.n1)
	}

	n.applied = true

	return
}

func (n *nni) Undo() (err error) {
	if !n.applied {
		return
	}
	var n12index, n1index int
	var n1n2index int
	var n11index, n2index int
	var e1, e2 *Edge
	var n11node *Node

	if n1n2index, err = n.n1.NodeIndex(n.n2); err != nil {
		err = fmt.Errorf("Cannot create NNI with unconnected nodes n1 n2")
	}

	// We first get n12 index for n2 node
	if n12index, err = n.n2.NodeIndex(n.n1_2); err != nil {
		err = fmt.Errorf("Cannot apply NNI with unconnected nodes n2 n1_2")
		return
	}
	// then we get n2 index for n12 node
	if n2index, err = n.n1_2.NodeIndex(n.n2); err != nil {
		err = fmt.Errorf("Cannot apply NNI with unconnected nodes n1_2 n2")
		return
	}

	// Same for n2 and n12 or n21 nodes
	if n.cross {
		n11node = n.n2_1
	} else {
		n11node = n.n2_2
	}
	if n11index, err = n.n1.NodeIndex(n11node); err != nil {
		err = fmt.Errorf("Cannot apply NNI with unconnected nodes n1 n2_1")
		return
	}
	if n1index, err = n11node.NodeIndex(n.n1); err != nil {
		err = fmt.Errorf("Cannot apply NNI with unconnected nodes n2_1 n1")
		return
	}

	e1 = n.n1.Edges()[n11index]
	e2 = n.n2.Edges()[n12index]

	// The root is somwhere in the
	// clade on the n1_2 side (connected to n2)
	if e2.Right() == n.n2 {
		// Reorient n1-n2 edge
		n.n1.Edges()[n1n2index].Inverse()
	}

	n.n1.Edges()[n11index] = e2
	n.n2.Edges()[n12index] = e1

	n.n1.Neigh()[n11index] = n.n1_2
	n.n1_2.Neigh()[n2index] = n.n1

	n.n2.Neigh()[n12index] = n11node
	n11node.Neigh()[n1index] = n.n2

	if e1.Left() == n.n1 {
		e1.setLeft(n.n2)
	} else {
		e1.setRight(n.n2)
	}
	if e2.Left() == n.n2 {
		e2.setLeft(n.n1)
	} else {
		e2.setRight(n.n1)
	}

	n.applied = false

	return
}
