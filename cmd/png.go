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

var pngwidth int
var pngheight int
var pngradial bool
var pngcircular bool
var pngfillbackground bool

// pngCmd represents the png command
var pngCmd = &cobra.Command{
	Use:   "png",
	Short: "Draw trees in png files",
	Long:  `Draw trees in png files.`,
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
				if extension == ".png" {
					fname = fname[0 : len(fname)-len(extension)]
				}
				fname = fmt.Sprintf(fname+"_%03d.png", ntree)
			}
			if f, err = openWriteFile(fname); err != nil {
				io.LogError(err)
				return
			}
			if pngradial {
				if err = t.Tree.ReinitIndexes(); err != nil {
					io.LogError(err)
					return
				}

				d = draw.NewPngTreeDrawer(f, pngwidth, pngheight, 30, 30, 30, 30, pngfillbackground)
				l = draw.NewRadialLayout(d, !drawNoBranchLengths, !drawNoTipLabels, drawInternalNodeLabels, drawSupport)
			} else if pngcircular {
				d = draw.NewPngTreeDrawer(f, min(pngwidth, pngheight), min(pngwidth, pngheight), 30, 30, 30, 30, pngfillbackground)
				l = draw.NewCircularLayout(d, !drawNoBranchLengths, !drawNoTipLabels, drawInternalNodeLabels, drawSupport)
			} else {
				d = draw.NewPngTreeDrawer(f, pngwidth, pngheight, 30, 30, 30, 30, pngfillbackground)
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
	drawCmd.AddCommand(pngCmd)
	pngCmd.PersistentFlags().IntVarP(&pngwidth, "width", "w", 200, "Width of png image in pixels")
	pngCmd.PersistentFlags().IntVarP(&pngheight, "height", "H", 200, "Height of png image in pixels")
	pngCmd.PersistentFlags().BoolVarP(&pngradial, "radial", "r", false, "Radial layout (default : normal)")
	pngCmd.PersistentFlags().BoolVar(&pngfillbackground, "fill-background", false, "If true, then background is white, otherwise transparent")
	pngCmd.PersistentFlags().BoolVarP(&pngcircular, "circular", "c", false, "Circular/Polar layout (default : normal)")
}
