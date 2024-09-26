package cmd

import (
	"github.com/spf13/cobra"
)

var cutdateCmd = &cobra.Command{
	Use:   "cut",
	Short: "Cut the tree",
	Long:  ``,
	RunE:  func(cmd *cobra.Command, args []string) (err error) { return },
}

func init() {
	RootCmd.AddCommand(cutdateCmd)
}
