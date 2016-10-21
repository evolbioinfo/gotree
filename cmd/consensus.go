package cmd

import (
	"github.com/fredericlemoine/gotree/io"
	"github.com/fredericlemoine/gotree/io/utils"
	"github.com/fredericlemoine/gotree/tree"
	"github.com/spf13/cobra"
	"os"
)

var consensinfile, consensoutfile string
var consensCutoff float64
var consensIntrees chan tree.Trees
var consensOut *os.File

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
`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		var err error
		var nbtrees int = 0
		consensIntrees = make(chan tree.Trees, 15)
		/* Read ref tree(s) */
		go func() {
			if nbtrees, err = utils.ReadCompTrees(consensinfile, consensIntrees); err != nil {
				io.ExitWithMessage(err)
			}
		}()
		consensOut = openWriteFile(statsoutfile)
	},
	Run: func(cmd *cobra.Command, args []string) {
		consensus := tree.Consensus(consensIntrees, consensCutoff)
		consensOut.WriteString(consensus.Newick() + "\n")
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		consensOut.Close()
	},
}

func init() {
	computeCmd.AddCommand(consensusCmd)

	consensusCmd.PersistentFlags().StringVarP(&consensinfile, "input", "i", "stdin", "Input tree")
	consensusCmd.PersistentFlags().StringVarP(&consensoutfile, "output", "o", "stdout", "Output file")
	consensusCmd.PersistentFlags().Float64VarP(&consensCutoff, "freq-min", "f", 0.5, "Minimum frequency to keep the bipartitions")
}
