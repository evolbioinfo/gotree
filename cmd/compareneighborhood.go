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

// compareCmd represents the compare command
var compareNeighborhoodCmd = &cobra.Command{
	Use:   "neighborhood",
	Short: "Compare tip neighborhoods of a reference tree to a compared tree",
	Long: `Compare tip neighborhoods of a reference tree to a compared tree.
`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		var treefile goio.Closer
		var treechan <-chan tree.Trees
		var refTree *tree.Tree
		var stats []tree.TipNeighborhoodStats

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

		fmt.Printf("Tree\tTip\tPercent\tJacquard\tInter\tUnion\n")
		for comptree := range treechan {
			if err = comptree.Tree.ReinitIndexes(); err != nil {
				io.LogError(err)
				return
			}

			if stats, err = tree.CompareTipNeighborhood(refTree, comptree.Tree); err != nil {
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
}
