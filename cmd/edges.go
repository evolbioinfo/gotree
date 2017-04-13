package cmd

import (
	"fmt"

	"github.com/fredericlemoine/gotree/io"
	"github.com/spf13/cobra"
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

var locality bool

// edgesCmd represents the edges command
var edgesCmd = &cobra.Command{
	Use:   "edges",
	Short: "Displays statistics on edges of input tree",
	Long: `Displays statistics on edges of input tree

Statistics are displayed in text format (tab separated):
1 - Id of edge
2 - Length
3 - Support
4 - Terminal (true/false)
5 - Depth (Shortest path to a tip)
6 - Topo depth (Number of tips on the lightest side)
7 - Name of the Right node

Example of usage:

gotree stats edges -i t.nw

`,
	Run: func(cmd *cobra.Command, args []string) {
		f := openWriteFile(outtreefile)
		f.WriteString("tree\tbrid\tlength\tsupport\tterminal\tdepth\ttopodepth\trightname")
		if locality {
			f.WriteString("\tlocalitymin")
			f.WriteString("\tlocalitymax")
			f.WriteString("\tlocalityavg")
			f.WriteString("\tlocalityhx")
			f.WriteString("\tlocalityhy")
		}
		f.WriteString("\n")
		treefile, trees := readTrees(intreefile)
		defer treefile.Close()

		for t := range trees {
			if t.Err != nil {
				io.ExitWithMessage(t.Err)
			}
			t.Tree.ComputeDepths()
			for i, e := range t.Tree.Edges() {
				f.WriteString(
					fmt.Sprintf("%d\t%d\t%s",
						t.Id, i, e.ToStatsString()))
				if locality {
					if e.Right().Tip() {
						f.WriteString("\tN/A\tN/A\tN/A\tN/A\tN/A")
					} else {
						/**
						hx: 1 if exists a neighbor branch with suppt > 0.8
						hy: 1 if the current branch has suppt > 0.8
						*/
						avg, min, max, hx, hy := e.Locality(1, cutoff)
						f.WriteString(fmt.Sprintf("\t%f", min))
						f.WriteString(fmt.Sprintf("\t%f", max))
						f.WriteString(fmt.Sprintf("\t%f", avg))
						f.WriteString(fmt.Sprintf("\t%t", hx))
						f.WriteString(fmt.Sprintf("\t%t", hy))
					}
				}
				f.WriteString("\n")
			}
		}
		f.Close()
	},
}

func init() {
	statsCmd.AddCommand(edgesCmd)
	edgesCmd.PersistentFlags().BoolVarP(&locality, "locality", "l", false, "If locality measure must be computed (average difference between supports of edges and their neighbors)")
	//edgesCmd.PersistentFlags().IntVarP(&localitymaxdist, "locality-max-dist", "d", 1, "If locality measure is true, sets a cutoff to the neighborhood of a branch (number of edges)")
	edgesCmd.PersistentFlags().Float64VarP(&cutoff, "support-cutoff", "s", 0.8, "Cutoff to consider a branch (or its neighbor) as above the cutoff")
}
