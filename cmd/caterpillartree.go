package cmd

import (
	"os"

	"github.com/evolbioinfo/gotree/io"
	"github.com/evolbioinfo/gotree/tree"
	"github.com/spf13/cobra"
)

func caterpilarTree(nbtrees int, nbtips int, output string, rooted bool) error {
	var f *os.File
	var err error
	var t *tree.Tree

	if output != "stdout" && output != "-" {
		f, err = os.Create(output)
		defer f.Close()
	} else {
		f = os.Stdout
	}
	if err != nil {
		return err
	}

	for i := 0; i < nbtrees; i++ {
		t, err = tree.RandomCaterpillarBinaryTree(nbtips, rooted)
		if err != nil {
			return err
		}
		f.WriteString(t.Newick() + "\n")
	}
	return nil
}

// binarytreeCmd represents the binarytree command
var caterpilartreeCmd = &cobra.Command{
	Use:   "caterpillartree",
	Short: "Generates a random caterpilar binary tree",
	Long:  `Generates a random caterpilar binary tree.`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		if err = caterpilarTree(generateNbTrees, generateNbTips, generateOutputfile, generateRooted); err != nil {
			io.LogError(err)
		}
		return
	},
}

func init() {
	generateCmd.AddCommand(caterpilartreeCmd)
	caterpilartreeCmd.PersistentFlags().IntVarP(&generateNbTips, "nbtips", "l", 10, "Number of tips/leaves of the tree to generate")
}
