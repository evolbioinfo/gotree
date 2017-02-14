package cmd

import (
	"github.com/spf13/cobra"
)

// collapseCmd represents the collapse command
var collapseCmd = &cobra.Command{
	Use:   "collapse",
	Short: "Collapse branches of input trees",
	Long:  `Collapse branches of input trees.`,
}

func init() {
	RootCmd.AddCommand(collapseCmd)
	collapseCmd.PersistentFlags().StringVarP(&intreefile, "input", "i", "stdin", "Input tree")
	collapseCmd.PersistentFlags().StringVarP(&outtreefile, "output", "o", "stdout", "Collapsed tree output file")
}
