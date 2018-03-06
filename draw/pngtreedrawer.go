package draw

import (
	"bufio"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"math"

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
func NewPngTreeDrawer(w io.Writer, width, height int, leftmargin, rightmargin, topmargin, bottommargin int) TreeDrawer {
	ptd := &pngTreeDrawer{
		w,
		width,
		height,
		leftmargin,
		rightmargin,
		topmargin,
		bottommargin,
		nil,
		nil,
		20.0,
	}
	ptd.img = image.NewRGBA(image.Rect(0, 0, width+leftmargin+rightmargin, height+bottommargin+topmargin))
	ptd.gc = draw2dimg.NewGraphicContext(ptd.img)
	ptd.initFonts()
	ptd.gc.SetFontData(draw2d.FontData{Name: "goregular"})
	ptd.gc.SetFontSize(10.0)
	return ptd
}

/*
Draw a tree in a png file.
*/
type pngTreeDrawer struct {
	outwriter    io.Writer                 // Output Writer
	width        int                       // Width of the ascii canvas
	height       int                       // Height of the ascii canvas
	leftmargin   int                       // Left margin of the canvas (in addition to the width)
	rightmargin  int                       // Right margin of the canvas (in addition to the width)
	topmargin    int                       // Top margin of the canvas (in addition to the height)
	bottommargin int                       // Bottom margin of the canvas (in addition to the height)
	img          *image.RGBA               // Image
	gc           *draw2dimg.GraphicContext // Graphic context to draw on the image
	dTip         float64                   // Distance from tip tolabel
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

func (ptd *pngTreeDrawer) DrawLine(x1, y1, x2, y2, maxlength, maxheight float64) {
	y1pos := float64(ptd.height)*y1/maxheight + float64(ptd.topmargin)
	y2pos := float64(ptd.height)*y2/maxheight + float64(ptd.topmargin)
	x1pos := float64(ptd.width)*x1/maxlength + float64(ptd.leftmargin)
	x2pos := float64(ptd.width)*x2/maxlength + float64(ptd.leftmargin)

	ptd.gc.SetFillColor(color.RGBA{0x00, 0x00, 0x00, 0xff})
	ptd.gc.SetStrokeColor(color.RGBA{0x00, 0x00, 0x00, 0xff})
	ptd.gc.SetLineWidth(2)
	ptd.gc.MoveTo(x1pos, y1pos)
	ptd.gc.LineTo(x2pos, y2pos)
	ptd.gc.Close()
	ptd.gc.FillStroke()
}

func (ptd *pngTreeDrawer) DrawCurve(centerx, centery float64, middlex, middley float64, radius float64, startAngle, endAngle float64, maxlength, maxheight float64) {
	centerx2 := centerx*float64(ptd.width)/maxlength + float64(ptd.topmargin)
	centery2 := centery*float64(ptd.height)/maxheight + float64(ptd.leftmargin)
	middlex2 := middlex*float64(ptd.width)/maxlength + float64(ptd.topmargin)
	middley2 := middley*float64(ptd.height)/maxheight + float64(ptd.leftmargin)
	radiusscaled := math.Sqrt(math.Pow((middley2-centery2), 2) + math.Pow((middlex2-centerx2), 2))

	ptd.gc.SetFillColor(color.RGBA{0x00, 0x00, 0x00, 0x00})
	ptd.gc.SetStrokeColor(color.RGBA{0x00, 0x00, 0x00, 0xff})
	ptd.gc.SetLineWidth(2)
	ptd.gc.ArcTo(centerx2, centery2, radiusscaled, radiusscaled, startAngle, endAngle-startAngle)
	ptd.gc.Stroke()
}

func (ptd *pngTreeDrawer) DrawCircle(x, y float64, maxlength, maxheight float64) {
	centerx2 := x*float64(ptd.width)/maxlength + float64(ptd.topmargin)
	centery2 := y*float64(ptd.height)/maxheight + float64(ptd.leftmargin)

	ptd.gc.SetFillColor(color.RGBA{0x77, 0xca, 0xff, 0xff})
	ptd.gc.SetStrokeColor(color.RGBA{0x00, 0x00, 0x00, 0xff})
	ptd.gc.SetLineWidth(1)
	ptd.gc.ArcTo(centerx2, centery2, 5, 5, 0, 2*math.Pi)
	ptd.gc.Close()
	ptd.gc.FillStroke()
}

/* angle:  incoming branch angle */
func (ptd *pngTreeDrawer) DrawName(x, y float64, name string, maxlength, maxheight float64, angle float64) {
	ptd.gc.SetFillColor(color.RGBA{0x00, 0x00, 0x00, 0xff})
	ptd.gc.SetStrokeColor(color.RGBA{0x00, 0x00, 0x00, 0xff})
	left, top, right, bottom := ptd.gc.GetStringBounds(name)
	ypos := float64(ptd.height)*y/maxheight + float64(ptd.topmargin)
	xpos := float64(ptd.width)*x/maxlength + float64(ptd.leftmargin)

	// We rotate the other way (text not upside down)
	if angle < 3*math.Pi/2.0 && angle > math.Pi/2.0 {
		ptd.gc.Translate(xpos, ypos)
		ptd.gc.Rotate(angle - math.Pi)
		ptd.gc.FillStringAt(name, -(right-left)-ptd.dTip, (bottom-top)/2.0)
		ptd.gc.Rotate(-angle + math.Pi)
		ptd.gc.Translate(-xpos, -ypos)

	} else {
		ptd.gc.Translate(xpos, ypos)
		ptd.gc.Rotate(angle)
		ptd.gc.FillStringAt(name, ptd.dTip, (bottom-top)/2.0)
		ptd.gc.Rotate(-angle)
		ptd.gc.Translate(-xpos, -ypos)
	}
}

func (ptd *pngTreeDrawer) Write() {
	// Create Writer from file
	b := bufio.NewWriter(ptd.outwriter)
	// Write the image into the buffer
	_ = png.Encode(b, ptd.img)
	_ = b.Flush()
}

func (ptd *pngTreeDrawer) Bounds() (width, height int) {
	width, height = ptd.width, ptd.height
	return
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

func (ptd *pngTreeDrawer) initFonts() {
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
	//draw2d.SetFontCache(fontCache)
	ptd.gc.FontCache = fontCache
}
