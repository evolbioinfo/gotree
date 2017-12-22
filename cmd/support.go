package cmd

import (
	"github.com/spf13/cobra"
)

// clearCmd represents the clear command
var supportCmd = &cobra.Command{
	Use:   "support",
	Short: "Modify supports of branches",
	Long:  `Modify supports of branches from input trees`,
}

func init() {
	RootCmd.AddCommand(supportCmd)
	supportCmd.PersistentFlags().StringVarP(&intreefile, "input", "i", "stdin", "Input tree")
	supportCmd.PersistentFlags().StringVarP(&outtreefile, "output", "o", "stdout", "Cleared tree output file")
}
