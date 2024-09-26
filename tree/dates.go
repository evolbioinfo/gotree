package tree

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"

	"github.com/evolbioinfo/gotree/io"
)

// LTTData describes a Lineage to Time data point
type LTTData struct {
	X float64 // Time or Mutations
	Y int     // Number of lineages
}

// Get Node dates
// Returns a slice of float correspsponding to all node dates (internal and external)
// Node IDs are their index in the slice.
// If one node does not have date or a malformed date, returns an error
func (t *Tree) NodeDates() (ndates []float64, err error) {
	var date float64
	var pattern *regexp.Regexp
	var matches []string

	ndates = make([]float64, 0)
	pattern = regexp.MustCompile(`(?i)&date="(.+)"`)
	nnodes := 0
	t.PreOrder(func(cur *Node, prev *Node, e *Edge) (keep bool) {
		keep = true
		if cur.Id() != nnodes {
			err = fmt.Errorf("node id does not correspond to postorder traversal: %d vs %d", cur.Id(), nnodes)
			keep = false
		} else if len(cur.Comments()) > 0 {
			keep = false
			for _, c := range cur.Comments() {
				matches = pattern.FindStringSubmatch(c)
				if len(matches) < 2 {
					err = fmt.Errorf("no date found: %s", c)
				} else if date, err = strconv.ParseFloat(matches[1], 64); err != nil {
					err = fmt.Errorf("one of the node date is malformed: %s", c)
				} else {
					ndates = append(ndates, date)
					err = nil
					keep = true
				}
			}
		} else {
			err = fmt.Errorf("a node with no date found")
			keep = false
		}
		nnodes += 1
		return
	})
	return
}

// LTTData describes a Lineage to Time data point
func (t *Tree) LTT() (lttdata []LTTData) {
	var lttdatadup []LTTData
	var dists []float64
	var err error

	// We compute distance from root to all nodes
	// If the field [&date=] exists, then takes it
	// Otherwise, computes the distance to the root
	if dists, err = t.NodeDates(); err != nil {
		io.LogWarning(err)
		io.LogWarning(fmt.Errorf("using mutations instead of dates"))
		dists = t.NodeRootDistance()
	}

	// This initializes
	lttdatadup = make([]LTTData, len(dists))
	// Version with one point per x, already summed up
	lttdata = make([]LTTData, len(dists))

	t.PreOrder(func(cur, prev *Node, e *Edge) (keep bool) {
		lttdatadup[cur.Id()].X = dists[cur.Id()]
		lttdatadup[cur.Id()].Y = cur.Nneigh()
		if prev != nil {
			lttdatadup[cur.Id()].Y -= 2
		}
		return true
	})
	sort.Slice(lttdatadup, func(i, j int) bool {
		return lttdatadup[i].X < lttdatadup[j].X
	})

	lasti := 0
	for i, l := range lttdatadup {
		if i == 0 {
			lttdata[i] = l
		} else {
			if lttdata[lasti].X == l.X {
				lttdata[lasti].Y += l.Y
			} else {
				lasti++
				lttdata[lasti] = l
			}
		}
	}
	lttdata = lttdata[:lasti+1]

	dists = nil
	total := 0
	for i := range lttdata {
		total += lttdata[i].Y
		lttdata[i].Y = total
	}
	return
}

// RTTData describes a Root To Tip Regression
type RTTData struct {
	X float64 // Date of the tip
	Y float64 // Distance to root
}

// RTTData describes a Root To Tip Regression
func (t *Tree) RTT() (rttdata []RTTData, err error) {
	var dists []float64
	var dates []float64

	// We compute distance from root to all nodes
	// If the field [&date=] exists, then takes it
	// Otherwise, computes the distance to the root
	if dates, err = t.NodeDates(); err != nil {
		io.LogWarning(err)
		err = fmt.Errorf("using mutations instead of dates")
		io.LogWarning(err)
		return
	}

	dists = t.NodeRootDistance()

	if len(dists) != len(dates) {
		err = fmt.Errorf("length of dates differs from length of distances")
		io.LogWarning(err)
		return
	}

	rttdata = make([]RTTData, 0, len(dates))
	for i, v := range dists {
		rttdata = append(rttdata, RTTData{dates[i], v})
	}

	return
}

// CutTreeMinDate traverses the tree, and only keep subtree starting at the given min date
//
// If a node has the exact same date as mindate: it becomes the root of a new tree
// If a node has a date > mindate and its parent has a date < mindate: a new node is added as a the root of a new tree, with one child, the currrent node.
// The output is a forest
func (t *Tree) CutTreeMinDate(mindate float64) (forest []*Tree, err error) {
	var dates []float64
	forest = make([]*Tree, 0, 10)
	var tmpforest []*Tree

	// If the field [&date=] exists, then takes it
	// Otherwise, returns an error
	if dates, err = t.NodeDates(); err != nil {
		io.LogWarning(err)
		err = fmt.Errorf("no dates provided in in the tree, of the form &date=")
		io.LogWarning(err)
		return
	}

	if tmpforest, err = cutTreeMinDateRecur(t.Root(), nil, nil, mindate, dates); err != nil {
		return
	}
	forest = append(forest, tmpforest...)

	return
}

func cutTreeMinDateRecur(cur, prev *Node, e *Edge, mindate float64, dates []float64) (forest []*Tree, err error) {
	// We take the branches/nodes >= min-date
	var tmptree *Tree
	var tmpnode *Node
	var tmpedge *Edge
	var tmpforest []*Tree

	forest = make([]*Tree, 0)
	// The current node is at the exact min date: we keep the subtree starting at this node
	// And disconnect the current node from its parent
	if dates[cur.Id()] == mindate || (prev == nil && dates[cur.Id()] >= mindate) {
		tmptree = NewTree()
		tmptree.SetRoot(cur)
		prev.delNeighbor(cur)
		cur.delNeighbor(prev)
		tmptree.ReinitIndexes()
		forest = append(forest, tmptree)
		return
	} else if prev != nil && dates[cur.Id()] > mindate && dates[prev.Id()] < mindate {
		tmptree = NewTree()
		tmpnode = tmptree.NewNode()
		tmptree.SetRoot(tmpnode)
		prev.delNeighbor(cur)
		cur.delNeighbor(prev)
		tmpedge = tmptree.ConnectNodes(tmpnode, cur)
		tmpnode.AddComment(fmt.Sprintf("&date=\"%f\"", mindate))
		tmpedge.SetLength(e.Length() * (dates[cur.Id()] - mindate) / (dates[cur.Id()] - dates[prev.Id()]))
		//tmptree.ReinitIndexes()
		forest = append(forest, tmptree)
		return
	}

	edges := make([]*Edge, len(cur.Edges()))
	copy(edges, cur.Edges())
	neigh := make([]*Node, len(cur.neigh))
	copy(neigh, cur.neigh)
	for i, n := range neigh {
		if n != prev {
			tmpforest, err = cutTreeMinDateRecur(n, cur, edges[i], mindate, dates)
			forest = append(forest, tmpforest...)
		}
	}

	return
}
