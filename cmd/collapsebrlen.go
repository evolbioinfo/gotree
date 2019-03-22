package cmd

import (
	goio "io"
	"os"

	"github.com/evolbioinfo/gotree/io"
	"github.com/evolbioinfo/gotree/tree"
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
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		var f *os.File
		var treefile goio.Closer
		var treechan <-chan tree.Trees

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
			t.Tree.CollapseShortBranches(shortbranchesThreshold)
			f.WriteString(t.Tree.Newick() + "\n")
		}
		return
	},
}

func init() {
	collapseCmd.AddCommand(collapsebrlenCmd)
	collapsebrlenCmd.Flags().Float64VarP(&shortbranchesThreshold, "length", "l", 0.0, "Length cutoff to collapse branches")
}
