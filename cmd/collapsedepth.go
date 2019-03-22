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

// collapsedepthCmd represents the collapsedepth command
var collapsedepthCmd = &cobra.Command{
	Use:   "depth",
	Short: "Collapse branches having a given depth",
	Long: `Collapse branches having a given depth.

Branches having depth (number of taxa on the lightest side of 
the bipartition) d such that:

min-depth<=d<=max-depth

will be collapsed.

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
			t.Tree.ReinitIndexes()
			if t.Err != nil {
				io.LogError(t.Err)
				return t.Err
			}
			t.Tree.CollapseTopoDepth(mindepthThreshold, maxdepthThreshold)
			f.WriteString(t.Tree.Newick() + "\n")
		}
		return
	},
}

func init() {
	collapseCmd.AddCommand(collapsedepthCmd)

	collapsedepthCmd.Flags().IntVarP(&mindepthThreshold, "min-depth", "m", 0, "Min depth cutoff to collapse branches")
	collapsedepthCmd.Flags().IntVarP(&maxdepthThreshold, "max-depth", "M", 0, "Max Depth cutoff to collapse branches")
}
