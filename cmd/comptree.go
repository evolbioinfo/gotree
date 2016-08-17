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
	"flag"
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
	tree2  int
}

func comptrees(tree1 string, tree2 string, tips bool, cpus int) {
	maxcpus := runtime.NumCPU()
	var refTree *tree.Tree
	var err error
	var nbtips int
	var edges []*tree.Edge

	compTreesChannel := make(chan utils.Trees, 100)
	statsChannel := make(chan stats, 100)

	flag.Parse()

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
	edges = refTree.Edges()
	if nbtips, err = refTree.NbTips(); err != nil {
		panic(err)
	}

	go func() {
		if err = utils.ReadCompTrees(tree2, compTreesChannel); err != nil {
			panic(err)
		}
	}()

	var wg sync.WaitGroup
	for cpu := 0; cpu < cpus; cpu++ {
		wg.Add(1)
		go func(cpu int) {
			for treeV := range compTreesChannel {
				var tree1, common, tree2 int
				var err error
				edges2 := treeV.Tree.Edges()

				if tree1, common, tree2, err = tree.CommonEdges(edges, edges2); err != nil {
					panic(err)
				}
				if !tips {
					common -= nbtips
				}
				statsChannel <- stats{
					treeV.Id,
					tree1,
					common,
					tree2,
				}
			}
			wg.Done()
		}(cpu)
	}

	go func() {
		wg.Wait()
		close(statsChannel)
	}()

	fmt.Printf("Tree\tspecref\tcommon\tspeccomp\n")
	for stats := range statsChannel {
		fmt.Printf("%d\t%d\t%d\t%d\n", stats.id, stats.tree1, stats.common, stats.tree2)
	}
}

var comptreeCpus int
var comptreeTree1 string
var comptreeTree2 string
var comptreeTips bool

// comptreeCmd represents the comptree command
var comptreeCmd = &cobra.Command{
	Use:   "comptree",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Work your own magic here
		comptrees(comptreeTree1, comptreeTree2, comptreeTips, comptreeCpus)
	},
}

func init() {
	RootCmd.AddCommand(comptreeCmd)
	maxcpus := runtime.NumCPU()

	comptreeCmd.Flags().StringVarP(&comptreeTree1, "reftree", "i", "stdin", "Reference tree input file")
	comptreeCmd.Flags().StringVarP(&comptreeTree2, "comptrees", "c", "none", "Compared trees input file")
	comptreeCmd.Flags().BoolVarP(&comptreeTips, "tips", "l", false, "Compared trees input file")
	comptreeCmd.Flags().IntVarP(&comptreeCpus, "threads", "t", 1, "Number of threads (Max=)"+strconv.Itoa(maxcpus))

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// comptreeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// comptreeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
