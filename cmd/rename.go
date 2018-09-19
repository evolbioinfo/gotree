package cmd

import (
	"errors"

	"github.com/fredericlemoine/gotree/io"
	"github.com/spf13/cobra"
)

var autorename bool
var autorenamelength int
var renameInternalNodes bool
var renameTips bool
var renameRegex string
var renameReplaceBy string

// renameCmd represents the rename command
var renameCmd = &cobra.Command{
	Use:   "rename",
	Short: "Rename nodes/tips of the input tree",
	Long: `Rename nodes/tips of the input tree.

* In default mode, only tips are renamed (--tips=true by default), 
  and a map file must be given (-m), and must be tab separated with columns:
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



* If -a is given, then tips/nodes are renamed using automatically generated identifiers 
  of length 10 Correspondance between old names and new names is written in the map file 
  given with -m. 
  In this mode, --revert has no effect.
  --length  allows to customize length of generated id. It is min 5.
  If several trees in input has different tip names, it does not matter, a new identifier is still
  generated for each new tip name, and same names are reused if needed.

* If -e (--regexp) and -b (--replace) is given, then  will replace matching strings in tip/node 
  names by string given by -b, ex:
  gotree rename -i tree.nh --regexp 'Tip(\d+)' --replace 'Leaf$1' -m map.txt
  this will replace all matches of 'Tip(\d+)' with 'Leaf$1', with $1 being the matched string 
  inside ().


Warning: If after this rename, several tips/nodes have the same name, subsequent commands may 
fail.


If --internal is specified, then internal nodes are renamed;
--tips is true by default. To inactivate it, you must specify --tips=false .
`,
	Run: func(cmd *cobra.Command, args []string) {
		var namemap map[string]string = nil
		var err error
		var setregex, setreplace bool
		setregex = cmd.Flags().Changed("regexp")
		setreplace = cmd.Flags().Changed("replace")

		if !(renameTips || renameInternalNodes) {
			io.ExitWithMessage(errors.New("You should rename at least internal nodes (--internal) or tips (--tips)"))
		}
		if setregex && !setreplace {
			io.ExitWithMessage(errors.New("--replace must be given with --regexp"))
		}

		if !autorename && !setregex {
			// Read map file
			if mapfile == "none" {
				io.ExitWithMessage(errors.New("map file is not given"))
			}
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
		defer closeWriteFile(f, outtreefile)

		// Read ref Trees and rename them
		treefile, trees := readTrees(intreefile)
		defer treefile.Close()

		curid := 1
		for tr := range trees {
			if tr.Err != nil {
				io.ExitWithMessage(tr.Err)
			}

			if autorename {
				err = tr.Tree.RenameAuto(renameInternalNodes, renameTips, autorenamelength, &curid, namemap)
				if err != nil {
					io.ExitWithMessage(err)
				}
			} else if setregex {
				err = tr.Tree.RenameRegexp(renameInternalNodes, renameTips, renameRegex, renameReplaceBy, namemap)
				if err != nil {
					io.ExitWithMessage(err)
				}
			} else {
				err = tr.Tree.Rename(namemap)
				if err != nil {
					io.ExitWithMessage(err)
				}
			}

			f.WriteString(tr.Tree.Newick() + "\n")
		}

		if (autorename || setregex) && mapfile != "none" {
			writeNameMap(namemap, mapfile)
		}
	},
}

func init() {
	RootCmd.AddCommand(renameCmd)
	renameCmd.Flags().StringVarP(&outtreefile, "output", "o", "stdout", "Renamed tree output file")
	renameCmd.Flags().StringVarP(&intreefile, "input", "i", "stdin", "Input tree")
	renameCmd.Flags().BoolVar(&renameInternalNodes, "internal", false, "Internal nodes are taken into account")
	renameCmd.Flags().BoolVar(&renameTips, "tips", true, "Tips are taken into account (--tips=false to cancel)")
	renameCmd.Flags().StringVarP(&mapfile, "map", "m", "none", "Tip name map file")
	renameCmd.Flags().StringVarP(&renameRegex, "regexp", "e", "none", "Regexp to get matching tip/node names")
	renameCmd.Flags().StringVarP(&renameReplaceBy, "replace", "b", "none", "String replacement to the given regexp")
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
