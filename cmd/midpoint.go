package cmd

import (
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
		for t2 := range readTrees(intreefile) {
			t2.Tree.RerootMidPoint()
			f.WriteString(t2.Tree.Newick() + "\n")
		}
		f.Close()
	},
}

func init() {
	rerootCmd.AddCommand(midpointCmd)
}
