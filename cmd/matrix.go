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
var matrixavg bool

// matrixCmd represents the matrix command
var matrixCmd = &cobra.Command{
	Use:   "matrix",
	Short: "Prints distance matrix associated to the input tree",
	Long: `Prints distance matrix associated to the input tree.
	
	The distance matrix can be computed in several ways, depending on the "metric" option:
	* --metric brlen : distances correspond to the sum of branch lengths between the tips (patristic distance). If there is no length for a given branch, 0.0 is the default.
	* --metric boot : distances correspond to the sum of supports of the internal branches separating the tips. If there is no support for a given branch (e.g. for a tip), 1.0 is the default. If branch supports range from 0 to 100, you may consider to use gotree support scale -f 0.01 first.
	* --metric none : distances correspond to the sum of the branches separating the tips, but each individual branch is counted as having a length of 1 (topological distance)
	`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		var f *os.File
		var treefile goio.Closer
		var treechan <-chan tree.Trees
		var distmetric int
		var mat [][]float64
		var tips []*tree.Node

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

		if matrixavg {
			if mat, tips, err = tree.AvgDistanceMatrix(distmetric, treechan); err != nil {
				io.LogError(err)
				return
			}
			f.WriteString(fmt.Sprintf("%d\n", len(tips)))
			for i, t := range tips {
				f.WriteString(t.Name())
				for j := range tips {
					f.WriteString("\t" + fmt.Sprintf("%.12f", mat[i][j]))
				}
				f.WriteString("\n")
			}
		} else {
			for t := range treechan {
				if t.Err != nil {
					io.LogError(t.Err)
					return t.Err
				}
				mat, tips = t.Tree.ToDistanceMatrix(distmetric)
				f.WriteString(fmt.Sprintf("%d\n", len(tips)))
				for i, t := range tips {
					f.WriteString(t.Name())
					for j := range tips {
						f.WriteString("\t" + fmt.Sprintf("%.12f", mat[i][j]))
					}
					f.WriteString("\n")
				}
			}
		}
		return
	},
}

func init() {
	RootCmd.AddCommand(matrixCmd)
	matrixCmd.PersistentFlags().StringVarP(&intreefile, "input", "i", "stdin", "Input tree")
	matrixCmd.PersistentFlags().StringVarP(&metric, "metric", "m", "brlen", "Distance metric (brlen|boot|none)")
	matrixCmd.PersistentFlags().BoolVar(&matrixavg, "avg", false, "Average the distance matrices of all input trees")
	matrixCmd.PersistentFlags().StringVarP(&outtreefile, "output", "o", "stdout", "Matrix output file")
}
