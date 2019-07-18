package cmd

import (
	goio "io"
	"os"

	"github.com/evolbioinfo/gotree/io"
	"github.com/evolbioinfo/gotree/tree"
	"github.com/fredericlemoine/gostats"
	"github.com/spf13/cobra"
)

var setlengthmean float64
var setlengthmeanMin float64
var setlengthmeanMax float64

// randbrlenCmd represents the randbrlen command
var randbrlenCmd = &cobra.Command{
	Use:   "setrand",
	Short: "Assign a random length to edges of input trees",
	Long: `Assign a random length to edges of input trees.

Branch lengths are drawn in an exponential distribution of parameter lambda=1/mean.
Two possibilities for the mean:

1) If --mean-min and --mean-max are given, and mean-min < mean-max and are both > 0 then 
"mean" is drawn uniformly in the interval [mean-min,mean-max]

2) Otherwise, 'mean' is set to --mean value.

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

		if treefile, treechan, err = readTrees(intreefile); err != nil {
			io.LogError(err)
			return
		}
		defer treefile.Close()

		lmean := setlengthmean
		for tr := range treechan {
			if cmd.Flags().Changed("min-mean") && cmd.Flags().Changed("max-mean") &&
				setlengthmeanMin < setlengthmeanMax && setlengthmeanMin >= 0 && setlengthmeanMax > 0 {
				lmean = gostats.Float64RangeF(setlengthmeanMin, setlengthmeanMax)
			}
			if tr.Err != nil {
				io.LogError(tr.Err)
				return tr.Err
			}

			for _, e := range tr.Tree.Edges() {
				e.SetLength(gostats.Exp(1.0 / lmean))
			}
			f.WriteString(tr.Tree.Newick() + "\n")
		}
		return
	},
}

func init() {
	brlenCmd.AddCommand(randbrlenCmd)

	randbrlenCmd.PersistentFlags().StringVarP(&intreefile, "input", "i", "stdin", "Input tree")
	randbrlenCmd.Flags().Float64VarP(&setlengthmean, "mean", "m", 0.1, "Mean of the exponential distribution of branch lengths")
	randbrlenCmd.Flags().Float64Var(&setlengthmeanMin, "min-mean", 0.001, "Mean of the exponential distribution of branch lengths will be drawn uniformly in the interval [min-mean,max-mean]")
	randbrlenCmd.Flags().Float64Var(&setlengthmeanMax, "max-mean", 0.05, "Mean of the exponential distribution of branch lengths will be drawn uniformly in the interval [min-mean,max-mean]")
	randbrlenCmd.PersistentFlags().StringVarP(&outtreefile, "output", "o", "stdout", "Random length output tree file")
}
