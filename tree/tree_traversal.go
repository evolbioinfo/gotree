package tree

func (t *Tree) PostOrder(f func(cur *Node, prev *Node, e *Edge) (keep bool)) {
	t.postOrderRecur(t.Root(), nil, nil, f)
}

func (t *Tree) postOrderRecur(cur *Node, prev *Node, e *Edge, f func(cur *Node, prev *Node, e *Edge) (keep bool)) (keep bool) {
	keep = true
	for i, n := range cur.neigh {
		if n != prev {
			if keep = t.postOrderRecur(n, cur, cur.Edges()[i], f); !keep {
				return
			}
		}
	}
	keep = keep && f(cur, prev, e)
	return
}

func (t *Tree) PreOrder(f func(cur *Node, prev *Node, e *Edge) (keep bool)) {
	t.preOrderRecur(t.Root(), nil, nil, f)
}

func (t *Tree) preOrderRecur(cur *Node, prev *Node, e *Edge, f func(cur *Node, prev *Node, e *Edge) (keep bool)) (keep bool) {
	keep = true
	if keep = f(cur, prev, e); !keep {
		return
	}
	for i, n := range cur.neigh {
		if n != prev {
			if keep = t.preOrderRecur(n, cur, cur.Edges()[i], f); !keep {
				return
			}
		}
	}
	return
}
