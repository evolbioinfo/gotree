package cmd

import (
	"github.com/fredericlemoine/gotree/io"
	"github.com/fredericlemoine/gotree/io/nexus"
	"github.com/spf13/cobra"
)

// nexusCmd represents the nexus command
var nexusCmd = &cobra.Command{
	Use:   "nexus",
	Short: "Reformats an input tree file into Nexus format",
	Long: `Reformats an input tree file into Nexus format.

- Input formats: Newick, Nexus,
- Output format: Nexus.`,
	Run: func(cmd *cobra.Command, args []string) {
		f := openWriteFile(outtreefile)
		defer f.Close()
		treefile, treechan := readTrees(intreefile)
		defer treefile.Close()
		nex, err := nexus.WriteNexus(treechan)
		if err != nil {
			io.ExitWithMessage(err)
		}
		f.WriteString(nex)
	},
}

func init() {
	reformatCmd.AddCommand(nexusCmd)
}
