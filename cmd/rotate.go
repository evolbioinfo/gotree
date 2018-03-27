package cmd

import (
	"github.com/spf13/cobra"
)

// rotateCmd represents the shuffletips command
var rotateCmd = &cobra.Command{
	Use:   "rotate",
	Short: "Rotates children of internal nodes",
	Long: `Rotates children of internal nodes by different means.

Either randomly with "rand" subcommand, either sorting by number of tips
with "sort" subcommand.

It does not change the topology, but just the order of neighbors 
of all node and thus the newick representation.

             ------C                    ------A
       x     |z	   	          x     |z	    
   A---------*ROOT     =>     B---------*ROOT  
             |t	   	                |t	    	 
             ------B 	                ------C

Example of usage:

gotree rotate rand -i t.nw
gotree rotate sort -i t.nw
`,
}

func init() {
	RootCmd.AddCommand(rotateCmd)
	rotateCmd.PersistentFlags().StringVarP(&intreefile, "input", "i", "stdin", "Input tree")
	rotateCmd.PersistentFlags().StringVarP(&outtreefile, "output", "o", "stdout", "Rotated tree output file")
}
