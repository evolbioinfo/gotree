package cmd

import (
	"github.com/fredericlemoine/gotree/io"
	"github.com/fredericlemoine/gotree/io/utils"
	"github.com/fredericlemoine/gotree/tree"
	"github.com/spf13/cobra"
	"os"
)

var clearInputTree string
var clearOutputTree string
var clearIntrees chan tree.Trees
var clearOutTrees *os.File

// clearCmd represents the clear command
var clearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clear lengths or supports from input trees",
	Long:  `Clear lengths or supports from input trees`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		var err error
		var nbtrees int = 0
		clearIntrees = make(chan tree.Trees, 15)

		/* Read ref tree(s) */
		go func() {
			if nbtrees, err = utils.ReadCompTrees(clearInputTree, clearIntrees); err != nil {
				io.ExitWithMessage(err)
			}
		}()

		clearOutTrees = openWriteFile(clearOutputTree)
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		clearOutTrees.Close()
	},
}

func init() {
	RootCmd.AddCommand(clearCmd)
	clearCmd.PersistentFlags().StringVarP(&clearInputTree, "input", "i", "stdin", "Input tree")
	clearCmd.PersistentFlags().StringVarP(&clearOutputTree, "output", "o", "stdout", "Cleared tree output file")
}
