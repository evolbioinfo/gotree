package cmd

import (
	goio "io"
	"os"

	"github.com/evolbioinfo/gotree/io"
	"github.com/evolbioinfo/gotree/tree"
	"github.com/spf13/cobra"
)

// collapseCmd represents the collapse command
var collapsesingleCmd = &cobra.Command{
	Use:   "single",
	Short: "Collapse branches that connect single nodes",
	Long: `Collapse branches that connect single nodes.

Single nodes are defined as nodes that:
- Connect only 2 neighbors
- Are not the root

* Branch lengths are added
* Max branch support is assigned to remaining branch

Ex:
           t1           t1
           /	       /
 n0--n1--n2   => n0--n2
           \	       \
            t2          t2

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
			t.Tree.RemoveSingleNodes()
			f.WriteString(t.Tree.Newick() + "\n")
		}
		return
	},
}

func init() {
	collapseCmd.AddCommand(collapsesingleCmd)
}
