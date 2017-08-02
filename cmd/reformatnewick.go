package cmd

import (
	"github.com/fredericlemoine/gotree/io"

	"github.com/spf13/cobra"
)

// newickCmd represents the newick command
var newickCmd = &cobra.Command{
	Use:   "newick",
	Short: "Reformats an input tree file into Newick format",
	Long: `Reformats an input tree file into Newick format.

- Input formats: Newick, Nexus,
- Output format: Newick.`,
	Run: func(cmd *cobra.Command, args []string) {
		f := openWriteFile(outtreefile)
		treefile, treechan := readTrees(intreefile)
		defer treefile.Close()

		for t := range treechan {
			if t.Err != nil {
				io.ExitWithMessage(t.Err)
			}
			f.WriteString(t.Tree.Newick() + "\n")
		}
		f.Close()
	},
}

func init() {
	reformatCmd.AddCommand(newickCmd)
}
