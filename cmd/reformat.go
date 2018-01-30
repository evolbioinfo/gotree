package cmd

import (
	"github.com/spf13/cobra"
)

// reformatCmd represents the reformat command
var reformatCmd = &cobra.Command{
	Use:   "reformat",
	Short: "Reformats an input tree file into different formats",
	Long: `Reformats an input tree file into different formats.

So far, it can be :
- Input formats: Newick, Nexus
- Output formats: Newick, Nexus.`,
}

func init() {
	RootCmd.AddCommand(reformatCmd)
	reformatCmd.PersistentFlags().StringVarP(&rootInputFormat, "input-format", "f", "newick", "Input tree format (newick, nexus, or phyloxml), alias to --format")
	reformatCmd.PersistentFlags().StringVarP(&intreefile, "input", "i", "stdin", "Input tree")
	reformatCmd.PersistentFlags().StringVarP(&outtreefile, "output", "o", "stdout", "Output file")

}
