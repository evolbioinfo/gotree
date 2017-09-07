package cmd

import (
	"errors"
	"fmt"

	"github.com/fredericlemoine/gotree/io"
	"github.com/spf13/cobra"
)

var autorename bool
var autorenamelength int

// renameCmd represents the rename command
var renameCmd = &cobra.Command{
	Use:   "rename",
	Short: "Rename tips of the input tree",
	Long: `Rename tips of the input tree.

In default mode, a map file must be given (-m), and must be tab separated with columns:
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



If -a is given, then tips are renamed using automatically generated identifiers of length 10
Correspondance between old names and new names is written in the map file given with -m. 
In this mode, --revert has no effect.
--length  allows to customize length of generated id. It is min 5.
If several trees in input has different tip names, it does not matter, a new identifier is still
generated for each new tip name.

`,
	Run: func(cmd *cobra.Command, args []string) {

		if mapfile == "none" {
			io.ExitWithMessage(errors.New("map file is not given"))
		}
		var namemap map[string]string = nil
		var err error

		if !autorename {
			// Read map file
			namemap, err = readMapFile(mapfile, revert)
			if err != nil {
				io.ExitWithMessage(err)
			}
		} else {
			if autorenamelength < 5 {
				autorenamelength = 5
			}
			namemap = make(map[string]string)
		}

		f := openWriteFile(outtreefile)
		// Read ref Trees and rename them
		treefile, trees := readTrees(intreefile)
		defer treefile.Close()

		curid := 1
		for tr := range trees {
			if tr.Err != nil {
				io.ExitWithMessage(tr.Err)
			}

			if autorename {
				for _, t := range tr.Tree.Tips() {
					if _, ok := namemap[t.Name()]; !ok {
						newname := fmt.Sprintf(fmt.Sprintf("T%%0%dd", (autorenamelength-1)), curid)
						if len(newname) != autorenamelength {
							io.ExitWithMessage(fmt.Errorf("Id length %d does not allow to generate as much ids: %d (%s)", autorenamelength, curid, newname))
						}
						namemap[t.Name()] = newname
						curid++
					}
				}
			}

			err = tr.Tree.Rename(namemap)
			if err != nil {
				io.ExitWithMessage(err)
			}

			f.WriteString(tr.Tree.Newick() + "\n")
		}

		if autorename {
			writeNameMap(namemap, mapfile)
		}

		f.Close()
	},
}

func init() {
	RootCmd.AddCommand(renameCmd)
	renameCmd.Flags().StringVarP(&outtreefile, "output", "o", "stdout", "Renamed tree output file")
	renameCmd.Flags().StringVarP(&intreefile, "input", "i", "stdin", "Input tree")
	renameCmd.Flags().StringVarP(&mapfile, "map", "m", "none", "Tip name map file")
	renameCmd.Flags().BoolVarP(&autorename, "auto", "a", false, "Renames automatically tips with auto generated id of length 10.")
	renameCmd.Flags().IntVarP(&autorenamelength, "length", "l", 10, "Length of automatically generated id. Only with --auto")
	renameCmd.Flags().BoolVarP(&revert, "revert", "r", false, "Revert orientation of map file")
}

func writeNameMap(namemap map[string]string, outfile string) {
	f := openWriteFile(outfile)
	for old, new := range namemap {
		f.WriteString(old)
		f.WriteString("\t")
		f.WriteString(new)
		f.WriteString("\n")
	}
	f.Close()
}
