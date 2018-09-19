package cmd

import (
	"math/rand"

	"github.com/fredericlemoine/gotree/io"
	"github.com/spf13/cobra"
)

// rotateCmd represents the shuffletips command
var rotateSortCmd = &cobra.Command{
	Use:   "sort",
	Short: "Sorts children of internal nodes by number of tips",
	Long: `Sorts children of internal nodes by number of tips.

It does not change the topology, but just the order of neighbors 
of all node and thus the newick representation.

             ------C                    ------A
       x     |z	   	          x     |z	    
   A---------*ROOT     =>     B---------*ROOT  
             |t	   	                |t	    	 
             ------B 	                ------C

Example of usage:

gotree rotate sort -i t.nw
`,
	Run: func(cmd *cobra.Command, args []string) {
		// Read Tree
		rand.Seed(seed)
		f := openWriteFile(outtreefile)
		defer closeWriteFile(f, outtreefile)

		treefile, treechan := readTrees(intreefile)
		defer treefile.Close()
		for t := range treechan {
			if t.Err != nil {
				io.ExitWithMessage(t.Err)
			}
			t.Tree.SortNeighborsByTips()
			f.WriteString(t.Tree.Newick() + "\n")
		}
	},
}

func init() {
	rotateCmd.AddCommand(rotateSortCmd)
}
