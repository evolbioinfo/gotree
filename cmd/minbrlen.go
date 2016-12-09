package cmd

import (
	"github.com/fredericlemoine/gotree/io"
	"github.com/fredericlemoine/gotree/io/utils"
	"github.com/fredericlemoine/gotree/tree"
	"github.com/spf13/cobra"
	"os"
)

var minbrlenCutoff float64
var minbrlenInputTree string
var minbrlenOutputTree string

// minbrlenCmd represents the minbrlen command
var minbrlenCmd = &cobra.Command{
	Use:   "minbrlen",
	Short: "Set a min branch length to all branches with length < cutoff",
	Long: `Set a min branch length to all branches with length < cutoff

Example of usage:

gotree minbrlen -i tree.nw -o out.nw -l 0.001

`,
	Run: func(cmd *cobra.Command, args []string) {
		// Read Tree
		var t *tree.Tree
		var err error
		t, err = utils.ReadRefTree(minbrlenInputTree)
		if err != nil {
			io.ExitWithMessage(err)
		}
		var f *os.File
		if minbrlenOutputTree != "stdout" {
			f, err = os.Create(minbrlenOutputTree)
		} else {
			f = os.Stdout
		}
		if err != nil {
			io.ExitWithMessage(err)
		}

		for _, e := range t.Edges() {
			if e.Length() < minbrlenCutoff {
				e.SetLength(minbrlenCutoff)
			}
		}
		if err != nil {
			io.ExitWithMessage(err)
		}

		f.WriteString(t.Newick() + "\n")
		f.Close()

	},
}

func init() {
	RootCmd.AddCommand(minbrlenCmd)
	minbrlenCmd.Flags().Float64VarP(&minbrlenCutoff, "length", "l", 0.0, "Min Length cutoff")
	minbrlenCmd.PersistentFlags().StringVarP(&minbrlenInputTree, "input", "i", "stdin", "Input tree")
	minbrlenCmd.PersistentFlags().StringVarP(&minbrlenOutputTree, "output", "o", "stdout", "Length corrected tree output file")
}
