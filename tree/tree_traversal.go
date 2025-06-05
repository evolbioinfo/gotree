package tree

func (t *Tree) PostOrder(f func(cur *Node, prev *Node, e *Edge) (keep bool)) {
	postOrderRecur(t.Root(), nil, nil, f)
}

func postOrderRecur(cur *Node, prev *Node, e *Edge, f func(cur *Node, prev *Node, e *Edge) (keep bool)) (keep bool) {
	keep = true
	for i, n := range cur.neigh {
		if n != prev {
			if keep = postOrderRecur(n, cur, cur.Edges()[i], f); !keep {
				return
			}
		}
	}
	keep = keep && f(cur, prev, e)
	return
}

func (t *Tree) PreOrder(f func(cur *Node, prev *Node, e *Edge) (keep bool)) {
	preOrderRecur(t.Root(), nil, nil, f)
}

func preOrderRecur(cur *Node, prev *Node, e *Edge, f func(cur *Node, prev *Node, e *Edge) (keep bool)) (keep bool) {
	keep = true
	if keep = f(cur, prev, e); !keep {
		return
	}
	for i, n := range cur.neigh {
		if n != prev {
			if keep = preOrderRecur(n, cur, cur.Edges()[i], f); !keep {
				return
			}
		}
	}
	return
}

func (n *Node) PostOrder(f func(cur *Node, prev *Node, e *Edge) (keep bool)) (err error) {
	var parent *Node
	if parent, err = n.Parent(); err != nil {
		postOrderRecur(n, parent, nil, f)
	}
	return
}

func (n *Node) PreOrder(f func(cur *Node, prev *Node, e *Edge) (keep bool)) (err error) {
	var parent *Node
	if parent, err = n.Parent(); err != nil {
		return
	}
	preOrderRecur(n, parent, nil, f)
	return
}
