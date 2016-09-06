package cmd

import (
	"errors"
	"fmt"
	"github.com/fredericlemoine/gotree/io"
	"github.com/fredericlemoine/gotree/io/utils"
	"github.com/fredericlemoine/gotree/tree"
	"github.com/spf13/cobra"
	"os"
	"runtime"
	"sync"
)

// Type for channel of tree stats
type stats struct {
	id     int
	tree1  int
	common int
}

func compare(tree1 string, tree2 string, tips bool, cpus int) {
	maxcpus := runtime.NumCPU()
	var refTree *tree.Tree
	var err error
	var edges []*tree.Edge

	compareChannel := make(chan utils.Trees, 15)
	statsChannel := make(chan stats, 15)

	if tree2 == "none" {
		io.ExitWithMessage(errors.New("You must provide a file containing compared trees"))
	}
	if cpus > maxcpus {
		cpus = maxcpus
	}

	fmt.Fprintf(os.Stderr, "Reference : %s\n", tree1)
	fmt.Fprintf(os.Stderr, "Compared  : %s\n", tree2)
	fmt.Fprintf(os.Stderr, "With tips : %b\n", tips)
	fmt.Fprintf(os.Stderr, "Threads   : %d\n", cpus)

	if refTree, err = utils.ReadRefTree(tree1); err != nil {
		io.ExitWithMessage(err)
	}

	var nbtrees int
	go func() {
		if nbtrees, err = utils.ReadCompTrees(tree2, compareChannel); err != nil {
			io.ExitWithMessage(err)
		}
	}()

	edges = refTree.Edges()
	var wg sync.WaitGroup
	for cpu := 0; cpu < cpus; cpu++ {
		wg.Add(1)
		go func(cpu int) {
			for treeV := range compareChannel {
				var tree1, common int
				var err error

				edges2 := treeV.Tree.Edges()

				// Check wether the 2 trees have the same set of tip names
				if err = refTree.CompareTipIndexes(treeV.Tree); err != nil {
					io.ExitWithMessage(err)
				}

				// Then compare edges
				if tree1, common, err = tree.CommonEdges(edges, edges2, tips); err != nil {
					io.ExitWithMessage(err)
				}

				statsChannel <- stats{
					treeV.Id,
					tree1,
					common,
				}
			}
			wg.Done()
		}(cpu)
	}

	go func() {
		wg.Wait()
		close(statsChannel)
	}()

	fmt.Printf("Tree\tspecref\tcommon\n")
	for stats := range statsChannel {
		fmt.Printf("%d\t%d\t%d\n", stats.id, stats.tree1, stats.common)
	}
}

var compareTree1 string
var compareTree2 string
var compareTips bool

// compareCmd represents the compare command
var compareCmd = &cobra.Command{
	Use:   "compare",
	Short: "Compares a reference tree with a set of trees",
	Long: `Compares a reference tree to a set of trees.

for each trees in the compared tree file, it prints the number of common edges
between it and the reference tree, as well as the number of specific edges.
`,
	Run: func(cmd *cobra.Command, args []string) {
		compare(compareTree1, compareTree2, compareTips, rootCpus)
	},
}

func init() {
	RootCmd.AddCommand(compareCmd)
	compareCmd.Flags().StringVarP(&compareTree1, "reftree", "i", "stdin", "Reference tree input file")
	compareCmd.Flags().StringVarP(&compareTree2, "compared", "c", "none", "Compared trees input file")
	compareCmd.Flags().BoolVarP(&compareTips, "tips", "l", false, "Compared trees input file")
}
