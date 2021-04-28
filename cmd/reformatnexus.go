package cmd

import (
	goio "io"
	"os"

	"github.com/evolbioinfo/gotree/io"
	"github.com/evolbioinfo/gotree/io/nexus"
	"github.com/evolbioinfo/gotree/tree"
	"github.com/spf13/cobra"
)

var nexusTranslate bool

// nexusCmd represents the nexus command
var nexusCmd = &cobra.Command{
	Use:   "nexus",
	Short: "Reformats an input tree file into Nexus format",
	Long: `Reformats an input tree file into Nexus format.

- Input formats: Newick, Nexus,
- Output format: Nexus.`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		var f *os.File
		var treefile goio.Closer
		var treechan <-chan tree.Trees
		var nex string

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
		if nex, err = nexus.WriteNexus(treechan, nexusTranslate); err != nil {
			io.LogError(err)
			return
		}
		f.WriteString(nex)
		return
	},
}

func init() {
	reformatCmd.AddCommand(nexusCmd)
	nexusCmd.PersistentFlags().BoolVar(&nexusTranslate, "translate", false, "Renames tip names with indices and add a translate table to the Nexus format")
}
