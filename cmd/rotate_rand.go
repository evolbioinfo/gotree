package cmd

import (
	"math/rand"
	"time"

	"github.com/fredericlemoine/gotree/io"
	"github.com/spf13/cobra"
)

// rotateCmd represents the shuffletips command
var rotateRandCmd = &cobra.Command{
	Use:   "rand",
	Short: "Randomly rotates children of internal nodes",
	Long: `Randomly rotates children of internal nodes.

It does not change the topology, but just the order of neighbors 
of all node and thus the newick representation.

             ------C                    ------A
       x     |z	   	          x     |z	    
   A---------*ROOT     =>     B---------*ROOT  
             |t	   	                |t	    	 
             ------B 	                ------C

Example of usage:

gotree rotate rand -i t.nw
`,
	Run: func(cmd *cobra.Command, args []string) {
		// Read Tree
		rand.Seed(seed)
		f := openWriteFile(outtreefile)
		treefile, treechan := readTrees(intreefile)
		defer treefile.Close()
		for t := range treechan {
			if t.Err != nil {
				io.ExitWithMessage(t.Err)
			}
			t.Tree.RotateInternalNodes()
			f.WriteString(t.Tree.Newick() + "\n")
		}
		f.Close()
	},
}

func init() {
	rotateCmd.AddCommand(rotateRandCmd)
	rotateRandCmd.Flags().Int64VarP(&seed, "seed", "s", time.Now().UTC().UnixNano(), "Initial Random Seed")
}
