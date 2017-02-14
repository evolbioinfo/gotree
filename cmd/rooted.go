package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

// rootedCmd represents the rooted command
var rootedCmd = &cobra.Command{
	Use:   "rooted",
	Short: "Tells wether the tree is rooted or unrooted",
	Long: `Tells wether the tree is rooted or unrooted

Example of usage:

gotree stats rooted -i t.nw

`,
	Run: func(cmd *cobra.Command, args []string) {
		f := openWriteFile(outtreefile)
		f.WriteString("tree\trooted\n")
		for statsintree := range readTrees(intreefile) {
			f.WriteString(fmt.Sprintf("%d\t", statsintree.Id))
			if statsintree.Tree.Rooted() {
				f.WriteString("rooted\n")
			} else {
				f.WriteString("unrooted\n")
			}
		}
		f.Close()
	},
}

func init() {
	statsCmd.AddCommand(rootedCmd)
}
