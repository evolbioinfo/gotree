package draw

import (
	"math"

	"github.com/fredericlemoine/gotree/tree"
)

type radialLayout struct {
	drawer           TreeDrawer
	spread           float64
	hasBranchLengths bool
	hasTipLabels     bool
	cache            *radialCache
}

type radialPoint struct {
	x       float64
	y       float64
	brAngle float64 // Angle of the incoming branch
	name    string
}

type radialLine struct {
	p1 *radialPoint
	p2 *radialPoint
}

type radialCache struct {
	tipLabelPoints []*radialPoint
	branchPaths    []*radialLine
	nodePoints     []*radialPoint
}

func newRadialCache() *radialCache {
	return &radialCache{
		make([]*radialPoint, 0),
		make([]*radialLine, 0),
		make([]*radialPoint, 0),
	}

}

func NewRadialLayout(td TreeDrawer, withBranchLengths, withTipLabels bool) TreeLayout {
	return &radialLayout{
		td,
		0.0,
		withBranchLengths,
		withTipLabels,
		newRadialCache(),
	}
}

/*
Draw the tree on the specific drawer. Does not close the file. The caller must do it.
This layout is an adaptation in Go of the figtree radial layout : figtree/treeviewer/treelayouts/RadialTreeLayout.java
( https://github.com/rambaut/figtree/ )
*/
func (layout *radialLayout) DrawTree(t *tree.Tree) error {
	root := t.Root()
	layout.spread = 0.0
	layout.constructNode(t, root, nil, 0.0, math.Pi*2, 0.0, 0.0, 0.0)
	layout.drawTree()
	layout.drawer.Write()
	return nil
}

func (layout *radialLayout) constructNode(t *tree.Tree, node *tree.Node, prev *tree.Node, angleStart, angleFinish, xPosition, yPosition, length float64) *radialPoint {
	branchAngle := (angleStart + angleFinish) / 2.0
	directionX := math.Cos(branchAngle)
	directionY := math.Sin(branchAngle)

	nodePoint := &radialPoint{xPosition + (length * directionX), yPosition + (length * directionY), branchAngle, node.Name()}

	if !node.Tip() {
		leafCounts := make([]int, 0)
		sumLeafCount := 0
		i := 0
		for num, child := range node.Neigh() {
			if child != prev {
				numT := int(node.Edges()[num].NumTips())
				leafCounts = append(leafCounts, numT)
				sumLeafCount += numT
				i++
			}
		}
		span := (angleFinish - angleStart)
		if node != t.Root() {
			span *= 1.0 + (layout.spread / 10.0)
			angleStart = branchAngle - (span / 2.0)
			angleFinish = branchAngle + (span / 2.0)
		}
		a2 := angleStart
		rotate := false
		i = 0
		for num, child := range node.Neigh() {
			if child != prev {
				index := i
				if rotate {
					index = len(node.Neigh()) - i - 1
				}
				brLen := node.Edges()[num].Length()
				if !layout.hasBranchLengths {
					brLen = 1.0
				}
				a1 := a2
				a2 = a1 + (span * float64(leafCounts[index]) / float64(sumLeafCount))
				childPoint := layout.constructNode(t, child, node, a1, a2, nodePoint.x, nodePoint.y, brLen)
				branchLine := &radialLine{childPoint, nodePoint}
				//add the branchLine to the map of branch paths
				layout.cache.branchPaths = append(layout.cache.branchPaths, branchLine)
				i++
			}
		}
	} else {
		layout.cache.tipLabelPoints = append(layout.cache.tipLabelPoints, nodePoint)
	}
	//layout.cache.nodePoints = append(layout.cache.nodePoints, nodePoints)
	return nodePoint
}

func (layout *radialLayout) drawTree() {
	xmin, ymin, xmax, ymax := layout.borders()
	xoffset := 0.0
	if xmin < 0 {
		xoffset = -xmin
	}
	yoffset := 0.0
	if ymin < 0 {
		yoffset = -ymin
	}

	for _, l := range layout.cache.branchPaths {
		layout.drawer.DrawLine(l.p1.x+xoffset, l.p1.y+yoffset, l.p2.x+xoffset, l.p2.y+yoffset, xmax+xoffset, ymax+yoffset)
	}
	if layout.hasTipLabels {
		for _, p := range layout.cache.tipLabelPoints {
			layout.drawer.DrawName(p.x+xoffset, p.y+yoffset, p.name, xmax+xoffset, ymax+yoffset, p.brAngle)
		}
	}
}

func (layout *radialLayout) borders() (xmin, ymin, xmax, ymax float64) {
	xmin, ymin = 100000.0, 100000.0
	xmax, ymax = 0.0, 0.0
	for _, line := range layout.cache.branchPaths {
		xmin = math.Min(math.Min(xmin, line.p1.x), line.p2.x)
		ymin = math.Min(math.Min(ymin, line.p1.y), line.p2.y)
		xmax = math.Max(math.Max(xmax, line.p1.x), line.p2.x)
		ymax = math.Max(math.Max(ymax, line.p1.y), line.p2.y)
	}
	return
}
