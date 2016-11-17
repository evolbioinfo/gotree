package cmd

import (
	"fmt"
	"github.com/fredericlemoine/gotree/io"
	"github.com/fredericlemoine/gotree/io/utils"
	"github.com/fredericlemoine/gotree/tree"
	"github.com/spf13/cobra"
	"os"
	"sync"
)

var edgeInTree string
var edgeOutFile string

type EdgeStruct struct {
	e   *tree.Edge
	idx int
}

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

		edges := make(chan EdgeStruct, 1000)

		go func() {
			for i, e := range t.Edges() {
				edges <- EdgeStruct{e, i}
			}
		}()

		var wg sync.WaitGroup
		for cpu := 0; cpu < rootCpus; cpu++ {
			wg.Add(1)
			go func() {
				for edgeS := range edges {
					if !edgeS.e.Right().Tip() {
						var edgeOut *os.File

						if edgeOutFile == "stdout" {
							edgeOut = openWriteFile(edgeOutFile)
						} else {
							edgeOut = openWriteFile(fmt.Sprintf("%s_%04d.nw", edgeOutFile, edgeS.idx))
						}
						edgeTree := tree.EdgeTree(t, edgeS.e, alltips)

						// We build a new Tree with a single edge
						edgeOut.WriteString(edgeTree.Newick() + "\n")
						if edgeOutFile != "stdout" {
							edgeOut.Close()
						}
					}
				}
				wg.Done()
			}()
		}
		wg.Wait()
	},
}

func init() {
	computeCmd.AddCommand(edgeTreesCmd)
	edgeTreesCmd.PersistentFlags().StringVarP(&edgeInTree, "reftree", "i", "stdin", "Reference tree input file")
	edgeTreesCmd.PersistentFlags().StringVarP(&edgeOutFile, "out", "o", "stdout", "Output tree files prefix")
}
