package cmd

import (
	"github.com/spf13/cobra"

	"github.com/fredericlemoine/gotree/io"
	"github.com/fredericlemoine/gotree/tree"
)

var multiplylengthfactor float64

// clearlengthCmd represents the clearlength command
var multiplylengthCmd = &cobra.Command{
	Use:   "multiply",
	Short: "Multiply lengths from input trees by a given factor",
	Long:  `Multiply lengths from input trees by a given factor.`,
	Run: func(cmd *cobra.Command, args []string) {
		f := openWriteFile(outtreefile)
		treefile, treechan := readTrees(intreefile)
		defer treefile.Close()
		for tr := range treechan {
			if tr.Err != nil {
				io.ExitWithMessage(tr.Err)
			}
			for _, e := range tr.Tree.Edges() {
				if e.Length() != tree.NIL_LENGTH {
					e.SetLength(e.Length() * multiplylengthfactor)
				}
			}
			f.WriteString(tr.Tree.Newick() + "\n")
		}
		f.Close()
	},
}

func init() {
	brlenCmd.AddCommand(multiplylengthCmd)
	multiplylengthCmd.Flags().Float64VarP(&multiplylengthfactor, "factor", "f", 1.0, "Branch length multiplication factor")
}
