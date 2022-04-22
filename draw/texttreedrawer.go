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
		0.0,
		0.0,
		0.0,
		0.0,
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

func (ttd *textTreeDrawer) SetMaxValues(maxLength, maxHeight float64, maxNameLength, maxNameHeight int) {
	ttd.maxLength = maxLength
	ttd.maxHeight = maxHeight
	ttd.maxNameLength = maxNameLength
	ttd.maxNameHeight = maxNameHeight
}

/*
Draw a tree as ASCII in any file (stdout/stderr, etc.).
*/
type textTreeDrawer struct {
	outwriter     io.Writer // Output file
	width         int       // Width of the ascii canvas
	rightmargin   int       // Right margin of the canvas (in addition to the width)
	height        int       // Height of the ascii canvas
	textCanvas    [][]rune  // ascii canvas
	maxHeight     float64   // Maximum height of object to draw (in original scale)
	maxLength     float64   // Maximum length of object to draw (in original scale)
	maxNameLength int       // Maximum length of species names / horitzontal
	maxNameHeight int       // Maximum length of species names / vertical
}

func (ttd *textTreeDrawer) DrawHLine(x1, x2, y float64) {
	min := float64(ttd.width-ttd.maxNameLength) * x1 / ttd.maxLength
	max := float64(ttd.width-ttd.maxNameLength) * x2 / ttd.maxLength
	ypos := y * float64(ttd.height-ttd.maxNameHeight) / ttd.maxHeight
	for i := int(min); float64(i) < max-1; i++ {
		if i == int(min) {
			ttd.textCanvas[int(ypos)][i] = '+'
		} else {
			ttd.textCanvas[int(ypos)][i] = '-'
		}
	}
}

func (ttd *textTreeDrawer) DrawVLine(x, y1, y2 float64) {
	min := float64(ttd.height-ttd.maxNameHeight) * y1 / ttd.maxHeight
	max := float64(ttd.height-ttd.maxNameHeight) * y2 / ttd.maxHeight
	xpos := float64(ttd.width-ttd.maxNameLength) * x / ttd.maxLength
	for i := int(min); float64(i) < max; i++ {
		if i == int(min) || i == int(max) {
			ttd.textCanvas[i][int(xpos)] = '+'
		} else {
			ttd.textCanvas[i][int(xpos)] = '|'
		}
	}
}

func (ttd *textTreeDrawer) DrawLine(x1, x2, y1, y2 float64) {
	log.Print("Method DrawLine cannot be called on textTreeDrawer: The line will not be drawn")
}

func (ttd *textTreeDrawer) DrawCurve(centerx, centery float64, middlex, middley float64, radius float64, startAngle, endAngle float64) {
	log.Print("Method DrawCurve cannot be called on textTreeDrawer: The curve will not be drawn")
}

func (ttd *textTreeDrawer) DrawCircle(x, y float64) {
	ypos := float64(ttd.height-ttd.maxNameHeight) * y / ttd.maxHeight
	xpos := float64(ttd.width-ttd.maxNameLength) * x / ttd.maxLength
	ttd.textCanvas[int(ypos)][int(xpos)] = '*'
}

func (ttd *textTreeDrawer) DrawColoredCircle(x, y float64, r, g, b, a uint8) {
	ttd.DrawCircle(x, y)
}

func (ttd *textTreeDrawer) DrawName(x, y float64, name string, angle float64) {
	ypos := float64(ttd.height-ttd.maxNameHeight) * y / ttd.maxHeight
	xpos := float64(ttd.width-ttd.maxNameLength) * x / ttd.maxLength
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
