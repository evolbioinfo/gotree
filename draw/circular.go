package draw

import (
	"log"
	"math"

	"github.com/evolbioinfo/gotree/tree"
)

type circularLayout struct {
	drawer                 TreeDrawer
	hasBranchLengths       bool
	hasTipLabels           bool
	hasInternalNodeLabels  bool
	hasInternalNodeSymbols bool
	hasNodeComments        bool
	hasSupport             bool
	supportCutoff          float64
	cache                  *layoutCache
}

/*
If withSuppportCircles is true, then it will draw circles on branches whose support is > 0.7. The cutoff may be set with layout.SetSupportCutoff()
*/
func NewCircularLayout(td TreeDrawer, withBranchLengths, withTipLabels, withInternalNodeLabel, withSupportCircles bool) TreeLayout {
	w, h := td.Bounds()
	if w != h {
		log.Print("Width!=Height : This is not advised with circular layout")
	}
	return &circularLayout{
		td,
		withBranchLengths,
		withTipLabels,
		withInternalNodeLabel,
		false,
		false,
		withSupportCircles,
		0.7,
		newLayoutCache(),
	}
}

func (layout *circularLayout) SetSupportCutoff(c float64) {
	layout.supportCutoff = c
}

func (layout *circularLayout) SetDisplayInternalNodes(s bool) {
	layout.hasInternalNodeSymbols = s
}

func (layout *circularLayout) SetDisplayNodeComments(s bool) {
	layout.hasNodeComments = s
}

/*
Draw the tree on the specific drawer. Does not close the file. The caller must do it.
*/
func (layout *circularLayout) DrawTree(t *tree.Tree) error {
	var err error = nil
	root := t.Root()
	ntips := len(t.Tips())
	curNbTips := 0
	maxLength := layout.maxLength(t)
	layout.drawTreeRecur(root, nil, tree.NIL_SUPPORT, 0, 0, maxLength, &curNbTips, ntips)
	layout.drawTree()
	layout.drawer.Write()
	return err
}

/*
Recursive function that draws the tree. Returns the angle of the current node
*/
func (layout *circularLayout) drawTreeRecur(n *tree.Node, prev *tree.Node, support, prevDistToRoot, distToRoot float64, maxLength float64, curtip *int, nbtips int) float64 {
	angle := 0.0
	if n.Tip() {
		angle = float64(*curtip)*2*math.Pi/float64(nbtips) + math.Pi/2
		x3 := distToRoot * math.Cos(angle)
		y3 := distToRoot * math.Sin(angle)
		node := &layoutPoint{x3, y3, angle, n.Name(), n.CommentsString()}
		layout.cache.tipLabelPoints = append(layout.cache.tipLabelPoints, node)
		*curtip++
	} else {
		minangle := -1.0
		maxangle := -1.0
		for i, child := range n.Neigh() {
			if child != prev {
				len := n.Edges()[i].Length()
				supp := n.Edges()[i].Support()
				if !layout.hasBranchLengths || len == tree.NIL_LENGTH {
					len = 1.0
				}
				tempangle := layout.drawTreeRecur(child, n, supp, distToRoot, distToRoot+len, maxLength, curtip, nbtips)
				if minangle == -1 || minangle > tempangle {
					minangle = tempangle
				}
				if maxangle == -1 || maxangle < tempangle {
					maxangle = tempangle
				}
			}
		}
		angle = (minangle + maxangle) / 2.0

		x4 := distToRoot * math.Cos(angle)
		y4 := distToRoot * math.Sin(angle)
		inode := &layoutPoint{x4, y4, angle, n.Name(), n.CommentsString()}
		layout.cache.nodePoints = append(layout.cache.nodePoints, inode)
		curve := &layoutCurve{&layoutPoint{0, 0, 0.0, "", ""}, inode, distToRoot, minangle, maxangle}
		layout.cache.curvePaths = append(layout.cache.curvePaths, curve)
	}
	x1 := prevDistToRoot * math.Cos(angle)
	y1 := prevDistToRoot * math.Sin(angle)
	x2 := distToRoot * math.Cos(angle)
	y2 := distToRoot * math.Sin(angle)
	line := &layoutLine{&layoutPoint{x1, y1, angle, "", ""}, &layoutPoint{x2, y2, angle, "", ""}, support}
	layout.cache.branchPaths = append(layout.cache.branchPaths, line)
	return angle
}

func (layout *circularLayout) maxLength(t *tree.Tree) float64 {
	maxlength := 0.0
	curlength := 0.0
	root := t.Root()
	layout.maxLengthRecur(root, nil, curlength, &maxlength)
	return maxlength
}

func (layout *circularLayout) maxLengthRecur(n *tree.Node, prev *tree.Node, curlength float64, maxlength *float64) {
	if curlength > *maxlength {
		*maxlength = curlength
	}
	for i, child := range n.Neigh() {
		if child != prev {
			brlen := n.Edges()[i].Length()
			if brlen == tree.NIL_LENGTH || !layout.hasBranchLengths {
				brlen = 1.0
			}
			layout.maxLengthRecur(child, n, curlength+brlen, maxlength)
		}
	}
}

func (layout *circularLayout) drawTree() {
	xmin, ymin, xmax, ymax := layout.cache.borders()
	xoffset := 0.0
	if xmin < 0 {
		xoffset = -xmin
	}
	yoffset := 0.0
	if ymin < 0 {
		yoffset = -ymin
	}

	max := math.Max(xmax+xoffset, ymax+yoffset)

	for _, l := range layout.cache.branchPaths {
		layout.drawer.DrawLine(l.p1.x+xoffset, l.p1.y+yoffset, l.p2.x+xoffset, l.p2.y+yoffset, max, max)
	}
	for _, c := range layout.cache.curvePaths {
		layout.drawer.DrawCurve(c.center.x+xoffset, c.center.y+yoffset, c.middlepoint.x+xoffset, c.middlepoint.y+yoffset, c.radius, c.startAngle, c.endAngle, max, max)
	}

	if layout.hasTipLabels {
		for _, p := range layout.cache.tipLabelPoints {
			if layout.hasNodeComments {
				layout.drawer.DrawName(p.x+xoffset, p.y+yoffset, p.name+p.comment, max, max, p.brAngle)
			} else {
				layout.drawer.DrawName(p.x+xoffset, p.y+yoffset, p.name, max, max, p.brAngle)
			}
		}
	}
	if layout.hasInternalNodeLabels {
		for _, p := range layout.cache.nodePoints {
			layout.drawer.DrawName(p.x+xoffset, p.y+yoffset, p.name, max, max, p.brAngle)
		}
	} else if layout.hasNodeComments {
		for _, p := range layout.cache.nodePoints {
			layout.drawer.DrawName(p.x+xoffset, p.y+yoffset, p.comment, max, max, p.brAngle)
		}
	}

	if layout.hasInternalNodeSymbols {
		for _, p := range layout.cache.nodePoints {
			layout.drawer.DrawCircle(p.x+xoffset, p.y+yoffset, max, max)
		}
	}
	for _, l := range layout.cache.branchPaths {
		middlex := (l.p1.x + l.p2.x + 2*xoffset) / 2.0
		middley := (l.p1.y + l.p2.y + 2*yoffset) / 2.0
		if layout.hasSupport && l.support != tree.NIL_SUPPORT && l.support >= layout.supportCutoff {
			layout.drawer.DrawCircle(middlex, middley, max, max)
		}
	}
}
