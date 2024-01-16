package mutations

import (
	"fmt"

	"github.com/evolbioinfo/goalign/align"
	"github.com/evolbioinfo/gotree/io"
	"github.com/evolbioinfo/gotree/tree"
)

func CountEEMs(t *tree.Tree, a align.Alignment) (mutations *MutationList, err error) {
	var sitemutations *MutationList
	mutations = NewMutationList()

	// We set branch ids
	nbranches := 0
	for _, e := range t.Edges() {
		e.SetId(nbranches)
		nbranches++
	}

	// We set nodes ids identical to index in alignment
	// and we check that sequences correspond to all tree nodes
	// and all tree nodes have a name
	for _, n := range t.Nodes() {
		var i int

		if n.Name() == "" {
			err = fmt.Errorf("all nodes of the phylogeny must have a name")
			io.LogError(err)
			return
		}
		if i = a.GetSequenceIdByName(n.Name()); i < 0 {
			err = fmt.Errorf("node %s of the phylogeny does not have an associated sequence in the alignment", n.Name())
			io.LogError(err)
			return
		}
		n.SetId(i)
	}

	// We iterate over alignment sites
	for i := 0; i < a.Length(); i++ {
		if sitemutations, err = countEEMsSite(t, a, i); err != nil {
			io.LogError(err)
			return
		}
		for _, v := range sitemutations.Mutations {
			id := fmt.Sprintf("%d-%c-%c", v.AlignmentSite,
				rune(v.ParentCharacter), rune(v.ChildCharacter))
			m, ok := mutations.Mutations[id]
			if ok {
				m.NumEEM = m.NumEEM + 1
				mutations.Mutations[id] = m
			} else {
				mutations.Mutations[id] = v
			}
		}
	}

	return
}

func countEEMsSite(t *tree.Tree, a align.Alignment, site int) (mutations *MutationList, err error) {
	mutations = NewMutationList()
	err = countEEMSiteBranch(t, nil, t.Root(), nil, a, site, mutations, nil)
	return
}

func countEEMSiteBranch(t *tree.Tree, prevNode *tree.Node, currentNode *tree.Node, currentBranch *tree.Edge,
	a align.Alignment, site int, mutations *MutationList, curMutation *Mutation) (err error) {
	var nextNode *tree.Node
	var nextIndex int
	var prevChar, curChar uint8
	var prevSeq, curSeq []uint8

	curSeq, _ = a.GetSequenceCharById(currentNode.Id())
	curChar = curSeq[site]

	if prevNode != nil {
		prevSeq, _ = a.GetSequenceCharById(prevNode.Id())
		prevChar = prevSeq[site]

		if prevChar != curChar {
			curMutation = &Mutation{
				AlignmentSite:   site,
				BranchIndex:     currentBranch.Id(),
				ChildNodeName:   currentNode.Name(),
				ParentCharacter: prevChar,
				ChildCharacter:  curChar,
				NumEEM:          1,
			}
		}
	}

	if currentNode.Tip() {
		if curMutation != nil {
			id := fmt.Sprintf("%d-%d-%c-%c", curMutation.AlignmentSite, curMutation.BranchIndex,
				rune(curMutation.ParentCharacter), rune(curMutation.ChildCharacter))
			mutations.Mutations[id] = *curMutation
		}
	} else {
		for nextIndex, nextNode = range currentNode.Neigh() {
			if nextNode != prevNode {
				if err = countEEMSiteBranch(t, currentNode, nextNode, currentNode.Edges()[nextIndex], a,
					site, mutations, curMutation); err != nil {
					return
				}
			}
		}
	}

	return
}
