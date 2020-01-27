package cmd

import (
	"errors"
	"fmt"
	goio "io"
	"os"

	"github.com/evolbioinfo/gotree/io"
	"github.com/evolbioinfo/gotree/tree"
	"github.com/spf13/cobra"
)

// outgroupCmd represents the outgroup command
var monoCmd = &cobra.Command{
	Use:   "monophyletic",
	Short: "Tells wether input tips form a monophyletic group in each of the input trees",
	Long: `Tells wether input tips form a monophyletic group in each of the input trees.

Returns true for each tree in which the given tips form a monophyletic group (form a clade containing no other tips).
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

		fmt.Fprintf(f, "Tree\tMonophyletic\n")
		for t := range treechan {
			if t.Err != nil {
				io.LogError(t.Err)
				return t.Err
			}

			var monophyletic bool

			if !t.Tree.Rooted() {
				if _, _, monophyletic, err = t.Tree.LeastCommonAncestorUnrooted(nil, tips...); err != nil {
					io.LogError(err)
					return
				}
			} else {
				if _, _, monophyletic, err = t.Tree.LeastCommonAncestorRooted(nil, tips...); err != nil {
					io.LogError(err)
					return
				}
			}

			fmt.Fprintf(f, "%d\t%t\n", t.Id, monophyletic)
		}
		return
	},
}

func init() {
	statsCmd.AddCommand(monoCmd)
	monoCmd.PersistentFlags().StringVarP(&tipfile, "tip-file", "l", "none", "File containing names of tips of the outgroup")
}
