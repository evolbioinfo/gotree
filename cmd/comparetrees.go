package cmd

import (
	"errors"
	"fmt"
	"github.com/fredericlemoine/gotree/io"
	"github.com/fredericlemoine/gotree/tree"
	"github.com/spf13/cobra"
	"os"
	"runtime"
	"sync"
)

// Type for channel of tree stats
type stats struct {
	id       int
	tree1    int
	tree2    int
	common   int
	sametree bool
}

var comparetreeidentical bool

func compare(tree1 string, tree2 string, tips bool, cpus int) {
	maxcpus := runtime.NumCPU()
	var edges []*tree.Edge

	statsChannel := make(chan stats, 15)

	if tree2 == "none" {
		io.ExitWithMessage(errors.New("You must provide a file containing compared trees"))
	}
	if cpus > maxcpus {
		cpus = maxcpus
	}

	fmt.Fprintf(os.Stderr, "Reference : %s\n", tree1)
	fmt.Fprintf(os.Stderr, "Compared  : %s\n", tree2)
	fmt.Fprintf(os.Stderr, "With tips : %t\n", tips)
	fmt.Fprintf(os.Stderr, "Threads   : %d\n", cpus)

	refTree := readTree(intreefile)
	compareChannel := readTrees(intree2file)

	edges = refTree.Edges()
	index := tree.NewEdgeIndex(int64(len(edges)*2), 0.75)
	total := 0
	for i, e := range edges {
		index.PutEdgeValue(e, i, e.Length())
		if tips || !e.Right().Tip() {
			total++
		}
	}
	var wg sync.WaitGroup
	for cpu := 0; cpu < cpus; cpu++ {
		wg.Add(1)
		go func(cpu int) {
			for treeV := range compareChannel {
				total2 := 0
				common := 0
				var err error

				edges2 := treeV.Tree.Edges()
				// Check wether the 2 trees have the same set of tip names
				if err = refTree.CompareTipIndexes(treeV.Tree); err != nil {
					io.ExitWithMessage(err)
				}

				sametree := true
				for _, e2 := range edges2 {
					ok := true
					if tips || !e2.Right().Tip() {
						total2++
					}
					if !e2.Right().Tip() {
						_, ok = index.Value(e2)
					}
					if !ok && comparetreeidentical {
						sametree = false
						break
					}
					if ok && (tips || !e2.Right().Tip()) {
						common++
					}
				}

				statsChannel <- stats{
					treeV.Id,
					total - common,
					total2 - common,
					common,
					sametree,
				}
			}
			wg.Done()
		}(cpu)
	}

	go func() {
		wg.Wait()
		close(statsChannel)
	}()

	if comparetreeidentical {
		fmt.Printf("tree\tidentical\n")
		for stats := range statsChannel {
			fmt.Printf("%d\t%v\n", stats.id, stats.sametree)
		}
	} else {
		fmt.Printf("tree\treference\tcommon\tcompared\n")
		for stats := range statsChannel {
			fmt.Printf("%d\t%d\t%d\t%d\n", stats.id, stats.tree1, stats.common, stats.tree2)
		}
	}
}

// compareCmd represents the compare command
var compareTreesCmd = &cobra.Command{
	Use:   "trees",
	Short: "Compare a reference tree with a set of trees",
	Long: `Compare a reference tree with a set of trees.

If --binary is given:
For each trees in the compared tree file, it will print tab separated values with:
1) The index of the compared tree in the file
2) "true" if the tree is identical, 
   "false" otherwise

Otherwise:
For each trees in the compared tree file, it will print tab separated values with:
1) The index of the compared tree in the file
2) The number of branches that are specific to the reference tree
3) The number of branches that are common to both trees
4) The number of branches that are specific to the compared tree

`,
	Run: func(cmd *cobra.Command, args []string) {
		compare(intreefile, intree2file, compareTips, rootCpus)
	},
}

func init() {
	compareCmd.AddCommand(compareTreesCmd)
	compareTreesCmd.Flags().BoolVarP(&compareTips, "tips", "l", false, "Include tips in the comparison")
	compareTreesCmd.Flags().BoolVar(&comparetreeidentical, "binary", false, "If true, then just print true (identical tree) or false (different tree) for each compared tree")
}
