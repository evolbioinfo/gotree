package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/fredericlemoine/gotree/draw"
	"github.com/spf13/cobra"
)

var svgwidth int
var svgheight int
var svgradial bool

// svgCmd represents the svg command
var svgCmd = &cobra.Command{
	Use:   "svg",
	Short: "Draw trees in svg files",
	Long:  `Draw trees in svg files.`,
	Run: func(cmd *cobra.Command, args []string) {
		var d draw.TreeDrawer
		var l draw.TreeLayout
		ntree := 0
		for tr := range readTrees(intreefile) {
			fname := outtreefile
			if ntree > 0 {
				extension := filepath.Ext(fname)
				if extension == ".svg" {
					fname = fname[0 : len(fname)-len(extension)]
				}
				fname = fmt.Sprintf(fname+"_%03d.svg", ntree)
			}
			f := openWriteFile(fname)
			d = draw.NewSvgTreeDrawer(f, svgwidth, svgheight, 50, 50, 50, 50)
			if svgradial {
				l = draw.NewRadialLayout(d, !drawNoBranchLengths, !drawNoTipLabels, drawInternalNodeLabels)
			} else {
				l = draw.NewNormalLayout(d, !drawNoBranchLengths, !drawNoTipLabels, drawInternalNodeLabels)
			}
			l.DrawTree(tr.Tree)
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
}
