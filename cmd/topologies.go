package cmd

import (
	"github.com/fredericlemoine/gotree/io"
	"github.com/fredericlemoine/gotree/tree"
	"github.com/spf13/cobra"
)

// binarytreeCmd represents the binarytree command
var topologiesCmd = &cobra.Command{
	Use:   "topologies",
	Short: "Generates all possible tree topologies",
	Long:  `Generates all possible tree topologies.`,
	Run: func(cmd *cobra.Command, args []string) {
		if trees, err := tree.AllTopologies(generateNbTips, generateRooted); err != nil {
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
}
