package draw

import (
	"github.com/fredericlemoine/gotree/tree"
)

type normalLayout struct {
	drawer                TreeDrawer
	hasBranchLengths      bool
	hasTipLabels          bool
	hasInternalNodeLabels bool
}

func NewNormalLayout(td TreeDrawer, withBranchLengths, withTipLabels, withInternalNodeLabel bool) TreeLayout {
	return &normalLayout{
		td,
		withBranchLengths,
		withTipLabels,
		withInternalNodeLabel,
	}
}

/*
Draw the tree on the specific drawer. Does not close the file. The caller must do it.
*/
func (layout *normalLayout) DrawTree(t *tree.Tree) error {
	var err error = nil
	root := t.Root()
	ntips := len(t.Tips())
	curNbTips := 0
	maxLength := layout.maxLength(t)
	layout.drawTreeRecur(root, nil, 0, 0, maxLength, &curNbTips, ntips, maxLength)
	layout.drawer.Write()
	return err
}

/*
Recursive function that draws the tree. Returns the yposition of the current node
*/
func (layout *normalLayout) drawTreeRecur(n *tree.Node, prev *tree.Node, prevDistToRoot, distToRoot float64, maxLength float64, curtip *int, nbtips int, maxlength float64) float64 {
	ypos := 0.0
	nbchild := 0.0
	if n.Tip() {
		ypos = float64(*curtip)
		nbchild = 1.0
		if layout.hasTipLabels {
			layout.drawer.DrawName(distToRoot, ypos, n.Name(), maxlength, float64(nbtips), 0.0)
		}
		*curtip++
	} else {
		minpos := -1.0
		maxpos := -1.0
		for i, child := range n.Neigh() {
			if child != prev {
				len := n.Edges()[i].Length()
				if !layout.hasBranchLengths || len == tree.NIL_LENGTH {
					len = 1.0
				}
				temppos := layout.drawTreeRecur(child, n, distToRoot, distToRoot+len, maxLength, curtip, nbtips, maxlength)
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
		layout.drawer.DrawVLine(distToRoot, minpos, maxpos, maxlength, float64(nbtips))
		if layout.hasInternalNodeLabels {
			layout.drawer.DrawName(distToRoot, ypos, n.Name(), maxlength, float64(nbtips), 0.0)
		}
	}
	layout.drawer.DrawHLine(prevDistToRoot, distToRoot, ypos, maxlength, float64(nbtips))
	return ypos
}

func (layout *normalLayout) maxLength(t *tree.Tree) float64 {
	maxlength := 0.0
	curlength := 0.0
	root := t.Root()
	layout.maxLengthRecur(root, nil, curlength, &maxlength)
	return maxlength
}

func (layout *normalLayout) maxLengthRecur(n *tree.Node, prev *tree.Node, curlength float64, maxlength *float64) {
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
