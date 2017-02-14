package cmd

import (
	"github.com/spf13/cobra"
)

// clearsupportCmd represents the clearsupport command
var clearsupportCmd = &cobra.Command{
	Use:   "supports",
	Short: "Clear supports from input trees",
	Long:  `Clear supports from input trees.`,
	Run: func(cmd *cobra.Command, args []string) {
		f := openWriteFile(outtreefile)
		for t := range readTrees(intreefile) {
			t.Tree.ClearSupports()
			f.WriteString(t.Tree.Newick() + "\n")
		}
		f.Close()
	},
}

func init() {
	clearCmd.AddCommand(clearsupportCmd)
}
