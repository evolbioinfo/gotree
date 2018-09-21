package cmd

import (
	goio "io"
	"os"

	"github.com/fredericlemoine/gotree/io"
	"github.com/fredericlemoine/gotree/io/phyloxml"
	"github.com/fredericlemoine/gotree/tree"
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
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		var f *os.File
		var treefile goio.Closer
		var treechan <-chan tree.Trees
		var xml string

		if f, err = openWriteFile(outtreefile); err != nil {
			io.LogError(err)
			return
		}
		defer closeWriteFile(f, outtreefile)

		if treefile, treechan, err = readTrees(intreefile); err != nil {
			io.LogError(err)
			return
		}
		defer treefile.Close()
		if xml, err = phyloxml.WritePhyloXML(treechan); err != nil {
			io.LogError(err)
			return
		}
		f.WriteString(xml)
		return
	},
}

func init() {
	reformatCmd.AddCommand(phyloxmlCmd)
}
