package cmd

import (
	"fmt"
	goio "io"
	"os"

	"github.com/evolbioinfo/gotree/io"
	"github.com/evolbioinfo/gotree/tree"
	"github.com/spf13/cobra"
)

var metric string

// matrixCmd represents the matrix command
var matrixCmd = &cobra.Command{
	Use:   "matrix",
	Short: "Prints distance matrix associated to the input tree",
	Long: `Prints distance matrix associated to the input tree.
	
	The distance matrix can be computed in several ways, depending on the "metric" option:
	* --metric brlen : distances correspond to the sum of branch lengths between the tips (patristic distance)
	* --metric boot : distances correspond to the sum of supports of the internal branches separating the tips
	* --metric none : distances correspond to the sum of the branches separating the tips, but each individual branch is counted as having a length of 1 (topological distance)
	`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		var f *os.File
		var treefile goio.Closer
		var treechan <-chan tree.Trees
		var distmetric int

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

		switch metric {
		case "brlen":
			distmetric = tree.DISTANCE_METRIC_BRLEN
		case "boot":
			distmetric = tree.DISTANCE_METRIC_BOOTS
		case "none":
			distmetric = tree.DISTANCE_METRIC_NONE
		default:
			err = fmt.Errorf("distance metric %s in not supported", metric)
			io.LogError(err)
			return
		}

		for t := range treechan {
			if t.Err != nil {
				io.LogError(t.Err)
				return t.Err
			}
			tips := t.Tree.Tips()
			f.WriteString(fmt.Sprintf("%d\n", len(tips)))
			mat := t.Tree.ToDistanceMatrix(distmetric)
			for i, t := range tips {
				f.WriteString(t.Name())
				for j := range tips {
					f.WriteString("\t" + fmt.Sprintf("%.12f", mat[i][j]))
				}
				f.WriteString("\n")
			}
		}
		return
	},
}

func init() {
	RootCmd.AddCommand(matrixCmd)
	matrixCmd.PersistentFlags().StringVarP(&intreefile, "input", "i", "stdin", "Input tree")
	matrixCmd.PersistentFlags().StringVarP(&metric, "metric", "m", "brlen", "Distance metric (brlen|boot|none)")
	matrixCmd.PersistentFlags().StringVarP(&outtreefile, "output", "o", "stdout", "Matrix output file")
}
