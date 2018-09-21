package cmd

import (
	goio "io"
	"os"

	"github.com/fredericlemoine/gotree/io"
	"github.com/fredericlemoine/gotree/tree"
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
			t.Tree.SortNeighborsByTips()
			f.WriteString(t.Tree.Newick() + "\n")
		}
		return
	},
}

func init() {
	rotateCmd.AddCommand(rotateSortCmd)
}
