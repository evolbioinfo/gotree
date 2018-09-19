package cmd

import (
	goio "io"
	"os"

	"github.com/fredericlemoine/gotree/io"
	"github.com/fredericlemoine/gotree/tree"
	"github.com/spf13/cobra"
)

var roundlengthprecision int

// roundlengthCmd represents the roundlength command
var roundlengthCmd = &cobra.Command{
	Use:   "round",
	Short: "Rounds branch lengths of input trees",
	Long: `Rounds branch lengths of input trees.

The precision is given by -p|--precision option, and is expressed in 1/10^precision

if -p 5 is given, precision of 10⁻5 is considered.


Does not do anything if precision is <=0;
Takes precision=15 if precision>15.

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
		for t := range treechan {
			if t.Err != nil {
				io.LogError(t.Err)
				return t.Err
			}
			t.Tree.RoundLengths(roundlengthprecision)
			f.WriteString(t.Tree.Newick() + "\n")
		}
		return
	},
}

func init() {
	brlenCmd.AddCommand(roundlengthCmd)
	roundlengthCmd.Flags().IntVarP(&roundlengthprecision, "precision", "p", 3, "Rounding length precision (x means 10^-x)")
	roundlengthCmd.PersistentFlags().StringVarP(&outtreefile, "output", "o", "stdout", "Rounded length output tree file")
}
