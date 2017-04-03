package cmd

import (
	"fmt"

	"github.com/fredericlemoine/gotree/draw"
	"github.com/spf13/cobra"
)

var pngwidth int
var pngheight int

// pngCmd represents the png command
var pngCmd = &cobra.Command{
	Use:   "png",
	Short: "Draw trees in png files",
	Long:  `Draw trees in png files.`,
	Run: func(cmd *cobra.Command, args []string) {
		var d draw.TreeDrawer
		var l draw.TreeLayout
		ntree := 0
		for tr := range readTrees(intreefile) {
			fname := outtreefile
			if ntree > 0 {
				fname = fmt.Sprintf(outtreefile+"_%03d.png", ntree)
			}
			f := openWriteFile(fname)
			d = draw.NewPngTreeDrawer(f, pngwidth, pngheight, 30, 30, 30, 30)
			l = draw.NewNormalLayout(d)
			l.DrawTree(tr.Tree)
			f.Close()
			ntree++
		}
	},
}

func init() {
	drawCmd.AddCommand(pngCmd)
	pngCmd.PersistentFlags().IntVarP(&pngwidth, "width", "w", 200, "Width of png image in pixels")
	pngCmd.PersistentFlags().IntVarP(&pngheight, "height", "H", 200, "Height of png image in pixels")
}
