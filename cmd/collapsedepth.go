package cmd

import (
	goio "io"
	"os"

	"github.com/evolbioinfo/gotree/io"
	"github.com/evolbioinfo/gotree/tree"
	"github.com/spf13/cobra"
)

var maxdepthThreshold int
var mindepthThreshold int
var collapseDepthRoot bool
var collapseDepthTips bool

// collapsedepthCmd represents the collapsedepth command
var collapsedepthCmd = &cobra.Command{
	Use:   "depth",
	Short: "Collapse branches having a given depth",
	Long: `Collapse branches having a given depth.

Removes internal branches (not connected to the root in case of rooted trees) 
having depth (number of taxa on the lightest side of the bipartition) d such that:

min-depth<=d<=max-depth

will be collapsed.

If --root is given, then it applies also to internal branches connected to the root in the case 
of rooted trees. This may unroot the tree.

If --tips is given, then it applies also to external branches (if min-depth<=1), just by setting their length to 0.0
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
			if err = t.Tree.ReinitIndexes(); err != nil {
				io.LogError(err)
				return
			}

			t.Tree.CollapseTopoDepth(mindepthThreshold, maxdepthThreshold, collapseDepthRoot, collapseDepthTips)
			f.WriteString(t.Tree.Newick() + "\n")
		}
		return
	},
}

func init() {
	collapseCmd.AddCommand(collapsedepthCmd)

	collapsedepthCmd.Flags().IntVarP(&mindepthThreshold, "min-depth", "m", 0, "Min depth cutoff to collapse branches")
	collapsedepthCmd.Flags().IntVarP(&maxdepthThreshold, "max-depth", "M", 0, "Max Depth cutoff to collapse branches")
	collapsedepthCmd.Flags().BoolVar(&collapseDepthRoot, "root", false, "Applies also to branches connected to the root (may unroot the tree)")
	collapsedepthCmd.Flags().BoolVar(&collapseDepthTips, "tips", false, "Applies also to tips (keeps a 0.0 length tip)")
}
