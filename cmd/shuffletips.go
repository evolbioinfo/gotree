package cmd

import (
	"github.com/fredericlemoine/gotree/io"
	"github.com/fredericlemoine/gotree/io/utils"
	"github.com/fredericlemoine/gotree/tree"
	"github.com/spf13/cobra"
	"math/rand"
	"os"
	"time"
)

var shuffletipsSeed int64
var shuffleTipsInputTree string
var shuffleTipsOutputTree string

// shuffletipsCmd represents the shuffletips command
var shuffletipsCmd = &cobra.Command{
	Use:   "shuffletips",
	Short: "Shuffles the tip names of an input tree",
	Long: `Shuffles the tip names of an input tree.


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

		rand.Seed(shuffletipsSeed)

		var t *tree.Tree
		var err error
		t, err = utils.ReadRefTree(shuffleTipsInputTree)
		if err != nil {
			io.ExitWithMessage(err)
		}
		var f *os.File
		if shuffleTipsOutputTree != "stdout" {
			f, err = os.Create(shuffleTipsOutputTree)
		} else {
			f = os.Stdout
		}
		if err != nil {
			io.ExitWithMessage(err)
		}

		t.ShuffleTips()

		if err != nil {
			io.ExitWithMessage(err)
		}

		f.WriteString(t.Newick() + "\n")
		f.Close()
	},
}

func init() {
	RootCmd.AddCommand(shuffletipsCmd)
	shuffletipsCmd.Flags().Int64VarP(&shuffletipsSeed, "seed", "s", time.Now().UTC().UnixNano(), "Initial Random Seed")
	shuffletipsCmd.PersistentFlags().StringVarP(&shuffleTipsInputTree, "input", "i", "stdin", "Input tree")
	shuffletipsCmd.PersistentFlags().StringVarP(&shuffleTipsOutputTree, "output", "o", "stdout", "Shuffled tree output file")

}
