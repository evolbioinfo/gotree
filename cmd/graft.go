package cmd

import (
	"errors"
	"os"

	"github.com/spf13/cobra"

	"github.com/evolbioinfo/gotree/io"
	"github.com/evolbioinfo/gotree/tree"
)

var tipname string

// compareCmd represents the compare command
var graftCmd = &cobra.Command{
	Use:   "graft",
	Short: "Graft a tree t2 on a tree t1, at the position of a given tip",
	Long: `Graft a tree t2 on a tree t1, at the position of a given tip.

	The root of t2 will replace the given tip of t2.

	Example: grafting t2 on t1, at tip l1

	t1:      t2:
	/--- l1  /---l4
	|----l2  |---l5
	\---l3   \---l6

	result:
	     /---l4
	/--- |---l5
	|    \---l6
	|---l2
	\---l3
`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		var refTree, graftTree *tree.Tree
		var f *os.File

		if f, err = openWriteFile(outtreefile); err != nil {
			io.LogError(err)
			return
		}
		defer closeWriteFile(f, outtreefile)

		if intree2file == "none" {
			err = errors.New("You must provide a file containing a tree to graft")
			io.LogError(err)
			return
		}

		if refTree, err = readTree(intreefile); err != nil {
			io.LogError(err)
			return
		}

		if err = refTree.UpdateTipIndex(); err != nil {
			io.LogError(err)
		}

		if graftTree, err = readTree(intree2file); err != nil {
			io.LogError(err)
			return
		}

		refTree.GraftTreeOnTip(tipname, graftTree)

		f.WriteString(refTree.Newick() + "\n")
		return
	},
}

func init() {
	RootCmd.AddCommand(graftCmd)
	graftCmd.PersistentFlags().StringVarP(&intreefile, "reftree", "i", "stdin", "Reference tree input file")
	graftCmd.PersistentFlags().StringVarP(&intree2file, "graft", "c", "none", "Tree to graft")
	graftCmd.PersistentFlags().StringVarP(&tipname, "tip", "l", "none", "Name of the tip to graft the second tree at")
	graftCmd.Flags().StringVarP(&outtreefile, "output", "o", "stdout", "Output tree")
}
