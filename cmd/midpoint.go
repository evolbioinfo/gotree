package cmd

import (
	"github.com/fredericlemoine/gotree/io"
	"github.com/spf13/cobra"
)

// rerootCmd represents the reroot command
var midpointCmd = &cobra.Command{
	Use:   "midpoint",
	Short: "Reroot trees at midpoint",
	Long: `Reroot tree at midpoint.

Example:

gotree reroot midpoint  -i tree.nw > reroot.nw
`,
	Run: func(cmd *cobra.Command, args []string) {
		f := openWriteFile(outtreefile)
		treefile, trees := readTrees(intreefile)
		defer treefile.Close()

		for t := range trees {
			if t.Err != nil {
				io.ExitWithMessage(t.Err)
			}
			err := t.Tree.RerootMidPoint()
			if err != nil {
				io.ExitWithMessage(err)
			}
			f.WriteString(t.Tree.Newick() + "\n")
		}
		f.Close()
	},
}

func init() {
	rerootCmd.AddCommand(midpointCmd)
}
