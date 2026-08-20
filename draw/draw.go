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
	"math"

	"github.com/evolbioinfo/gotree/tree"
)

const (
	// metaBaseGap is the pixel distance between a tip and the first metadata marker.
	metaBaseGap = 8.0
	// metaCircleSpacing is the pixel distance between consecutive metadata markers.
	metaCircleSpacing = 11.0
	// metaLabelExtraGap is added after the last metadata marker before the tip label starts.
	metaLabelExtraGap = 5.0

	// legendPadding is the padding (px) between the legend box border and its content.
	legendPadding = 6.0
	// legendRowHeight is the vertical space (px) allotted to each legend row.
	legendRowHeight = 13.0
	// legendSwatchGap is the horizontal space (px) reserved for a row's marker before its label starts.
	legendSwatchGap = 14.0
	// legendGap is the pixel gap left between the tree area and the legend band reserved below it.
	legendGap = 10.0
	// legendColumnGap is the horizontal pixel gap between two field columns in the legend.
	legendColumnGap = 10.0
)

// metaLabelOffset returns the pixel offset at which the tip label should
// start, to leave room for nFields metadata markers (0 if there are none).
func metaLabelOffset(nFields int) float64 {
	if nFields == 0 {
		return 0
	}
	return metaBaseGap + metaCircleSpacing*float64(nFields) + metaLabelExtraGap
}

// metaExtraNameChars converts the metadata marker row's pixel width into an
// equivalent number of "name characters", so it can inflate the
// character-based tip-label margin passed to TreeDrawer.SetMaxValues
// (which reserves 5px/char) and keep both markers and label text inside
// the image bounds.
func metaExtraNameChars(nFields int) int {
	return int(math.Ceil(metaLabelOffset(nFields) / 5.0))
}

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
	   (in pixels) from (x,y) along angle at which the shape is centered.
	   filled : colored fill + black stroke if true, else unfilled with a
	   grey stroke (used for missing metadata values). */
	DrawColoredShapeAtOffset(x, y float64, angle, offsetPixels float64, shape Shape, r, g, b, a uint8, filled bool)
	/* Sets the pixel offset at which tip names start being drawn (to make
	   room for metadata circles drawn closer to the tip). */
	SetTipLabelOffset(px float64)
	/* angle : angle of the tip incoming branch */
	DrawName(x, y float64, name string, angle float64)
	/* Draws a legend box (field name, marker shape, and color-coded value
	   labels for each metadata field) anchored to a corner of the image,
	   in absolute pixel coordinates independent of the tree's data-space
	   transform. No-op if entries is empty. */
	DrawLegend(entries []LegendEntry)
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
	SetTipMetadata(fields []string, shapes []Shape, values map[string][]TipMetaColor, legend []LegendEntry)
}

// metaShapeAt returns shapes[i], defaulting to ShapeCircle if shapes is too short.
func metaShapeAt(shapes []Shape, i int) Shape {
	if i < len(shapes) {
		return shapes[i]
	}
	return ShapeCircle
}

// legendRow is one line of a rendered legend: either a field-name header
// (isHeader, no swatch) or a color-coded value row (hasSwatch).
type legendRow struct {
	text       string
	isHeader   bool
	hasSwatch  bool
	shape      Shape
	r, g, b, a uint8
}

// legendColumns lays out legend entries as one column per metadata field
// (landscape legend): each column is a header row (the field name) followed
// by that field's value rows (marker + label), plus a "..." row when the
// field's values were truncated (see maxLegendDiscreteValues). Columns are
// meant to be drawn left to right, each internally top to bottom.
func legendColumns(entries []LegendEntry) [][]legendRow {
	cols := make([][]legendRow, 0, len(entries))
	for _, e := range entries {
		col := make([]legendRow, 0, len(e.Values)+2)
		col = append(col, legendRow{text: e.Field, isHeader: true})
		for _, v := range e.Values {
			col = append(col, legendRow{text: v.Label, hasSwatch: true, shape: e.Shape, r: v.R, g: v.G, b: v.B, a: v.A})
		}
		if e.Truncated {
			col = append(col, legendRow{text: "..."})
		}
		cols = append(cols, col)
	}
	return cols
}

// TextWidthFunc measures the rendered pixel width of text in a legend row
// (bold is true for a field-name header row), for a specific TreeDrawer's
// font/size. Implementations: SvgTextWidth (estimate, no real font metrics
// available for SVG output) and PngTextWidth (exact, via draw2d font
// metrics, the same font pngTreeDrawer draws with).
type TextWidthFunc func(text string, bold bool) float64

// legendColumnWidth returns a column's rendered content width in pixels
// (marker + label text, or header text), excluding padding/gaps: for each
// row, the swatch (if any) plus the row's measured text width, maxed
// across all rows in the column.
func legendColumnWidth(col []legendRow, widthFn TextWidthFunc) float64 {
	max := 0.0
	for _, r := range col {
		w := widthFn(r.text, r.isHeader)
		if r.hasSwatch {
			w += legendSwatchGap
		}
		if w > max {
			max = w
		}
	}
	return max
}

// LegendSize returns the pixel (width, height) of the legend box that
// DrawLegend would draw for entries, or (0, 0) if entries is empty.
// Callers that want the image canvas to grow to fit the legend (rather
// than have it overlaid on the tree) should compute this before
// constructing a TreeDrawer and reserve the returned space, e.g. via
// NewSvgTreeDrawer/NewPngTreeDrawer's legendWidth/legendHeight parameters.
// widthFn must measure text the same way the eventual TreeDrawer's
// DrawLegend will (SvgTextWidth for an svg drawer, PngTextWidth for a png
// one), or the reserved space and the actual rendering will disagree.
func LegendSize(entries []LegendEntry, widthFn TextWidthFunc) (width, height int) {
	cols := legendColumns(entries)
	if len(cols) == 0 {
		return 0, 0
	}
	w := 2 * legendPadding
	maxRows := 0
	for i, col := range cols {
		if i > 0 {
			w += legendColumnGap
		}
		w += legendColumnWidth(col, widthFn)
		if len(col) > maxRows {
			maxRows = len(col)
		}
	}
	h := float64(maxRows)*legendRowHeight + 2*legendPadding
	return int(math.Ceil(w)), int(math.Ceil(h))
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
