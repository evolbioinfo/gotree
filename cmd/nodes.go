package cmd

import (
	"fmt"
	goio "io"
	"os"

	"github.com/evolbioinfo/gotree/io"
	"github.com/evolbioinfo/gotree/tree"
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
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		var f *os.File
		var treefile goio.Closer
		var treechan <-chan tree.Trees

		if f, err = openWriteFile(outtreefile); err != nil {
			io.LogError(err)
			return
		}
		defer closeWriteFile(f, outtreefile)

		f.WriteString("tree\tnid\tnneigh\tname\tdepth\tcomments\n")
		var depth int
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
			t.Tree.ComputeDepths()
			for i, n := range t.Tree.Nodes() {
				if depth, err = n.Depth(); err != nil {
					io.LogError(err)
					return
				}
				f.WriteString(fmt.Sprintf("%d\t%d\t%d\t%s\t%d\t%s\n", t.Id, i, n.Nneigh(), n.Name(), depth, n.CommentsString()))
			}
		}
		return
	},
}

func init() {
	statsCmd.AddCommand(nodesCmd)
}
