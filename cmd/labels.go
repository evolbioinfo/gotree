package cmd

import (
	"fmt"
	goio "io"
	"os"

	"github.com/evolbioinfo/gotree/io"
	"github.com/evolbioinfo/gotree/tree"
	"github.com/spf13/cobra"
)

var labelsNodes bool
var labelsTips bool

// labelsCmd represents the labels command
var labelsCmd = &cobra.Command{
	Use:   "labels",
	Short: "Lists labels of all tree tips",
	Long: `Lists labels of all tree tips

Example of usage:

gotree labels -i t.mw

If several trees are given in the input file, labels of all trees are listed.

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
			for _, n := range t.Tree.Nodes() {
				if (n.Tip() && labelsTips) || (!n.Tip() && labelsNodes && n.Name() != "") {
					f.WriteString(fmt.Sprintln(n.Name()))

				}
			}
		}
		return
	},
}

func init() {
	RootCmd.AddCommand(labelsCmd)
	labelsCmd.PersistentFlags().StringVarP(&intreefile, "input", "i", "stdin", "Input tree")
	labelsCmd.Flags().BoolVar(&labelsNodes, "internal", false, "Internal node labels are listed")
	labelsCmd.Flags().BoolVar(&labelsTips, "tips", true, "Tip labels are listed (--tips=false to cancel)")
}
