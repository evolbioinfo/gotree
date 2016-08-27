package cmd

import (
	"github.com/fredericlemoine/gotree/io"
	"github.com/fredericlemoine/gotree/io/utils"
	"github.com/fredericlemoine/gotree/tree"
	"github.com/spf13/cobra"
	"os"
)

// unrootCmd represents the unroot command
var unrootCmd = &cobra.Command{
	Use:   "unroot",
	Short: "Unroots the input tree",
	Long: `Unroots the input tree.

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

gotree tranform unroot -i tree.nw -o tree_u.nw

`,
	Run: func(cmd *cobra.Command, args []string) {
		// Read Tree
		var t *tree.Tree
		var err error
		t, err = utils.ReadRefTree(transformInputTree)
		if err != nil {
			io.ExitWithMessage(err)
		}
		var f *os.File
		if renameouttree != "stdout" {
			f, err = os.Create(transformOutputTree)
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
	transformCmd.AddCommand(unrootCmd)
}
