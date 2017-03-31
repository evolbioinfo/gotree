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
		f := openWriteFile(outtreefile)
		d = draw.NewTextTreeDrawer(f, termwidth, 10)
		for tr := range readTrees(intreefile) {
			d.DrawTree(tr.Tree)
		}
		f.Close()
	},
}

func init() {
	drawCmd.AddCommand(textCmd)
	drawCmd.PersistentFlags().IntVarP(&termwidth, "width", "w", 200, "Width of tree/terminal (in characters)")

}
