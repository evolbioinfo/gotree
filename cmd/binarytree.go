package cmd

import (
	"github.com/fredericlemoine/gotree/io"
	"github.com/fredericlemoine/gotree/tree"
	"github.com/spf13/cobra"
	"math/rand"
	"os"
	"time"
)

var binarytreeNbTips int
var binarytreeNbTrees int
var binarytreeOutputfile string
var binarytreeSeed int64

func binarytree(nbtrees int, nbtips int, output string, binarytreeSeed int64) error {
	var f *os.File
	var err error
	var t *tree.Tree

	rand.Seed(binarytreeSeed)

	if output != "stdout" {
		f, err = os.Create(output)
	} else {
		f = os.Stdout
	}
	if err != nil {
		return err
	}

	for i := 0; i < nbtrees; i++ {
		t, err = tree.RandomBinaryTree(nbtips)
		if err != nil {
			return err
		}
		f.WriteString(t.Newick() + "\n")
	}
	f.Close()
	return nil
}

// binarytreeCmd represents the binarytree command
var binarytreeCmd = &cobra.Command{
	Use:   "binarytree",
	Short: "Generates a random binary tree",
	Long: `Generates a random binary tree
`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := binarytree(binarytreeNbTrees, binarytreeNbTips, binarytreeOutputfile, binarytreeSeed); err != nil {
			io.ExitWithMessage(err)
		}
	},
}

func init() {
	generateCmd.AddCommand(binarytreeCmd)
	binarytreeCmd.Flags().IntVarP(&binarytreeNbTips, "nbtips", "l", 10, "Number of tips/leaves of the tree to generate")
	binarytreeCmd.Flags().IntVarP(&binarytreeNbTrees, "nbtrees", "n", 1, "Number of trees to generate")
	binarytreeCmd.Flags().Int64VarP(&binarytreeSeed, "seed", "s", time.Now().UTC().UnixNano(), "Initial Random Seed")
	binarytreeCmd.Flags().StringVarP(&binarytreeOutputfile, "output", "o", "stdout", "Number of tips of the tree to generate")
}
