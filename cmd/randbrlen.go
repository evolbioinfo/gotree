package cmd

import (
	"github.com/fredericlemoine/gostats"
	"github.com/fredericlemoine/gotree/io"
	"github.com/fredericlemoine/gotree/io/utils"
	"github.com/fredericlemoine/gotree/tree"
	"github.com/spf13/cobra"
	"math/rand"
	"os"
	"time"
)

var setlengthinput string
var setlengthintrees chan tree.Trees
var setlengthout string
var setlengthseed int64
var setlengthmean float64
var setlengthoutfile *os.File

// randbrlenCmd represents the randbrlen command
var randbrlenCmd = &cobra.Command{
	Use:   "randbrlen",
	Short: "Assign a  random length to edges of input trees",
	Long: `Assign a  random length to edges of input trees.

Length follows an exponential distribution of parameter lambda=1/0.1 
(Mean=0.1)
`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		var err error
		rand.Seed(shuffletipsSeed)
		setlengthintrees = make(chan tree.Trees, 15)
		/* Read ref tree(s) */
		go func() {
			if _, err = utils.ReadCompTrees(setlengthinput, setlengthintrees); err != nil {
				io.ExitWithMessage(err)
			}
		}()
		setlengthoutfile = openWriteFile(setlengthout)
	},
	Run: func(cmd *cobra.Command, args []string) {
		for tr := range setlengthintrees {
			for _, e := range tr.Tree.Edges() {
				e.SetLength(gostats.Exp(1 / setlengthmean))
			}
			setlengthoutfile.WriteString(tr.Tree.Newick() + "\n")
		}
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		setlengthoutfile.Close()
	},
}

func init() {
	RootCmd.AddCommand(randbrlenCmd)

	randbrlenCmd.PersistentFlags().StringVarP(&setlengthinput, "input", "i", "stdin", "Input tree")
	randbrlenCmd.PersistentFlags().StringVarP(&setlengthout, "output", "o", "stdout", "Output file")
	randbrlenCmd.Flags().Int64VarP(&setlengthseed, "seed", "s", time.Now().UTC().UnixNano(), "Initial Random Seed")
	randbrlenCmd.Flags().Float64VarP(&setlengthmean, "mean", "m", 0.1, "Mean of the exponential distribution of branch lengths")
}
