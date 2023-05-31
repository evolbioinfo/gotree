package cmd

import (
	goio "io"
	"os"

	"github.com/evolbioinfo/gotree/io"
	"github.com/evolbioinfo/gotree/tree"
	"github.com/spf13/cobra"
)

var scalelengthfactor float64

// clearlengthCmd represents the clearlength command
var scalelengthCmd = &cobra.Command{
	Use:   "scale",
	Short: "Scale lengths from input trees by a given factor",
	Long: `Scale lengths from input trees by a given factor.
	
	if --internal=false is given, it won't apply to internal branches (only external)
	if --external=false is given, it won't apply to external branches (only internal)

`,
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
			tr.Tree.ScaleLengths(scalelengthfactor, brleninternal, brlenexternal)
			f.WriteString(tr.Tree.Newick() + "\n")
		}
		return
	},
}

func init() {
	brlenCmd.AddCommand(scalelengthCmd)
	scalelengthCmd.Flags().Float64VarP(&scalelengthfactor, "factor", "f", 1.0, "Branch length scaling factor")
	scalelengthCmd.PersistentFlags().StringVarP(&outtreefile, "output", "o", "stdout", "Scaled length output tree file")
}
