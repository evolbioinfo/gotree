package cmd

import (
	"os"

	"github.com/evolbioinfo/gotree/io"
	"github.com/evolbioinfo/gotree/tree"
	"github.com/spf13/cobra"
)

func balancedTree(nbtrees int, depth int, output string, rooted bool) error {
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
		t, err = tree.RandomBalancedBinaryTree(depth, rooted)
		if err != nil {
			return err
		}
		f.WriteString(t.Newick() + "\n")
	}
	return nil
}

// binarytreeCmd represents the binarytree command
var balancedtreeCmd = &cobra.Command{
	Use:   "balancedtree",
	Short: "Generates a random balanced binary tree",
	Long: `Generates a random balanced binary tree
`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		if err = balancedTree(generateNbTrees, generateDepth, generateOutputfile, generateRooted); err != nil {
			io.LogError(err)
		}
		return
	},
}

func init() {
	generateCmd.AddCommand(balancedtreeCmd)
	balancedtreeCmd.PersistentFlags().IntVarP(&generateDepth, "depth", "d", 3, "Depth of the balanced binary tree")

}
