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
	"github.com/fredericlemoine/gotree/tree"
)

/*
Generic struct to draw on different supports:
 * ascii in terminal
 * png
 * svg
*/
type TreeDrawer interface {
	DrawHLine(x1, x2, y, maxlength, maxheight float64)
	DrawVLine(x, y1, y2, maxlength, maxheight float64)
	DrawLine(x1, y1, x2, y2, maxlength, maxheight float64)
	DrawCurve(centerx, centery float64, middlex, middley float64, radius float64, startAngle, endAngle float64, maxlength, maxheight float64)
	DrawCircle(x, y float64, maxwidth, maxheight float64)
	/* angle : angle of the tip incoming branch */
	DrawName(x, y float64, name string, maxlength, maxheight float64, angle float64)
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
}
