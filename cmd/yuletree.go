package cmd

import (
	"github.com/fredericlemoine/gotree/io"
	"github.com/fredericlemoine/gotree/tree"
	"github.com/spf13/cobra"
	"math/rand"
	"os"
)

func yuleTree(nbtrees int, nbtips int, output string, seed int64, rooted bool) error {
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
		t, err = tree.RandomYuleBinaryTree(nbtips, rooted)
		if err != nil {
			return err
		}
		f.WriteString(t.Newick() + "\n")
	}

	return nil
}

// binarytreeCmd represents the binarytree command
var yuletreeCmd = &cobra.Command{
	Use:   "yuletree",
	Short: "Generates a random yule binary tree",
	Long:  `Generates a random yule binary tree.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := yuleTree(generateNbTrees, generateNbTips, generateOutputfile, generateSeed, generateRooted); err != nil {
			io.ExitWithMessage(err)
		}
	},
}

func init() {
	generateCmd.AddCommand(yuletreeCmd)
	yuletreeCmd.PersistentFlags().IntVarP(&generateNbTips, "nbtips", "l", 10, "Number of tips/leaves of the tree to generate")
}
