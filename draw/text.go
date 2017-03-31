package draw

import (
	"fmt"
	"math"
	"os"

	"github.com/fredericlemoine/gotree/tree"
)

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

/*
  TextTreeDrawer initializer. TextTreeDraws draws tree as ASCII on stdout or any file.
  So far: Does not take into account branch lengths.
*/
func NewTextTreeDrawer(file *os.File, width int, rightmargin int) *textTreeDrawer {
	return &textTreeDrawer{
		file,
		width,
		rightmargin,
		0,
		nil,
	}
}

/*
  Draw the tree on ttd.outfile. Does not close the file. The caller must do it.
*/
func (ttd *textTreeDrawer) DrawTree(t *tree.Tree) error {
	var err error = nil
	root := t.Root()
	ntips := len(t.Tips())
	curNbTips := 0
	maxLength := ttd.maxLength(t)
	ttd.height = ntips * 2
	ttd.textCanvas = make([][]rune, ttd.height)
	for i := 0; i < len(ttd.textCanvas); i++ {
		ttd.textCanvas[i] = make([]rune, ttd.width+ttd.rightmargin)
		for j := 0; j < len(ttd.textCanvas[i]); j++ {
			ttd.textCanvas[i][j] = ' '
		}
	}
	ttd.drawTreeRecur(root, nil, 0, 0, maxLength, &curNbTips, ntips, maxLength)
	ttd.dumpCanvas()
	return err
}

/*
  Recursive function that draws the tree. Returns the yposition of the current node
*/
func (ttd *textTreeDrawer) drawTreeRecur(n *tree.Node, prev *tree.Node, prevDistToRoot, distToRoot float64, maxLength float64, curtip *int, nbtips int, maxlength float64) float64 {
	ypos := 0.0
	nbchild := 0.0
	if n.Tip() {
		ypos = float64(*curtip)
		nbchild = 1.0
		ttd.drawName(distToRoot, ypos, n.Name(), maxlength, nbtips)
		*curtip++
	} else {
		minpos := -1.0
		maxpos := -1.0
		for i, child := range n.Neigh() {
			if child != prev {
				len := n.Edges()[i].Length()
				temppos := ttd.drawTreeRecur(child, n, distToRoot, distToRoot+len, maxLength, curtip, nbtips, maxlength)
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
		ttd.drawVLine(distToRoot, minpos, maxpos, maxlength, nbtips)
	}
	ttd.drawHLine(prevDistToRoot, distToRoot, ypos, maxlength, nbtips)
	return ypos
}

func (ttd *textTreeDrawer) dumpCanvas() {
	for _, l := range ttd.textCanvas {
		for _, c := range l {
			ttd.outfile.WriteString(fmt.Sprintf("%c", c))
		}
		ttd.outfile.WriteString("\n")
	}
}

func (ttd *textTreeDrawer) maxLength(t *tree.Tree) float64 {
	maxlength := 0.0
	curlength := 0.0
	root := t.Root()
	ttd.maxLengthRecur(root, nil, curlength, &maxlength)
	return maxlength
}

func (ttd *textTreeDrawer) maxLengthRecur(n *tree.Node, prev *tree.Node, curlength float64, maxlength *float64) {
	if curlength > *maxlength {
		*maxlength = curlength
	}
	for i, child := range n.Neigh() {
		if child != prev {
			brlen := n.Edges()[i].Length()
			if brlen == tree.NIL_LENGTH {
				brlen = 1.0
			}
			ttd.maxLengthRecur(child, n, curlength+brlen, maxlength)
		}
	}
}

func (ttd *textTreeDrawer) drawHLine(x1, x2, y, maxlength float64, nbtips int) {
	min := float64(ttd.width) * x1 / maxlength
	max := float64(ttd.width) * x2 / maxlength
	ypos := y * float64(ttd.height) / float64(nbtips)
	for i := int(min); float64(i) < max-1; i++ {
		if i == int(min) {
			ttd.textCanvas[int(ypos)][i] = '+'
		} else {
			ttd.textCanvas[int(ypos)][i] = '-'
		}
	}
}

func (ttd *textTreeDrawer) drawVLine(x, y1, y2, maxlength float64, nbtips int) {
	min := float64(ttd.height) * y1 / float64(nbtips)
	max := float64(ttd.height) * y2 / float64(nbtips)
	xpos := float64(ttd.width) * x / maxlength
	for i := int(min); float64(i) < max; i++ {
		if i == int(min) || i == int(max) {
			ttd.textCanvas[i][int(xpos)] = '+'
		} else {
			ttd.textCanvas[i][int(xpos)] = '|'
		}
	}
}

func (ttd *textTreeDrawer) drawName(x, y float64, name string, maxlength float64, nbtips int) {
	ypos := float64(ttd.height) * y / float64(nbtips)
	xpos := float64(ttd.width) * x / maxlength
	for i, c := range []rune(name) {
		if int(math.Ceil(xpos))+i < len(ttd.textCanvas[int(ypos)]) {
			ttd.textCanvas[int(ypos)][int(math.Ceil(xpos))+i] = c
		}
	}
}
