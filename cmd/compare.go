package cmd

import (
	"github.com/spf13/cobra"
)

// compareCmd represents the compare command
var compareCmd = &cobra.Command{
	Use:   "compare",
	Short: "Compare full trees, edges, or tips",
	Long: `Compare full trees, edges, or tips.
`,
}

func init() {
	RootCmd.AddCommand(compareCmd)
	compareCmd.PersistentFlags().StringVarP(&intreefile, "reftree", "i", "stdin", "Reference tree input file")
	compareCmd.PersistentFlags().StringVarP(&intree2file, "compared", "c", "none", "Compared trees input file")
}
