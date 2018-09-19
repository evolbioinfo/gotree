package cmd

import (
	"github.com/fredericlemoine/gotree/io"
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
	Run: func(cmd *cobra.Command, args []string) {
		f := openWriteFile(outtreefile)
		defer closeWriteFile(f, outtreefile)

		treefile, treechan := readTrees(intreefile)
		defer treefile.Close()
		for t := range treechan {
			if t.Err != nil {
				io.ExitWithMessage(t.Err)
			}
			t.Tree.RemoveSingleNodes()
			f.WriteString(t.Tree.Newick() + "\n")
		}
	},
}

func init() {
	collapseCmd.AddCommand(collapsesingleCmd)
}
