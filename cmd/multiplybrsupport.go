package cmd

import (
	"github.com/spf13/cobra"

	"github.com/fredericlemoine/gotree/io"
)

var scalesupportfactor float64

// clearsupportCmd represents the support scale command
var scalesupportCmd = &cobra.Command{
	Use:   "scale",
	Short: "Scale branch supports from input trees by a given factor",
	Long:  `Scale branch supports from input trees by a given factor.`,
	Run: func(cmd *cobra.Command, args []string) {
		f := openWriteFile(outtreefile)
		defer closeWriteFile(f, outtreefile)

		treefile, treechan := readTrees(intreefile)
		defer treefile.Close()
		for tr := range treechan {
			if tr.Err != nil {
				io.ExitWithMessage(tr.Err)
			}
			tr.Tree.ScaleSupports(scalesupportfactor)
			f.WriteString(tr.Tree.Newick() + "\n")
		}
	},
}

func init() {
	supportCmd.AddCommand(scalesupportCmd)
	scalesupportCmd.Flags().Float64VarP(&scalesupportfactor, "factor", "f", 1.0, "Branch support scaling factor")
}
