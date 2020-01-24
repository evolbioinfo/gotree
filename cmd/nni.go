package cmd

import (
	goio "io"
	"os"

	"github.com/evolbioinfo/gotree/io"
	"github.com/evolbioinfo/gotree/tree"
	"github.com/spf13/cobra"
)

// rerootCmd represents the reroot command
var nniCmd = &cobra.Command{
	Use:   "nni",
	Short: "Generates all NNI neighbors from a given tree",
	Long: `Generates all NNI neighbors from a given tree.
`,

	RunE: func(cmd *cobra.Command, args []string) (err error) {
		var f *os.File
		var treefile goio.Closer
		var treechan <-chan tree.Trees

		if treefile, treechan, err = readTrees(intreefile); err != nil {
			io.LogError(err)
			return
		}
		defer treefile.Close()

		if f, err = openWriteFile(outtreefile); err != nil {
			io.LogError(err)
			return
		}
		defer closeWriteFile(f, outtreefile)

		r := &tree.NNIRearranger{}

		for t := range treechan {
			r.Rearrange(t.Tree, func(re tree.Rearrangement) bool {
				if err = re.Apply(); err != nil {
					return false
				}
				if err = t.Tree.CheckTreePostOrder(); err != nil {
					return false
				}

				f.WriteString(t.Tree.Newick() + "\n")

				if err = re.Undo(); err != nil {
					return false
				}
				if err = t.Tree.CheckTreePostOrder(); err != nil {
					return false
				}
				return true
			})

			if err != nil {
				io.LogError(err)
				return
			}
		}

		return
	},
}

func init() {
	RootCmd.AddCommand(nniCmd)
	nniCmd.PersistentFlags().StringVarP(&intreefile, "input", "i", "stdin", "Input Tree")
	nniCmd.PersistentFlags().StringVarP(&outtreefile, "output", "o", "stdout", "NNI output tree file")
}
