package cmd

import (
	"errors"
	goio "io"
	"os"

	"github.com/evolbioinfo/gotree/io"
	"github.com/evolbioinfo/gotree/tree"
	"github.com/spf13/cobra"
)

// outgroupCmd represents the outgroup command
var outgroupCmd = &cobra.Command{
	Use:   "outgroup",
	Short: "Reroot trees using an outgroup",
	Long: `Reroot the tree using an outgroup given in argument or in stdin.

Example:

Reroot on 1 tip named "Tip10" using stdin:
echo "Tip10" | gotree reroot outgroup -i tree.nw -l - > reroot.nw

Reroot using an outgroup defined by 3 tips using stdin:
echo "Tip1,Tip2,Tip10" | gotree reroot outgroup -i tree.nw -l - > reroot.nw

Reroot using an outgroup defined by 3 tips using command args:

gotree reroot outgroup -i tree.nw Tip1 Tip2 Tip3 > reroot.nw

If the outgroup includes a tip that is not present in the tree,
this tip will not be taken into account for the reroot. A warning
will be issued.

By default (--strict=false), if the outgroup is not monophyletic it will
take all the descendant of the LCA to reroot and print a warning.If the
outgroup is not monophyletic and if --strict is given, it exits with an 
error.

If the option -r|--remove-outgroup is given, then the outgroup is
removed after reroot.
`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		var f *os.File
		var treefile goio.Closer
		var treechan <-chan tree.Trees

		var tips []string
		if tipfile != "none" {
			if tips, err = parseTipsFile(tipfile); err != nil {
				io.LogError(err)
				return
			}
		} else if len(args) > 0 {
			tips = args
		} else {
			err = errors.New("Not group given")
			io.LogError(err)
			return
		}

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
			err = t.Tree.RerootOutGroup(removeoutgroup, rerootstrict, tips...)
			if err != nil {
				io.LogError(err)
				return
			}

			f.WriteString(t.Tree.Newick() + "\n")
		}
		return
	},
}

func init() {
	rerootCmd.AddCommand(outgroupCmd)
	outgroupCmd.PersistentFlags().StringVarP(&tipfile, "tip-file", "l", "none", "File containing names of tips of the outgroup")
	outgroupCmd.PersistentFlags().BoolVarP(&removeoutgroup, "remove-outgroup", "r", false, "Removes the outgroup after reroot")
	outgroupCmd.PersistentFlags().BoolVar(&rerootstrict, "strict", false, "Enforce the outgroup to be monophyletic (else throw an error)")
}
