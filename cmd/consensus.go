package cmd

import (
	goio "io"
	"os"

	"github.com/fredericlemoine/gotree/io"
	"github.com/fredericlemoine/gotree/tree"
	"github.com/spf13/cobra"
)

// consensusCmd represents the consensus command
var consensusCmd = &cobra.Command{
	Use:   "consensus",
	Short: "Computes the consensus of a set of trees",
	Long: `Computes the consensus of a set of input trees
Trees must have the same tip names.

Two parameters:
-i : Input file containing several trees
-f : Percentage threshold to keep a bipartition in the consensus 
     It must be >=0.5 && <=1

In the output consensus tree:
1) Branch supports are computed as the proportion of trees in which
   the bipartition is present
2) Branch lengths are computed as the average length of the same branch
   over all the trees where it is present

`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		var f *os.File
		var treefile goio.Closer
		var treechan <-chan tree.Trees
		var consensus *tree.Tree

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
		consensus, err = tree.Consensus(treechan, cutoff)
		if err != nil {
			io.LogError(err)
			return
		}
		f.WriteString(consensus.Newick() + "\n")
		return
	},
}

func init() {
	computeCmd.AddCommand(consensusCmd)
	consensusCmd.PersistentFlags().StringVarP(&intreefile, "input", "i", "stdin", "Input tree")
	consensusCmd.PersistentFlags().StringVarP(&outtreefile, "output", "o", "stdout", "Output file")
	consensusCmd.PersistentFlags().Float64VarP(&cutoff, "freq-min", "f", 0.5, "Minimum frequency to keep the bipartitions")
}
