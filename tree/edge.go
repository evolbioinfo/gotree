package tree

import (
	"bytes"
	"errors"
	"fmt"
	"math"
	"strconv"

	"github.com/evolbioinfo/gotree/io"
	"github.com/evolbioinfo/gotree/mutils"
	"github.com/fredericlemoine/bitset"
)

// Structure of an edge
type Edge struct {
	left, right *Node    // Left and right nodes
	length      float64  // length of branch
	comment     []string // Comment if any in the newick file
	support     float64  // -1 if no support
	pvalue      float64  // -1 if no pvalue
	// a Bit at index i in the bitset corresponds to the position of the tip i
	//left:0/right:1 .
	// i is the index of the tip in the sorted tip name array
	bitset *bitset.BitSet // Bitset of length Number of taxa each
	id     int            // this field is used at discretion of the user to store information
}

// Constant for uninitialized values
const (
	NIL_SUPPORT = -1.0
	NIL_LENGTH  = -1.0
	NIL_PVALUE  = -1.0
	NIL_ID      = -1.0
)

/* Edge functions */
/******************/

// Sets left node (parent)
func (e *Edge) setLeft(left *Node) {
	e.left = left
}

// Sets right node (child)
func (e *Edge) setRight(right *Node) {
	e.right = right
}

// Inverse Edge orientation:
// left becomes right and
// right becomes left
func (e *Edge) Inverse() {
	e.left, e.right = e.right, e.left

}

// Sets the pvalue of this edge (if not null, pvalue
// is stored/parsed as "/pvalue" in the bootstrap value
// field.
func (e *Edge) SetPValue(pval float64) {
	e.pvalue = pval
}

// Sets the length of the branch
func (e *Edge) SetLength(length float64) {
	e.length = length
}

// Sets the branch support
func (e *Edge) SetSupport(support float64) {
	e.support = support
}

// returns the length of the branch
func (e *Edge) Length() float64 {
	return e.length
}

// Returns the length as a string representing the
// right precision float (not 0.010000000 but
// 0.01 for example)
func (e *Edge) LengthString() string {
	length := "N/A"
	if e.Length() != NIL_LENGTH {
		length = fmt.Sprintf("%s", strconv.FormatFloat(e.Length(), 'f', -1, 64))
	}
	return length
}

// Returns the support of that branch
func (e *Edge) Support() float64 {
	return e.support
}

// Returns the support as a string representing the
// right precision float (not 0.90000000 but
// 0.9 for example)
func (e *Edge) SupportString() string {
	support := "N/A"
	if e.Support() != NIL_SUPPORT {
		support = fmt.Sprintf("%s", strconv.FormatFloat(e.Support(), 'f', -1, 64))
	}
	return support
}

// Returns the Pvalue of that branch
func (e *Edge) PValue() float64 {
	return e.pvalue
}

// Returns the node at the right side of the edge (child)
func (e *Edge) Right() *Node {
	return e.right
}

// Returns the node at the left side of the edge (parent)
func (e *Edge) Left() *Node {
	return e.left
}

// Returns the BitSet of that edge. It may be nil if not initialized.
//
// the ith bit corresponds position of tip i around the branch (left:0/right:1).
//
// i is the index of the tip in the sorted tip name array
func (e *Edge) Bitset() *bitset.BitSet {
	return e.bitset
}

// Returns the Id of the branch.
//
// Returns an error if not initialized.
func (e *Edge) Id() int {
	if e.id == NIL_ID {
		io.ExitWithMessage(errors.New("Id has not been set"))
	}
	return e.id
}

// Sets the id of the branch
func (e *Edge) SetId(id int) {
	e.id = id
}

// Returns the name associated to this Edge.
//
// If rooted, the output clade name is the name of the
// descendent node.
//
// Else, the clade name is the name of the node on the
// lightest side. In that case bitsets need to be initialized.
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

// Returns the size (number of tips) of the light side
// (smallest number of tips) of the given branch.
//
// Bitsets must be initialized otherwise returns an error.
func (e *Edge) TopoDepth() (int, error) {
	if e.bitset == nil {
		return -1, errors.New("Cannot compute topodepth, Bitset is nil")
	}
	if e.bitset.None() {
		return -1, errors.New("Cannot compute topodepth, Bitset is 000...0")
	}
	count := int(e.bitset.Count())
	total := int(e.bitset.Len())
	return mutils.Min(count, total-count), nil
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
	6 - name of right node if any
        7 - comments associated to the edge
        8 - name of left node if any
        9 - comment of right node if any
       10 - comment of left node if any
*/
func (e *Edge) ToStatsString(withedgecomments bool) string {
	var err error
	length := e.LengthString()
	support := e.SupportString()

	var depth, leftdepth, rightdepth int

	if leftdepth, err = e.Left().Depth(); err != nil {
		io.ExitWithMessage(err)
	}
	if rightdepth, err = e.Right().Depth(); err != nil {
		io.ExitWithMessage(err)
	}
	depth = mutils.Min(leftdepth, rightdepth)
	var topodepth int
	topodepth, err = e.TopoDepth()
	if err != nil {
		io.ExitWithMessage(err)
	}

	rightname := e.Right().Name()

	comment := ""
	if withedgecomments {
		rightcomment := e.Right().CommentsString()
		leftname := e.Left().Name()
		leftcomment := e.Left().CommentsString()
		comment = "\t" + e.CommentsString() + "\t" + leftname + "\t" + rightcomment + "\t" + leftcomment
	}

	return fmt.Sprintf("%s\t%s\t%t\t%d\t%d\t%s%s",
		length, support, e.Right().Tip(),
		depth, topodepth, rightname, comment)
}

// Returns true if this edge defines the same biparition of the tips
// than the edge in argument.
//
// Bitsets must be initialized
func (e *Edge) SameBipartition(e2 *Edge) bool {
	return e.bitset.EqualOrComplement(e2.bitset)
}

// Tests wether the tip with index id in the bitset
// is Set or not.
//
// The index corresponds to tree.Tipindex(tipname)
func (e *Edge) TipPresent(id uint) bool {
	return e.bitset.Test(id)
}

// Number of tips on the right side of the bipartition
// Used by "TopoDepth" function for example.
//
// Bitsets must be initialized, otherwise returns an error.
func (e *Edge) NumTipsRight() (int, error) {
	if e.bitset == nil {
		return -1, errors.New("Cannot count right tips, Bitset is nil")
	}
	if e.bitset.None() {
		return -1, errors.New("Cannot count right tips, Bitset is 000...0")
	}

	return int(e.bitset.Count()), nil
}

// Number of tips on the left side of the bipartition
// Used by "TopoDepth" function for example.
//
// Bitsets must be initialized, otherwise returns an error.
func (e *Edge) NumTipsLeft() (int, error) {
	if e.bitset == nil {
		return -1, errors.New("Cannot count left tips, Bitset is nil")
	}
	if e.bitset.None() {
		return -1, errors.New("Cannot count left tips, Bitset is 000...0")
	}
	return int(e.bitset.Len() - e.bitset.Count()), nil
}

// Return the given edge in the array of edges comparing bitsets fields
// Return nil if not found.
//
// Bitsets must be initialized otherwise returns an error.
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

// Returns the average difference and the max difference in support between
// the current edge and its neighbors.
//
// The neighbors are defined by the branches located in a area defined
// by number of branches separating them (<d).
//
//	* cutoff: Cutoff to consider hx=true or hy=true
//	* hx=true if exists a neighbor branch with suppt > cutoff
//	* hy=true if the current branch has suppt > cutoff */
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
//
// The neighbors are defined by the branches located in a area defined by
// number of branches separating them (<d).
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

// Adds a comment to the edge. It will be coded by a list of []
// In the Newick format.
func (e *Edge) AddComment(comment string) {
	e.comment = append(e.comment, comment)
}

// Returns the list of comments associated to the edge.
func (e *Edge) Comments() []string {
	return e.comment
}

// Returns the string of comma separated comments
// surounded by [].
func (e *Edge) CommentsString() string {
	var buf bytes.Buffer
	buf.WriteRune('[')
	for i, c := range e.comment {
		if i > 0 {
			buf.WriteRune(',')
		}
		buf.WriteString(c)
	}
	buf.WriteRune(']')
	return buf.String()
}

// Removes all comments associated to the node
func (e *Edge) ClearComments() {
	e.comment = e.comment[:0]
}
