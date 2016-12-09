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
	Run: func(cmd *cobra.Command, args []string) {
		// Read Tree
		var t *tree.Tree
		var err error
		t, err = utils.ReadRefTree(unrootInputTree)
		if err != nil {
			io.ExitWithMessage(err)
		}
		var f *os.File
		if unrootOutputTree != "stdout" {
			f, err = os.Create(unrootOutputTree)
		} else {
			f = os.Stdout
		}
		if err != nil {
			io.ExitWithMessage(err)
		}

		t.UnRoot()

		f.WriteString(t.Newick() + "\n")
		f.Close()
	},
}

func init() {
	RootCmd.AddCommand(unrootCmd)

	unrootCmd.PersistentFlags().StringVarP(&unrootInputTree, "input", "i", "stdin", "Input tree")
	unrootCmd.PersistentFlags().StringVarP(&unrootOutputTree, "output", "o", "stdout", "Collapsed tree output file")

}
