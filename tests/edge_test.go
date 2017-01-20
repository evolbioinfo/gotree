package tests

import (
	"fmt"
	"github.com/fredericlemoine/gotree/tree"
	"math"
	"testing"
)

// Tests the function to get neighboring edges of a given edges
// Testing the edges of distance 1
func TestEdgeNeighbor(t *testing.T) {
	for i := 0; i < 10; i++ {
		tr, err := tree.RandomYuleBinaryTree(200, true)
		if err != nil {
			t.Error(err)
		}
		edges := tr.Edges()
		for _, e := range edges {
			neighbors := e.NeigborEdges(1)
			if e.Left().Tip() || e.Right().Tip() {
				if len(neighbors) != 2 {
					t.Error(fmt.Sprintf("Number of neighbors of this edge should be 2 and is %d", len(neighbors)))
				}
			} else if tr.Root() == e.Left() || tr.Root() == e.Right() {
				if len(neighbors) != 3 {
					t.Error(fmt.Sprintf("Number of neighbors of this edge should be 3 and is %d", len(neighbors)))
				}
			} else {
				if len(neighbors) != 4 {
					t.Error(fmt.Sprintf("Number of neighbors of this edge should be 4 and is %d", len(neighbors)))
				}
			}
		}
	}
}

// Tests the function to get neighboring edges of a given edges
// Testing the edges of distance x on a balanced rooted tree
func TestEdgeNeighbor2(t *testing.T) {
	// Expected number of neighbors of root connected edges
	expected := 0
	// Random balanced binary tree
	tr, err := tree.RandomBalancedBinaryTree(15, true)
	if err != nil {
		t.Error(err)
	}

	var d uint
	for d = 1; d < 15; d++ {
		// we add 2^d neighbors on one side
		expected += (1 << d)
		// we add 2^(d-1) neighbors on the other side
		expected += (1 << (d - 1))
		// We only test branches connected to the root
		for _, e := range tr.Root().Edges() {
			neighbors := e.NeigborEdges(int(d))
			if len(neighbors) != expected {
				t.Error(fmt.Sprintf("Number of neighbors of depth %d of this edge should be %d and is %d", d, expected, len(neighbors)))
			}
		}
	}
}

// Tests locality
// Of distance 1
func TestLocality(t *testing.T) {
	cutoff := 0.8
	for i := 0; i < 10; i++ {
		tr, err := tree.RandomYuleBinaryTree(100, true)
		if err != nil {
			t.Error(err)
		}
		edges := tr.Edges()
		for _, e := range edges {
			e.SetSupport(1.0)
		}
		for _, e := range edges {
			for dist := 1; dist < 5; dist++ {
				avgloc, minloc, maxloc, hx, hy := e.Locality(dist, cutoff)

				if avgloc != 0 {
					t.Error(fmt.Sprintf("Avg locality should be 0 is and is %f", avgloc))
				}
				if minloc != 0 {
					t.Error(fmt.Sprintf("Min locality should be 0 is and is %f", minloc))
				}
				if maxloc != 0 {
					t.Error(fmt.Sprintf("Max locality should be 0 is and is %f", maxloc))
				}
				if !hy {
					t.Error(fmt.Sprintf("Branch entropy for cutoff > 0.8 should be true and is %t", hx))
				}
				if !hx {
					t.Error(fmt.Sprintf("Neighbor entropy for cutoff > 0.8 should be true and is %t", hx))
				}
			}
		}
	}
}

// Tests locality 2
// Of distance 1 with binary tree
// and alternating 0 and 1 supports from tips to root
func TestLocality2(t *testing.T) {
	expected := 3.0 / 4.0
	cutoff := 0.8
	tr, err := tree.RandomBalancedBinaryTree(10, true)
	if err != nil {
		t.Error(err)
	}
	edges := tr.Edges()
	for _, e := range edges {
		d, err := e.TopoDepth()
		if err != nil {
			t.Error(err)
		}
		if int(math.Log2(float64(d)))%2 == 0 {
			e.SetSupport(1.0)
		} else {
			e.SetSupport(0.0)
		}
	}
	// fmt.Println(tr.Newick())
	for _, e := range edges {
		if tr.Root() != e.Left() && !e.Right().Tip() {
			avgloc, minloc, maxloc, hx, hy := e.Locality(1, cutoff)
			if avgloc != expected {
				t.Error(fmt.Sprintf("Avg locality should be %f is and is %f", expected, avgloc))
			}
			if minloc != 0.0 {
				t.Error(fmt.Sprintf("Min locality should be %f is and is %f", 0.0, minloc))
			}
			if maxloc != 1.0 {
				t.Error(fmt.Sprintf("Max locality should be %f is and is %f", 1.0, maxloc))
			}
			if e.Support() > cutoff && !hy {
				t.Error(fmt.Sprintf("Branch entropy for cutoff > 0.8 should be true and is %t", hx))
			}
			if !hx {
				t.Error(fmt.Sprintf("Neighbor entropy for cutoff > 0.8 should be true and is %t", hx))
			}
		}
	}
}
