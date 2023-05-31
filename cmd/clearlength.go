package cmd

import (
	goio "io"
	"os"

	"github.com/evolbioinfo/gotree/io"
	"github.com/evolbioinfo/gotree/tree"
	"github.com/spf13/cobra"
)

// clearlengthCmd represents the clearlength command
var clearlengthCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clear lengths from input trees",
	Long: `Clear lengths from input trees.
	
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
		for t := range treechan {
			if t.Err != nil {
				io.LogError(t.Err)
				return t.Err
			}
			t.Tree.ClearLengths(brleninternal, brlenexternal)
			f.WriteString(t.Tree.Newick() + "\n")
		}
		return
	},
}

func init() {
	brlenCmd.AddCommand(clearlengthCmd)
	clearlengthCmd.PersistentFlags().StringVarP(&outtreefile, "output", "o", "stdout", "Cleared tree output file")
}
