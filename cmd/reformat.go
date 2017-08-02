package cmd

import (
	"fmt"
	"strings"

	"github.com/fredericlemoine/gotree/io"
	"github.com/fredericlemoine/gotree/io/utils"
	"github.com/spf13/cobra"
)

var reformatinputformat string

// reformatCmd represents the reformat command
var reformatCmd = &cobra.Command{
	Use:   "reformat",
	Short: "Reformats an input tree file into different formats",
	Long: `Reformats an input tree file into different formats.

So far, it can be :
- Input formats: Newick, Nexus
- Output formats: Newick, Nexus.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		switch strings.ToLower(reformatinputformat) {
		case "newick":
			treeformat = utils.FORMAT_NEWICK
		case "nexus":
			treeformat = utils.FORMAT_NEXUS
		default:
			io.ExitWithMessage(fmt.Errorf("Tree input format is not supported : %q", reformatinputformat))
		}
	},
}

func init() {
	RootCmd.AddCommand(reformatCmd)
	reformatCmd.PersistentFlags().StringVarP(&reformatinputformat, "format", "f", "newick", "Input format (newick, nexus)")
	reformatCmd.PersistentFlags().StringVarP(&intreefile, "input", "i", "stdin", "Input tree")
	reformatCmd.PersistentFlags().StringVarP(&outtreefile, "output", "o", "stdout", "Output file")

}
