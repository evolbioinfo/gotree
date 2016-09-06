package cmd

import (
	"github.com/fredericlemoine/gotree/io"
	"github.com/fredericlemoine/gotree/io/utils"
	"github.com/spf13/cobra"
)

var shortbranchesThreshold float64
var collapseInputTree string
var collapseOutputTree string

// collapseCmd represents the collapse command
var collapsebrlenCmd = &cobra.Command{
	Use:   "collapsebrlen",
	Short: "Collapse short branches of the input tree",
	Long: `Collapse short branches of the input tree

`,
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		var nbtrees int = 0
		intrees := make(chan utils.Trees, 15)

		/* Read ref tree(s) */
		go func() {
			if nbtrees, err = utils.ReadCompTrees(collapseInputTree, intrees); err != nil {
				io.ExitWithMessage(err)
			}
		}()

		/* Collapsing branches */
		f := openWriteFile(collapseOutputTree)
		for t := range intrees {
			t.Tree.CollapseShortBranches(shortbranchesThreshold)
			f.WriteString(t.Tree.Newick() + "\n")
		}
		f.Close()
	},
}

func init() {
	RootCmd.AddCommand(collapsebrlenCmd)
	collapsebrlenCmd.Flags().Float64VarP(&shortbranchesThreshold, "length", "l", 0.0, "Length cutoff to collapse the branch")
	collapsebrlenCmd.PersistentFlags().StringVarP(&collapseInputTree, "input", "i", "stdin", "Input tree")
	collapsebrlenCmd.PersistentFlags().StringVarP(&collapseOutputTree, "output", "o", "stdout", "Collapsed tree output file")

}
