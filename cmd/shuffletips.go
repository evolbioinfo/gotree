package cmd

import (
	"math/rand"
	"time"

	"github.com/fredericlemoine/gotree/io"
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
			t.Tree.ShuffleTips()
			f.WriteString(t.Tree.Newick() + "\n")
		}
	},
}

func init() {
	RootCmd.AddCommand(shuffletipsCmd)
	shuffletipsCmd.Flags().Int64VarP(&seed, "seed", "s", time.Now().UTC().UnixNano(), "Initial Random Seed")
	shuffletipsCmd.PersistentFlags().StringVarP(&intreefile, "input", "i", "stdin", "Input tree")
	shuffletipsCmd.PersistentFlags().StringVarP(&outtreefile, "output", "o", "stdout", "Shuffled tree output file")
}
