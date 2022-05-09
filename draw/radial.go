package draw

import (
	"math"

	"github.com/evolbioinfo/gotree/tree"
)

type radialLayout struct {
	drawer                 TreeDrawer
	spread                 float64
	hasBranchLengths       bool
	hasTipLabels           bool
	hasTipSymbols          bool
	hasInternalNodeLabels  bool
	hasInternalNodeSymbols bool
	hasNodeComments        bool
	hasSupport             bool
	supportCutoff          float64
	cache                  *layoutCache
	tipColors              map[string][]uint8
}

func NewRadialLayout(td TreeDrawer, withBranchLengths, withTipLabels, withInternalNodeLabels, withSuppportCircles bool) TreeLayout {
	return &radialLayout{
		td,
		0.0,
		withBranchLengths,
		withTipLabels,
		false,
		withInternalNodeLabels,
		false,
		false,
		withSuppportCircles,
		0.7,
		newLayoutCache(),
		make(map[string][]uint8),
	}
}

func (layout *radialLayout) SetSupportCutoff(c float64) {
	layout.supportCutoff = c
}

func (layout *radialLayout) SetDisplayInternalNodes(s bool) {
	layout.hasInternalNodeSymbols = s
}
func (layout *radialLayout) SetDisplayNodeComments(s bool) {
	layout.hasNodeComments = s
}

func (layout *radialLayout) SetTipColors(colors map[string][]uint8) {
	layout.hasTipSymbols = true
	layout.tipColors = colors
}

/*
Draw the tree on the specific drawer. Does not close the file. The caller must do it.
This layout is an adaptation in Go of the figtree radial layout : figtree/treeviewer/treelayouts/RadialTreeLayout.java
( https://github.com/rambaut/figtree/ )
Tree indexes must have been set with t.ReinitIndexes()
*/
func (layout *radialLayout) DrawTree(t *tree.Tree) error {
	root := t.Root()
	layout.spread = 0.0
	layout.constructNode(t, root, nil, 0.0, 0.0, math.Pi*2, 0.0, 0.0, 0.0)
	_, maxNameLength := maxLength(t, layout.hasBranchLengths, layout.hasTipLabels, layout.hasNodeComments)
	layout.drawTree(maxNameLength)
	layout.drawer.Write()
	return nil
}

func (layout *radialLayout) constructNode(t *tree.Tree, node *tree.Node, prev *tree.Node, support, angleStart, angleFinish, xPosition, yPosition, length float64) *layoutPoint {
	branchAngle := (angleStart + angleFinish) / 2.0
	directionX := math.Cos(branchAngle)
	directionY := math.Sin(branchAngle)

	nodePoint := &layoutPoint{xPosition + (length * directionX), yPosition + (length * directionY), branchAngle, node.Name(), node.CommentsString()}

	if !node.Tip() {
		leafCounts := make([]int, 0)
		sumLeafCount := 0
		i := 0
		for num, child := range node.Neigh() {
			if child != prev {

				numT := node.Edges()[num].NumTipsRight()
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
				supp := node.Edges()[num].Support()

				if !layout.hasBranchLengths || brLen == tree.NIL_LENGTH {
					brLen = 1.0
				}
				a1 := a2
				a2 = a1 + (span * float64(leafCounts[index]) / float64(sumLeafCount))
				childPoint := layout.constructNode(t, child, node, supp, a1, a2, nodePoint.x, nodePoint.y, brLen)
				branchLine := &layoutLine{childPoint, nodePoint, supp}
				//add the branchLine to the map of branch paths
				layout.cache.branchPaths = append(layout.cache.branchPaths, branchLine)
				i++
			}
		}
		layout.cache.nodePoints = append(layout.cache.nodePoints, nodePoint)
	} else {
		layout.cache.tipLabelPoints = append(layout.cache.tipLabelPoints, nodePoint)
	}
	return nodePoint
}

func (layout *radialLayout) drawTree(maxNameLength int) {
	xmin, ymin, xmax, ymax := layout.cache.borders()
	xoffset := 0.0
	if xmin < 0 {
		xoffset = -xmin
	}
	yoffset := 0.0
	if ymin < 0 {
		yoffset = -ymin
	}

	layout.drawer.SetMaxValues(xmax+xoffset, ymax+yoffset, maxNameLength, maxNameLength)

	for _, l := range layout.cache.branchPaths {
		layout.drawer.DrawLine(l.p1.x+xoffset, l.p1.y+yoffset, l.p2.x+xoffset, l.p2.y+yoffset)
	}
	if layout.hasTipLabels {
		for _, p := range layout.cache.tipLabelPoints {
			// Add space to label so it's not hidden by node symbol
			// There is probably a better way to do this
			spc := ""
			if layout.hasTipSymbols {
				if _, ok := layout.tipColors[p.name]; ok {
					spc = "  "
				}
			}
			if layout.hasNodeComments {
				layout.drawer.DrawName(p.x+xoffset, p.y+yoffset, spc+p.name+p.comment+spc, p.brAngle)
			} else {
				layout.drawer.DrawName(p.x+xoffset, p.y+yoffset, spc+p.name+spc, p.brAngle)
			}
		}
	}

	if layout.hasInternalNodeLabels {
		for _, p := range layout.cache.nodePoints {
			layout.drawer.DrawName(p.x+xoffset, p.y+yoffset, p.name, p.brAngle)
		}
	} else if layout.hasNodeComments {
		for _, p := range layout.cache.nodePoints {
			layout.drawer.DrawName(p.x+xoffset, p.y+yoffset, p.comment, p.brAngle)
		}
	}

	if layout.hasInternalNodeSymbols {
		for _, p := range layout.cache.nodePoints {
			layout.drawer.DrawCircle(p.x+xoffset, p.y+yoffset)
		}
	}

	if layout.hasTipSymbols {
		for _, p := range layout.cache.tipLabelPoints {
			if col, ok := layout.tipColors[p.name]; ok {
				layout.drawer.DrawColoredCircle(p.x+xoffset, p.y+yoffset, col[0], col[1], col[2], 0xff)
			}
		}
	}

	for _, l := range layout.cache.branchPaths {
		middlex := (l.p1.x + l.p2.x + 2*xoffset) / 2.0
		middley := (l.p1.y + l.p2.y + 2*yoffset) / 2.0
		if layout.hasSupport && l.support != tree.NIL_SUPPORT && l.support >= layout.supportCutoff {
			layout.drawer.DrawCircle(middlex, middley)
		}
	}
}
