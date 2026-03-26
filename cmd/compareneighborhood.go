package cmd

import (
	"errors"
	"fmt"
	goio "io"
	"runtime"

	"github.com/spf13/cobra"

	"github.com/evolbioinfo/gotree/io"
	"github.com/evolbioinfo/gotree/tree"
)

var compareNeighborhoodMetric string

// compareCmd represents the compare command
var compareNeighborhoodCmd = &cobra.Command{
	Use:   "neighborhood",
	Short: "Compare tip neighborhoods of a reference tree to a compared tree",
	Long: `Compare tip neighborhoods of a reference tree to compared trees.

For each tip of the reference tree and for each compared tree, this command
computes Jaccard similarity between neighborhoods defined by increasing
percentages of closest tips (from 1% to 100%).

Neighborhood distance can be computed with:
* --metric brlen : patristic distance (sum of branch lengths)
* --metric boot  : sum of branch supports
* --metric none  : topological distance (all branches have distance 1)

Output columns (tab-separated):
Tree, Tip, Percent, Jacquard, Inter, Union
`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		var treefile goio.Closer
		var treechan <-chan tree.Trees
		var refTree *tree.Tree
		var stats []tree.TipNeighborhoodStats
		var distMetric int

		if intree2file == "none" {
			err = errors.New("You must provide a file containing compared trees")
			io.LogError(err)
			return
		}

		maxcpus := runtime.NumCPU()
		if rootCpus > maxcpus {
			rootCpus = maxcpus
		}
		if refTree, err = readTree(intreefile); err != nil {
			io.LogError(err)
			return
		}

		if err = refTree.ReinitIndexes(); err != nil {
			io.LogError(err)
		}

		if treefile, treechan, err = readTrees(intree2file); err != nil {
			io.LogError(err)
			return
		}
		defer treefile.Close()

		switch compareNeighborhoodMetric {
		case "brlen":
			distMetric = tree.DISTANCE_METRIC_BRLEN
		case "boot":
			distMetric = tree.DISTANCE_METRIC_BOOTS
		case "none":
			distMetric = tree.DISTANCE_METRIC_NONE
		default:
			err = fmt.Errorf("distance metric %s in not supported", compareNeighborhoodMetric)
			io.LogError(err)
			return
		}

		fmt.Printf("Tree\tTip\tPercent\tJacquard\tInter\tUnion\n")
		for comptree := range treechan {
			if err = comptree.Tree.ReinitIndexes(); err != nil {
				io.LogError(err)
				return
			}

			if stats, err = tree.CompareTipNeighborhood(refTree, comptree.Tree, distMetric); err != nil {
				io.LogError(err)
				return
			}
			for _, s := range stats {
				for v := range s.PercentTips {
					fmt.Printf("%d\t%s\t%d\t%f\t%d\t%d\n", comptree.Id, s.Id, s.PercentTips[v], s.Jacquard[v], s.Common[v], s.Union[v])
				}
			}
		}

		return
	},
}

func init() {
	compareCmd.AddCommand(compareNeighborhoodCmd)
	compareNeighborhoodCmd.Flags().StringVarP(&compareNeighborhoodMetric, "metric", "m", "brlen", "Distance metric (brlen|boot|none)")
}
