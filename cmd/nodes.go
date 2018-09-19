package cmd

import (
	"fmt"
	"github.com/fredericlemoine/gotree/io"
	"github.com/spf13/cobra"
)

// nodesCmd represents the nodes command
var nodesCmd = &cobra.Command{
	Use:   "nodes",
	Short: "Displays statistics on nodes of input tree",
	Long: `Displays statistics on nodes of input tree.

Statistics are displayed in text format (tab separated):
1 - Id of node
2 - Nb neighbors
3 - Name of node
4 - depth of node (shortest path to a tip)

Example of usage:

gotree stats nodes -i t.nw

`,
	Run: func(cmd *cobra.Command, args []string) {
		f := openWriteFile(outtreefile)
		defer closeWriteFile(f, outtreefile)

		f.WriteString("tree\tnid\tnneigh\tname\tdepth\tcomments\n")
		var depth int
		var err error
		treefile, trees := readTrees(intreefile)
		defer treefile.Close()

		for t := range trees {
			if t.Err != nil {
				io.ExitWithMessage(t.Err)
			}
			t.Tree.ComputeDepths()
			for i, n := range t.Tree.Nodes() {
				if depth, err = n.Depth(); err != nil {
					io.ExitWithMessage(err)
				}
				f.WriteString(fmt.Sprintf("%d\t%d\t%d\t%s\t%d\t%s\n", t.Id, i, n.Nneigh(), n.Name(), depth, n.CommentsString()))
			}
		}
	},
}

func init() {
	statsCmd.AddCommand(nodesCmd)
}
