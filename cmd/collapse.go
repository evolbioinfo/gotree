package cmd

import (
	"github.com/fredericlemoine/gotree/io"
	"github.com/fredericlemoine/gotree/io/utils"
	"github.com/fredericlemoine/gotree/tree"
	"github.com/spf13/cobra"
	"os"
)

var shortbranchesThreshold float64

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
		t, err = utils.ReadRefTree(transformInputTree)
		if err != nil {
			io.ExitWithMessage(err)
		}
		var f *os.File
		if renameouttree != "stdout" {
			f, err = os.Create(transformOutputTree)
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
	transformCmd.AddCommand(collapsebrlenCmd)
	collapsebrlenCmd.Flags().Float64VarP(&shortbranchesThreshold, "length", "l", 0.0, "Length cutoff to collapse the branch")

}
