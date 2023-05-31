package cmd

import (
	"github.com/spf13/cobra"
)

// brlenCmd represents the brlen command
var brlenCmd = &cobra.Command{
	Use:   "brlen",
	Short: "Modify branch lengths",
	Long: `Commands to modify lengths of branches:
Set a minimum branch length, or set random branch lengths, or multiply branch lengths by a factor.
`,
}
var brlenexternal, brleninternal bool

func init() {
	RootCmd.AddCommand(brlenCmd)
	brlenCmd.PersistentFlags().StringVarP(&intreefile, "input", "i", "stdin", "Input tree")
	brlenCmd.PersistentFlags().BoolVar(&brlenexternal, "external", true, "Applies to external branches")
	brlenCmd.PersistentFlags().BoolVar(&brleninternal, "internal", true, "Applies to internal branches")
}
