package cmd

import (
	"fmt"
	"github.com/fredericlemoine/gotree/io"
	"github.com/fredericlemoine/gotree/io/utils"
	"github.com/fredericlemoine/gotree/tree"
	"github.com/spf13/cobra"
	"os"
)

var statsintreestr string
var statsoutfile string
var statsintree *tree.Tree
var statsout *os.File

// statsCmd represents the stats command
var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Displays different statistics about the tree",
	Long: `Displays different statistics about the tree

For exemple:
- Edge informations
- Node informations
- Tips informations

`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		var err error
		statsintree, err = utils.ReadRefTree(statsintreestr)
		if err != nil {
			io.ExitWithMessage(err)
		}
		statsintree.ComputeDepths()
		if statsoutfile != "stdout" {
			statsout, err = os.Create(statsoutfile)
		} else {
			statsout = os.Stdout
		}
		if err != nil {
			io.ExitWithMessage(err)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		statsout.WriteString("nodes\t" + fmt.Sprintf("%d", len(statsintree.Nodes())) + "\n")
		statsout.WriteString("tips\t" + fmt.Sprintf("%d", len(statsintree.Tips())) + "\n")
		statsout.WriteString("edges\t" + fmt.Sprintf("%d", len(statsintree.Edges())) + "\n")
		statsout.WriteString("meanbrlen\t" + fmt.Sprintf("%.4f", statsintree.MeanBrLength()) + "\n")
		statsout.WriteString("meansupport\t" + fmt.Sprintf("%.4f", statsintree.MeanSupport()) + "\n")
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		statsout.Close()
	},
}

func init() {
	RootCmd.AddCommand(statsCmd)
	statsCmd.PersistentFlags().StringVarP(&statsintreestr, "input", "i", "stdin", "Input tree")
	statsCmd.PersistentFlags().StringVarP(&statsoutfile, "output", "o", "stdout", "Output file")
}
