package cmd

import (
	"github.com/fredericlemoine/gotree/io"
	"github.com/fredericlemoine/gotree/io/phyloxml"
	"github.com/spf13/cobra"
)

// phyloxmlCmd represents the phyloxml command
var phyloxmlCmd = &cobra.Command{
	Use:   "phyloxml",
	Short: "Reformats an input tree file into PhyloXML format",
	Long: `Reformats an input tree file into PhyloXML format.

- Input formats: Newick, Nexus, PhyloXML
- Output format: PhyloXML.

Note that only toplogical information, node names, branch lengths and 
branch supports are kept.
`,
	Run: func(cmd *cobra.Command, args []string) {
		f := openWriteFile(outtreefile)
		defer closeWriteFile(f, outtreefile)

		treefile, treechan := readTrees(intreefile)
		defer treefile.Close()
		xml, err := phyloxml.WritePhyloXML(treechan)
		if err != nil {
			io.ExitWithMessage(err)
		}
		f.WriteString(xml)
	},
}

func init() {
	reformatCmd.AddCommand(phyloxmlCmd)
}
