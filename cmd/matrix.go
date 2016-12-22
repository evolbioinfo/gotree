package cmd

import (
	"fmt"
	"github.com/fredericlemoine/gotree/io"
	"github.com/fredericlemoine/gotree/io/utils"
	"github.com/fredericlemoine/gotree/tree"
	"github.com/spf13/cobra"
	"os"
)

var matrixInputTree string
var matrixOutput string
var matrixIntrees chan tree.Trees
var matrixOut *os.File

// matrixCmd represents the matrix command
var matrixCmd = &cobra.Command{
	Use:   "matrix",
	Short: "Prints distance matrix associated to the input tree",
	Long:  `Prints distance matrix associated to the input tree.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		var err error
		var nbtrees int = 0
		matrixIntrees = make(chan tree.Trees, 15)

		/* Read ref tree(s) */
		go func() {
			if nbtrees, err = utils.ReadCompTrees(matrixInputTree, matrixIntrees); err != nil {
				io.ExitWithMessage(err)
			}
		}()

		matrixOut = openWriteFile(matrixOutput)
	},
	Run: func(cmd *cobra.Command, args []string) {
		for t := range matrixIntrees {
			tips := t.Tree.Tips()
			mat := t.Tree.ToDistanceMatrix()
			for _, t := range tips {
				matrixOut.WriteString("\t" + t.Name())
			}
			matrixOut.WriteString("\n")
			for i, t := range tips {
				matrixOut.WriteString(t.Name())
				for j, _ := range tips {
					matrixOut.WriteString("\t" + fmt.Sprintf("%.4f", mat[i][j]))
				}
				matrixOut.WriteString("\n")
			}
		}
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		matrixOut.Close()
	},
}

func init() {
	RootCmd.AddCommand(matrixCmd)
	matrixCmd.PersistentFlags().StringVarP(&matrixInputTree, "input", "i", "stdin", "Input tree")
	matrixCmd.PersistentFlags().StringVarP(&matrixOutput, "output", "o", "stdout", "Matrix output file")
}
