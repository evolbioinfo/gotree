package cmd

import (
	"fmt"
	goio "io"
	"log"
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
var rttrate float64
var rttminrate float64
var rttmaxrate float64
var rttrootdate float64
var rttminrootdate float64
var rttmaxrootdate float64

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
		var line, line2, line3 *plotter.Line

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
				fmt.Fprintf(f, "%d\t%f\t%f\n", tr.Id, l.X, l.Y)
			}

			if rttoutimagefile != "none" {
				var mindate, maxdate float64 = 10000.0, 0.0
				pts := make(plotter.XYs, len(rtt))
				for i, l := range rtt {
					pts[i].X = l.X
					pts[i].Y = l.Y
					if l.X > maxdate {
						maxdate = l.X
					}
					if l.X < mindate {
						mindate = l.X
					}
				}
				if rttrootdate < 0 {
					rttrootdate = mindate
				}

				point, err = plotter.NewScatter(pts)
				point.GlyphStyleFunc = func(i int) draw.GlyphStyle {
					var col color.RGBA
					if rtt[i].Tip {
						col = color.RGBA{R: 0, G: 0, B: 255, A: 255}
					} else {
						col = color.RGBA{R: 255, G: 210, B: 52, A: 255}
					}
					return draw.GlyphStyle{Color: col, Radius: 1.5, Shape: draw.CircleGlyph{}}
				}
				p.Add(point)

				if rttrate > 0 && rttrootdate > 0 {
					pts2 := plotter.XYs{{X: rttrootdate, Y: 0}, {X: maxdate, Y: rttrate * (maxdate - rttrootdate)}}
					line, err = plotter.NewLine(pts2)
					line.Width = 2
					if err != nil {
						log.Panic(err)
					}
					p.Add(line)
				}

				if rttminrate > 0 && rttmaxrootdate > 0 {
					pts3 := plotter.XYs{{X: rttmaxrootdate, Y: 0}, {X: maxdate, Y: rttminrate * (maxdate - rttmaxrootdate)}}
					line2, err = plotter.NewLine(pts3)
					line2.Color = color.RGBA{R: 100, G: 100, B: 100, A: 255}
					line2.Width = 2
					if err != nil {
						log.Panic(err)
					}
					p.Add(line2)
				}

				if rttmaxrate > 0 && rttminrootdate > 0 {
					pts4 := plotter.XYs{{X: rttminrootdate, Y: 0}, {X: maxdate, Y: rttmaxrate * (maxdate - rttminrootdate)}}
					line3, err = plotter.NewLine(pts4)
					line3.Color = color.RGBA{R: 100, G: 100, B: 100, A: 255}
					line3.Width = 2
					if err != nil {
						log.Panic(err)
					}
					p.Add(line3)
				}
				//plotutil.AddLinePoints(p, fmt.Sprintf("LTT_%d", tr.Id), pts, draw.CircleGlyph{})

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
	rttCmd.PersistentFlags().Float64Var(&rttmaxrate, "max-rate", -1.0, "Mutation rate higher bound")
	rttCmd.PersistentFlags().Float64Var(&rttminrootdate, "min-root-date", -1.0, "Root date")
	rttCmd.PersistentFlags().Float64Var(&rttrate, "rate", -1.0, "Mutation rate to display on the figure")
	rttCmd.PersistentFlags().Float64Var(&rttmaxrootdate, "max-root-date", -1.0, "Root date")
	rttCmd.PersistentFlags().Float64Var(&rttminrate, "min-rate", -1.0, "Mutation rate lower bound")
	rttCmd.PersistentFlags().Float64Var(&rttrootdate, "root-date", -1.0, "Root date")
	rttCmd.PersistentFlags().IntVar(&rttoutimagewidth, "image-width", 4, "RTT plot image image output width")
	rttCmd.PersistentFlags().IntVar(&rttoutimageheight, "image-height", 4, "RTT plot image output heigh")
}
