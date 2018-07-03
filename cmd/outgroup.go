package cmd

import (
	"errors"

	"github.com/fredericlemoine/gotree/io"
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
	Run: func(cmd *cobra.Command, args []string) {
		var tips []string
		if tipfile != "none" {
			tips = parseTipsFile(tipfile)
		} else if len(args) > 0 {
			tips = args
		} else {
			io.ExitWithMessage(errors.New("Not group given"))
		}

		f := openWriteFile(outtreefile)
		treefile, trees := readTrees(intreefile)
		defer treefile.Close()

		for t := range trees {
			if t.Err != nil {
				io.ExitWithMessage(t.Err)
			}
			err := t.Tree.RerootOutGroup(removeoutgroup, rerootstrict, tips...)
			if err != nil {
				io.ExitWithMessage(err)
			}

			f.WriteString(t.Tree.Newick() + "\n")
		}

		f.Close()
	},
}

func init() {
	rerootCmd.AddCommand(outgroupCmd)
	outgroupCmd.PersistentFlags().StringVarP(&tipfile, "tip-file", "l", "none", "File containing names of tips of the outgroup")
	outgroupCmd.PersistentFlags().BoolVarP(&removeoutgroup, "remove-outgroup", "r", false, "Removes the outgroup after reroot")
	outgroupCmd.PersistentFlags().BoolVar(&rerootstrict, "strict", false, "Enforce the outgroup to be monophyletic (else throw an error)")
}
