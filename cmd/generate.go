package cmd

import (
	"github.com/spf13/cobra"
)

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Command to genereate random trees",
	Long: `Command to generate random trees
`,
}

func init() {
	RootCmd.AddCommand(generateCmd)
}
