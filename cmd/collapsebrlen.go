package cmd

import (
	"github.com/fredericlemoine/gotree/io"
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
		treefile, treechan := readTrees(intreefile)
		defer treefile.Close()
		for t := range treechan {
			if t.Err != nil {
				io.ExitWithMessage(t.Err)
			}
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
