package cmd

import (
	"github.com/fredericlemoine/gotree/io"
	"github.com/spf13/cobra"
)

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
		f := openWriteFile(outtreefile)
		treefile, treechan := readTrees(intreefile)
		defer treefile.Close()
		for t := range treechan {
			if t.Err != nil {
				io.ExitWithMessage(t.Err)
			}
			t.Tree.UnRoot()
			f.WriteString(t.Tree.Newick() + "\n")
		}
		f.Close()
	},
}

func init() {
	RootCmd.AddCommand(unrootCmd)

	unrootCmd.PersistentFlags().StringVarP(&intreefile, "input", "i", "stdin", "Input tree")
	unrootCmd.PersistentFlags().StringVarP(&outtreefile, "output", "o", "stdout", "Collapsed tree output file")

}
