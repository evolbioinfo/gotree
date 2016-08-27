package cmd

import (
	"github.com/spf13/cobra"
)

// computeCmd represents the compute command
var computeCmd = &cobra.Command{
	Use:   "compute",
	Short: "Command to do different computations",
	Long: `Command to do different computations such
as support.
`,
}

func init() {
	RootCmd.AddCommand(computeCmd)
}
