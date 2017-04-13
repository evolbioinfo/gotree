package cmd

import (
	"fmt"

	"github.com/fredericlemoine/gotree/io"
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
		f := openWriteFile(outtreefile)
		f.WriteString("tree\tid\tnneigh\tname\n")
		treefile, treechan := readTrees(intreefile)
		defer treefile.Close()
		for t := range treechan {
			if t.Err != nil {
				io.ExitWithMessage(t.Err)
			}
			for i, n := range t.Tree.Nodes() {
				if n.Nneigh() == 1 {
					f.WriteString(fmt.Sprintf("%d\t%d\t%d\t%s\n", t.Id, i, n.Nneigh(), n.Name()))
				}
			}
		}
		f.Close()
	},
}

func init() {
	statsCmd.AddCommand(tipsCmd)
}
