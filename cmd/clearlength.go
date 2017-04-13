package cmd

import (
	"github.com/fredericlemoine/gotree/io"
	"github.com/spf13/cobra"
)

// clearlengthCmd represents the clearlength command
var clearlengthCmd = &cobra.Command{
	Use:   "lengths",
	Short: "Clear lengths from input trees",
	Long:  `Clear lengths from input trees.`,
	Run: func(cmd *cobra.Command, args []string) {
		f := openWriteFile(outtreefile)
		treefile, treechan := readTrees(intreefile)
		defer treefile.Close()
		for t := range treechan {
			if t.Err != nil {
				io.ExitWithMessage(t.Err)
			}
			t.Tree.ClearLengths()
			f.WriteString(t.Tree.Newick() + "\n")
		}
		f.Close()
	},
}

func init() {
	clearCmd.AddCommand(clearlengthCmd)
}
