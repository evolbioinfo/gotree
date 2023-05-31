package cmd

import (
	goio "io"
	"os"

	"github.com/evolbioinfo/gotree/io"
	"github.com/evolbioinfo/gotree/tree"
	"github.com/spf13/cobra"
)

// minbrlenCmd represents the minbrlen command
var minbrlenCmd = &cobra.Command{
	Use:   "setmin",
	Short: "Set a min branch length to all branches with length < cutoff",
	Long: `Set a min branch length to all branches with length < cutoff

Example of usage:

gotree brlen setmin -i tree.nw -o out.nw -l 0.001

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
			for _, e := range t.Tree.Edges() {
				if ((e.Right().Tip() && brlenexternal) || (!e.Right().Tip() && brleninternal)) && e.Length() < cutoff {
					e.SetLength(cutoff)
				}
			}
			f.WriteString(t.Tree.Newick() + "\n")
		}
		return
	},
}

func init() {
	brlenCmd.AddCommand(minbrlenCmd)
	minbrlenCmd.Flags().Float64VarP(&cutoff, "length", "l", 0.0, "Min Length cutoff")
	minbrlenCmd.PersistentFlags().StringVarP(&intreefile, "input", "i", "stdin", "Input tree")
	minbrlenCmd.PersistentFlags().StringVarP(&outtreefile, "output", "o", "stdout", "Min length output tree file")
}
