package cmd

import (
	"bytes"
	"fmt"
	"os"
	"sync"

	"github.com/fredericlemoine/gotree/io"
	"github.com/fredericlemoine/gotree/tree"
	"github.com/spf13/cobra"
)

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
		t := readTree(intreefile)
		t.ReinitIndexes()
		alltips := t.AllTipNames()
		edges := make(chan EdgeStruct, 1000)

		go func() {
			if deepestedge {
				// We take the deepest edge and give it to the channel
				maxdepth := 0
				var maxedge *tree.Edge = nil
				maxid := -1
				for i, e := range t.Edges() {
					if d, er := e.TopoDepth(); er != nil {
						io.ExitWithMessage(er)
					} else {
						if d > maxdepth {
							maxdepth = d
							maxedge = e
							maxid = i
						}
					}
				}
				edges <- EdgeStruct{maxedge, maxid}
			} else {
				for i, e := range t.Edges() {
					edges <- EdgeStruct{e, i}
				}
			}
			close(edges)
		}()

		var wg sync.WaitGroup
		for cpu := 0; cpu < rootCpus; cpu++ {
			wg.Add(1)
			go func() {
				for edgeS := range edges {
					if !edgeS.e.Right().Tip() {
						var edgeOut *os.File

						if outtreefile == "stdout" || outtreefile == "-" {
							edgeOut = openWriteFile("stdout")
						} else {
							edgeOut = openWriteFile(fmt.Sprintf("%s_%06d.nw", outtreefile, edgeS.idx))
						}

						if edgeformattext {
							var leftbuffer bytes.Buffer
							var rightbuffer bytes.Buffer
							leftnb := 0
							rightnb := 0
							for _, n := range alltips {
								bitsetindex, err := t.TipIndex(n)
								if err != nil {
									io.ExitWithMessage(err)
								}
								if edgeS.e.TipPresent(bitsetindex) {
									if leftnb > 0 {
										leftbuffer.WriteRune(',')
									}
									leftbuffer.WriteString(n)
									leftnb++
								} else {
									if rightnb > 0 {
										rightbuffer.WriteRune(',')
									}
									rightbuffer.WriteString(n)
									rightnb++
								}
							}
							if leftnb > rightnb {
								edgeOut.WriteString(leftbuffer.String())
								edgeOut.WriteString("|")
								edgeOut.WriteString(rightbuffer.String() + "\n")
							} else {
								edgeOut.WriteString(rightbuffer.String())
								edgeOut.WriteString("|")
								edgeOut.WriteString(leftbuffer.String() + "\n")
							}
						} else {
							edgeTree := tree.EdgeTree(t, edgeS.e, alltips)
							// We build a new Tree with a single edge
							edgeOut.WriteString(edgeTree.Newick() + "\n")
						}
						closeWriteFile(edgeOut, outtreefile)
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
	edgeTreesCmd.PersistentFlags().StringVarP(&intreefile, "reftree", "i", "stdin", "Reference tree input file")
	edgeTreesCmd.PersistentFlags().StringVarP(&outtreefile, "out", "o", "stdout", "Output tree files prefix")
	edgeTreesCmd.PersistentFlags().BoolVar(&deepestedge, "deepest", false, "Output a tree only for the deepest bipartition")
	edgeTreesCmd.PersistentFlags().BoolVar(&edgeformattext, "text-format", false, "Output bipartitions in the form t1,t2,...,tp|tp+1,...tn instead of newick")

}
