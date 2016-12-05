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
		for t := range clearIntrees {
			t.Tree.ClearSupports()
			clearOutTrees.WriteString(t.Tree.Newick() + "\n")
		}
	},
}

func init() {
	clearCmd.AddCommand(clearsupportCmd)
}
