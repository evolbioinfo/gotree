package cmd

import (
	"github.com/fredericlemoine/gotree/io"
	"github.com/fredericlemoine/gotree/io/utils"
	"github.com/fredericlemoine/gotree/tree"
	"github.com/spf13/cobra"
	"os"
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
		// Read Tree
		var t *tree.Tree
		var err error
		t, err = utils.ReadRefTree(collapseInputTree)
		if err != nil {
			io.ExitWithMessage(err)
		}
		var f *os.File
		if collapseOutputTree != "stdout" {
			f, err = os.Create(collapseOutputTree)
		} else {
			f = os.Stdout
		}
		if err != nil {
			io.ExitWithMessage(err)
		}

		t.CollapseShortBranches(shortbranchesThreshold)

		if err != nil {
			io.ExitWithMessage(err)
		}

		f.WriteString(t.Newick() + "\n")
		f.Close()
	},
}

func init() {
	RootCmd.AddCommand(collapsebrlenCmd)
	collapsebrlenCmd.Flags().Float64VarP(&shortbranchesThreshold, "length", "l", 0.0, "Length cutoff to collapse the branch")
	collapsebrlenCmd.PersistentFlags().StringVarP(&collapseInputTree, "input", "i", "stdin", "Input tree")
	collapsebrlenCmd.PersistentFlags().StringVarP(&collapseOutputTree, "output", "o", "stdout", "Collapsed tree output file")

}
