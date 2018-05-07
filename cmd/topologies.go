package cmd

import (
	"github.com/fredericlemoine/gotree/io"
	"github.com/fredericlemoine/gotree/tree"
	"github.com/spf13/cobra"
)

var generateIntreeFile string

// binarytreeCmd represents the binarytree command
var topologiesCmd = &cobra.Command{
	Use:   "topologies",
	Short: "Generates all possible tree topologies",
	Long:  `Generates all possible tree topologies.`,
	Run: func(cmd *cobra.Command, args []string) {
		var tipNames []string = make([]string, 0)
		// If we have an input tree, we get the tip names of that tree
		if generateIntreeFile != "none" {
			treefile, trees := readTrees(generateIntreeFile)
			defer treefile.Close()
			t := <-trees
			if t.Err != nil {
				io.ExitWithMessage(t.Err)
			}
			tipNames = t.Tree.AllTipNames()
			generateNbTips = len(tipNames)
		}
		if trees, err := tree.AllTopologies(generateNbTips, generateRooted, tipNames...); err != nil {
			io.ExitWithMessage(err)
		} else {
			f := openWriteFile(generateOutputfile)
			for _, t := range trees {
				f.WriteString(t.Newick() + "\n")
			}
			f.Close()
		}
	},
}

func init() {
	generateCmd.AddCommand(topologiesCmd)
	topologiesCmd.PersistentFlags().IntVarP(&generateNbTips, "nbtips", "l", 10, "Number of tips/leaves of the trees to generate")
	topologiesCmd.PersistentFlags().StringVarP(&generateIntreeFile, "input", "i", "none", "Input Tree: Tip names of generate trees are taken from it")
}
