package cmd

import (
	"github.com/fredericlemoine/gostats"
	"github.com/fredericlemoine/gotree/io"
	"github.com/spf13/cobra"
	"math/rand"
	"time"
)

var setlengthmean float64

// randbrlenCmd represents the randbrlen command
var randbrlenCmd = &cobra.Command{
	Use:   "setrand",
	Short: "Assign a random length to edges of input trees",
	Long: `Assign a random length to edges of input trees.

Length follows an exponential distribution of parameter lambda=1/0.1 
(Mean=0.1)
`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		RootCmd.PersistentPreRun(cmd, args)
		rand.Seed(seed)
	},
	Run: func(cmd *cobra.Command, args []string) {
		f := openWriteFile(outtreefile)
		defer closeWriteFile(f, outtreefile)

		treefile, trees := readTrees(intreefile)
		defer treefile.Close()

		for tr := range trees {
			if tr.Err != nil {
				io.ExitWithMessage(tr.Err)
			}

			for _, e := range tr.Tree.Edges() {
				e.SetLength(gostats.Exp(1 / setlengthmean))
			}
			f.WriteString(tr.Tree.Newick() + "\n")
		}
	},
}

func init() {
	brlenCmd.AddCommand(randbrlenCmd)

	randbrlenCmd.PersistentFlags().StringVarP(&intreefile, "input", "i", "stdin", "Input tree")
	randbrlenCmd.Flags().Int64VarP(&seed, "seed", "s", time.Now().UTC().UnixNano(), "Initial Random Seed")
	randbrlenCmd.Flags().Float64VarP(&setlengthmean, "mean", "m", 0.1, "Mean of the exponential distribution of branch lengths")
	randbrlenCmd.PersistentFlags().StringVarP(&outtreefile, "output", "o", "stdout", "Random length output tree file")
}
