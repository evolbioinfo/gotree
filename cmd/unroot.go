package cmd

import (
	goio "io"
	"os"

	"github.com/fredericlemoine/gotree/io"
	"github.com/fredericlemoine/gotree/tree"
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
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		var f *os.File
		var treefile goio.Closer
		var treechan <-chan tree.Trees

		if f, err = openWriteFile(outtreefile); err != nil {
			io.LogError(err)
			return
		}
		defer closeWriteFile(f, outtreefile)

		if treefile, treechan, err = readTrees(intreefile); err != nil {
			io.LogError(err)
			return
		}
		defer treefile.Close()
		for t := range treechan {
			if t.Err != nil {
				io.LogError(t.Err)
				return t.Err
			}
			t.Tree.UnRoot()
			f.WriteString(t.Tree.Newick() + "\n")
		}
		return
	},
}

func init() {
	RootCmd.AddCommand(unrootCmd)

	unrootCmd.PersistentFlags().StringVarP(&intreefile, "input", "i", "stdin", "Input tree")
	unrootCmd.PersistentFlags().StringVarP(&outtreefile, "output", "o", "stdout", "Collapsed tree output file")

}
