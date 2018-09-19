package cmd

import (
	goio "io"
	"math/rand"
	"os"
	"time"

	"github.com/fredericlemoine/gotree/io"
	"github.com/fredericlemoine/gotree/tree"
	"github.com/spf13/cobra"
)

// shuffletipsCmd represents the shuffletips command
var shuffletipsCmd = &cobra.Command{
	Use:   "shuffletips",
	Short: "Shuffle tip names of an input tree",
	Long: `Shuffle tip names of an input tree.


             ------C                    ------A
       x     |z	   	          x     |z	    
   A---------*ROOT     =>     B---------*ROOT  
             |t	   	                |t	    	 
             ------B 	                ------C

Example of usage:

gotree shuffletips -i t.nw

`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		var f *os.File
		var treefile goio.Closer
		var treechan <-chan tree.Trees

		// Read Tree
		rand.Seed(seed)
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
			t.Tree.ShuffleTips()
			f.WriteString(t.Tree.Newick() + "\n")
		}
		return
	},
}

func init() {
	RootCmd.AddCommand(shuffletipsCmd)
	shuffletipsCmd.Flags().Int64VarP(&seed, "seed", "s", time.Now().UTC().UnixNano(), "Initial Random Seed")
	shuffletipsCmd.PersistentFlags().StringVarP(&intreefile, "input", "i", "stdin", "Input tree")
	shuffletipsCmd.PersistentFlags().StringVarP(&outtreefile, "output", "o", "stdout", "Shuffled tree output file")
}
