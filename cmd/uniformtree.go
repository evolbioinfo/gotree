package cmd

import (
	"os"

	"github.com/evolbioinfo/gotree/io"
	"github.com/evolbioinfo/gotree/tree"
	"github.com/spf13/cobra"
)

func uniformTree(nbtrees int, nbtips int, output string, rooted bool) error {
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
		t, err = tree.RandomUniformBinaryTree(nbtips, rooted)
		if err != nil {
			return err
		}
		f.WriteString(t.Newick() + "\n")
	}
	return nil
}

// binarytreeCmd represents the binarytree command
var uniformtreeCmd = &cobra.Command{
	Use:   "uniformtree",
	Short: "Generates a random uniform binary tree",
	Long: `Generates a random uniform binary tree
`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		if err = uniformTree(generateNbTrees, generateNbTips, generateOutputfile, generateRooted); err != nil {
			io.LogError(err)
		}
		return
	},
}

func init() {
	generateCmd.AddCommand(uniformtreeCmd)
	uniformtreeCmd.PersistentFlags().IntVarP(&generateNbTips, "nbtips", "l", 10, "Number of tips/leaves of the tree to generate")
}
