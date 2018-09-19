package cmd

import (
	"github.com/fredericlemoine/gotree/io"
	"github.com/spf13/cobra"
)

// clearsupportCmd represents the clearsupport command
var clearsupportCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clear supports from input trees",
	Long:  `Clear supports from input trees.`,
	Run: func(cmd *cobra.Command, args []string) {
		f := openWriteFile(outtreefile)
		defer closeWriteFile(f, outtreefile)

		treefile, treechan := readTrees(intreefile)
		defer treefile.Close()
		for t := range treechan {
			if t.Err != nil {
				io.ExitWithMessage(t.Err)
			}
			t.Tree.ClearSupports()
			f.WriteString(t.Tree.Newick() + "\n")
		}
	},
}

func init() {
	supportCmd.AddCommand(clearsupportCmd)
}
