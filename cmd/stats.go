package cmd

import (
	"fmt"

	"github.com/fredericlemoine/gotree/io"
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
		f.WriteString("tree\tnodes\ttips\tedges\tmeanbrlen\tsumbrlen\tmeansupport\tmediansupport\trooted\tnbcherries\tcolless\tsackin\n")
		treefile, treechan := readTrees(intreefile)
		defer treefile.Close()
		for t := range treechan {
			if t.Err != nil {
				io.ExitWithMessage(t.Err)
			}
			t.Tree.ComputeDepths()
			f.WriteString(fmt.Sprintf("%d", t.Id))
			f.WriteString(fmt.Sprintf("\t%d", len(t.Tree.Nodes())))
			f.WriteString(fmt.Sprintf("\t%d", len(t.Tree.Tips())))
			f.WriteString(fmt.Sprintf("\t%d", len(t.Tree.Edges())))
			f.WriteString(fmt.Sprintf("\t%.8f", t.Tree.MeanBranchLength()))
			f.WriteString(fmt.Sprintf("\t%.8f", t.Tree.SumBranchLengths()))
			f.WriteString(fmt.Sprintf("\t%.8f", t.Tree.MeanSupport()))
			f.WriteString(fmt.Sprintf("\t%.8f", t.Tree.MedianSupport()))
			if t.Tree.Rooted() {
				f.WriteString(fmt.Sprintf("\trooted"))
			} else {
				f.WriteString(fmt.Sprintf("\tunrooted"))
			}
			f.WriteString(fmt.Sprintf("\t%d", t.Tree.NbCherries()))
			f.WriteString(fmt.Sprintf("\t%d", t.Tree.CollessIndex()))
			f.WriteString(fmt.Sprintf("\t%d\n", t.Tree.SackinIndex()))
		}
		f.Close()
	},
}

func init() {
	RootCmd.AddCommand(statsCmd)
	statsCmd.PersistentFlags().StringVarP(&intreefile, "input", "i", "stdin", "Input tree")
	statsCmd.PersistentFlags().StringVarP(&outtreefile, "output", "o", "stdout", "Output file")
}
