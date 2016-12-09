package cmd

import (
	"github.com/spf13/cobra"
)

// computeCmd represents the compute command
var computeCmd = &cobra.Command{
	Use:   "compute",
	Short: "Computations such as consensus and supports",
	Long: `Computations such as consensus and supports.
`,
}

func init() {
	RootCmd.AddCommand(computeCmd)
}
