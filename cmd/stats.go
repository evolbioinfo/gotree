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
		statsout.WriteString(fmt.Sprintf("nodes\t%d", len(statsintree.Nodes())))
		statsout.WriteString(fmt.Sprintf("tips\t%d\n", len(statsintree.Tips())))
		statsout.WriteString(fmt.Sprintf("edges\t%d\n", len(statsintree.Edges())))
		statsout.WriteString(fmt.Sprintf("meanbrlen\t%.4f\n", statsintree.MeanBrLength()))
		statsout.WriteString(fmt.Sprintf("meansupport\t%.4f\n", statsintree.MeanSupport()))
		if statsintree.Rooted() {
			statsout.WriteString(fmt.Sprintf("root\trooted\n"))
		} else {
			statsout.WriteString(fmt.Sprintf("root\tunrooted\n"))
		}
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
