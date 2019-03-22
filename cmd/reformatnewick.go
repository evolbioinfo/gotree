package cmd

import (
	goio "io"
	"os"

	"github.com/evolbioinfo/gotree/io"
	"github.com/evolbioinfo/gotree/tree"
	"github.com/spf13/cobra"
)

// newickCmd represents the newick command
var newickCmd = &cobra.Command{
	Use:   "newick",
	Short: "Reformats an input tree file into Newick format",
	Long: `Reformats an input tree file into Newick format.

- Input formats: Newick, Nexus,
- Output format: Newick.`,
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
			f.WriteString(t.Tree.Newick() + "\n")
		}
		return
	},
}

func init() {
	reformatCmd.AddCommand(newickCmd)
}
