package draw

import (
	"math"
	"os"

	"github.com/ajstarks/svgo"
)

/*
TextTreeDrawer initializer. TextTreeDraws draws tree as ASCII on stdout or any file.
So far: Does not take into account branch lengths.
*/
func NewSvgTreeDrawer(file *os.File, width, height int, leftmargin, rightmargin, topmargin, bottommargin int) TreeDrawer {
	svgtd := &svgTreeDrawer{
		file,
		width,
		height,
		leftmargin,
		rightmargin,
		topmargin,
		bottommargin,
		nil,
		20.0,
	}
	svgtd.canvas = svg.New(file)
	svgtd.canvas.Start(width+leftmargin+rightmargin, height+topmargin+bottommargin)
	return svgtd
}

/*
Draw a tree in a svg file.
*/
type svgTreeDrawer struct {
	outfile      *os.File // Output file
	width        int      // Width of the ascii canvas
	height       int      // Height of the ascii canvas
	leftmargin   int      // Left margin of the canvas (in addition to the width)
	rightmargin  int      // Right margin of the canvas (in addition to the width)
	topmargin    int      // Top margin of the canvas (in addition to the height)
	bottommargin int      // Bottom margin of the canvas (in addition to the height)
	canvas       *svg.SVG // SVN Canvas
	dTip         float64  // Distance from tip to label
}

func (svgtd *svgTreeDrawer) DrawHLine(x1, x2, y, maxlength, maxheight float64) {
	min := int(float64(svgtd.width)*x1/maxlength + float64(svgtd.leftmargin))
	max := int(float64(svgtd.width)*x2/maxlength + float64(svgtd.leftmargin))
	ypos := int(float64(svgtd.height)*y/maxheight + float64(svgtd.topmargin))
	svgtd.canvas.Line(min, ypos, max, ypos, "stroke-width:2; fill:black; stroke: black;")
}

func (svgtd *svgTreeDrawer) DrawVLine(x, y1, y2, maxlength, maxheight float64) {
	min := int(float64(svgtd.height)*y1/maxheight + float64(svgtd.topmargin))
	max := int(float64(svgtd.height)*y2/maxheight + float64(svgtd.topmargin))
	xpos := int(float64(svgtd.width)*x/maxlength + float64(svgtd.leftmargin))
	svgtd.canvas.Line(xpos, min, xpos, max, "stroke-width:2; fill:black; stroke: black;")
}

func (svgtd *svgTreeDrawer) DrawLine(x1, y1, x2, y2, maxlength, maxheight float64) {
	y1pos := int(float64(svgtd.height)*y1/maxheight + float64(svgtd.topmargin))
	y2pos := int(float64(svgtd.height)*y2/maxheight + float64(svgtd.topmargin))
	x1pos := int(float64(svgtd.width)*x1/maxlength + float64(svgtd.leftmargin))
	x2pos := int(float64(svgtd.width)*x2/maxlength + float64(svgtd.leftmargin))
	svgtd.canvas.Line(x1pos, y1pos, x2pos, y2pos, "stroke-width:2; fill:black; stroke: black;")
}

/* angle:  incoming branch angle */
func (svgtd *svgTreeDrawer) DrawName(x, y float64, name string, maxlength, maxheight float64, angle float64) {
	degree := angle * 180.0 / math.Pi
	//left, top, right, bottom := ptd.gc.GetStringBounds(name)
	// Text width: Not very elegant so far...
	textsize := 10 * len(name)
	ypos := int(float64(svgtd.height)*y/maxheight + float64(svgtd.topmargin))
	xpos := int(float64(svgtd.width)*x/maxlength + float64(svgtd.leftmargin))

	// We rotate the other way (text not upside down)
	if angle < 3*math.Pi/2.0 && angle > math.Pi/2.0 {
		svgtd.canvas.Translate(xpos, ypos)
		svgtd.canvas.Rotate(degree - 180)
		svgtd.canvas.Text(-(textsize)-int(svgtd.dTip), 0, name, "font-family: sans-serif;")
		svgtd.canvas.Gend()
		svgtd.canvas.Gend()
	} else {
		svgtd.canvas.Translate(xpos, ypos)
		svgtd.canvas.Rotate(degree)
		svgtd.canvas.Text(int(svgtd.dTip), 0, name, "font-family: sans-serif;")
		svgtd.canvas.Gend()
		svgtd.canvas.Gend()
	}
}

func (svgtd *svgTreeDrawer) Write() {
	svgtd.canvas.End()
}
