package cmd

import (
	"github.com/fredericlemoine/gostats"
	"github.com/spf13/cobra"
	"math/rand"
	"time"
)

var setlengthmean float64

// randbrlenCmd represents the randbrlen command
var randbrlenCmd = &cobra.Command{
	Use:   "randbrlen",
	Short: "Assign a random length to edges of input trees",
	Long: `Assign a random length to edges of input trees.

Length follows an exponential distribution of parameter lambda=1/0.1 
(Mean=0.1)
`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		rand.Seed(seed)
	},
	Run: func(cmd *cobra.Command, args []string) {
		f := openWriteFile(outtreefile)
		for tr := range readTrees(intreefile) {
			for _, e := range tr.Tree.Edges() {
				e.SetLength(gostats.Exp(1 / setlengthmean))
			}
			f.WriteString(tr.Tree.Newick() + "\n")
		}
		f.Close()
	},
}

func init() {
	RootCmd.AddCommand(randbrlenCmd)

	randbrlenCmd.PersistentFlags().StringVarP(&intreefile, "input", "i", "stdin", "Input tree")
	randbrlenCmd.PersistentFlags().StringVarP(&outtreefile, "output", "o", "stdout", "Output file")
	randbrlenCmd.Flags().Int64VarP(&seed, "seed", "s", time.Now().UTC().UnixNano(), "Initial Random Seed")
	randbrlenCmd.Flags().Float64VarP(&setlengthmean, "mean", "m", 0.1, "Mean of the exponential distribution of branch lengths")
}
