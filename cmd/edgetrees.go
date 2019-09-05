package cmd

import (
	"bytes"
	"fmt"
	"os"
	"sync"

	"github.com/evolbioinfo/gotree/io"
	"github.com/evolbioinfo/gotree/tree"
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
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		var t *tree.Tree

		if t, err = readTree(intreefile); err != nil {
			io.LogError(err)
			return
		}
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
						err = er
						return
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
						var err2 error
						var bitsetindex uint

						if outtreefile == "stdout" || outtreefile == "-" {
							if edgeOut, err2 = openWriteFile("stdout"); err2 != nil {
								err = err2
								return
							}
						} else {
							if edgeOut, err2 = openWriteFile(fmt.Sprintf("%s_%06d.nw", outtreefile, edgeS.idx)); err2 != nil {
								err = err2
								return
							}
						}

						if edgeformattext {
							var leftbuffer bytes.Buffer
							var rightbuffer bytes.Buffer
							leftnb := 0
							rightnb := 0
							for _, n := range alltips {
								if bitsetindex, err2 = t.TipIndex(n); err2 != nil {
									err = err2
									return
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
		return
	},
}

func init() {
	computeCmd.AddCommand(edgeTreesCmd)
	edgeTreesCmd.PersistentFlags().StringVarP(&intreefile, "reftree", "i", "stdin", "Reference tree input file")
	edgeTreesCmd.PersistentFlags().StringVarP(&outtreefile, "out", "o", "stdout", "Output tree files prefix")
	edgeTreesCmd.PersistentFlags().BoolVar(&deepestedge, "deepest", false, "Output a tree only for the deepest bipartition")
	edgeTreesCmd.PersistentFlags().BoolVar(&edgeformattext, "text-format", false, "Output bipartitions in the form t1,t2,...,tp|tp+1,...tn instead of newick")

}
