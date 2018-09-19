package cmd

import (
	"fmt"

	"github.com/fredericlemoine/gotree/io"
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
	Run: func(cmd *cobra.Command, args []string) {
		f := openWriteFile(outtreefile)
		defer closeWriteFile(f, outtreefile)

		treefile, treechan := readTrees(intreefile)
		defer treefile.Close()
		for t := range treechan {
			if t.Err != nil {
				io.ExitWithMessage(t.Err)
			}
			bags, err := t.Tree.CutEdgesMaxLength(cutlengthmax)
			if err != nil {
				io.ExitWithMessage(t.Err)
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
	},
}

func init() {
	brlenCmd.AddCommand(cutCmd)
	cutCmd.PersistentFlags().Float64VarP(&cutlengthmax, "max-length", "l", 0.5, "Length cutoff. Branches with length greater than or equal to this cutoff are considered removed")
	cutCmd.PersistentFlags().StringVarP(&outtreefile, "output", "o", "stdout", "Output file with groups of tips/connected components")
}
