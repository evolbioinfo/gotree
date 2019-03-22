package cmd

import (
	"fmt"
	goio "io"
	"os"

	"github.com/evolbioinfo/gotree/io"
	"github.com/evolbioinfo/gotree/tree"
	"github.com/spf13/cobra"
)

// splitsCmd represents the splits command
var splitsCmd = &cobra.Command{
	Use:   "splits",
	Short: "Prints all the splits from an input tree",
	Long: `Prints all the splits from an input tree.

First line : List of taxa
Then: One line per branch, and 0/1 
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
			if t.Err != nil {
				io.LogError(t.Err)
				return t.Err
			}
			t.Tree.ReinitIndexes()
			f.WriteString("Tree\t")
			names := t.Tree.SortedTips()
			for i := len(names) - 1; i >= 0; i-- {
				if i < len(names)-1 {
					f.WriteString("|")
				}
				f.WriteString(names[i])
			}
			f.WriteString("\n")
			for _, e := range t.Tree.Edges() {
				f.WriteString(fmt.Sprintf("%d\t", t.Id))
				f.WriteString(e.DumpBitSet() + "\n")
			}
		}
		return
	},
}

func init() {
	statsCmd.AddCommand(splitsCmd)
}
