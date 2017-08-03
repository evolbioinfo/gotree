package cmd

import (
	"bytes"
	"strconv"

	"github.com/fredericlemoine/gotree/io"
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
		taxlabels := false
		f := openWriteFile(outtreefile)
		treefile, treechan := readTrees(intreefile)
		defer treefile.Close()
		var buffer bytes.Buffer
		buffer.WriteString("#NEXUS\n")

		for t := range treechan {
			if t.Err != nil {
				io.ExitWithMessage(t.Err)
			}

			if !taxlabels {
				buffer.WriteString("BEGIN TAXA;\n")
				buffer.WriteString(" TAXLABELS")
				for _, tip := range t.Tree.Tips() {
					buffer.WriteString(" " + tip.Name())
				}
				buffer.WriteString(";\n")
				buffer.WriteString("END;\n")
				buffer.WriteString("BEGIN TREES;\n")
				taxlabels = true
			}
			buffer.WriteString("  TREE tree")
			buffer.WriteString(strconv.Itoa(t.Id))
			buffer.WriteString(" = ")
			buffer.WriteString(t.Tree.Newick())
			buffer.WriteString("\n")
		}
		buffer.WriteString("END;\n")
		f.WriteString(buffer.String())
		f.Close()
	},
}

func init() {
	reformatCmd.AddCommand(nexusCmd)
}
