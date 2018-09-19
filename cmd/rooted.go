package cmd

import (
	"fmt"

	"github.com/fredericlemoine/gotree/io"
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
	Run: func(cmd *cobra.Command, args []string) {
		f := openWriteFile(outtreefile)
		defer closeWriteFile(f, outtreefile)

		f.WriteString("tree\trooted\n")
		treefile, treechan := readTrees(intreefile)
		defer treefile.Close()
		for t := range treechan {
			if t.Err != nil {
				io.ExitWithMessage(t.Err)
			}
			f.WriteString(fmt.Sprintf("%d\t", t.Id))
			if t.Tree.Rooted() {
				f.WriteString("rooted\n")
			} else {
				f.WriteString("unrooted\n")
			}
		}
	},
}

func init() {
	statsCmd.AddCommand(rootedCmd)
}
