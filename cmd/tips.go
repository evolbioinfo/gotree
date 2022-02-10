package cmd

import (
	"fmt"
	goio "io"
	"os"

	"github.com/evolbioinfo/gotree/io"
	"github.com/evolbioinfo/gotree/tree"
	"github.com/spf13/cobra"
)

func tipInfoRecur(t *tree.Tree, f *os.File, id int, cur *tree.Node, prev *tree.Node, prevEdge *tree.Edge, height float64) {
	if cur == nil {
		cur = t.Root()
	}
	if cur.Tip() {
		f.WriteString(fmt.Sprintf("%d\t%d\t%d\t%s\t%.8f\t%.8f\n", id, cur.Id(), cur.Nneigh(), cur.Name(), prevEdge.Length(), height))
	}
	for i, n := range cur.Neigh() {
		if n != prev {
			e := cur.Edges()[i]
			tipInfoRecur(t, f, id, n, cur, e, height+e.Length())
		}
	}
}

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

		f.WriteString("tree\tid\tnneigh\tname\tExternalBranch\tRootToTip\n")
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
			tipInfoRecur(t.Tree, f, t.Id, nil, nil, nil, 0.0)
		}
		return
	},
}

func init() {
	statsCmd.AddCommand(tipsCmd)
}
