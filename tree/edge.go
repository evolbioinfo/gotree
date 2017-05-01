package tree

import (
	"errors"
	"fmt"
	"math"
	"strconv"

	"github.com/fredericlemoine/bitset"
	"github.com/fredericlemoine/gotree/io"
)

type Edge struct {
	left, right *Node   // Left and right nodes
	length      float64 // length of branch
	support     float64 // -1 if no support
	pvalue      float64 // -1 if no pvalue
	// a Bit at index i in the bitset corresponds to the position of the tip i
	//left:0/right:1 .
	// i is the index of the tip in the sorted tip name array
	bitset *bitset.BitSet // Bitset of length Number of taxa each
	id     int            // this field is used at discretion of the user to store information
}

const (
	NIL_SUPPORT = -1.0
	NIL_LENGTH  = -1.0
	NIL_PVALUE  = -1.0
	NIL_ID      = -1.0
)

/* Edge functions */
/******************/

func (e *Edge) setLeft(left *Node) {
	e.left = left
}
func (e *Edge) setRight(right *Node) {
	e.right = right
}

func (e *Edge) SetPValue(pval float64) {
	e.pvalue = pval
}

func (e *Edge) SetLength(length float64) {
	e.length = length
}

func (e *Edge) SetSupport(support float64) {
	e.support = support
}

func (e *Edge) Length() float64 {
	return e.length
}

func (e *Edge) Support() float64 {
	return e.support
}

func (e *Edge) PValue() float64 {
	return e.pvalue
}

func (e *Edge) Right() *Node {
	return e.right
}

func (e *Edge) Left() *Node {
	return e.left
}

func (e *Edge) Bitset() *bitset.BitSet {
	return e.bitset
}

func (e *Edge) Id() int {
	if e.id == NIL_ID {
		io.ExitWithMessage(errors.New("Id has not been set"))
	}
	return e.id
}

func (e *Edge) SetId(id int) {
	e.id = id
}

/*
If rooted, the output clade name is the name of the
descendent node.

if not rooted, then the clade name is the name of the node on
the lightest side
*/
func (e *Edge) Name(rooted bool) (nodename string) {
	//If rooted, the clade name is the name of the
	// descendent node
	if rooted || e.bitset.Count() <= e.bitset.Len() {
		nodename = e.Right().Name()
	} else {
		nodename = e.Left().Name()
	}
	return
}

// Returns the size (number of tips) of the smallest subtree
// between the two subtrees connected to this edge
func (e *Edge) TopoDepth() (int, error) {
	if e.bitset == nil {
		return -1, errors.New("Cannot compute topodepth, Bitset is nil")
	}
	if e.bitset.None() {
		return -1, errors.New("Cannot compute topodepth, Bitset is 000...0")
	}
	count := int(e.bitset.Count())
	total := int(e.bitset.Len())
	return min(count, total-count), nil
}

// Returns a string representing the bitset (bipartition)
// defined by this edge
func (e *Edge) DumpBitSet() string {
	if e.bitset == nil {
		return "nil"
	}
	s := e.bitset.DumpAsBits()
	return s[len(s)-int(e.bitset.Len())-1 : len(s)]
}

/* Returns a string containing informations about the edge:
Tab delimited:
1 - length
2 - support
3 - istip?
4 - depth
5 - topo depth
6 - name of node if any
*/
func (e *Edge) ToStatsString() string {
	var err error
	var length = "N/A"
	if e.Length() != NIL_LENGTH {
		length = fmt.Sprintf("%s", strconv.FormatFloat(e.Length(), 'f', -1, 64))
	}
	var support = "N/A"
	if e.Support() != NIL_SUPPORT {
		support = fmt.Sprintf("%s", strconv.FormatFloat(e.Support(), 'f', -1, 64))
	}

	var depth, leftdepth, rightdepth int

	if leftdepth, err = e.Left().Depth(); err != nil {
		io.ExitWithMessage(err)
	}
	if rightdepth, err = e.Right().Depth(); err != nil {
		io.ExitWithMessage(err)
	}
	depth = min(leftdepth, rightdepth)
	var topodepth int
	topodepth, err = e.TopoDepth()
	if err != nil {
		io.ExitWithMessage(err)
	}

	name := ""
	if e.PValue() != NIL_PVALUE {
		name = fmt.Sprintf("%s/%s", strconv.FormatFloat(e.Support(), 'f', -1, 64), strconv.FormatFloat(e.PValue(), 'f', -1, 64))
	} else {
		name = e.Right().Name()
	}

	return fmt.Sprintf("%s\t%s\t%t\t%d\t%d\t%s",
		length, support, e.Right().Tip(),
		depth, topodepth, name)

}

// Returns true if this edge defines the same biparition of the tips
// than the edge in argument
func (e *Edge) SameBipartition(e2 *Edge) bool {
	return e.bitset.EqualOrComplement(e2.bitset)
}

// Tests wether the tip with index id in the bitset
// is Set or not
// The index corresponds to tree.Tipindex(tipname)
func (e *Edge) TipPresent(id uint) bool {
	return e.bitset.Test(id)
}

// Number of tips on one side of the bipartition
// Used by "TopoDepth" function for example
func (e *Edge) NumTips() uint {
	return e.bitset.Count()
}

// Return the given edge in the array of edges comparing bitsets fields
// Return nil if not found
func (e *Edge) FindEdge(edges []*Edge) (*Edge, error) {
	if e.bitset == nil {
		return nil, errors.New("BitSets has not been initialized with tree.clearBitSetsRecur(nil, nil, uint(len(tree.tipIndex)))")
	}
	if e.bitset.None() {
		return nil, errors.New("One edge has a bitset of 0...000 : May be BitSets have not been updated with tree.UpdateBitSet()?")
	}
	for _, e2 := range edges {
		if e2.bitset == nil {
			return nil, errors.New("BitSets has not been initialized with tree.clearBitSetsRecur(nil, nil, uint(len(tree.tipIndex)))")
		}

		if e.Right().Tip() != e2.Right().Tip() {
			continue
		}
		// If we take all the edges, or if both edges are not tips
		if e.bitset.EqualOrComplement(e2.bitset) {
			if e2.bitset.None() {
				return nil, errors.New("One edge has a bitset of 0...000 : May be BitSets have not been updated with tree.UpdateBitSet()?")
			}
			return e, nil
		}
	}
	return nil, nil
}

// Returns the average difference and the max difference in support between the current edge and its neighbors
// The neighbors are defined by the branches with length of the path separating the branches < d
// cutoff: Cutoff to consider hx=true or hy=true
// hx=true if exists a neighbor branch with suppt > cutoff
// hy=true if the current branch has suppt > cutoff */
// Returns (avg diff, min diff, max diff, hx, hy)
func (e *Edge) Locality(maxdist int, cutoff float64) (float64, float64, float64, bool, bool) {
	neighbors := e.NeigborEdges(maxdist)

	avgdiff := 0.0 /* Avg diff of br sup and neighb sup */
	maxdiff := 0.0 /* max diff of br sup and neighb sup */
	mindiff := 0.0 /* min diff of br sup and neighb sup */
	hx := false    /* hx: true if exists a neighbor branch with suppt > cutoff */
	hy := false    /* hy: true if the current branch has suppt > cutoff */
	nbe := 0       /* nb neigh branches with support */

	hy = (e.Support() != NIL_SUPPORT && e.Support() > cutoff)
	for _, n := range neighbors {
		if n.Support() != NIL_SUPPORT {
			if n.Support() != NIL_SUPPORT && n.Support() > cutoff {
				hx = true
			}

			diff := math.Abs(e.Support() - n.Support())
			avgdiff += diff
			maxdiff = math.Max(maxdiff, diff)
			if nbe == 0 {
				mindiff = diff
			} else {
				mindiff = math.Min(mindiff, diff)
			}
			nbe++
		}
	}
	return avgdiff / float64(nbe), mindiff, maxdiff, hx, hy
}

// Returns the neighbors of the given edge.
// Neighbors are defined as branches separated of given branch by a path whose length < maxdist
func (e *Edge) NeigborEdges(maxdist int) []*Edge {
	edges := make([]*Edge, 0, 0)

	neigborEdgesRecur(e.Left(), e, e.Right(), &edges, maxdist, 0)
	neigborEdgesRecur(e.Right(), e, e.Left(), &edges, maxdist, 0)

	return edges
}

func neigborEdgesRecur(cur *Node, curEdge *Edge, prev *Node, e *[]*Edge, maxdist, curdist int) {
	if curdist <= maxdist {
		// We do not take the first edge as its own neighbor
		if curdist > 0 {
			*e = append((*e), curEdge)
		}
		for i, n := range cur.neigh {
			if n != prev {
				neigborEdgesRecur(n, cur.br[i], cur, e, maxdist, curdist+1)
			}
		}
	}
}
