package cmd

import (
	goio "io"
	"os"

	"github.com/evolbioinfo/gotree/io"
	"github.com/evolbioinfo/gotree/tree"
	"github.com/spf13/cobra"
)

var scalesupportfactor float64

// clearsupportCmd represents the support scale command
var scalesupportCmd = &cobra.Command{
	Use:   "scale",
	Short: "Scale branch supports from input trees by a given factor",
	Long:  `Scale branch supports from input trees by a given factor.`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		var f *os.File
		var treefile goio.Closer
		var treechan <-chan tree.Trees

		if f, err = openWriteFile(outtreefile); err != nil {
			io.LogError(err)
			return
		}
		defer closeWriteFile(f, outtreefile)

		if treefile, treechan, err = readTrees(intreefile); err != nil {
			io.LogError(err)
			return
		}
		defer treefile.Close()
		for tr := range treechan {
			if tr.Err != nil {
				io.LogError(tr.Err)
				return tr.Err
			}
			tr.Tree.ScaleSupports(scalesupportfactor)
			f.WriteString(tr.Tree.Newick() + "\n")
		}
		return
	},
}

func init() {
	supportCmd.AddCommand(scalesupportCmd)
	scalesupportCmd.Flags().Float64VarP(&scalesupportfactor, "factor", "f", 1.0, "Branch support scaling factor")
}
