package cmd

import (
	"fmt"
	goio "io"
	"os"

	"github.com/fredericlemoine/gotree/io"
	"github.com/fredericlemoine/gotree/tree"
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
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		var f *os.File
		var treefile goio.Closer
		var treechan <-chan tree.Trees

		/* Dividing trees */
		i := 0
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
			if f, err = openWriteFile(fmt.Sprintf("%s_%03d.nw", outtreefile, i)); err != nil {
				io.LogError(err)
				return
			}
			f.WriteString(t.Tree.Newick() + "\n")
			f.Close()
			i++
		}
		return
	},
}

func init() {
	RootCmd.AddCommand(divideCmd)
	divideCmd.PersistentFlags().StringVarP(&intreefile, "input", "i", "stdin", "Input tree(s) file")
	divideCmd.PersistentFlags().StringVarP(&outtreefile, "output", "o", "prefix", "Divided trees output file prefix")
}
