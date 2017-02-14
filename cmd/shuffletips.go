package cmd

import (
	"github.com/spf13/cobra"
	"math/rand"
	"time"
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
		for t2 := range readTrees(intreefile) {
			t2.Tree.ShuffleTips()
			f.WriteString(t2.Tree.Newick() + "\n")
		}
		f.Close()
	},
}

func init() {
	RootCmd.AddCommand(shuffletipsCmd)
	shuffletipsCmd.Flags().Int64VarP(&seed, "seed", "s", time.Now().UTC().UnixNano(), "Initial Random Seed")
	shuffletipsCmd.PersistentFlags().StringVarP(&intreefile, "input", "i", "stdin", "Input tree")
	shuffletipsCmd.PersistentFlags().StringVarP(&outtreefile, "output", "o", "stdout", "Shuffled tree output file")
}
