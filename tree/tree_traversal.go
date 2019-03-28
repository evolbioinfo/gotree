package tree

func (t *Tree) postOrderRecur(cur *Node, prev *Node, f func(cur *Node, prev *Node)) {
	for _, n := range cur.neigh {
		if n != prev {
			t.postOrderRecur(n, cur, f)
		}
	}
	f(cur, prev)
}

func (t *Tree) PostOrder(f func(*Node, *Node)) {
	t.postOrderRecur(t.Root(), nil, f)
}

func (t *Tree) preOrderRecur(cur *Node, prev *Node, f func(cur *Node, prev *Node)) {
	f(cur, prev)
	for _, n := range cur.neigh {
		if n != prev {
			t.preOrderRecur(n, cur, f)
		}
	}
}

func (t *Tree) PreOrder(f func(cur *Node, prev *Node)) {
	t.preOrderRecur(t.Root(), nil, f)
}
