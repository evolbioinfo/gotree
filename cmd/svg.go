package cmd

import (
	"fmt"
	goio "io"
	"os"
	"path/filepath"

	"github.com/evolbioinfo/gotree/draw"
	"github.com/evolbioinfo/gotree/io"
	"github.com/evolbioinfo/gotree/tree"
	"github.com/spf13/cobra"
)

var svgwidth int
var svgheight int
var svgradial bool
var svgcircular bool

// svgCmd represents the svg command
var svgCmd = &cobra.Command{
	Use:   "svg",
	Short: "Draw trees in svg files",
	Long:  `Draw trees in svg files.`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		var f *os.File
		var treefile goio.Closer
		var treechan <-chan tree.Trees
		var d draw.TreeDrawer
		var l draw.TreeLayout

		ntree := 0
		if treefile, treechan, err = readTrees(intreefile); err != nil {
			io.LogError(err)
			return
		}
		defer treefile.Close()
		for t := range treechan {
			if t.Err != nil {
				io.LogError(t.Err)
				return t.Err
			}
			fname := outtreefile
			if ntree > 0 {
				extension := filepath.Ext(fname)
				if extension == ".svg" {
					fname = fname[0 : len(fname)-len(extension)]
				}
				fname = fmt.Sprintf(fname+"_%03d.svg", ntree)
			}
			if f, err = openWriteFile(fname); err != nil {
				io.LogError(err)
				return
			}

			margin := 30
			// for _, n := range t.Tree.AllTipNames() {
			// 	if len(n)*10 > margin {
			// 		margin = len(n) * 8
			// 	}
			// }

			if svgradial {
				if err = t.Tree.ReinitIndexes(); err != nil {
					io.LogError(err)
					return
				}

				d = draw.NewSvgTreeDrawer(f, svgwidth, svgheight, margin, margin, margin, margin)
				l = draw.NewRadialLayout(d, !drawNoBranchLengths, !drawNoTipLabels, drawInternalNodeLabels, drawSupport)
				l.SetDisplayInternalNodes(drawInternalNodeSymbols)
			} else if svgcircular {
				d = draw.NewSvgTreeDrawer(f, min(svgwidth, svgheight), min(svgwidth, svgheight), margin, margin, margin, margin)
				l = draw.NewCircularLayout(d, !drawNoBranchLengths, !drawNoTipLabels, drawInternalNodeLabels, drawSupport)
			} else {
				d = draw.NewSvgTreeDrawer(f, svgwidth, svgheight, 30, margin, 30, 30)
				l = draw.NewNormalLayout(d, !drawNoBranchLengths, !drawNoTipLabels, drawInternalNodeLabels, drawSupport)
			}
			l.SetDisplayInternalNodes(drawInternalNodeSymbols)
			l.SetDisplayNodeComments(drawNodeComment)
			l.SetSupportCutoff(drawSupportCutoff)
			l.DrawTree(t.Tree)
			closeWriteFile(f, fname)
			ntree++
		}
		return
	},
}

func init() {
	drawCmd.AddCommand(svgCmd)
	svgCmd.PersistentFlags().IntVarP(&svgwidth, "width", "w", 200, "Width of svg image in pixels")
	svgCmd.PersistentFlags().IntVarP(&svgheight, "height", "H", 200, "Height of svg image in pixels")
	svgCmd.PersistentFlags().BoolVarP(&svgradial, "radial", "r", false, "Radial layout (default : normal)")
	svgCmd.PersistentFlags().BoolVarP(&svgcircular, "circular", "c", false, "Circular/Polar layout (default : normal)")
}
