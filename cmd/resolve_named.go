package cmd

import (
	goio "io"
	"os"

	"github.com/evolbioinfo/gotree/io"
	"github.com/evolbioinfo/gotree/tree"
	"github.com/spf13/cobra"
)

// resolveCmd represents the resolve command
var resolveNamedCmd = &cobra.Command{
	Use:   "named",
	Short: "Resolve named internal nodes into tip + 0 length branches",
	Long: `Resolve named internal nodes into tip + 0 length branches

	Example:

	-------T1      -------T1
	|              |
	*N1        =>  *---N1      
	|              |
	-------T2      -------T2

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
			tr.Tree.ResolveNamedInternalNodes()
			f.WriteString(tr.Tree.Newick() + "\n")
		}
		return
	},
}

func init() {
	resolveCmd.AddCommand(resolveNamedCmd)

}
