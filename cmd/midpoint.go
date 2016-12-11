package cmd

import (
	"github.com/fredericlemoine/gotree/io"
	"github.com/fredericlemoine/gotree/io/utils"
	"github.com/fredericlemoine/gotree/tree"
	"github.com/spf13/cobra"
	"os"
)

// rerootCmd represents the reroot command
var midpointCmd = &cobra.Command{
	Use:   "midpoint",
	Short: "Reroot trees at midpoint",
	Long: `Reroot tree at midpoint.

Example:

gotree reroot midpoint  -i tree.nw > reroot.nw
`,
	Run: func(cmd *cobra.Command, args []string) {

		var err error
		var nbtrees int

		compareChannel := make(chan tree.Trees, 15)

		go func() {
			if nbtrees, err = utils.ReadCompTrees(rerootinputfile, compareChannel); err != nil {
				io.ExitWithMessage(err)
			}
		}()

		var f *os.File
		if rerootoutputfile != "stdout" {
			f, err = os.Create(rerootoutputfile)
		} else {
			f = os.Stdout
		}
		if err != nil {
			io.ExitWithMessage(err)
		}

		for t2 := range compareChannel {
			t2.Tree.RerootMidPoint()
			f.WriteString(t2.Tree.Newick() + "\n")
		}

		f.Close()
	},
}

func init() {
	rerootCmd.AddCommand(midpointCmd)
}
