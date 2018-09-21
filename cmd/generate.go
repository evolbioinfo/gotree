package cmd

import (
	"github.com/spf13/cobra"
)

var generateNbTips int
var generateDepth int
var generateNbTrees int
var generateOutputfile string
var generateRooted bool

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate random trees",
	Long: `Generate random trees
`,
}

func init() {
	RootCmd.AddCommand(generateCmd)
	generateCmd.PersistentFlags().IntVarP(&generateNbTrees, "nbtrees", "n", 1, "Number of trees to generate")
	generateCmd.PersistentFlags().StringVarP(&generateOutputfile, "output", "o", "stdout", "Number of tips of the tree to generate")
	generateCmd.PersistentFlags().BoolVarP(&generateRooted, "rooted", "r", false, "Generate rooted trees")
}
