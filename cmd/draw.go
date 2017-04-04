package cmd

import (
	"github.com/spf13/cobra"
)

var drawNoTipLabels bool
var drawNoBranchLengths bool
var drawInternalNodeLabels bool

// drawCmd represents the draw command
var drawCmd = &cobra.Command{
	Use:   "draw",
	Short: "Draw trees",
	Long:  `Draw trees `,
}

func init() {
	RootCmd.AddCommand(drawCmd)

	drawCmd.PersistentFlags().StringVarP(&intreefile, "input", "i", "stdin", "Input tree")
	drawCmd.PersistentFlags().StringVarP(&outtreefile, "output", "o", "stdout", "Output file")
	drawCmd.PersistentFlags().BoolVar(&drawNoTipLabels, "no-tip-labels", false, "Draw the tree without tip labels")
	drawCmd.PersistentFlags().BoolVar(&drawNoBranchLengths, "no-branch-lengths", false, "Draw the tree without branch lengths (all the same length)")
	drawCmd.PersistentFlags().BoolVar(&drawInternalNodeLabels, "with-node-labels", false, "Draw the tree with internal node labels")
}
