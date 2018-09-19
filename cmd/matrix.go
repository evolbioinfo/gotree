package cmd

import (
	"fmt"
	goio "io"
	"os"

	"github.com/fredericlemoine/gotree/io"
	"github.com/fredericlemoine/gotree/tree"
	"github.com/spf13/cobra"
)

// matrixCmd represents the matrix command
var matrixCmd = &cobra.Command{
	Use:   "matrix",
	Short: "Prints distance matrix associated to the input tree",
	Long:  `Prints distance matrix associated to the input tree.`,
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
		}
		defer treefile.Close()

		for t := range treechan {
			if t.Err != nil {
				io.LogError(t.Err)
				return t.Err
			}
			tips := t.Tree.Tips()
			f.WriteString(fmt.Sprintf("%d\n", len(tips)))
			mat := t.Tree.ToDistanceMatrix()
			for i, t := range tips {
				f.WriteString(t.Name())
				for j, _ := range tips {
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
	matrixCmd.PersistentFlags().StringVarP(&outtreefile, "output", "o", "stdout", "Matrix output file")
}
