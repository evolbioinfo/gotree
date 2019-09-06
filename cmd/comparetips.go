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
	Short: "Print diff between tip names of two trees or between a tree and a list of tips",
	Long: `Print diff between tip names of two trees or between a tree and a list of tips.

* Example between 2 trees:
t1.nh : (t1,t2,(t3,t4));
t2.nh : (t10,t2,(t3,t4));

gotree difftips -i t1.nh -c t2.nh

should produce the following output:
< t1
> t10
= 3

* Example between a tree and a list of tips:
t1.nh : (t1,t2,(t3,t4));
t2.txt : 

t10
t2
t3
t4

gotree difftips -i t1.nh -f t2.txt

should produce the following output:
< t1
> t10
= 3

* Options
-c has priority over -f

`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		var refTree *tree.Tree
		var treefile goio.Closer
		var treechan <-chan tree.Trees
		var ok bool
		var tips map[string]bool

		if intree2file == "none" && tipfile == "none" {
			err = errors.New("At least a compare tree file or a tip list file must be provided with -c or -f")
			io.LogError(err)
			return
		}

		if refTree, err = readTree(intreefile); err != nil {
			io.LogError(err)
			return
		}
		refTree.UpdateTipIndex()

		if intree2file != "none" {
			// We compare with a tree
			if treefile, treechan, err = readTrees(intree2file); err != nil {
				io.LogError(err)
				return
			}
			defer treefile.Close()
		} else {
			//We compare with a tip list
			tips = make(map[string]bool)
			var tipstmp []string
			if tipstmp, err = parseTipsFile(tipfile); err != nil {
				io.LogError(err)
				return
			}
			for _, v := range tipstmp {
				tips[v] = true
			}
		}

		// Compare to an other tree
		if treechan != nil {
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
					}
					if !ok {
						fmt.Fprintf(os.Stdout, "(Tree %d) < %s\n", compTree.Id, t.Name())
					} else {
						eq++
					}
				}

				for _, t := range compTree.Tree.Tips() {
					if ok, err = refTree.ExistsTip(t.Name()); err != nil {
						io.LogError(err)
						return
					}
					if !ok {
						fmt.Fprintf(os.Stdout, "(Tree %d) > %s\n", compTree.Id, t.Name())
					}
				}
				fmt.Fprintf(os.Stdout, "(Tree %d) = %d\n", compTree.Id, eq)
			}
		} else {
			eq := 0
			for _, t := range refTree.Tips() {
				if _, ok = tips[t.Name()]; !ok {
					fmt.Fprintf(os.Stdout, "(Tree %d) < %s\n", 0, t.Name())
				} else {
					eq++
				}
			}

			for k, _ := range tips {
				if ok, err = refTree.ExistsTip(k); err != nil {
					io.LogError(err)
					return
				}
				if !ok {
					fmt.Fprintf(os.Stdout, "(Tree %d) > %s\n", 0, k)
				}
			}
			fmt.Fprintf(os.Stdout, "(Tree %d) = %d\n", 0, eq)
		}
		return
	},
}

func init() {
	compareCmd.AddCommand(difftipsCmd)
	difftipsCmd.PersistentFlags().StringVarP(&tipfile, "tipfile", "f", "none", "Tip File (Optional)")
}
