package cmd

import (
	"fmt"
	goio "io"
	"os"

	"github.com/evolbioinfo/gotree/io"
	"github.com/evolbioinfo/gotree/tree"
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
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		var f *os.File
		var treefile goio.Closer
		var treechan <-chan tree.Trees

		if f, err = openWriteFile(outtreefile); err != nil {
			io.LogError(err)
			return
		}
		defer closeWriteFile(f, outtreefile)

		f.WriteString("tree\tid\tnneigh\tname\n")
		if treefile, treechan, err = readTrees(intreefile); err != nil {
			io.LogError(err)
			return
		}
		defer treefile.Close()
		for t := range treechan {
			if t.Err != nil {
				io.LogError(t.Err)
				return t.Err
			}
			for i, n := range t.Tree.Nodes() {
				if n.Nneigh() == 1 {
					f.WriteString(fmt.Sprintf("%d\t%d\t%d\t%s\n", t.Id, i, n.Nneigh(), n.Name()))
				}
			}
		}
		return
	},
}

func init() {
	statsCmd.AddCommand(tipsCmd)
}
