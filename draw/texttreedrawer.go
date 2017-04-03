package draw

import (
	"fmt"
	"math"
	"os"
)

/*
TextTreeDrawer initializer. TextTreeDraws draws tree as ASCII on stdout or any file.
So far: Does not take into account branch lengths.
*/
func NewTextTreeDrawer(file *os.File, width, height int, rightmargin int) TreeDrawer {
	ttd := &textTreeDrawer{
		file,
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
	outfile     *os.File // Output file
	width       int      // Width of the ascii canvas
	rightmargin int      // Right margin of the canvas (in addition to the width)
	height      int      // Height of the ascii canvas
	textCanvas  [][]rune // ascii canvas
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

func (ttd *textTreeDrawer) DrawName(x, y float64, name string, maxlength, maxheight float64) {
	ypos := float64(ttd.height) * y / maxheight
	xpos := float64(ttd.width) * x / maxlength
	for i, c := range []rune(name) {
		if int(math.Ceil(xpos))+i < len(ttd.textCanvas[int(ypos)]) {
			ttd.textCanvas[int(ypos)][int(math.Ceil(xpos))+i] = c
		}
	}
}
func (ttd *textTreeDrawer) Write() {
	for _, l := range ttd.textCanvas {
		for _, c := range l {
			ttd.outfile.WriteString(fmt.Sprintf("%c", c))
		}
		ttd.outfile.WriteString("\n")
	}
}
