package cmd

import (
	"errors"
	"github.com/fredericlemoine/gotree/io"
	"github.com/spf13/cobra"
)

// renameCmd represents the rename command
var renameCmd = &cobra.Command{
	Use:   "rename",
	Short: "Rename tips of the input tree, given a map file",
	Long: `Rename tips of the input tree, given a map file.

Map file must be tab separated with columns:
1) Current name of the tip
2) Desired new name of the tip
(if --revert then it is the other way)

If a tip name does not appear in the map file, it will not be renamed. 
If a name that does not exist appears in the map file, it will not throw an error.

Example :

MapFile :
A   A2
B   B2
C   C2

gotree rename -m MapFile -i t.nw

             ------C                   ------C2
       x     |z	     	        x      |z	    
   A---------*ROOT    =>    A2---------*ROOT  
             |t	     	               |t	    
             ------B 	               ------B2

`,
	Run: func(cmd *cobra.Command, args []string) {

		if mapfile == "none" {
			io.ExitWithMessage(errors.New("map file is not given"))
		}

		// Read map file
		namemap, err := readMapFile(mapfile, revert)
		if err != nil {
			io.ExitWithMessage(err)
		}

		f := openWriteFile(outtreefile)
		// Read ref Trees and rename them
		for reftree := range readTrees(intreefile) {
			err = reftree.Tree.Rename(namemap)
			if err != nil {
				io.ExitWithMessage(err)
			}

			f.WriteString(reftree.Tree.Newick() + "\n")
		}
		f.Close()
	},
}

func init() {
	RootCmd.AddCommand(renameCmd)
	renameCmd.Flags().StringVarP(&outtreefile, "output", "o", "stdout", "Renamed tree output file")
	renameCmd.Flags().StringVarP(&intreefile, "input", "i", "stdin", "Input tree")
	renameCmd.Flags().StringVarP(&mapfile, "map", "m", "none", "Tip name map file")
	renameCmd.Flags().BoolVarP(&revert, "revert", "r", false, "Revert orientation of map file")

}
