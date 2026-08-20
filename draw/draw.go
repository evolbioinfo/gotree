/*
Package intended to draw phylogenetic trees on different devices :
 - Terminal,
 - Images (svg, png)
 - ...
And with different drawing algorithms. So far, only ASCII form in terminal.
 - Circular
 - Normal
 - Unrooted
*/
package draw

import (
	"github.com/evolbioinfo/gotree/tree"
)

const (
	// metaBaseGap is the pixel distance between a tip and the first metadata circle.
	metaBaseGap = 8.0
	// metaCircleSpacing is the pixel distance between consecutive metadata circles.
	metaCircleSpacing = 11.0
	// metaLabelExtraGap is added after the last metadata circle before the tip label starts.
	metaLabelExtraGap = 5.0
)

/*
Generic struct to draw on different supports:
 * ascii in terminal
 * png
 * svg
*/
type TreeDrawer interface {
	SetMaxValues(maxObjectWidth, maxObjectHeight float64, maxNameLength, maxNameHeight int)
	DrawHLine(x1, x2, y float64)
	DrawVLine(x, y1, y float64)
	DrawLine(x1, y1, x2, y2 float64)
	DrawCurve(centerx, centery float64, middlex, middley float64, radius float64, startAngle, endAngle float64)
	DrawCircle(x, y float64)
	/* angle : angle of the tip incoming branch. offsetPixels : distance
	   (in pixels) from (x,y) along angle at which the circle is centered.
	   filled : colored fill + black stroke if true, else unfilled with a
	   grey stroke (used for missing metadata values). */
	DrawColoredCircleAtOffset(x, y float64, angle, offsetPixels float64, r, g, b, a uint8, filled bool)
	/* Sets the pixel offset at which tip names start being drawn (to make
	   room for metadata circles drawn closer to the tip). */
	SetTipLabelOffset(px float64)
	/* angle : angle of the tip incoming branch */
	DrawName(x, y float64, name string, angle float64)
	Write()
	Bounds() (int, int) /* width, height*/
}

/*
Generic struct that represents tree layout:
 * circular
 * normal
 * unrooted
*/
type TreeLayout interface {
	DrawTree(t *tree.Tree) error
	SetSupportCutoff(float64)
	SetDisplayInternalNodes(bool)
	SetDisplayNodeComments(bool)
	SetTipMetadata(fields []string, values map[string][]TipMetaColor)
}

func maxLength(t *tree.Tree, hasBranchLengths, hasTipNames, hasNodeComments bool) (float64, int) {
	maxlength := 0.0
	curlength := 0.0
	maxname := 0
	root := t.Root()
	maxLengthRecur(root, nil, curlength, &maxlength, &maxname, hasBranchLengths, hasTipNames, hasNodeComments)
	return maxlength, maxname
}

func maxLengthRecur(n *tree.Node, prev *tree.Node, curlength float64, maxlength *float64, maxname *int, hasBranchLengths, hasTipNames, hasNodeComments bool) {
	if curlength > *maxlength {
		*maxlength = curlength
	}
	if n.Tip() {
		if hasTipNames && hasNodeComments {
			if len(n.Name()+n.CommentsString()) > *maxname {
				*maxname = len(n.Name() + n.CommentsString())
			}
		} else if hasTipNames {
			if len(n.Name()) > *maxname {
				*maxname = len(n.Name())
			}
		} else if hasNodeComments {
			if len(n.CommentsString()) > *maxname {
				*maxname = len(n.CommentsString())
			}
		}
	}
	for i, child := range n.Neigh() {
		if child != prev {
			brlen := n.Edges()[i].Length()
			if brlen == tree.NIL_LENGTH || !hasBranchLengths {
				brlen = 1.0
			}
			maxLengthRecur(child, n, curlength+brlen, maxlength, maxname, hasBranchLengths, hasTipNames, hasNodeComments)
		}
	}
}
