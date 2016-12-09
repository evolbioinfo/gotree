package cmd

import (
	"github.com/fredericlemoine/gotree/io"
	"github.com/fredericlemoine/gotree/io/utils"
	"github.com/fredericlemoine/gotree/tree"
	"github.com/spf13/cobra"
	"os"
)

var prunereftree string
var prunecomptree string
var pruneouttree string

func specificTips(ref *tree.Tree, comp *tree.Tree) []string {
	compmap := make(map[string]*tree.Node)
	spectips := make([]string, 0)
	for _, n := range comp.Nodes() {
		if n.Nneigh() == 1 {
			compmap[n.Name()] = n
		}
	}

	for _, n := range ref.Nodes() {
		if n.Nneigh() == 1 {
			_, ok := compmap[n.Name()]
			if !ok {
				spectips = append(spectips, n.Name())
			}
		}
	}
	return spectips
}

// pruneCmd represents the prune command
var pruneCmd = &cobra.Command{
	Use:   "prune",
	Short: "Remove tips of the input tree that are not in the compared tree",
	Long: `This tool removes tips of the input reference tree that 
are not present in the compared tree.

In output, we have a tree containing only tips that are common to both trees.

If several trees are present in the file given by -i, they are all analyzed and 
written in the output.

`,
	Run: func(cmd *cobra.Command, args []string) {
		var comptree *tree.Tree
		var err error
		var specificTipNames []string

		var f *os.File
		if pruneouttree != "stdout" {
			f, err = os.Create(pruneouttree)
		} else {
			f = os.Stdout
		}
		if err != nil {
			io.ExitWithMessage(err)
		}

		// Read comp Tree : Only one tree in input
		comptree, err = utils.ReadRefTree(prunecomptree)
		if err != nil {
			io.ExitWithMessage(err)
		}

		intreesChan := make(chan tree.Trees, 15)
		/* Read ref tree(s) */
		go func() {
			if _, err = utils.ReadCompTrees(prunereftree, intreesChan); err != nil {
				io.ExitWithMessage(err)
			}
		}()

		// Read ref Trees
		for reftree := range intreesChan {
			specificTipNames = specificTips(reftree.Tree, comptree)
			err = reftree.Tree.RemoveTips(specificTipNames...)
			if err != nil {
				io.ExitWithMessage(err)
			}
			f.WriteString(reftree.Tree.Newick() + "\n")
		}
		f.Close()
	},
}

func init() {
	RootCmd.AddCommand(pruneCmd)
	pruneCmd.Flags().StringVarP(&prunereftree, "ref", "i", "stdin", "Input reference tree")
	pruneCmd.Flags().StringVarP(&prunecomptree, "comp", "c", "none", "Input compared tree ")
	pruneCmd.Flags().StringVarP(&pruneouttree, "output", "o", "stdout", "Output tree")
}
