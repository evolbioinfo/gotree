package cmd

import (
	"github.com/spf13/cobra"
)

// minbrlenCmd represents the minbrlen command
var minbrlenCmd = &cobra.Command{
	Use:   "minbrlen",
	Short: "Set a min branch length to all branches with length < cutoff",
	Long: `Set a min branch length to all branches with length < cutoff

Example of usage:

gotree minbrlen -i tree.nw -o out.nw -l 0.001

`,
	Run: func(cmd *cobra.Command, args []string) {
		f := openWriteFile(outtreefile)
		for tr := range readTrees(intreefile) {
			for _, e := range tr.Tree.Edges() {
				if e.Length() < cutoff {
					e.SetLength(cutoff)
				}
			}
			f.WriteString(tr.Tree.Newick() + "\n")
		}
		f.Close()

	},
}

func init() {
	RootCmd.AddCommand(minbrlenCmd)
	minbrlenCmd.Flags().Float64VarP(&cutoff, "length", "l", 0.0, "Min Length cutoff")
	minbrlenCmd.PersistentFlags().StringVarP(&intreefile, "input", "i", "stdin", "Input tree")
	minbrlenCmd.PersistentFlags().StringVarP(&outtreefile, "output", "o", "stdout", "Length corrected tree output file")
}
