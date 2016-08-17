// Copyright © 2016 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"github.com/fredericlemoine/gotree/io/utils"
	"github.com/fredericlemoine/gotree/tree"
	"github.com/spf13/cobra"
	"os"
	"runtime"
	"strconv"
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

	compareChannel := make(chan utils.Trees, 100)
	statsChannel := make(chan stats, 100)

	if tree2 == "none" {
		panic("You must provide a file containing compared trees")
	}
	if cpus > maxcpus {
		cpus = maxcpus
	}

	fmt.Fprintf(os.Stderr, "Reference : %s\n", tree1)
	fmt.Fprintf(os.Stderr, "Compared  : %s\n", tree2)
	fmt.Fprintf(os.Stderr, "With tips : %b\n", tips)
	fmt.Fprintf(os.Stderr, "Threads   : %d\n", cpus)

	if refTree, err = utils.ReadRefTree(tree1); err != nil {
		panic(err)
	}

	var nbtrees int
	go func() {
		if nbtrees, err = utils.ReadCompTrees(tree2, compareChannel); err != nil {
			panic(err)
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
					panic(err)
				}

				// Then compare edges
				if tree1, common, err = tree.CommonEdges(edges, edges2, tips); err != nil {
					panic(err)
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

var compareCpus int
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
		compare(compareTree1, compareTree2, compareTips, compareCpus)
	},
}

func init() {
	RootCmd.AddCommand(compareCmd)
	maxcpus := runtime.NumCPU()

	compareCmd.Flags().StringVarP(&compareTree1, "reftree", "i", "stdin", "Reference tree input file")
	compareCmd.Flags().StringVarP(&compareTree2, "compared", "c", "none", "Compared trees input file")
	compareCmd.Flags().BoolVarP(&compareTips, "tips", "l", false, "Compared trees input file")
	compareCmd.Flags().IntVarP(&compareCpus, "threads", "t", 1, "Number of threads (Max=)"+strconv.Itoa(maxcpus))

}
