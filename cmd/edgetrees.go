package cmd

import (
	"errors"
	"fmt"
	"github.com/fredericlemoine/gotree/io"
	"github.com/fredericlemoine/gotree/io/utils"
	"github.com/fredericlemoine/gotree/tree"
	"github.com/spf13/cobra"
	"os"
)

var edgeInTree string
var edgeOutFile string

// edgeTreesCmd represents the edgeTrees command
var edgeTreesCmd = &cobra.Command{
	Use:   "edgetrees",
	Short: "For each edge of the input tree, builds a tree with only this edge",
	Long: `For each edge of the input tree, builds a tree with only this edge.

The resulting trees are star trees to which we added one biparition. All branch lengths are set to 1.
`,
	Run: func(cmd *cobra.Command, args []string) {
		var t *tree.Tree
		var err error
		t, err = utils.ReadRefTree(edgeInTree)
		if err != nil {
			io.ExitWithMessage(err)
		}
		alltips := t.AllTipNames()
		for i, e := range t.Edges() {
			var edgeOut *os.File
			// We build a star Tree
			if startree, err := tree.StarTreeFromName(alltips...); err != nil {
				io.ExitWithMessage(err)
			} else {
				if !e.Right().Tip() {
					if edgeOutFile == "stdout" {
						edgeOut = openWriteFile(edgeOutFile)
					} else {
						edgeOut = openWriteFile(fmt.Sprintf("%s_%04d.nw", edgeOutFile, i))
					}

					nodeindex := tree.NewNodeIndex(startree)
					names := make([]string, 0, e.Bitset().Count())
					for _, n := range alltips {
						if idx, err := startree.TipIndex(n); err != nil {
							io.ExitWithMessage(err)
						} else {
							if e.Bitset().Test(idx) {
								names = append(names, n)
							}
						}
					}
					// We add the bipartition to the tree and write the tree into a file
					node, edges, monophyletic := startree.LeastCommonAncestorUnrooted(nodeindex, names...)
					if node == nil {
						io.ExitWithMessage(errors.New("EdgeTree error: No common ancestor found for biparition"))
					}
					if edges == nil || len(edges) == 0 {
						io.ExitWithMessage(errors.New("EdgeTree error: No common ancestor Edges found"))
					}
					if !monophyletic {
						io.ExitWithMessage(errors.New("The group should be monophyletic"))
					}
					// We add the bipartition with a support value corresponding to the percentage of
					// trees in which it appears
					// TODO: Average branch length : Need to change the data structure
					startree.AddBipartition(node, edges, 1.0, -1.0)
					edgeOut.WriteString(startree.Newick() + "\n")
					if edgeOutFile != "stdout" {
						edgeOut.Close()
					}
				}
			}
		}
	},
}

func init() {
	computeCmd.AddCommand(edgeTreesCmd)
	edgeTreesCmd.PersistentFlags().StringVarP(&edgeInTree, "reftree", "i", "stdin", "Reference tree input file")
	edgeTreesCmd.PersistentFlags().StringVarP(&edgeOutFile, "out", "o", "stdout", "Output tree files prefix")
}
