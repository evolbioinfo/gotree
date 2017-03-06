package cmd

import (
	"fmt"

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
	Run: func(cmd *cobra.Command, args []string) {
		f := openWriteFile(outtreefile)
		for statsintree := range readTrees(intreefile) {
			f.WriteString("Tree\t")
			names := statsintree.Tree.SortedTips()
			for i := len(names) - 1; i >= 0; i-- {
				if i < len(names)-1 {
					f.WriteString("|")
				}
				f.WriteString(names[i])
			}
			f.WriteString("\n")
			for _, e := range statsintree.Tree.Edges() {
				f.WriteString(fmt.Sprintf("%d\t", statsintree.Id))
				f.WriteString(e.DumpBitSet() + "\n")
			}
		}
		f.Close()
	},
}

func init() {
	statsCmd.AddCommand(splitsCmd)
}
