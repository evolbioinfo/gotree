package cmd

import (
	"github.com/spf13/cobra"
)

// clearlengthCmd represents the clearlength command
var clearlengthCmd = &cobra.Command{
	Use:   "lengths",
	Short: "Clear lengths from input trees",
	Long:  `Clear lengths from input trees.`,
	Run: func(cmd *cobra.Command, args []string) {
		for t := range clearIntrees {
			t.Tree.ClearLengths()
			clearOutTrees.WriteString(t.Tree.Newick() + "\n")
		}
	},
}

func init() {
	clearCmd.AddCommand(clearlengthCmd)
}
