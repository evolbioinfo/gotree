package cmd

import (
	"fmt"
	goio "io"
	"os"

	"github.com/evolbioinfo/gotree/io"
	"github.com/evolbioinfo/gotree/tree"
	"github.com/spf13/cobra"
)

var brnamefile string
var brid bool

// collapseCmd represents the collapse command
var collapsebrnameCmd = &cobra.Command{
	Use:   "name",
	Short: "Collapse branches having given name or ID",
	Long: `Collapse branches having given name or ID.

	Names (or ID) are defined in an input file (-b)

	If an external branch name/id is given, then does not do anything.
`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		var f *os.File
		var treefile goio.Closer
		var treechan <-chan tree.Trees
		var brnames []string
		var brids []int
		var toremove []*tree.Edge

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

		if brid {

			if brids, err = parseIntFile(brnamefile); err != nil {
				io.LogError(err)
				return
			}
		} else {
			if brnames, err = parseStringFile(brnamefile); err != nil {
				io.LogError(err)
				return
			}
		}

		for t := range treechan {
			if t.Err != nil {
				io.LogError(t.Err)
				return t.Err
			}
			t.Tree.ReinitIndexes()
			alledges := t.Tree.Edges()
			if brid {
				for _, i := range brids {
					if i < 0 || i >= len(alledges) {
						err = fmt.Errorf("branch index is not in the tree (<0 or >#branches)")
						io.LogError(err)
						return
					}
					toremove = append(toremove, alledges[i])
				}
			} else {
				for _, n := range brnames {
					found := false
					for _, e := range alledges {
						if e.Name(t.Tree.Rooted()) == n {
							toremove = append(toremove, e)
							found = true
						}
					}
					if !found {
						err = fmt.Errorf("branch name %s not found in the tree", n)
						io.LogError(err)
						return
					}
				}
			}
			t.Tree.RemoveEdges(true, false, toremove...)
			f.WriteString(t.Tree.Newick() + "\n")
		}
		return
	},
}

func init() {
	collapseCmd.AddCommand(collapsebrnameCmd)
	collapsebrnameCmd.Flags().StringVarP(&brnamefile, "brfile", "b", "none", "File with one branch name/id per line")
	collapsebrnameCmd.Flags().BoolVar(&brid, "id", false, "Input file contains branch ids (otherwise, branch names)")
}
