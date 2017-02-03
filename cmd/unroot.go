package cmd

import (
	"github.com/fredericlemoine/gotree/io"
	"github.com/fredericlemoine/gotree/io/utils"
	"github.com/fredericlemoine/gotree/tree"
	"github.com/spf13/cobra"
	"os"
)

var unrootInputTree string
var unrootOutputTree string
var unrootIntrees chan tree.Trees
var unrootOutTrees *os.File

// unrootCmd represents the unroot command
var unrootCmd = &cobra.Command{
	Use:   "unroot",
	Short: "Unroot input tree",
	Long: `Unroot input tree.

If the tree is already unrooted does nothing
Otherwise places the root on a trifurcated node and removes
old root.
br length : Take the sum
br support: Take the max

             ------C         
             |z	         
    ---------*	                       ------C 
    |x       |t	                 x+y   |z	   
ROOT*        ------B   =>    A---------*ROOT   
    |y		                       |t	   		 
    ---*A                              ------B 

Example of usage:

gotree unroot -i tree.nw -o tree_u.nw

`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		var err error
		var nbtrees int = 0
		unrootIntrees = make(chan tree.Trees, 15)
		/* Read ref tree(s) */
		go func() {
			if nbtrees, err = utils.ReadCompTrees(unrootInputTree, unrootIntrees); err != nil {
				io.ExitWithMessage(err)
			}
		}()
		unrootOutTrees = openWriteFile(unrootOutputTree)
	},
	Run: func(cmd *cobra.Command, args []string) {
		for t := range unrootIntrees {
			t.Tree.UnRoot()
			unrootOutTrees.WriteString(t.Tree.Newick() + "\n")
		}
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		unrootOutTrees.Close()
	},
}

func init() {
	RootCmd.AddCommand(unrootCmd)

	unrootCmd.PersistentFlags().StringVarP(&unrootInputTree, "input", "i", "stdin", "Input tree")
	unrootCmd.PersistentFlags().StringVarP(&unrootOutputTree, "output", "o", "stdout", "Collapsed tree output file")

}
