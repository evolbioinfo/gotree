package cmd

import (
	goio "io"
	"os"

	"github.com/fredericlemoine/gotree/io"
	"github.com/fredericlemoine/gotree/tree"
	"github.com/spf13/cobra"
)

var lowSupportThreshold float64

// collapsesupportCmd represents the collapsesupport command
var collapsesupportCmd = &cobra.Command{
	Use:   "support",
	Short: "Collapse lowly supported branches of the input tree",
	Long: `Collapse lowly supported branches of the input tree.

Lowly supported branches are defined by a threshold (-s). All branches 
with support < threshold are removed.
`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		var f *os.File
		var treefile goio.Closer
		var treechan <-chan tree.Trees

		if f, err = openWriteFile(outtreefile); err != nil {
			io.LogError(err)
			return
		}
		defer closeWriteFile(f, outtreefile)

		treefile, treechan, err = readTrees(intreefile)
		defer treefile.Close()
		for t := range treechan {
			if t.Err != nil {
				io.LogError(t.Err)
				return t.Err
			}
			t.Tree.CollapseLowSupport(lowSupportThreshold)
			f.WriteString(t.Tree.Newick() + "\n")
		}
		return
	},
}

func init() {
	collapseCmd.AddCommand(collapsesupportCmd)
	collapsesupportCmd.Flags().Float64VarP(&lowSupportThreshold, "support", "s", 0.0, "Support cutoff to collapse branches")
}
