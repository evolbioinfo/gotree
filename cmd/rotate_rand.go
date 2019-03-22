package cmd

import (
	goio "io"
	"os"

	"github.com/evolbioinfo/gotree/io"
	"github.com/evolbioinfo/gotree/tree"
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
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		var f *os.File
		var treefile goio.Closer
		var treechan <-chan tree.Trees

		// Read Tree
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
			t.Tree.RotateInternalNodes()
			f.WriteString(t.Tree.Newick() + "\n")
		}
		return
	},
}

func init() {
	rotateCmd.AddCommand(rotateRandCmd)
}
