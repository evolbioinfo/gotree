package tree

func (t *Tree) PostOrder(f func(cur *Node, prev *Node, e *Edge)) {
	t.postOrderRecur(t.Root(), nil, nil, f)
}

func (t *Tree) postOrderRecur(cur *Node, prev *Node, e *Edge, f func(cur *Node, prev *Node, e *Edge)) {
	for i, n := range cur.neigh {
		if n != prev {
			t.postOrderRecur(n, cur, cur.Edges()[i], f)
		}
	}
	f(cur, prev, e)
}

func (t *Tree) PreOrder(f func(cur *Node, prev *Node, e *Edge)) {
	t.preOrderRecur(t.Root(), nil, nil, f)
}

func (t *Tree) preOrderRecur(cur *Node, prev *Node, e *Edge, f func(cur *Node, prev *Node, e *Edge)) {
	f(cur, prev, e)
	for i, n := range cur.neigh {
		if n != prev {
			t.preOrderRecur(n, cur, cur.Edges()[i], f)
		}
	}
}
