package cmd

import (
	"github.com/fredericlemoine/gotree/io"
	"github.com/fredericlemoine/gotree/io/utils"
	"github.com/fredericlemoine/gotree/tree"
	"github.com/spf13/cobra"
	"os"
)

var collapseInputTree string
var collapseOutputTree string
var collapseIntrees chan tree.Trees
var collapseOutTrees *os.File

// collapseCmd represents the collapse command
var collapseCmd = &cobra.Command{
	Use:   "collapse",
	Short: "Commands to collapse branches of input trees",
	Long:  `Commands to collapse branches of input trees.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		var err error
		var nbtrees int = 0
		collapseIntrees = make(chan tree.Trees, 15)

		/* Read ref tree(s) */
		go func() {
			if nbtrees, err = utils.ReadCompTrees(collapseInputTree, collapseIntrees); err != nil {
				io.ExitWithMessage(err)
			}
		}()

		collapseOutTrees = openWriteFile(collapseOutputTree)
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		collapseOutTrees.Close()
	},
}

func init() {
	RootCmd.AddCommand(collapseCmd)
	collapseCmd.PersistentFlags().StringVarP(&collapseInputTree, "input", "i", "stdin", "Input tree")
	collapseCmd.PersistentFlags().StringVarP(&collapseOutputTree, "output", "o", "stdout", "Collapsed tree output file")
}
