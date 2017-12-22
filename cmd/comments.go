package cmd

import (
	"github.com/spf13/cobra"
)

// brlenCmd represents the brlen command
var commentsCmd = &cobra.Command{
	Use:   "comments",
	Short: "Modify branch/node comments",
	Long:  `Modify branch/node comments`,
}

func init() {
	RootCmd.AddCommand(commentsCmd)
	commentsCmd.PersistentFlags().StringVarP(&intreefile, "input", "i", "stdin", "Input tree")
	commentsCmd.PersistentFlags().StringVarP(&outtreefile, "output", "o", "stdout", "Cleared tree output file")
}
