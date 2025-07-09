package cmd

import (
	goio "io"
	"os"

	"github.com/evolbioinfo/gotree/io"
	"github.com/evolbioinfo/gotree/tree"
	"github.com/spf13/cobra"
)

var resolverooted bool

// resolveCmd represents the resolve command
var resolveCmd = &cobra.Command{
	Use:   "resolve",
	Short: "Resolve multifurcations by adding 0 length branches",
	Long: `Resolve multifurcations by adding 0 length branches.

* If any node has more than 3 neighbors :
   Resolve neighbors randomly by adding 0 length 
   branches until it has 3 neighbors
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
			tr.Tree.Resolve(resolverooted)
			f.WriteString(tr.Tree.Newick() + "\n")
		}
		return
	},
}

func init() {
	RootCmd.AddCommand(resolveCmd)
	resolveCmd.PersistentFlags().StringVarP(&intreefile, "input", "i", "stdin", "Input tree(s) file")
	resolveCmd.PersistentFlags().StringVarP(&outtreefile, "output", "o", "stdout", "Resolved tree(s) output file")
	resolveCmd.PersistentFlags().BoolVar(&resolverooted, "rooted", false, "Considers the tree as rooted (will randomly resolve the root also if needed)")
}
