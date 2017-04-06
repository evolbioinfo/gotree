package draw

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"math"
)

/*
TextTreeDrawer initializer. TextTreeDraws draws tree as ASCII on stdout or any file.
So far: Does not take into account branch lengths.
*/
func NewTextTreeDrawer(w io.Writer, width, height int, rightmargin int) TreeDrawer {
	ttd := &textTreeDrawer{
		w,
		width,
		rightmargin,
		height,
		nil,
	}
	//ttd.height = ntips * 2
	ttd.textCanvas = make([][]rune, ttd.height)
	for i := 0; i < len(ttd.textCanvas); i++ {
		ttd.textCanvas[i] = make([]rune, ttd.width+ttd.rightmargin)
		for j := 0; j < len(ttd.textCanvas[i]); j++ {
			ttd.textCanvas[i][j] = ' '
		}
	}
	return ttd
}

/*
Draw a tree as ASCII in any file (stdout/stderr, etc.).
*/
type textTreeDrawer struct {
	outwriter   io.Writer // Output file
	width       int       // Width of the ascii canvas
	rightmargin int       // Right margin of the canvas (in addition to the width)
	height      int       // Height of the ascii canvas
	textCanvas  [][]rune  // ascii canvas
}

func (ttd *textTreeDrawer) DrawHLine(x1, x2, y, maxlength, maxheight float64) {
	min := float64(ttd.width) * x1 / maxlength
	max := float64(ttd.width) * x2 / maxlength
	ypos := y * float64(ttd.height) / maxheight
	for i := int(min); float64(i) < max-1; i++ {
		if i == int(min) {
			ttd.textCanvas[int(ypos)][i] = '+'
		} else {
			ttd.textCanvas[int(ypos)][i] = '-'
		}
	}
}

func (ttd *textTreeDrawer) DrawVLine(x, y1, y2, maxlength, maxheight float64) {
	min := float64(ttd.height) * y1 / maxheight
	max := float64(ttd.height) * y2 / maxheight
	xpos := float64(ttd.width) * x / maxlength
	for i := int(min); float64(i) < max; i++ {
		if i == int(min) || i == int(max) {
			ttd.textCanvas[i][int(xpos)] = '+'
		} else {
			ttd.textCanvas[i][int(xpos)] = '|'
		}
	}
}

func (ttd *textTreeDrawer) DrawLine(x1, x2, y1, y2, maxlength, maxheight float64) {
	log.Print("Method DrawLine cannot be called on textTreeDrawer: The line will not be drawn")
}

func (ttd *textTreeDrawer) DrawCurve(centerx, centery float64, middlex, middley float64, radius float64, startAngle, endAngle float64, maxlength, maxheight float64) {
	log.Print("Method DrawCurve cannot be called on textTreeDrawer: The curve will not be drawn")
}

func (ttd *textTreeDrawer) DrawCircle(x, y float64, maxwidth, maxheight float64) {
	ypos := float64(ttd.height) * y / maxheight
	xpos := float64(ttd.width) * x / maxwidth
	ttd.textCanvas[int(ypos)][int(xpos)] = '*'
}

func (ttd *textTreeDrawer) DrawName(x, y float64, name string, maxlength, maxheight float64, angle float64) {
	ypos := float64(ttd.height) * y / maxheight
	xpos := float64(ttd.width) * x / maxlength
	for i, c := range []rune(name) {
		if int(math.Ceil(xpos))+i < len(ttd.textCanvas[int(ypos)]) {
			ttd.textCanvas[int(ypos)][int(math.Ceil(xpos))+i] = c
		}
	}
}
func (ttd *textTreeDrawer) Write() {
	// Create Buffered Writer from io.writer
	b := bufio.NewWriter(ttd.outwriter)
	for _, l := range ttd.textCanvas {
		for _, c := range l {
			b.WriteString(fmt.Sprintf("%c", c))
		}
		b.WriteString("\n")
	}
	_ = b.Flush()
}

func (ttd *textTreeDrawer) Bounds() (width, height int) {
	width, height = ttd.width, ttd.height
	return
}
