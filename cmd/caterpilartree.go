package cmd

import (
	"github.com/fredericlemoine/gotree/io"
	"github.com/fredericlemoine/gotree/tree"
	"github.com/spf13/cobra"
	"math/rand"
	"os"
)

func caterpilarTree(nbtrees int, nbtips int, output string, seed int64, rooted bool) error {
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
		t, err = tree.RandomCaterpilarBinaryTree(nbtips, rooted)
		if err != nil {
			return err
		}
		f.WriteString(t.Newick() + "\n")
	}
	f.Close()
	return nil
}

// binarytreeCmd represents the binarytree command
var caterpilartreeCmd = &cobra.Command{
	Use:   "caterpilartree",
	Short: "Generates a random caterpilar binary tree",
	Long:  `Generates a random caterpilar binary tree.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := caterpilarTree(generateNbTrees, generateNbTips, generateOutputfile, generateSeed, generateRooted); err != nil {
			io.ExitWithMessage(err)
		}
	},
}

func init() {
	generateCmd.AddCommand(caterpilartreeCmd)
}
