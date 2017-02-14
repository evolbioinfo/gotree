package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

// divideCmd represents the divide command
var divideCmd = &cobra.Command{
	Use:   "divide",
	Short: "Divide an input tree file into several tree files",
	Long: `Divide an input tree file into several tree files

If the input file contains several trees, lets say 10, then 10 output files 
will be created, each containing 1 tree.

Example:

gotree divide -i trees.nw -o prefix_

`,
	Run: func(cmd *cobra.Command, args []string) {
		/* Dividing trees */
		i := 0
		for t := range readTrees(intreefile) {
			f := openWriteFile(fmt.Sprintf("%s_%03d.nw", outtreefile, i))
			f.WriteString(t.Tree.Newick() + "\n")
			f.Close()
			i++
		}
	},
}

func init() {
	RootCmd.AddCommand(divideCmd)
	divideCmd.PersistentFlags().StringVarP(&intreefile, "input", "i", "stdin", "Input tree(s) file")
	divideCmd.PersistentFlags().StringVarP(&outtreefile, "output", "o", "prefix", "Divided trees output file prefix")
}
