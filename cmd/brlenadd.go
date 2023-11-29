package cmd

import (
	goio "io"
	"os"

	"github.com/evolbioinfo/gotree/io"
	"github.com/evolbioinfo/gotree/tree"
	"github.com/spf13/cobra"
)

var addlengthfactor float64

// addCmd represents the cut command
var brlenAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add the given length to all branches of the tree.",
	Long: `Add the given length to all branches of the tree.

Example:

gotree brlen add -i tree.nwk -l <length>
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
			tr.Tree.AddLength(addlengthfactor, brleninternal, brlenexternal)
			f.WriteString(tr.Tree.Newick() + "\n")
		}
		return
	},
}

func init() {
	brlenCmd.AddCommand(brlenAddCmd)
	brlenAddCmd.Flags().Float64VarP(&addlengthfactor, "add-length", "l", 0.0, "Length to add to all branches")
	brlenAddCmd.PersistentFlags().StringVarP(&outtreefile, "output", "o", "stdout", "Output tree file")
}
