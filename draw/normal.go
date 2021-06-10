package draw

import (
	"github.com/evolbioinfo/gotree/tree"
)

type normalLayout struct {
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

func NewNormalLayout(td TreeDrawer, withBranchLengths, withTipLabels, withInternalNodeLabel, withSupportCircles bool) TreeLayout {
	return &normalLayout{
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

func (layout *normalLayout) SetSupportCutoff(c float64) {
	layout.supportCutoff = c
}

func (layout *normalLayout) SetDisplayInternalNodes(s bool) {
	layout.hasInternalNodeSymbols = s
}

func (layout *normalLayout) SetDisplayNodeComments(s bool) {
	layout.hasNodeComments = s
}

/*
Draw the tree on the specific drawer. Does not close the file. The caller must do it.
*/
func (layout *normalLayout) DrawTree(t *tree.Tree) error {
	var err error = nil
	root := t.Root()
	ntips := len(t.Tips())
	curNbTips := 0
	maxLength, maxName := maxLength(t, layout.hasBranchLengths, layout.hasTipLabels, layout.hasNodeComments)
	layout.drawer.SetMaxValues(maxLength, float64(ntips), maxName, 0)
	layout.drawTreeRecur(root, nil, tree.NIL_SUPPORT, 0, 0, &curNbTips)
	layout.drawTree()
	layout.drawer.Write()
	return err
}

/*
Recursive function that draws the tree. Returns the yposition of the current node
*/
func (layout *normalLayout) drawTreeRecur(n *tree.Node, prev *tree.Node, support, prevDistToRoot, distToRoot float64, curtip *int) float64 {
	ypos := 0.0
	nbchild := 0.0
	if n.Tip() {
		ypos = float64(*curtip)
		nbchild = 1.0
		if layout.hasTipLabels {
			node := &layoutPoint{distToRoot, ypos, 0.0, n.Name(), n.CommentsString()}
			layout.cache.tipLabelPoints = append(layout.cache.tipLabelPoints, node)
		}
		*curtip++
	} else {
		minpos := -1.0
		maxpos := -1.0
		for i, child := range n.Neigh() {
			if child != prev {
				len := n.Edges()[i].Length()
				supp := n.Edges()[i].Support()
				if !layout.hasBranchLengths || len == tree.NIL_LENGTH {
					len = 1.0
				}
				temppos := layout.drawTreeRecur(child, n, supp, distToRoot, distToRoot+len, curtip)
				if minpos == -1 || minpos > temppos {
					minpos = temppos
				}
				if maxpos == -1 || maxpos < temppos {
					maxpos = temppos
				}
				ypos += temppos
				nbchild += 1.0
			}
		}
		ypos /= nbchild
		line := &layoutVLine{distToRoot, minpos, maxpos, tree.NIL_SUPPORT}
		layout.cache.verticalPaths = append(layout.cache.verticalPaths, line)

		inode := &layoutPoint{distToRoot, ypos, 0.0, n.Name(), n.CommentsString()}
		layout.cache.nodePoints = append(layout.cache.nodePoints, inode)
	}

	line := &layoutHLine{prevDistToRoot, distToRoot, ypos, support}
	layout.cache.horizontalPaths = append(layout.cache.horizontalPaths, line)
	return ypos
}

func (layout *normalLayout) drawTree() {
	for _, l := range layout.cache.horizontalPaths {
		layout.drawer.DrawHLine(l.x1, l.x2, l.y)
	}
	for _, l := range layout.cache.verticalPaths {
		layout.drawer.DrawVLine(l.x, l.y1, l.y2)
	}
	if layout.hasTipLabels {
		for _, p := range layout.cache.tipLabelPoints {
			if layout.hasNodeComments {
				layout.drawer.DrawName(p.x, p.y, p.name+p.comment, 0.0)
			} else {
				layout.drawer.DrawName(p.x, p.y, p.name, 0.0)
			}
		}
	}
	if layout.hasInternalNodeLabels {
		for _, p := range layout.cache.nodePoints {
			layout.drawer.DrawName(p.x, p.y, p.name, 0.0)
		}
	} else if layout.hasNodeComments {
		for _, p := range layout.cache.nodePoints {
			layout.drawer.DrawName(p.x, p.y, p.comment, 0.0)
		}
	}

	if layout.hasInternalNodeSymbols {
		for _, p := range layout.cache.nodePoints {
			layout.drawer.DrawCircle(p.x, p.y)
		}
	}
	for _, l := range layout.cache.horizontalPaths {
		middlex := (l.x1 + l.x2) / 2.0
		middley := (l.y + l.y) / 2.0
		if layout.hasSupport && l.support != tree.NIL_SUPPORT && l.support >= layout.supportCutoff {
			layout.drawer.DrawCircle(middlex, middley)
		}
	}
}
