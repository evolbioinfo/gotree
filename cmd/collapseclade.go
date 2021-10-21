package cmd

import (
	"errors"
	goio "io"
	"os"

	"github.com/evolbioinfo/gotree/io"
	"github.com/evolbioinfo/gotree/tree"
	"github.com/spf13/cobra"
)

var cladestrict bool
var cladetipname string

// collapceClade represents the collapse command
var collapseClade = &cobra.Command{
	Use:   "clade",
	Short: "Collaps the clade defined by the given tip names",
	Long: `Collapse the clade defined by the given tip names, and replace it by a tip with a given name.

Example:

gotree collapse clade -i tree.nw -l tip.txt -n newtip
or
gotree collapse clade -i tree.nw -n newtip tip1 tip2 tip3
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
			if err = t.Tree.CollapseClade(cladestrict, cladetipname, tips...); err != nil {
				io.LogError(err)
				return
			}

			f.WriteString(t.Tree.Newick() + "\n")
		}
		return
	},
}

func init() {
	collapseCmd.AddCommand(collapseClade)
	collapseClade.PersistentFlags().StringVarP(&tipfile, "tip-file", "l", "none", "File containing names of tips of the outgroup")
	collapseClade.PersistentFlags().StringVarP(&cladetipname, "tip-name", "n", "none", "Name of the tip that will replace the clade")
	collapseClade.PersistentFlags().BoolVar(&cladestrict, "strict", false, "Enforce the outgroup to be monophyletic (else throw an error)")
}
