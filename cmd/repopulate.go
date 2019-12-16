package cmd

import (
	"fmt"
	goio "io"
	"os"

	"github.com/evolbioinfo/gotree/io"
	"github.com/evolbioinfo/gotree/tree"
	"github.com/spf13/cobra"
)

var groupfile string

// renameCmd represents the rename command
var repopulateCmd = &cobra.Command{
	Use:   "repopulate",
	Short: "Re populate the tree with identical tips",
	Long: `Re populate the tree with tips that have the same sequences.

When a tree is inferred, some tools first remove identical sequences.

However, it may be useful to keep them in the tree. To do so, this command takes:

1. A input tree
2. A file containing a list of tips that are identical, in the following format:
    Tip1,Tip2
    Tip3,Tip4
    Meaning that Tip1 is identical to Tip2, and Tip3 is identical to Tip4.

"repopulate" command then adds Tip2 next to Tip1 if Tip1 is present in the tree, or 
Tip1 next to Tip2 if Tip2 is present in the tree. To do so, it adds two 0.0 length
 branches. 

Example with Tip1,Tip2 :

 Before     |   After (if l>0.0)  |  After (if l=0.0)
------------+---------------------+-------------------
            |         *Tip1       |     *Tip1
    l       |    l   /.0          |    /0.0
 *----*Tip1 |   ----*	          |   *
            |        \.0          |    \0.0
            |         *Tip2       |     *Tip2

Each identical group must contain exactly 1 already present tip, otherwise it returns
 an error.

If a new tip is present in several groups, then returns and error.

The tree after "repopulate" command may contain polytomies.
`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		var f *os.File
		var treefile goio.Closer
		var treechan <-chan tree.Trees
		var identicalgroups [][]string
		var setgroups bool

		setgroups = cmd.Flags().Changed("id-groups")

		if !setgroups {
			err = fmt.Errorf("File with groups of identical tips must be provided")
			io.LogError(err)
			return
		}

		identicalgroups, err = readIdenticalGroupFile(groupfile)

		if f, err = openWriteFile(outtreefile); err != nil {
			io.LogError(err)
			return
		}
		defer closeWriteFile(f, outtreefile)

		// Read ref Trees and rename them
		if treefile, treechan, err = readTrees(intreefile); err != nil {
			io.LogError(err)
			return
		}
		defer treefile.Close()

		for tr := range treechan {
			if tr.Err != nil {
				io.LogError(tr.Err)
				return tr.Err
			}

			if err = tr.Tree.UpdateTipIndex(); err != nil {
				io.LogError(err)
				return
			}

			if err = tr.Tree.InsertIdenticalTips(identicalgroups); err != nil {
				io.LogError(err)
				return
			}

			f.WriteString(tr.Tree.Newick() + "\n")
		}
		return
	},
}

func init() {
	RootCmd.AddCommand(repopulateCmd)
	repopulateCmd.Flags().StringVarP(&outtreefile, "output", "o", "stdout", "Renamed tree output file")
	repopulateCmd.Flags().StringVarP(&intreefile, "input", "i", "stdin", "Input tree")
	repopulateCmd.Flags().StringVarP(&groupfile, "id-groups", "g", "none", "File with groups of identical tips")
}
