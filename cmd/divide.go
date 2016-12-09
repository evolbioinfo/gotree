package cmd

import (
	"fmt"
	"github.com/fredericlemoine/gotree/io"
	"github.com/fredericlemoine/gotree/io/utils"
	"github.com/fredericlemoine/gotree/tree"
	"github.com/spf13/cobra"
)

var divideInputTree string
var divideOutputTree string

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
		var err error
		var nbtrees int = 0
		intrees := make(chan tree.Trees, 15)

		/* Read ref tree(s) */
		go func() {
			if nbtrees, err = utils.ReadCompTrees(divideInputTree, intrees); err != nil {
				io.ExitWithMessage(err)
			}
		}()

		/* Dividing trees */
		i := 0
		for t := range intrees {
			f := openWriteFile(fmt.Sprintf("%s_%03d.nw", divideOutputTree, i))
			f.WriteString(t.Tree.Newick() + "\n")
			f.Close()
			i++
		}
	},
}

func init() {
	RootCmd.AddCommand(divideCmd)
	divideCmd.PersistentFlags().StringVarP(&divideInputTree, "input", "i", "stdin", "Input tree(s) file")
	divideCmd.PersistentFlags().StringVarP(&divideOutputTree, "output", "o", "prefix", "Divided trees output file prefix")
}
