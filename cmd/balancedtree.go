package cmd

import (
	"github.com/fredericlemoine/gotree/io"
	"github.com/fredericlemoine/gotree/tree"
	"github.com/spf13/cobra"
	"math/rand"
	"os"
)

func balancedTree(nbtrees int, depth int, output string, seed int64, rooted bool) error {
	var f *os.File
	var err error
	var t *tree.Tree

	rand.Seed(seed)

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
	Run: func(cmd *cobra.Command, args []string) {
		if err := balancedTree(generateNbTrees, generateDepth, generateOutputfile, generateSeed, generateRooted); err != nil {
			io.ExitWithMessage(err)
		}
	},
}

func init() {
	generateCmd.AddCommand(balancedtreeCmd)
	balancedtreeCmd.PersistentFlags().IntVarP(&generateDepth, "depth", "d", 3, "Depth of the balanced binary tree")

}
