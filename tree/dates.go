package tree

import (
	"fmt"
	"sort"

	"github.com/evolbioinfo/gotree/io"
)

// LTTData describes a Lineage to Time data point
type LTTData struct {
	X float64 // Time or Mutations
	Y int     // Number of lineages
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
