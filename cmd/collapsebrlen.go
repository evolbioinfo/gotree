package cmd

import (
	goio "io"
	"os"

	"github.com/evolbioinfo/gotree/io"
	"github.com/evolbioinfo/gotree/tree"
	"github.com/spf13/cobra"
)

var shortbranchesThreshold float64
var shortbranchesRemoveRoot bool
var shortbranchesRemoveTips bool

// collapseCmd represents the collapse command
var collapsebrlenCmd = &cobra.Command{
	Use:   "length",
	Short: "Collapse short branches of the input tree",
	Long: `Collapse short branches of the input tree.

Short branches are defined by a threshold (-l). All internal branches 
with length <= threshold are removed.

If --root is given, then it applies also to internal branches connected to the root in the case 
of rooted trees. This may unroot the tree. In that case, so far the two branches connected to the root are 
considered independently whereas it may be more useful to consider them as a single bipartition if the 
tree is going to be unrooted.

If --tips is given, then it applies also to external branches, just by setting their length to 0.0
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
			t.Tree.CollapseShortBranches(shortbranchesThreshold, shortbranchesRemoveRoot, shortbranchesRemoveTips)
			f.WriteString(t.Tree.Newick() + "\n")
		}
		return
	},
}

func init() {
	collapseCmd.AddCommand(collapsebrlenCmd)
	collapsebrlenCmd.Flags().Float64VarP(&shortbranchesThreshold, "length", "l", 0.0, "Length cutoff to collapse branches")
	collapsebrlenCmd.Flags().BoolVar(&shortbranchesRemoveRoot, "root", false, "Applies also to branches connected to the root (may unroot the tree)")
	collapsebrlenCmd.Flags().BoolVar(&shortbranchesRemoveTips, "tips", false, "Applies also to tips (keeps a 0.0 length tip)")
}
