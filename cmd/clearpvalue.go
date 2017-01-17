package cmd

import (
	"github.com/spf13/cobra"
)

// clearsupportCmd represents the clearsupport command
var clearpvalueCmd = &cobra.Command{
	Use:   "pvalues",
	Short: "Clear pvalues associated to supports from input trees",
	Long:  `Clear pvalues associated to supports from input trees.`,
	Run: func(cmd *cobra.Command, args []string) {
		for t := range clearIntrees {
			t.Tree.ClearPvalues()
			clearOutTrees.WriteString(t.Tree.Newick() + "\n")
		}
	},
}

func init() {
	clearCmd.AddCommand(clearpvalueCmd)
}
