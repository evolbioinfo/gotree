package cmd

import (
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
	Run: func(cmd *cobra.Command, args []string) {
		for t := range collapseIntrees {
			t.Tree.CollapseTopoDepth(mindepthThreshold, maxdepthThreshold)
			collapseOutTrees.WriteString(t.Tree.Newick() + "\n")
		}
	},
}

func init() {
	collapseCmd.AddCommand(collapsedepthCmd)

	collapsedepthCmd.Flags().IntVarP(&mindepthThreshold, "min-depth", "m", 0, "Min depth cutoff to collapse branches")
	collapsedepthCmd.Flags().IntVarP(&maxdepthThreshold, "max-depth", "M", 0, "Max Depth cutoff to collapse branches")
}
