package cmd

import (
	"github.com/fredericlemoine/gotree/io"
	"github.com/spf13/cobra"
)

// minbrlenCmd represents the minbrlen command
var minbrlenCmd = &cobra.Command{
	Use:   "setmin",
	Short: "Set a min branch length to all branches with length < cutoff",
	Long: `Set a min branch length to all branches with length < cutoff

Example of usage:

gotree minbrlen -i tree.nw -o out.nw -l 0.001

`,
	Run: func(cmd *cobra.Command, args []string) {
		f := openWriteFile(outtreefile)
		treefile, trees := readTrees(intreefile)
		defer treefile.Close()

		for t := range trees {
			if t.Err != nil {
				io.ExitWithMessage(t.Err)
			}
			for _, e := range t.Tree.Edges() {
				if e.Length() < cutoff {
					e.SetLength(cutoff)
				}
			}
			f.WriteString(t.Tree.Newick() + "\n")
		}
		f.Close()

	},
}

func init() {
	brlenCmd.AddCommand(minbrlenCmd)
	minbrlenCmd.Flags().Float64VarP(&cutoff, "length", "l", 0.0, "Min Length cutoff")
	minbrlenCmd.PersistentFlags().StringVarP(&intreefile, "input", "i", "stdin", "Input tree")
	minbrlenCmd.PersistentFlags().StringVarP(&outtreefile, "output", "o", "stdout", "Min length output tree file")
}
