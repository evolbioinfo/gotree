package cmd

import (
	"github.com/spf13/cobra"
)

var shortbranchesThreshold float64

// collapseCmd represents the collapse command
var collapsebrlenCmd = &cobra.Command{
	Use:   "length",
	Short: "Collapse short branches of the input tree",
	Long: `Collapse short branches of the input tree.

Short branches are defined by a threshold (-l). All branches 
with length <= threshold are removed.

`,
	Run: func(cmd *cobra.Command, args []string) {
		f := openWriteFile(outtreefile)
		for t := range readTrees(intreefile) {
			t.Tree.CollapseShortBranches(shortbranchesThreshold)
			f.WriteString(t.Tree.Newick() + "\n")
		}
		f.Close()
	},
}

func init() {
	collapseCmd.AddCommand(collapsebrlenCmd)
	collapsebrlenCmd.Flags().Float64VarP(&shortbranchesThreshold, "length", "l", 0.0, "Length cutoff to collapse branches")
}
