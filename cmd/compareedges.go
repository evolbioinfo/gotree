package cmd

import (
	"fmt"
	"github.com/fredericlemoine/gotree/io"
	"github.com/fredericlemoine/gotree/io/utils"
	"github.com/fredericlemoine/gotree/support"
	"github.com/fredericlemoine/gotree/tree"
	"github.com/spf13/cobra"
	"os"
)

/* If transfer dist should be also given in output */
var edgesMastDist bool

// compareedgesCmd represents the compareedges command
var compareedgesCmd = &cobra.Command{
	Use:   "edges",
	Short: "Compare edges of a reference tree with another tree",
	Long: `Compare edges of a reference tree with another tree

If the compared tree file contains several trees, it will take the first one only
`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Fprintf(os.Stderr, "Reference : %s\n", compareTree1)
		fmt.Fprintf(os.Stderr, "Compared  : %s\n", compareTree2)
		var err error
		var refTree *tree.Tree
		if refTree, err = utils.ReadRefTree(compareTree1); err != nil {
			io.ExitWithMessage(err)
		}
		refTree.ComputeDepths()

		nbtrees := 0
		compareChannel := make(chan tree.Trees, 15)

		go func() {
			if nbtrees, err = utils.ReadCompTrees(compareTree2, compareChannel); err != nil {
				io.ExitWithMessage(err)
			}
		}()

		edges1 := refTree.Edges()
		fmt.Printf("tree\tbrid\tlength\tsupport\tterminal\tdepth\ttopodepth\trightname\tfound")
		if edgesMastDist {
			fmt.Printf("\tmast")
		}
		fmt.Printf("\n")
		for t2 := range compareChannel {

			edges2 := t2.Tree.Edges()

			var min_dist []uint16
			var min_dist_edges []int
			if edgesMastDist {
				tips := refTree.Tips()
				min_dist = make([]uint16, len(edges1))
				min_dist_edges = make([]int, len(tips))
				var i_matrix [][]uint16 = make([][]uint16, len(edges1))
				var c_matrix [][]uint16 = make([][]uint16, len(edges1))
				var hamming [][]uint16 = make([][]uint16, len(edges1))

				for i, e := range edges1 {
					e.SetId(i)
					min_dist[i] = uint16(len(tips))
					i_matrix[i] = make([]uint16, len(edges2))
					c_matrix[i] = make([]uint16, len(edges2))
					hamming[i] = make([]uint16, len(edges2))
				}
				for i, e := range edges2 {
					e.SetId(i)
				}
				support.Update_all_i_c_post_order_ref_tree(refTree, &edges1, t2.Tree, &edges2, &i_matrix, &c_matrix)
				support.Update_all_i_c_post_order_boot_tree(refTree, uint(len(tips)), &edges1, t2.Tree, &edges2, &i_matrix, &c_matrix, &hamming, &min_dist, &min_dist_edges)
			}

			for i, e1 := range edges1 {
				found := false
				for _, e2 := range edges2 {
					if e1.SameBipartition(e2) {
						found = true
						break
					}
				}
				fmt.Printf("%d\t%d\t%s\t%t", t2.Id, i, e1.ToStatsString(), found)
				if edgesMastDist {
					fmt.Printf("\t%d", min_dist[e1.Id()])
				}
				fmt.Printf("\n")
			}
		}
	},
}

func init() {
	compareCmd.AddCommand(compareedgesCmd)
	compareedgesCmd.PersistentFlags().BoolVarP(&edgesMastDist, "mast-dist", "m", false, "If mast dist must be computed for each edge")
}
