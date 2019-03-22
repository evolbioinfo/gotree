package cmd

import (
	goio "io"
	"os"

	"github.com/evolbioinfo/gotree/draw"
	"github.com/evolbioinfo/gotree/io"
	"github.com/evolbioinfo/gotree/tree"
	"github.com/spf13/cobra"
)

var termwidth int

// textCmd represents the text command
var textCmd = &cobra.Command{
	Use:   "text",
	Short: "Print trees in ASCII",
	Long:  `Print trees in ASCII.`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		var f *os.File
		var treefile goio.Closer
		var treechan <-chan tree.Trees
		var d draw.TreeDrawer
		var l draw.TreeLayout
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
		for t := range treechan {
			if t.Err != nil {
				io.LogError(t.Err)
				return t.Err
			}
			d = draw.NewTextTreeDrawer(f, termwidth, len(t.Tree.Tips())*2, 10)
			l = draw.NewNormalLayout(d, !drawNoBranchLengths, !drawNoTipLabels, drawInternalNodeLabels, drawSupport)
			l.SetDisplayNodeComments(drawNodeComment)
			l.SetSupportCutoff(drawSupportCutoff)
			l.DrawTree(t.Tree)
		}
		return
	},
}

func init() {
	drawCmd.AddCommand(textCmd)
	textCmd.PersistentFlags().IntVarP(&termwidth, "width", "w", 200, "Width of tree/terminal (in characters)")
}
