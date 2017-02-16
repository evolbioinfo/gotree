package cmd

import (
	"github.com/spf13/cobra"
)

// uploadCmd represents the upload command
var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Upload a tree to a given server",
	Long:  `Upload a tree to a given server`,
}

func init() {
	RootCmd.AddCommand(uploadCmd)
	uploadCmd.PersistentFlags().StringVarP(&intreefile, "input", "i", "stdin", "Input tree")
}
