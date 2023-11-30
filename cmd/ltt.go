package cmd

import (
	"fmt"
	goio "io"
	"os"

	"image/color"

	"github.com/evolbioinfo/gotree/io"
	"github.com/evolbioinfo/gotree/tree"
	"github.com/spf13/cobra"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/font"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

var lttoutimagefile string
var lttoutimagewidth int
var lttoutimageheight int

// resolveCmd represents the resolve command
var lttCmd = &cobra.Command{
	Use:   "ltt",
	Short: "Lineage Through Time data",
	Long: `Compute Lineage Through Time data.

1) Will output data visualizable in statistical packages (R, python, etc.).
Set of x,y coordinates pairs: x: time (or mutations) and y: number of lineages.
2) If --image <image file> is specified, then a ltt plot is drawn in the given output file.
the format of the image depends on the extension (.png, .svg, .pdf, etc. 
	see https://github.com/gonum/plot/blob/342a5cee2153b051d94ae813861f9436c5584de2/plot.go#L525C17-L525C17).
Image width and height can be specified (in inch) with --image-width and --image-height.
`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		var f *os.File
		var treefile goio.Closer
		var treechan <-chan tree.Trees
		var p *plot.Plot
		var point *plotter.Scatter
		var line *plotter.Line

		if f, err = openWriteFile(outtreefile); err != nil {
			io.LogError(err)
			return
		}
		defer closeWriteFile(f, outtreefile)

		if treefile, treechan, err = readTrees(intreefile); err != nil {
			io.LogError(err)
			return
		}
		defer treefile.Close()

		if lttoutimagefile != "none" {
			p = plot.New()
			p.Title.Text = "Lineage through time plot"
			p.X.Label.Text = "Time"
			p.Y.Label.Text = "Number of lineages"
		}

		for tr := range treechan {
			if tr.Err != nil {
				io.LogError(tr.Err)
				return tr.Err
			}
			ltt := tr.Tree.LTT()
			for _, l := range ltt {
				f.WriteString(fmt.Sprintf("%d\t%f\t%d\n", tr.Id, l.X, l.Y))
			}

			if lttoutimagefile != "none" {
				pts := make(plotter.XYs, len(ltt)-1)
				for i, l := range ltt {
					// We do not take the last point where there is 0 lineage left
					if i < (len(ltt) - 1) {
						pts[i].X = float64(l.X)
						pts[i].Y = float64(l.Y)
					}
				}
				line, point, err = plotter.NewLinePoints(pts)
				point.Shape = draw.CircleGlyph{}
				point.Radius = 1
				point.Color = color.RGBA{R: 0, G: 0, B: 255, A: 255}

				//plotutil.AddLinePoints(p, fmt.Sprintf("LTT_%d", tr.Id), pts, draw.CircleGlyph{})
				p.Add(line, point)
				p.Legend.Add(fmt.Sprintf("LTT_%d", tr.Id))
			}
		}

		if lttoutimagefile != "none" {
			// Save the plot to a PNG file.
			if err = p.Save(font.Length(lttoutimagewidth)*vg.Inch, font.Length(lttoutimageheight)*vg.Inch, lttoutimagefile); err != nil {
				io.LogError(err)
				return
			}
		}
		return
	},
}

func init() {
	RootCmd.AddCommand(lttCmd)
	lttCmd.PersistentFlags().StringVarP(&intreefile, "input", "i", "stdin", "Input tree(s) file")
	lttCmd.PersistentFlags().StringVarP(&outtreefile, "output", "o", "stdout", "LTT output file")
	lttCmd.PersistentFlags().StringVar(&lttoutimagefile, "image", "none", "LTT plot image image output file")
	lttCmd.PersistentFlags().IntVar(&lttoutimagewidth, "image-width", 4, "LTT plot image image output width")
	lttCmd.PersistentFlags().IntVar(&lttoutimageheight, "image-height", 4, "LTT plot image output heigh")

}
