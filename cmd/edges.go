package cmd

import (
	"fmt"
	goio "io"
	"os"

	"github.com/evolbioinfo/gotree/io"
	"github.com/evolbioinfo/gotree/tree"
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
 0 - Tree id
 1 - Edge id
 2 - Length
 3 - Support
 4 - Terminal (true/false)
 5 - Depth (Shortest path to a tip)
 6 - Topo depth (Number of tips on the lightest side)
 7 - Name of the Right node
 8 - Comment of the edge if any
 9 - name of left node if any
10 - comment of right node if any
11 - comment of left node if any

Example of usage:

gotree stats edges -i t.nw

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

		f.WriteString("tree\tbrid\tlength\tsupport\tterminal\tdepth\ttopodepth\trightname\tcomments\tleftname\trightcomment\tleftcomment")
		f.WriteString("\n")
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
			if err = t.Tree.ReinitIndexes(); err != nil {
				io.LogError(err)
				return err
			}
			t.Tree.ComputeDepths()
			for i, e := range t.Tree.Edges() {
				f.WriteString(
					fmt.Sprintf("%d\t%d\t%s",
						t.Id, i, e.ToStatsString(true)))
				f.WriteString("\n")
			}
		}
		return
	},
}

func init() {
	statsCmd.AddCommand(edgesCmd)
}
