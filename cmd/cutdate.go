package cmd

import (
	goio "io"
	"os"

	"github.com/evolbioinfo/gotree/io"
	"github.com/evolbioinfo/gotree/tree"
	"github.com/spf13/cobra"
)

var pruneMinDate float64
var pruneMaxDate float64

// resolveCmd represents the resolve command
var pruneDateCmd = &cobra.Command{
	Use:   "date",
	Short: "Cut the input tree by keeping only parts in date window",
	Long: `Cut the input tree by keeping only parts in date window.

This command will extract part of the tree corresponding to >= min-date and <= max-date.

If min-date falls on an internal branch, it will create a new root node and will extract a tree starting at this node.
If max-date falls on an internal branch, we do not take this part of the tree, and we remove branches that end into these cases.

`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		var f *os.File
		var treefile goio.Closer
		var treechan <-chan tree.Trees
		var forest []*tree.Tree

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

		for tr := range treechan {
			if tr.Err != nil {
				io.LogError(tr.Err)
				return tr.Err
			}
			if forest, err = tr.Tree.CutTreeMinDate(pruneMinDate); err != nil {
				io.LogError(err)
				return
			}
			for _, t := range forest {
				if pruneMaxDate > 0 {
					if err = t.CutTreeMaxDate(pruneMaxDate); err != nil {
						io.LogError(err)
						return
					}
				}
				if len(t.Edges()) > 0 {
					f.WriteString(t.Newick() + "\n")
				}
			}
		}

		return
	},
}

func init() {
	cutdateCmd.AddCommand(pruneDateCmd)
	pruneDateCmd.PersistentFlags().StringVarP(&intreefile, "input", "i", "stdin", "Input tree(s) file")
	pruneDateCmd.PersistentFlags().StringVarP(&outtreefile, "output", "o", "stdout", "Forest output file")
	pruneDateCmd.PersistentFlags().Float64Var(&pruneMinDate, "min-date", 0, "Minimum date to cut the tree")
	pruneDateCmd.PersistentFlags().Float64Var(&pruneMaxDate, "max-date", 0, "Maximum date to cut the tree (0=no max date)")
}
