package cmd

import (
	"github.com/fredericlemoine/gotree/draw"
	"github.com/fredericlemoine/gotree/io"
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
		treefile, treechan := readTrees(intreefile)
		defer treefile.Close()
		for t := range treechan {
			if t.Err != nil {
				io.ExitWithMessage(t.Err)
			}
			d = draw.NewTextTreeDrawer(f, termwidth, len(t.Tree.Tips())*2, 10)
			l = draw.NewNormalLayout(d, !drawNoBranchLengths, !drawNoTipLabels, drawInternalNodeLabels, drawSupport)
			l.SetSupportCutoff(drawSupportCutoff)
			l.DrawTree(t.Tree)
		}
		f.Close()
	},
}

func init() {
	drawCmd.AddCommand(textCmd)
	textCmd.PersistentFlags().IntVarP(&termwidth, "width", "w", 200, "Width of tree/terminal (in characters)")
}
