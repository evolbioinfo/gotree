package tree

import (
	"errors"
	"fmt"
	"github.com/fredericlemoine/bitset"
	"github.com/fredericlemoine/gotree/io"
)

type Edge struct {
	left, right *Node   // Left and right nodes
	length      float64 // length of branch
	support     float64 // -1 if no support
	// a Bit at index i in the bitset corresponds to the position of the tip i
	//left:0/right:1 .
	// i is the index of the tip in the sorted tip name array
	bitset *bitset.BitSet // Bitset of length Number of taxa each
	id     int            // this field is used at discretion of the user to store information
}

/* Edge functions */
/******************/

func (e *Edge) setLeft(left *Node) {
	e.left = left
}
func (e *Edge) setRight(right *Node) {
	e.right = right
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
	if e.id == -1 {
		io.ExitWithMessage(errors.New("Id has not been set"))
	}
	return e.id
}

func (e *Edge) SetId(id int) {
	e.id = id
}

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

func (e *Edge) DumpBitSet() string {
	if e.bitset == nil {
		return "nil"
	}
	return e.bitset.DumpAsBits()
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
	if e.Length() != -1 {
		length = fmt.Sprintf("%f", e.Length())
	}
	var support = "N/A"
	if e.Support() != -1 {
		support = fmt.Sprintf("%f", e.Support())
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
	return fmt.Sprintf("%s\t%s\t%t\t%d\t%d\t%s\n",
		length, support, e.Right().Tip(),
		depth, topodepth, e.Right().Name())

}

func (e *Edge) SameBipartition(e2 *Edge) bool {
	return e.bitset.EqualOrComplement(e2.bitset)
}

// Tests wether the tip with index id in the bitset
// is Set or not
// The index corresponds to tree.Tipindex(tipname)
func (e *Edge) TipPresent(id uint) bool {
	return e.bitset.Test(id)
}

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
