package cmd

import (
	goio "io"
	"os"

	"github.com/evolbioinfo/gotree/io"
	"github.com/evolbioinfo/gotree/tree"
	"github.com/spf13/cobra"
)

// rerootCmd represents the reroot command
var midpointCmd = &cobra.Command{
	Use:   "midpoint",
	Short: "Reroot trees at midpoint",
	Long: `Reroot tree at midpoint.

Example:

gotree reroot midpoint  -i tree.nw > reroot.nw
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
			err = t.Tree.RerootMidPoint()
			if err != nil {
				io.LogError(err)
				return
			}
			f.WriteString(t.Tree.Newick() + "\n")
		}
		return
	},
}

func init() {
	rerootCmd.AddCommand(midpointCmd)
}
