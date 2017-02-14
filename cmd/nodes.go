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
		f.WriteString("tree\tnid\tnneigh\tname\tdepth\n")
		var depth int
		var err error
		for statsintree := range readTrees(intreefile) {
			statsintree.Tree.ComputeDepths()
			for i, n := range statsintree.Tree.Nodes() {
				if depth, err = n.Depth(); err != nil {
					io.ExitWithMessage(err)
				}
				f.WriteString(fmt.Sprintf("%d\t%d\t%d\t%s\t%d\n", statsintree.Id, i, n.Nneigh(), n.Name(), depth))
			}
		}
		f.Close()
	},
}

func init() {
	statsCmd.AddCommand(nodesCmd)
}
