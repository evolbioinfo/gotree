package cmd

import (
	"fmt"

	"github.com/fredericlemoine/gotree/io"
	"github.com/spf13/cobra"
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// edgesCmd represents the edges command
var edgesCmd = &cobra.Command{
	Use:   "edges",
	Short: "Displays statistics on edges of input tree",
	Long: `Displays statistics on edges of input tree

Statistics are displayed in text format (tab separated):
1 - Id of edge
2 - Length
3 - Support
4 - Terminal (true/false)
5 - Depth (Shortest path to a tip)
6 - Topo depth (Number of tips on the lightest side)
7 - Name of the Right node

Example of usage:

gotree stats edges -i t.nw

`,
	Run: func(cmd *cobra.Command, args []string) {
		f := openWriteFile(outtreefile)
		f.WriteString("tree\tbrid\tlength\tsupport\tterminal\tdepth\ttopodepth\trightname\tcomments")
		f.WriteString("\n")
		treefile, trees := readTrees(intreefile)
		defer treefile.Close()

		for t := range trees {
			if t.Err != nil {
				io.ExitWithMessage(t.Err)
			}
			t.Tree.ReinitIndexes()
			t.Tree.ComputeDepths()
			for i, e := range t.Tree.Edges() {
				f.WriteString(
					fmt.Sprintf("%d\t%d\t%s",
						t.Id, i, e.ToStatsString()))
				f.WriteString("\n")
			}
		}
		f.Close()
	},
}

func init() {
	statsCmd.AddCommand(edgesCmd)
}
