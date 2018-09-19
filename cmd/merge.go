package cmd

import (
	"github.com/fredericlemoine/gotree/io"
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
	Run: func(cmd *cobra.Command, args []string) {
		f := openWriteFile(outtreefile)
		defer closeWriteFile(f, outtreefile)

		refTree := readTree(intreefile)
		compTree := readTree(intree2file)
		refTree.UpdateTipIndex()
		compTree.UpdateTipIndex()
		err := refTree.Merge(compTree)
		if err != nil {
			io.ExitWithMessage(err)
		}
		f.WriteString(refTree.Newick() + "\n")
	},
}

func init() {
	RootCmd.AddCommand(mergeCmd)
	mergeCmd.PersistentFlags().StringVarP(&intreefile, "reftree", "i", "stdin", "Reference tree input file")
	mergeCmd.PersistentFlags().StringVarP(&intree2file, "compared", "c", "stdin", "Compared tree input file")
	mergeCmd.PersistentFlags().StringVarP(&outtreefile, "output", "o", "stdout", "Merged tree output file")
}
