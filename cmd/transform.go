package cmd

import (
	"github.com/spf13/cobra"
)

var transformInputTree string
var transformOutputTree string

// transformCmd represents the transform command
var transformCmd = &cobra.Command{
	Use:   "transform",
	Short: "Command to tranform input tree",
	Long: `Command to transform input tree

Collapse, shuffle, unroot, etc.

`,
}

func init() {
	RootCmd.AddCommand(transformCmd)

	transformCmd.PersistentFlags().StringVarP(&transformInputTree, "input", "i", "stdin", "Input tree")
	transformCmd.PersistentFlags().StringVarP(&transformOutputTree, "output", "o", "stdout", "Collapsed tree output file")
}
