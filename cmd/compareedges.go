package cmd

import (
	"fmt"
	"github.com/fredericlemoine/gotree/io"
	"github.com/fredericlemoine/gotree/io/utils"
	"github.com/fredericlemoine/gotree/tree"
	"github.com/spf13/cobra"
	"os"
)

// compareedgesCmd represents the compareedges command
var compareedgesCmd = &cobra.Command{
	Use:   "edges",
	Short: "Compare edges of a reference tree with another tree",
	Long: `Compare edges of a reference tree with another tree

If the compared tree file contains several trees, it will take the first one only
`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Fprintf(os.Stderr, "Reference : %s\n", compareTree1)
		fmt.Fprintf(os.Stderr, "Compared  : %s\n", compareTree2)
		var err error
		var refTree *tree.Tree
		if refTree, err = utils.ReadRefTree(compareTree1); err != nil {
			io.ExitWithMessage(err)
		}
		refTree.ComputeDepths()

		nbtrees := 0
		compareChannel := make(chan tree.Trees, 15)

		go func() {
			if nbtrees, err = utils.ReadCompTrees(compareTree2, compareChannel); err != nil {
				io.ExitWithMessage(err)
			}
		}()

		t2 := <-compareChannel

		edges1 := refTree.Edges()
		edges2 := t2.Tree.Edges()

		fmt.Printf("brid\tlength\tsupport\tterminal\tdepth\ttopodepth\trightname\tfound\n")
		for i, e1 := range edges1 {
			found := false
			for _, e2 := range edges2 {
				if e1.SameBipartition(e2) {
					found = true
					break
				}
			}
			fmt.Printf("%d\t%s\t%t\n", i, e1.ToStatsString(), found)
		}
	},
}

func init() {
	compareCmd.AddCommand(compareedgesCmd)
}
