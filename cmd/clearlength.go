package cmd

import (
	"github.com/fredericlemoine/gotree/io"
	"github.com/spf13/cobra"
)

// clearlengthCmd represents the clearlength command
var clearlengthCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clear lengths from input trees",
	Long:  `Clear lengths from input trees.`,
	Run: func(cmd *cobra.Command, args []string) {
		f := openWriteFile(outtreefile)
		defer closeWriteFile(f, outtreefile)

		treefile, treechan := readTrees(intreefile)
		defer treefile.Close()
		for t := range treechan {
			if t.Err != nil {
				io.ExitWithMessage(t.Err)
			}
			t.Tree.ClearLengths()
			f.WriteString(t.Tree.Newick() + "\n")
		}
	},
}

func init() {
	brlenCmd.AddCommand(clearlengthCmd)
	clearlengthCmd.PersistentFlags().StringVarP(&outtreefile, "output", "o", "stdout", "Cleared tree output file")
}
