package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

// tipsCmd represents the tips command
var tipsCmd = &cobra.Command{
	Use:   "tips",
	Short: "Displays statistics on tips of input tree",
	Long: `Displays statistics on tips of input tree

Statistics are displayed in text format (tab separated):
1 - Id of tip
2 - Nb neighbors
3 - Tip Name

Example of usage:

gotree stats tips -i t.mw

`,
	Run: func(cmd *cobra.Command, args []string) {
		statsout.WriteString("tree\tid\tnneigh\tname\n")
		for statsintree := range statintrees {
			for i, n := range statsintree.Tree.Nodes() {
				if n.Nneigh() == 1 {
					statsout.WriteString(fmt.Sprintf("%d\t%d\t%d\t%s\n", statsintree.Id, i, n.Nneigh(), n.Name()))
				}
			}
		}
	},
}

func init() {
	statsCmd.AddCommand(tipsCmd)
}
