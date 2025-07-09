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

var rttoutimagefile string
var rttoutimagewidth int
var rttoutimageheight int
var rttinternalnodes bool

// resolveCmd represents the resolve command
var rttCmd = &cobra.Command{
	Use:   "rtt",
	Short: "Root To Tip regression",
	Long: `Compute Root To Tip regression.

It considers input tree as rooted.
`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		var f *os.File
		var treefile goio.Closer
		var treechan <-chan tree.Trees
		var p *plot.Plot
		var point *plotter.Scatter
		var rtt []tree.RTTData

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

		if rttoutimagefile != "none" {
			p = plot.New()
			p.Title.Text = "Root to Tips"
			p.X.Label.Text = "Time"
			p.Y.Label.Text = "Distance to root"
		}

		for tr := range treechan {
			if tr.Err != nil {
				io.LogError(tr.Err)
				return tr.Err
			}
			rtt, err = tr.Tree.RTT(!rttinternalnodes)
			for _, l := range rtt {
				f.WriteString(fmt.Sprintf("%d\t%f\t%f\n", tr.Id, l.X, l.Y))
			}

			if rttoutimagefile != "none" {
				pts := make(plotter.XYs, len(rtt))
				for i, l := range rtt {
					pts[i].X = l.X
					pts[i].Y = l.Y
				}
				point, err = plotter.NewScatter(pts)
				point.Shape = draw.CircleGlyph{}
				point.Radius = 1
				point.GlyphStyleFunc = func(i int) draw.GlyphStyle {
					var col color.RGBA
					if rtt[i].Tip {
						col = color.RGBA{R: 0, G: 0, B: 255, A: 255}
					} else {
						col = color.RGBA{R: 255, G: 210, B: 52, A: 255}
					}
					return draw.GlyphStyle{Color: col, Radius: 1, Shape: draw.CircleGlyph{}}
				}

				//plotutil.AddLinePoints(p, fmt.Sprintf("LTT_%d", tr.Id), pts, draw.CircleGlyph{})
				p.Add(point)
				p.Legend.Add(fmt.Sprintf("TTT_%d", tr.Id))
			}
		}

		if rttoutimagefile != "none" {
			// Save the plot to a PNG file.
			if err = p.Save(font.Length(rttoutimagewidth)*vg.Inch, font.Length(rttoutimageheight)*vg.Inch, rttoutimagefile); err != nil {
				io.LogError(err)
				return
			}
		}
		return
	},
}

func init() {
	RootCmd.AddCommand(rttCmd)
	rttCmd.PersistentFlags().StringVarP(&intreefile, "input", "i", "stdin", "Input tree(s) file")
	rttCmd.PersistentFlags().StringVarP(&outtreefile, "output", "o", "stdout", "RTT output file")
	rttCmd.PersistentFlags().StringVar(&rttoutimagefile, "image", "none", "RTT plot image image output file")
	rttCmd.PersistentFlags().BoolVar(&rttinternalnodes, "internal-nodes", false, "include internal nodes")
	rttCmd.PersistentFlags().IntVar(&rttoutimagewidth, "image-width", 4, "RTT plot image image output width")
	rttCmd.PersistentFlags().IntVar(&rttoutimageheight, "image-height", 4, "RTT plot image output heigh")
}
