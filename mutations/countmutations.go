package mutations

import (
	"fmt"

	"github.com/evolbioinfo/goalign/align"
	"github.com/evolbioinfo/gotree/io"
	"github.com/evolbioinfo/gotree/tree"
)

func CountMutations(t *tree.Tree, a align.Alignment) (mutations *MutationList, err error) {
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
		if sitemutations, err = countMutationsSite(t, a, i); err != nil {
			io.LogError(err)
			return
		}
		if err = mutations.Append(sitemutations); err != nil {
			io.LogError(err)
			return
		}
	}

	return
}

func countMutationsSite(t *tree.Tree, a align.Alignment, site int) (mutations *MutationList, err error) {
	mutations, _, _, err = countMutationSiteBranch(t, nil, t.Root(), nil, a, site)
	return
}

func countMutationSiteBranch(t *tree.Tree, prevNode *tree.Node, currentNode *tree.Node, currentBranch *tree.Edge, a align.Alignment, site int) (mutations *MutationList, ntips int, characterDistribution map[uint8]int, err error) {
	var tmpMutations *MutationList
	var tmpntips, nidtips int
	var tmpCharacterDistribution map[uint8]int
	var nextNode *tree.Node
	var nextIndex int
	var prevChar, curChar uint8
	var prevSeq, curSeq []uint8

	mutations = NewMutationList()
	characterDistribution = make(map[uint8]int)

	curSeq, _ = a.GetSequenceCharById(currentNode.Id())
	curChar = curSeq[site]

	if currentNode.Tip() {
		ntips = 1
		characterDistribution[curChar] = 1
	} else {
		for nextIndex, nextNode = range currentNode.Neigh() {
			if nextNode != prevNode {
				if tmpMutations, tmpntips, tmpCharacterDistribution, err = countMutationSiteBranch(t, currentNode, nextNode, currentNode.Edges()[nextIndex], a, site); err != nil {
					return
				}
				for char, nb := range tmpCharacterDistribution {
					if _, exist := characterDistribution[char]; !exist {
						characterDistribution[char] = nb
					} else {
						characterDistribution[char] += nb
					}
				}
				ntips += tmpntips
				if err = mutations.Append(tmpMutations); err != nil {
					return
				}
			}
		}
	}

	if prevNode != nil {
		prevSeq, _ = a.GetSequenceCharById(prevNode.Id())
		prevChar = prevSeq[site]

		if n, exist := characterDistribution[curChar]; !exist {
			nidtips = 0
		} else {
			nidtips = n
		}

		if prevChar != curChar {
			k := fmt.Sprintf("%d-%d-%c-%c", site, currentBranch.Id(), rune(prevChar), rune(curChar))
			m := Mutation{
				AlignmentSite:             site,
				BranchIndex:               currentBranch.Id(),
				ChildNodeName:             currentNode.Name(),
				ParentCharacter:           prevChar,
				ChildCharacter:            curChar,
				NumTips:                   ntips,
				NumTipsWithChildCharacter: nidtips,
			}
			mutations.Mutations[k] = m
		}
	}
	return
}
