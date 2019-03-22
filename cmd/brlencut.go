package cmd

import (
	"fmt"
	goio "io"
	"os"

	"github.com/evolbioinfo/gotree/io"
	"github.com/evolbioinfo/gotree/tree"
	"github.com/spf13/cobra"
)

var cutlengthmax float64

// cutCmd represents the cut command
var cutCmd = &cobra.Command{
	Use:   "cut",
	Short: "Cut branches whose length is greater than or equal to the given length ",
	Long: `Cut branches whose length is greater than or equal to the given length.

As output, it prints groups of tips that are in connected components of the now disconnected tree.

Output format: One line per group/connected component. Each line contains id \t ntips \t t1,t2,t3, 
with id="id of the input tree", ntips="Number of tips in that group" and t1,t2,t3="a coma separated list of tips in the group".

Example:

gotree brlen cut -i tree.nhx -l 0.1 -o groups.txt

 `,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		var bags []*tree.TipBag
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
			if bags, err = t.Tree.CutEdgesMaxLength(cutlengthmax); err != nil {
				io.LogError(err)
				return
			}
			for _, b := range bags {
				f.WriteString(fmt.Sprintf("%d\t%d\t", t.Id, b.Size()))
				for i, tip := range b.Tips() {
					if i > 0 {
						f.WriteString(",")
					}
					f.WriteString(tip.Name())
				}
				f.WriteString("\n")
			}
		}
		return
	},
}

func init() {
	brlenCmd.AddCommand(cutCmd)
	cutCmd.PersistentFlags().Float64VarP(&cutlengthmax, "max-length", "l", 0.5, "Length cutoff. Branches with length greater than or equal to this cutoff are considered removed")
	cutCmd.PersistentFlags().StringVarP(&outtreefile, "output", "o", "stdout", "Output file with groups of tips/connected components")
}
