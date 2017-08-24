package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/fredericlemoine/gotree/draw"
	"github.com/fredericlemoine/gotree/io"
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
	Run: func(cmd *cobra.Command, args []string) {
		var d draw.TreeDrawer
		var l draw.TreeLayout
		ntree := 0
		treefile, treechan := readTrees(intreefile)
		defer treefile.Close()
		for t := range treechan {
			if t.Err != nil {
				io.ExitWithMessage(t.Err)
			}
			fname := outtreefile
			if ntree > 0 {
				extension := filepath.Ext(fname)
				if extension == ".svg" {
					fname = fname[0 : len(fname)-len(extension)]
				}
				fname = fmt.Sprintf(fname+"_%03d.svg", ntree)
			}
			f := openWriteFile(fname)
			if svgradial {
				t.Tree.ReinitIndexes()
				d = draw.NewSvgTreeDrawer(f, svgwidth, svgheight, 30, 30, 30, 30)
				l = draw.NewRadialLayout(d, !drawNoBranchLengths, !drawNoTipLabels, drawInternalNodeLabels, drawSupport)
			} else if svgcircular {
				d = draw.NewSvgTreeDrawer(f, min(svgwidth, svgheight), min(svgwidth, svgheight), 30, 30, 30, 30)
				l = draw.NewCircularLayout(d, !drawNoBranchLengths, !drawNoTipLabels, drawInternalNodeLabels, drawSupport)
			} else {
				d = draw.NewSvgTreeDrawer(f, svgwidth, svgheight, 30, 30, 30, 30)
				l = draw.NewNormalLayout(d, !drawNoBranchLengths, !drawNoTipLabels, drawInternalNodeLabels, drawSupport)
			}
			l.SetSupportCutoff(drawSupportCutoff)
			l.DrawTree(t.Tree)
			f.Close()
			ntree++
		}
	},
}

func init() {
	drawCmd.AddCommand(svgCmd)
	svgCmd.PersistentFlags().IntVarP(&svgwidth, "width", "w", 200, "Width of svg image in pixels")
	svgCmd.PersistentFlags().IntVarP(&svgheight, "height", "H", 200, "Height of svg image in pixels")
	svgCmd.PersistentFlags().BoolVarP(&svgradial, "radial", "r", false, "Radial layout (default : normal)")
	svgCmd.PersistentFlags().BoolVarP(&svgcircular, "circular", "c", false, "Circular/Polar layout (default : normal)")
}
