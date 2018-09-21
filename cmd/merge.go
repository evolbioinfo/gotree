package cmd

import (
	"os"

	"github.com/fredericlemoine/gotree/io"
	"github.com/fredericlemoine/gotree/tree"
	"github.com/spf13/cobra"
)

// mergeCmd represents the merge command
var mergeCmd = &cobra.Command{
	Use:   "merge",
	Short: "Merges two rooted trees",
	Long: `Merges two rooted trees by adding a new root connecting two former roots.

If one of the tree is not rooted, returns an error
Tip names must be different between the two trees, otherwise returns an error

Edges connecting new root with old roots have length of 1.0.

`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		var f *os.File
		var reftree, comptree *tree.Tree

		if f, err = openWriteFile(outtreefile); err != nil {
			io.LogError(err)
			return
		}
		defer closeWriteFile(f, outtreefile)

		if reftree, err = readTree(intreefile); err != nil {
			io.LogError(err)
			return
		}
		if comptree, err = readTree(intree2file); err != nil {
			io.LogError(err)
			return
		}
		reftree.UpdateTipIndex()
		comptree.UpdateTipIndex()
		if err = reftree.Merge(comptree); err != nil {
			io.LogError(err)
			return
		}
		f.WriteString(reftree.Newick() + "\n")
		return
	},
}

func init() {
	RootCmd.AddCommand(mergeCmd)
	mergeCmd.PersistentFlags().StringVarP(&intreefile, "reftree", "i", "stdin", "Reference tree input file")
	mergeCmd.PersistentFlags().StringVarP(&intree2file, "compared", "c", "stdin", "Compared tree input file")
	mergeCmd.PersistentFlags().StringVarP(&outtreefile, "output", "o", "stdout", "Merged tree output file")
}
