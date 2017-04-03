package cmd

import (
	"github.com/fredericlemoine/gotree/draw"
	"github.com/spf13/cobra"
)

var termwidth int

// textCmd represents the text command
var textCmd = &cobra.Command{
	Use:   "text",
	Short: "Print trees in ASCII",
	Long:  `Print trees in ASCII.`,
	Run: func(cmd *cobra.Command, args []string) {
		var d draw.TreeDrawer
		var l draw.TreeLayout
		f := openWriteFile(outtreefile)
		for tr := range readTrees(intreefile) {
			d = draw.NewTextTreeDrawer(f, termwidth, len(tr.Tree.Tips())*2, 10)
			l = draw.NewNormalLayout(d, !drawNoBranchLengths, !drawNoTipLabels)
			l.DrawTree(tr.Tree)
		}
		f.Close()
	},
}

func init() {
	drawCmd.AddCommand(textCmd)
	textCmd.PersistentFlags().IntVarP(&termwidth, "width", "w", 200, "Width of tree/terminal (in characters)")
}
