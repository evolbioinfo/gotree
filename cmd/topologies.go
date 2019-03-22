package cmd

import (
	goio "io"
	"os"

	"github.com/evolbioinfo/gotree/io"
	"github.com/evolbioinfo/gotree/tree"
	"github.com/spf13/cobra"
)

var generateIntreeFile string

// binarytreeCmd represents the binarytree command
var topologiesCmd = &cobra.Command{
	Use:   "topologies",
	Short: "Generates all possible tree topologies",
	Long:  `Generates all possible tree topologies.`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		var f *os.File
		var treefile goio.Closer
		var treechan <-chan tree.Trees
		var trees []*tree.Tree
		var tipNames []string = make([]string, 0)

		// If we have an input tree, we get the tip names of that tree
		if generateIntreeFile != "none" {
			if treefile, treechan, err = readTrees(generateIntreeFile); err != nil {
				io.LogError(err)
				return
			}
			defer treefile.Close()
			t := <-treechan
			if t.Err != nil {
				io.LogError(t.Err)
				return t.Err
			}
			tipNames = t.Tree.AllTipNames()
			generateNbTips = len(tipNames)
		}
		if trees, err = tree.AllTopologies(generateNbTips, generateRooted, tipNames...); err != nil {
			io.LogError(err)
			return
		} else {
			if f, err = openWriteFile(generateOutputfile); err != nil {
				io.LogError(err)
				return
			}
			for _, t := range trees {
				f.WriteString(t.Newick() + "\n")
			}
			closeWriteFile(f, generateOutputfile)
		}
		return
	},
}

func init() {
	generateCmd.AddCommand(topologiesCmd)
	topologiesCmd.PersistentFlags().IntVarP(&generateNbTips, "nbtips", "l", 10, "Number of tips/leaves of the trees to generate")
	topologiesCmd.PersistentFlags().StringVarP(&generateIntreeFile, "input", "i", "none", "Input Tree: Tip names of generate trees are taken from it")
}
