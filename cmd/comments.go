package cmd

import (
	"github.com/spf13/cobra"
)

// brlenCmd represents the brlen command
var commentCmd = &cobra.Command{
	Use:   "comment",
	Short: "Modify branch/node comments",
	Long:  `Modify branch/node comments`,
}

func init() {
	RootCmd.AddCommand(commentCmd)
	commentCmd.PersistentFlags().StringVarP(&intreefile, "input", "i", "stdin", "Input tree")
	commentCmd.PersistentFlags().StringVarP(&outtreefile, "output", "o", "stdout", "Cleared tree output file")
}
