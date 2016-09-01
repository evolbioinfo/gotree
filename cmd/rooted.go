package cmd

import (
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
		if statsintree.Rooted() {
			statsout.WriteString("rooted\n")
		} else {
			statsout.WriteString("unrooted\n")
		}
	},
}

func init() {
	statsCmd.AddCommand(rootedCmd)
}
