package cmd

import (
	"github.com/spf13/cobra"

	"github.com/fredericlemoine/gotree/io"
)

var scalelengthfactor float64

// clearlengthCmd represents the clearlength command
var scalelengthCmd = &cobra.Command{
	Use:   "scale",
	Short: "Scale lengths from input trees by a given factor",
	Long:  `Scale lengths from input trees by a given factor.`,
	Run: func(cmd *cobra.Command, args []string) {
		f := openWriteFile(outtreefile)
		defer closeWriteFile(f, outtreefile)

		treefile, treechan := readTrees(intreefile)
		defer treefile.Close()
		for tr := range treechan {
			if tr.Err != nil {
				io.ExitWithMessage(tr.Err)
			}
			tr.Tree.ScaleLengths(scalelengthfactor)
			f.WriteString(tr.Tree.Newick() + "\n")
		}
	},
}

func init() {
	brlenCmd.AddCommand(scalelengthCmd)
	scalelengthCmd.Flags().Float64VarP(&scalelengthfactor, "factor", "f", 1.0, "Branch length scaling factor")
	scalelengthCmd.PersistentFlags().StringVarP(&outtreefile, "output", "o", "stdout", "Scaled length output tree file")
}
