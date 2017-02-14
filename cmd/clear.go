package cmd

import (
	"github.com/spf13/cobra"
)

// clearCmd represents the clear command
var clearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clear lengths or supports from input trees",
	Long:  `Clear lengths or supports from input trees`,
}

func init() {
	RootCmd.AddCommand(clearCmd)
	clearCmd.PersistentFlags().StringVarP(&intreefile, "input", "i", "stdin", "Input tree")
	clearCmd.PersistentFlags().StringVarP(&outtreefile, "output", "o", "stdout", "Cleared tree output file")
}
