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

// difftipsCmd represents the difftips command
var difftipsCmd = &cobra.Command{
	Use:   "tips",
	Short: "Print diff between tip names of two trees",
	Long: `Print diff between tip names of two trees.

For example:
t1.nh : (t1,t2,(t3,t4));
t2.nh : (t10,t2,(t3,t4));

gotree difftips -i t1.nh -c t2.nh

should produce the following output:
< t1
> t10
= 3

`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		var refTree *tree.Tree
		var treefile goio.Closer
		var treechan <-chan tree.Trees
		var ok bool

		if intree2file == "none" {
			err = errors.New("Compare tree file must be provided with -c")
			io.LogError(err)
			return
		}

		if refTree, err = readTree(intreefile); err != nil {
			io.LogError(err)
			return
		}
		refTree.UpdateTipIndex()
		if treefile, treechan, err = readTrees(intree2file); err != nil {
			io.LogError(err)
			return
		}
		defer treefile.Close()
		for compTree := range treechan {
			if compTree.Err != nil {
				io.LogError(compTree.Err)
				return compTree.Err
			}
			compTree.Tree.UpdateTipIndex()

			eq := 0
			for _, t := range refTree.Tips() {
				if ok, err = compTree.Tree.ExistsTip(t.Name()); err != nil {
					io.LogError(err)
					return
				} else {
					if !ok {
						fmt.Fprintf(os.Stdout, "(Tree %d) < %s\n", compTree.Id, t.Name())
					} else {
						eq++
					}
				}
			}
			for _, t := range compTree.Tree.Tips() {
				if ok, err = refTree.ExistsTip(t.Name()); err != nil {
					io.LogError(err)
					return
				} else {
					if !ok {
						fmt.Fprintf(os.Stdout, "(Tree %d) > %s\n", compTree.Id, t.Name())
					}
				}
			}
			fmt.Fprintf(os.Stdout, "(Tree %d) = %d\n", compTree.Id, eq)
		}
		return
	},
}

func init() {
	compareCmd.AddCommand(difftipsCmd)
}
