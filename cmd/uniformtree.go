package cmd

import (
	"github.com/fredericlemoine/gotree/io"
	"github.com/fredericlemoine/gotree/tree"
	"github.com/spf13/cobra"
	"math/rand"
	"os"
)

func uniformTree(nbtrees int, nbtips int, output string, seed int64, rooted bool) error {
	var f *os.File
	var err error
	var t *tree.Tree

	rand.Seed(seed)

	if output != "stdout" {
		f, err = os.Create(output)
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
	f.Close()
	return nil
}

// binarytreeCmd represents the binarytree command
var uniformtreeCmd = &cobra.Command{
	Use:   "uniformtree",
	Short: "Generates a random uniform binary tree",
	Long: `Generates a random uniform binary tree
`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := uniformTree(generateNbTrees, generateNbTips, generateOutputfile, generateSeed, generateRooted); err != nil {
			io.ExitWithMessage(err)
		}
	},
}

func init() {
	generateCmd.AddCommand(uniformtreeCmd)
}
