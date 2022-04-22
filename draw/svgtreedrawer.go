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

func (svgtd *svgTreeDrawer) DrawColoredCircle(x, y float64, r, g, b, a uint8) {
	centerx2 := x*float64(svgtd.width-svgtd.maxNameLength)/svgtd.maxLength + float64(svgtd.topmargin)
	centery2 := y*float64(svgtd.height-svgtd.maxNameHeight)/svgtd.maxHeight + float64(svgtd.leftmargin)
	svgtd.canvas.Circle(round(centerx2), round(centery2), 5, fmt.Sprintf("stroke-width:1;fill: #%02x%02x%02x%02x;stroke: black;", r, g, b, a))
}

/* angle:  incoming branch angle */
func (svgtd *svgTreeDrawer) DrawName(x, y float64, name string, angle float64) {
	degree := angle * 180.0 / math.Pi
	//left, top, right, bottom := ptd.gc.GetStringBounds(name)
	// Text width: Not very elegant so far...
	textsize := 4 * len(name)
	ypos := int(float64(svgtd.height-svgtd.maxNameHeight)*y/svgtd.maxHeight + float64(svgtd.topmargin))
	xpos := int(float64(svgtd.width-svgtd.maxNameLength)*x/svgtd.maxLength + float64(svgtd.leftmargin))

	// We rotate the other way (text not upside down)
	if angle < 3*math.Pi/2.0 && angle > math.Pi/2.0 {
		svgtd.canvas.Translate(xpos, ypos)
		svgtd.canvas.Rotate(degree - 180)
		svgtd.canvas.Text(-(textsize)-int(svgtd.dTip), 0, name, "font-family: sans-serif;font-size:8px;")
		svgtd.canvas.Gend()
		svgtd.canvas.Gend()
	} else {
		svgtd.canvas.Translate(xpos, ypos)
		svgtd.canvas.Rotate(degree)
		svgtd.canvas.Text(int(svgtd.dTip), 0, name, "font-family: sans-serif;font-size:8px;")
		svgtd.canvas.Gend()
		svgtd.canvas.Gend()
	}
}

func (svgtd *svgTreeDrawer) Write() {
	svgtd.canvas.End()
}

func (svgtd *svgTreeDrawer) Bounds() (width, height int) {
	width, height = svgtd.width, svgtd.height
	return
}

func round(x float64) int {
	if x < 0 {
		return int(math.Ceil(x - .5))
	} else {
		return int(math.Floor(x + .5))
	}
}
