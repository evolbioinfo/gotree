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

// randbrlenCmd represents the randbrlen command
var randbrlenCmd = &cobra.Command{
	Use:   "setrand",
	Short: "Assign a random length to edges of input trees",
	Long: `Assign a random length to edges of input trees.

Length follows an exponential distribution of parameter lambda=1/0.1 
(Mean=0.1)
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

		for tr := range treechan {
			if tr.Err != nil {
				io.LogError(tr.Err)
				return tr.Err
			}

			for _, e := range tr.Tree.Edges() {
				e.SetLength(gostats.Exp(1.0 / setlengthmean))
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
	randbrlenCmd.PersistentFlags().StringVarP(&outtreefile, "output", "o", "stdout", "Random length output tree file")
}
