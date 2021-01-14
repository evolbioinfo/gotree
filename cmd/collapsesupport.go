package cmd

import (
	goio "io"
	"os"

	"github.com/evolbioinfo/gotree/io"
	"github.com/evolbioinfo/gotree/tree"
	"github.com/spf13/cobra"
)

var lowSupportThreshold float64
var supportRemoveRoot bool

// collapsesupportCmd represents the collapsesupport command
var collapsesupportCmd = &cobra.Command{
	Use:   "support",
	Short: "Collapse lowly supported branches of the input tree",
	Long: `Collapse lowly supported branches of the input tree.

Lowly supported branches are defined by a threshold (-s). All internal branches 
with support < threshold and that are not connected to the root in case of rooted tree
 are removed.

 If --root is given, then it applies also to internal branches connected to the root in the case 
 of rooted trees. This may unroot the tree.
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

		treefile, treechan, err = readTrees(intreefile)
		defer treefile.Close()
		for t := range treechan {
			if t.Err != nil {
				io.LogError(t.Err)
				return t.Err
			}
			t.Tree.CollapseLowSupport(lowSupportThreshold, supportRemoveRoot)
			f.WriteString(t.Tree.Newick() + "\n")
		}
		return
	},
}

func init() {
	collapseCmd.AddCommand(collapsesupportCmd)
	collapsesupportCmd.Flags().Float64VarP(&lowSupportThreshold, "support", "s", 0.0, "Support cutoff to collapse branches")
	collapsesupportCmd.Flags().BoolVar(&supportRemoveRoot, "root", false, "Applies also to branches connected to the root (may unroot the tree)")

}
