package cmd

import (
	"fmt"
	goio "io"
	"os"

	"github.com/evolbioinfo/gotree/io"
	"github.com/evolbioinfo/gotree/tree"
	"github.com/spf13/cobra"
)

// rootedCmd represents the rooted command
var rootedCmd = &cobra.Command{
	Use:   "rooted",
	Short: "Tells wether the tree is rooted or unrooted",
	Long: `Tells wether the tree is rooted or unrooted

Example of usage:

gotree stats rooted -i t.nw

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

		f.WriteString("tree\trooted\n")
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
			f.WriteString(fmt.Sprintf("%d\t", t.Id))
			if t.Tree.Rooted() {
				f.WriteString("rooted\n")
			} else {
				f.WriteString("unrooted\n")
			}
		}
		return
	},
}

func init() {
	statsCmd.AddCommand(rootedCmd)
}
