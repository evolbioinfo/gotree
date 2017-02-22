package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

// statsCmd represents the stats command
var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Print statistics about the tree",
	Long: `Print statistics about the tree

For example:
- Edge informations
- Node informations
- Tips informations

`,
	Run: func(cmd *cobra.Command, args []string) {
		/* Dividing trees */
		f := openWriteFile(outtreefile)
		f.WriteString("tree\tnodes\ttips\tedges\tmeanbrlen\tsumbrlen\tmeansupport\tmediansupport\trooted\n")
		for statsintree := range readTrees(intreefile) {
			statsintree.Tree.ComputeDepths()
			f.WriteString(fmt.Sprintf("%d", statsintree.Id))
			f.WriteString(fmt.Sprintf("\t%d", len(statsintree.Tree.Nodes())))
			f.WriteString(fmt.Sprintf("\t%d", len(statsintree.Tree.Tips())))
			f.WriteString(fmt.Sprintf("\t%d", len(statsintree.Tree.Edges())))
			f.WriteString(fmt.Sprintf("\t%.8f", statsintree.Tree.MeanBrLength()))
			f.WriteString(fmt.Sprintf("\t%.8f", statsintree.Tree.SumBranchLengths()))
			f.WriteString(fmt.Sprintf("\t%.8f", statsintree.Tree.MeanSupport()))
			f.WriteString(fmt.Sprintf("\t%.8f", statsintree.Tree.MedianSupport()))
			if statsintree.Tree.Rooted() {
				f.WriteString(fmt.Sprintf("\trooted\n"))
			} else {
				f.WriteString(fmt.Sprintf("\tunrooted\n"))
			}
		}
		f.Close()
	},
}

func init() {
	RootCmd.AddCommand(statsCmd)
	statsCmd.PersistentFlags().StringVarP(&intreefile, "input", "i", "stdin", "Input tree")
	statsCmd.PersistentFlags().StringVarP(&outtreefile, "output", "o", "stdout", "Output file")
}
