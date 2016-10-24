package tree

import (
	"errors"
	"fmt"
	"github.com/fredericlemoine/gotree/io"
	//"os"
)

/* Given a set of tip names
returns the node that is the common ancestor of them
and the edges that connects this node to the subtree
=> Considers the tree as unrooted
       e2---1
 ----a|
|      e1---2
|     ---3
 ----|
|     ---4
|     ---5
 ----|
      ---6
LeastCommonAncestorUnrooted(1,2) returns a,e1,e2,true
returned boolean value is true if the group is monophyletic
*/
func (t *Tree) LeastCommonAncestorUnrooted(nodeindex *nodeIndex, tips ...string) (*Node, []*Edge, bool) {
	if nodeindex == nil {
		nodeindex = NewNodeIndex(t)
	}
	tipindex := make(map[string]*Node, 0)
	for _, name := range tips {
		node, found := nodeindex.GetNode(name)
		if !found {
			io.ExitWithMessage(errors.New(fmt.Sprintf("Tip not found in the tree : %s", name)))
		}
		tipindex[name] = node
	}

	// We search a tip that is not in the input tips
	// It will serve as a temporary root for the tree
	var temproot *Node = nil
	for _, othertip := range t.Tips() {
		_, found := tipindex[othertip.Name()]
		if !found {
			temproot = othertip
			break
		}
	}

	// If temproot == nil : Means that the input tips consist of all the tips of the tree
	if temproot == nil {
		io.ExitWithMessage(errors.New("All tips of the tree given : Nothing to do"))
	}
	// otherwise we take the only child of the tip as first root
	ancestor, goodedges, _, diff, _ := t.LeastCommonAncestorUnrootedRecur(temproot.neigh[0], nil, tipindex)

	// fmt.Println("--")
	// fmt.Println(com)
	// fmt.Println(diff)
	// fmt.Println(found)
	// for _, s := range tips {
	// 	fmt.Println(s)
	// }
	// fmt.Println("--")
	return ancestor, goodedges, diff == 0
}

/* Returns for a given node ... */
func (t *Tree) LeastCommonAncestorUnrootedRecur(current *Node, prev *Node, tipIndex map[string]*Node) (*Node, []*Edge, int, int, bool) {
	common := 0
	edges := make([]*Edge, 0, 3)
	different := 0
	allFound := false

	// If current is a tip
	if current.Tip() {
		//fmt.Println(current.Name())
		_, found := tipIndex[current.Name()]
		if found {
			common++
			if idx, e := current.NodeIndex(prev); e == nil {
				edges = append(edges, current.br[idx])
			} else {
				io.ExitWithMessage(e)
			}
		} else {
			different = 1
		}
	}

	// If current is not a tip
	tmpdiff := 0
	for i, n := range current.neigh {
		if n != prev {
			node, succedges, com, diff, found := t.LeastCommonAncestorUnrootedRecur(n, current, tipIndex)
			if found {
				//fmt.Println("int found - diff:", diff)
				return node, succedges, com, diff, found
			} else if com > 0 {
				edges = append(edges, current.br[i])
				common += com
				different += diff
			} else {
				tmpdiff += diff
			}
		}
	}
	//fmt.Println("tmpdiff: ", tmpdiff)
	allFound = common == len(tipIndex)
	if allFound {
		//fmt.Println("found - diff:", different)
		return current, edges, common, different, allFound
	} else {
		different += tmpdiff
		//fmt.Println("diff:", different)
		return nil, nil, common, different, allFound
	}
}

/*
This function adds a bipartition at the given node and the given edges
Immagine a star tree with central node n,
     1
     |
     |
6----n-----2
    /|\
   / | \
 e5 e4  e3
if we call AddBipartition(n,{e3,e4,e5}) we end with:
     1
     |
     |
6----n-----2
     |
     |
     n2
    /|\
   / | \
 e5 e4  e3
*/
func (t *Tree) AddBipartition(n *Node, edges []*Edge, length, support float64) *Edge {
	n2 := t.NewNode()
	// Number of edges in direction n->e->other
	nbout := 0
	// Number of edges in direction n<-e<-other
	nbin := 0
	if len(edges) <= 1 || len(edges) >= len(n.br)-1 {
		io.ExitWithMessage(errors.New("We cannot add the bipartition, it already exists"))
	}
	for _, e := range edges {
		// We check if the edges are connected to the node
		// Else it exits with an error
		if e.Left() != n && e.Right() != n {
			io.ExitWithMessage(errors.New("Edges need to be connected to the node to add a bipartition"))
		}
		// Direction : true if n->e->other..., false if n<-e<-other
		// According to left / right
		dir := e.Left() == n
		var other *Node
		boot := e.Support()
		len := e.Length()
		var etmp *Edge
		if dir {
			nbout++
			other = e.Right()
			other.delNeighbor(n)
			n.delNeighbor(other)
			etmp = t.ConnectNodes(n2, other)
		} else {
			nbin++
			other = e.Left()
			other.delNeighbor(n)
			n.delNeighbor(other)
			etmp = t.ConnectNodes(other, n2)
		}
		etmp.SetLength(len)
		etmp.SetSupport(boot)
	}

	var e *Edge
	if nbin == 0 {
		e = t.ConnectNodes(n, n2)
	} else {
		e = t.ConnectNodes(n2, n)
	}
	e.SetLength(length)
	e.SetSupport(support)
	return e
}

/*
Builds the consensus of trees given in the input channel.
If the cutoff is 0.5 : The majority rule consensus is computed
If tht cutoff is 1   : The strict consensus is computed
In the output consensus tree:
1) Branch supports are computed as the proportion of trees in which the bipartition is present
2) Branch lengths are computed as the average length of the same branch over all the trees where it is present
There can be errors if:
* The cutoff <0.5 or >1
* The tip names are different in the different trees
* Incompatible bipartition are generated to build the consensus (It should not happen since cutoff should be >=0.5)
*/
func Consensus(trees <-chan Trees, cutoff float64) *Tree {
	if cutoff < 0.5 || cutoff > 1 {
		io.ExitWithMessage(errors.New("Min frequency for bipartition must be >=0.5 and <=1"))
	}
	nbtrees := 0
	edgeindex := NewEdgeIndex(128, .75)
	var nodeindex *nodeIndex
	var startree *Tree = nil
	nbtips := 0
	var alltips []string
	var err error
	// We fill the edge index with all the bipartition and their count
	for curtree := range trees {
		// If the star tree is not initialized, we create it with the tips of the first tree
		if startree == nil {
			alltips = curtree.Tree.AllTipNames()
			if startree, err = StarTreeFromTree(curtree.Tree); err != nil {
				io.ExitWithMessage(err)
			} else {
				nbtips = len(alltips)
				// We first build the node index
				nodeindex = NewNodeIndex(startree)
			}
		} else {
			// Compare tip names between star tree and current tree
			// Error if different sets (use already computed indexes)
			names := curtree.Tree.AllTipNames()
			if len(names) != nbtips {
				io.ExitWithMessage(errors.New("Trees do not have the same set of tips"))
			}
			for _, name := range names {
				if ok, err3 := startree.ExistsTip(name); err3 != nil {
					io.ExitWithMessage(err)
				} else if !ok {
					io.ExitWithMessage(errors.New("Trees do not have the same set of tips"))
				}
			}
		}
		// We add the edge into the index
		for _, e := range curtree.Tree.Edges() {
			edgeindex.AddEdgeCount(e)
		}
		nbtrees++
	}

	// We take the bipartitions that are present in more than cutoff trees and less
	// than or equal the number of trees
	// And we add it to the startree
	for _, bs := range edgeindex.BitSets(int(cutoff*float64(nbtrees)), nbtrees) {
		names := make([]string, 0, bs.key.Count())
		for _, n := range alltips {
			if idx, err := startree.TipIndex(n); err != nil {
				io.ExitWithMessage(err)
			} else {
				if bs.key.Test(idx) {
					names = append(names, n)
				}
			}
		}

		// Names of the tips in one side of the bipartition
		if len(names) < 2 {
			if len(names) == 1 {
				if t, ok := nodeindex.GetNode(names[0]); !ok || !t.Tip() {
					io.ExitWithMessage(errors.New(fmt.Sprintf("This taxon name does not exist in the consensus: %s", names[0])))
				} else {
					t.br[0].SetLength(float64(bs.val.Len) / float64(bs.val.Count))
				}
			} else {
				io.ExitWithMessage(errors.New("This bipartition has a side with no taxa"))
			}
		} else {
			node, edges, monophyletic := startree.LeastCommonAncestorUnrooted(nodeindex, names...)
			if node == nil {
				io.ExitWithMessage(errors.New("Consensus error: No common ancestor found for biparition"))
			}
			if edges == nil || len(edges) == 0 {
				io.ExitWithMessage(errors.New("Consensus error: No common ancestor Edges found"))
			}
			if !monophyletic {
				io.ExitWithMessage(errors.New("The group should be monophyletic"))
			}
			// We add the bipartition with a support value corresponding to the percentage of
			// trees in which it appears
			// TODO: Average branch length : Need to change the data structure
			startree.AddBipartition(node, edges, float64(bs.val.Len)/float64(bs.val.Count), float64(bs.val.Count)/float64(nbtrees))
		}
	}

	startree.UpdateTipIndex()
	startree.ClearBitSets()
	startree.UpdateBitSet()
	startree.ComputeDepths()

	return startree
}
