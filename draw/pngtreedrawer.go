package draw

import (
	"bufio"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"

	"github.com/golang/freetype/truetype"
	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dimg"
	"golang.org/x/image/font/gofont/gobold"
	"golang.org/x/image/font/gofont/goitalic"
	"golang.org/x/image/font/gofont/gomono"
	"golang.org/x/image/font/gofont/goregular"
)

/*
TextTreeDrawer initializer. TextTreeDraws draws tree as ASCII on stdout or any file.
So far: Does not take into account branch lengths.
*/
func NewPngTreeDrawer(file *os.File, width, height int, leftmargin, rightmargin, topmargin, bottommargin int) TreeDrawer {
	ptd := &pngTreeDrawer{
		file,
		width,
		height,
		leftmargin,
		rightmargin,
		topmargin,
		bottommargin,
		nil,
		nil,
	}
	ptd.img = image.NewRGBA(image.Rect(0, 0, width+leftmargin+rightmargin, height+bottommargin+topmargin))
	ptd.gc = draw2dimg.NewGraphicContext(ptd.img)
	initFonts()
	ptd.gc.SetFontData(draw2d.FontData{Name: "goregular"})
	return ptd
}

/*
Draw a tree in a png file.
*/
type pngTreeDrawer struct {
	outfile      *os.File                  // Output file
	width        int                       // Width of the ascii canvas
	height       int                       // Height of the ascii canvas
	leftmargin   int                       // Left margin of the canvas (in addition to the width)
	rightmargin  int                       // Right margin of the canvas (in addition to the width)
	topmargin    int                       // Top margin of the canvas (in addition to the height)
	bottommargin int                       // Bottom margin of the canvas (in addition to the height)
	img          *image.RGBA               // Image
	gc           *draw2dimg.GraphicContext // Graphic context to draw on the image
}

func (ptd *pngTreeDrawer) DrawHLine(x1, x2, y, maxlength, maxheight float64) {
	min := float64(ptd.width)*x1/maxlength + float64(ptd.leftmargin)
	max := float64(ptd.width)*x2/maxlength + float64(ptd.leftmargin)
	ypos := float64(ptd.height)*y/maxheight + float64(ptd.topmargin)
	ptd.gc.SetFillColor(color.RGBA{0x00, 0x00, 0x00, 0xff})
	ptd.gc.SetStrokeColor(color.RGBA{0x00, 0x00, 0x00, 0xff})
	ptd.gc.SetLineWidth(2)
	ptd.gc.MoveTo(min, ypos)
	ptd.gc.LineTo(max, ypos)
	ptd.gc.Close()
	ptd.gc.FillStroke()
}

func (ptd *pngTreeDrawer) DrawVLine(x, y1, y2, maxlength, maxheight float64) {
	min := float64(ptd.height)*y1/maxheight + float64(ptd.topmargin)
	max := float64(ptd.height)*y2/maxheight + float64(ptd.topmargin)
	xpos := float64(ptd.width)*x/maxlength + float64(ptd.leftmargin)
	ptd.gc.SetFillColor(color.RGBA{0x00, 0x00, 0x00, 0xff})
	ptd.gc.SetStrokeColor(color.RGBA{0x00, 0x00, 0x00, 0xff})
	ptd.gc.SetLineWidth(2)
	ptd.gc.MoveTo(xpos, min)
	ptd.gc.LineTo(xpos, max)
	ptd.gc.Close()
	ptd.gc.FillStroke()
}

func (ptd *pngTreeDrawer) DrawName(x, y float64, name string, maxlength, maxheight float64) {
	ptd.gc.SetFillColor(color.RGBA{0x00, 0x00, 0x00, 0xff})
	ptd.gc.SetStrokeColor(color.RGBA{0x00, 0x00, 0x00, 0xff})
	ypos := float64(ptd.height)*y/maxheight + float64(ptd.topmargin)
	xpos := float64(ptd.width)*x/maxlength + float64(ptd.leftmargin)
	ptd.gc.FillStringAt(name, xpos, ypos)
}
func (ptd *pngTreeDrawer) Write() {
	// Create Writer from file
	b := bufio.NewWriter(ptd.outfile)
	// Write the image into the buffer
	_ = png.Encode(b, ptd.img)
	_ = b.Flush()
}

type myFontCache map[string]*truetype.Font

func (fc myFontCache) Store(fd draw2d.FontData, font *truetype.Font) {
	fc[fd.Name] = font
}

func (fc myFontCache) Load(fd draw2d.FontData) (*truetype.Font, error) {
	font, stored := fc[fd.Name]
	if !stored {
		return nil, fmt.Errorf("Font %s is not stored in font cache.", fd.Name)
	}
	return font, nil
}

func initFonts() {
	fontCache := myFontCache{}

	TTFs := map[string]([]byte){
		"goregular": goregular.TTF,
		"gobold":    gobold.TTF,
		"goitalic":  goitalic.TTF,
		"gomono":    gomono.TTF,
	}

	for fontName, TTF := range TTFs {
		font, err := truetype.Parse(TTF)
		if err != nil {
			panic(err)
		}
		fontCache.Store(draw2d.FontData{Name: fontName}, font)
	}

	draw2d.SetFontCache(fontCache)
}
