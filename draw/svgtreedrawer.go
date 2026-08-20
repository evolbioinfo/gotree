package draw

import (
	"fmt"
	"io"
	"math"

	svg "github.com/ajstarks/svgo"
)

/*
TextTreeDrawer initializer. TextTreeDraws draws tree as ASCII on stdout or any file.
So far: Does not take into account branch lengths.
*/
func NewSvgTreeDrawer(w io.Writer, width, height int,
	leftmargin, rightmargin, topmargin, bottommargin int) TreeDrawer {
	svgtd := &svgTreeDrawer{
		w,
		width,
		height,
		leftmargin,
		rightmargin,
		topmargin,
		bottommargin,
		nil,
		1.0,
		0.0,
		0.0,
		0.0,
		0.0,
	}
	svgtd.canvas = svg.New(w)
	svgtd.canvas.Start(width+leftmargin+rightmargin, height+topmargin+bottommargin)
	return svgtd
}

func (svgtd *svgTreeDrawer) SetMaxValues(maxLength, maxHeight float64, maxNameLength, maxNameHeight int) {
	svgtd.maxLength = maxLength
	svgtd.maxHeight = maxHeight
	svgtd.maxNameLength = 5 * maxNameLength
	svgtd.maxNameHeight = 5 * maxNameHeight
}

/*
Draw a tree in a svg file.
*/
type svgTreeDrawer struct {
	outwriter     io.Writer // Output file
	width         int       // Width of the ascii canvas
	height        int       // Height of the ascii canvas
	leftmargin    int       // Left margin of the canvas (in addition to the width)
	rightmargin   int       // Right margin of the canvas (in addition to the width)
	topmargin     int       // Top margin of the canvas (in addition to the height)
	bottommargin  int       // Bottom margin of the canvas (in addition to the height)
	canvas        *svg.SVG  // SVN Canvas
	dTip          float64   // Distance from tip to label
	maxLength     float64   // Maximum length of object to draw (in original scale)
	maxHeight     float64   // Maximum height of object to draw (in original scale)
	maxNameLength int       // Maximum length of species names / horitzontal
	maxNameHeight int       // Maximum length of species names / vertical
}

func (svgtd *svgTreeDrawer) DrawHLine(x1, x2, y float64) {
	min := int(float64(svgtd.width-svgtd.maxNameLength)*x1/svgtd.maxLength + float64(svgtd.leftmargin))
	max := int(float64(svgtd.width-svgtd.maxNameLength)*x2/svgtd.maxLength + float64(svgtd.leftmargin))
	ypos := int(float64(svgtd.height-svgtd.maxNameHeight)*y/svgtd.maxHeight + float64(svgtd.topmargin))
	svgtd.canvas.Line(min, ypos, max, ypos, "stroke-width:2; fill:black; stroke: black;")
}

func (svgtd *svgTreeDrawer) DrawVLine(x, y1, y2 float64) {
	min := int(float64(svgtd.height-svgtd.maxNameHeight)*y1/svgtd.maxHeight + float64(svgtd.topmargin))
	max := int(float64(svgtd.height-svgtd.maxNameHeight)*y2/svgtd.maxHeight + float64(svgtd.topmargin))
	xpos := int(float64(svgtd.width-svgtd.maxNameLength)*x/svgtd.maxLength + float64(svgtd.leftmargin))
	svgtd.canvas.Line(xpos, min, xpos, max, "stroke-width:2; fill:black; stroke: black;")
}

func (svgtd *svgTreeDrawer) DrawLine(x1, y1, x2, y2 float64) {
	y1pos := int(float64(svgtd.height-svgtd.maxNameHeight)*y1/svgtd.maxHeight + float64(svgtd.topmargin))
	y2pos := int(float64(svgtd.height-svgtd.maxNameHeight)*y2/svgtd.maxHeight + float64(svgtd.topmargin))
	x1pos := int(float64(svgtd.width-svgtd.maxNameLength)*x1/svgtd.maxLength + float64(svgtd.leftmargin))
	x2pos := int(float64(svgtd.width-svgtd.maxNameLength)*x2/svgtd.maxLength + float64(svgtd.leftmargin))
	svgtd.canvas.Line(x1pos, y1pos, x2pos, y2pos, "stroke-width:2; fill:black; stroke: black;")
}

func (svgtd *svgTreeDrawer) DrawCurve(centerx, centery, middlex, middley float64, radius float64, startAngle, endAngle float64) {
	x1 := (radius*math.Cos(startAngle)+centerx)*float64(svgtd.width-svgtd.maxNameLength)/svgtd.maxLength + float64(svgtd.topmargin)
	y1 := (radius*math.Sin(startAngle)+centery)*float64(svgtd.height-svgtd.maxNameHeight)/svgtd.maxHeight + float64(svgtd.leftmargin)
	x2 := (radius*math.Cos(endAngle)+centerx)*float64(svgtd.width-svgtd.maxNameLength)/svgtd.maxLength + float64(svgtd.topmargin)
	y2 := (radius*math.Sin(endAngle)+centery)*float64(svgtd.height-svgtd.maxNameHeight)/svgtd.maxHeight + float64(svgtd.leftmargin)
	centerx2 := centerx*float64(svgtd.width-svgtd.maxNameLength)/svgtd.maxLength + float64(svgtd.topmargin)
	centery2 := centery*float64(svgtd.height-svgtd.maxNameHeight)/svgtd.maxHeight + float64(svgtd.leftmargin)
	// middlex2 := middlex*float64(svgtd.width)/maxlength + float64(svgtd.topmargin)
	// middley2 := middley*float64(svgtd.height)/maxheight + float64(svgtd.leftmargin)
	radiusscaled := round(math.Sqrt(math.Pow((y2-centery2), 2) + math.Pow((x2-centerx2), 2)))
	largeArcFlag := true
	if endAngle-startAngle < math.Pi {
		largeArcFlag = false
	}
	svgtd.canvas.Arc(round(x1), round(y1), radiusscaled, radiusscaled, 0, largeArcFlag, true, round(x2), round(y2), "stroke-width:2; fill:none;stroke: black;")
}

func (svgtd *svgTreeDrawer) DrawCircle(x, y float64) {
	centerx2 := x*float64(svgtd.width-svgtd.maxNameLength)/svgtd.maxLength + float64(svgtd.topmargin)
	centery2 := y*float64(svgtd.height-svgtd.maxNameHeight)/svgtd.maxHeight + float64(svgtd.leftmargin)
	svgtd.canvas.Circle(round(centerx2), round(centery2), 5, "stroke-width:1; fill:orange;stroke: black;")
}

func (svgtd *svgTreeDrawer) SetTipLabelOffset(px float64) {
	svgtd.dTip = px
}

/* angle : incoming branch angle. offsetPixels : distance from (x,y) along angle */
func (svgtd *svgTreeDrawer) DrawColoredShapeAtOffset(x, y float64, angle, offsetPixels float64, shape Shape, r, g, b, a uint8, filled bool) {
	degree := angle * 180.0 / math.Pi
	ypos := int(float64(svgtd.height-svgtd.maxNameHeight)*y/svgtd.maxHeight + float64(svgtd.topmargin))
	xpos := int(float64(svgtd.width-svgtd.maxNameLength)*x/svgtd.maxLength + float64(svgtd.leftmargin))

	style := svgFillStyle(r, g, b, a, "black")
	if !filled {
		style = "stroke-width:1;fill:none;stroke:#999999;"
	}

	svgtd.canvas.Translate(xpos, ypos)
	if angle < 3*math.Pi/2.0 && angle > math.Pi/2.0 {
		svgtd.canvas.Rotate(degree - 180)
		svgtd.drawShape(-round(offsetPixels), shape, style)
	} else {
		svgtd.canvas.Rotate(degree)
		svgtd.drawShape(round(offsetPixels), shape, style)
	}
	svgtd.canvas.Gend()
	svgtd.canvas.Gend()
}

// drawShape draws shape centered at local (localx, 0), in the current
// (already translated/rotated) canvas coordinate system.
func (svgtd *svgTreeDrawer) drawShape(localx int, shape Shape, style string) {
	switch shape {
	case ShapeSquare:
		svgtd.canvas.Rect(localx-4, -4, 8, 8, style)
	case ShapeTriangle:
		svgtd.canvas.Polygon(
			[]int{localx, localx - 5, localx + 5},
			[]int{-5, 4, 4},
			style,
		)
	case ShapeDiamond:
		svgtd.canvas.Polygon(
			[]int{localx, localx + 5, localx, localx - 5},
			[]int{-5, 0, 5, 0},
			style,
		)
	case ShapeStar:
		xs, ys := starPoints(localx, 0, 5, 2)
		svgtd.canvas.Polygon(xs, ys, style)
	default: // ShapeCircle
		svgtd.canvas.Circle(localx, 0, 4, style)
	}
}

// starPoints returns the 10 vertices (5 outer, 5 inner, alternating) of a
// 5-pointed star centered at (cx,cy), pointing up.
func starPoints(cx, cy int, outerR, innerR float64) (xs, ys []int) {
	xs = make([]int, 10)
	ys = make([]int, 10)
	for i := 0; i < 10; i++ {
		radius := outerR
		if i%2 == 1 {
			radius = innerR
		}
		angle := -math.Pi/2.0 + float64(i)*math.Pi/5.0
		xs[i] = cx + round(radius*math.Cos(angle))
		ys[i] = cy + round(radius*math.Sin(angle))
	}
	return
}

/* angle:  incoming branch angle */
func (svgtd *svgTreeDrawer) DrawName(x, y float64, name string, angle float64) {
	degree := angle * 180.0 / math.Pi
	//left, top, right, bottom := ptd.gc.GetStringBounds(name)
	// Text width: Not very elegant so far...
	ypos := int(float64(svgtd.height-svgtd.maxNameHeight)*y/svgtd.maxHeight + float64(svgtd.topmargin))
	xpos := int(float64(svgtd.width-svgtd.maxNameLength)*x/svgtd.maxLength + float64(svgtd.leftmargin))
	//fmt.Fprintf(os.Stderr, "Tip angle: %f (%f)\n", angle, degree)
	svgtd.canvas.Translate(xpos, ypos)
	// We rotate the other way (text not upside down)
	if angle < 3*math.Pi/2.0 && angle > math.Pi/2.0 {
		svgtd.canvas.Rotate(degree - 180)
		svgtd.canvas.Text(-int(svgtd.dTip), 0, name, "alignment-baseline:middle;text-anchor:end;font-family: sans-serif;font-size:8px;")
	} else {
		svgtd.canvas.Rotate(degree)
		svgtd.canvas.Text(int(svgtd.dTip), 0, name, "alignment-baseline:middle;text-anchor:start;font-family: sans-serif;font-size:8px;")
	}
	//svgtd.canvas.Rect(int(svgtd.dTip), 0, 100, 10)
	svgtd.canvas.Gend()
	svgtd.canvas.Gend()
}

// DrawLegend draws a legend box in the bottom-left corner of the image, in
// absolute pixel coordinates (independent of the tree's data-space
// transform), listing each metadata field's name, marker shape, and
// color-coded values.
func (svgtd *svgTreeDrawer) DrawLegend(entries []LegendEntry) {
	rows := legendRows(entries)
	if len(rows) == 0 {
		return
	}

	totalH := float64(svgtd.height + svgtd.topmargin + svgtd.bottommargin)

	maxChars := legendMaxChars(rows)
	legendW := legendSwatchGap + float64(maxChars)*legendCharWidth + 2*legendPadding
	legendH := float64(len(rows))*legendRowHeight + 2*legendPadding

	x0 := legendPadding
	y0 := totalH - legendPadding - legendH

	svgtd.canvas.Rect(round(x0), round(y0), round(legendW), round(legendH),
		"fill:white;fill-opacity:0.85;stroke:#999999;stroke-width:1;")

	for i, r := range rows {
		rowY := y0 + legendPadding + float64(i)*legendRowHeight + legendRowHeight/2.0
		textX := x0 + legendPadding
		textStyle := "alignment-baseline:middle;text-anchor:start;font-family:sans-serif;font-size:8px;"
		if r.isHeader {
			textStyle = "alignment-baseline:middle;text-anchor:start;font-family:sans-serif;font-size:8px;font-weight:bold;"
		}
		if r.hasSwatch {
			style := svgFillStyle(r.r, r.g, r.b, r.a, "black")
			svgtd.canvas.Translate(round(x0+legendPadding+4), round(rowY))
			svgtd.drawShape(0, r.shape, style)
			svgtd.canvas.Gend()
			textX = x0 + legendPadding + legendSwatchGap
		}
		svgtd.canvas.Text(round(textX), round(rowY), r.text, textStyle)
	}
}

func (svgtd *svgTreeDrawer) Write() {
	svgtd.canvas.End()
}

func (svgtd *svgTreeDrawer) Bounds() (width, height int) {
	width, height = svgtd.width, svgtd.height
	return
}

// svgFillStyle builds a fill style using a standard 6-digit hex color plus
// a separate fill-opacity, rather than the 8-digit #rrggbbaa shorthand:
// some SVG renderers (e.g. older Inkscape/librsvg builds) fail to parse
// the CSS Color Level 4 alpha-hex syntax and silently fall back to an
// opaque black fill instead of reporting an error.
func svgFillStyle(r, g, b, a uint8, stroke string) string {
	return fmt.Sprintf("stroke-width:1;fill:#%02x%02x%02x;fill-opacity:%.3f;stroke:%s;", r, g, b, float64(a)/255.0, stroke)
}

func round(x float64) int {
	if x < 0 {
		return int(math.Ceil(x - .5))
	} else {
		return int(math.Floor(x + .5))
	}
}
