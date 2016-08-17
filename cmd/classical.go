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
	"github.com/fredericlemoine/gotree/io/utils"
	"github.com/fredericlemoine/gotree/tree"
	"github.com/spf13/cobra"
	"os"
	"runtime"
	"sync"
)

func classical(reftreefile, boottreefile string, outfile string, cpus int) {
	var reftree *tree.Tree
	var err error
	var f *os.File
	if outfile != "stdout" {
		f, err = os.Create(outfile)
	} else {
		f = os.Stdout
	}
	if err != nil {
		panic(err)
	}

	maxcpus := runtime.NumCPU()
	if cpus > maxcpus {
		cpus = maxcpus
	}

	if reftree, err = utils.ReadRefTree(reftreefile); err != nil {
		panic(err)
	}

	if boottreefile == "none" {
		panic("You must provide a file containing bootstrap trees")
	}

	if cpus > maxcpus {
		cpus = maxcpus
	}

	var nbtrees int
	compareChannel := make(chan utils.Trees, 100)
	go func() {
		if nbtrees, err = utils.ReadCompTrees(boottreefile, compareChannel); err != nil {
			panic(err)
		}
	}()

	edges := reftree.Edges()
	foundEdges := make(chan int, 1000)
	foundBoot := make([]int, len(edges))
	var wg sync.WaitGroup
	for cpu := 0; cpu < cpus; cpu++ {
		wg.Add(1)
		go func(cpu int) {
			for treeV := range compareChannel {
				edges2 := treeV.Tree.Edges()
				for i, e := range edges {
					if !e.Right().Tip() {
						for _, e2 := range edges2 {
							if !e2.Right().Tip() && e.SameBipartition(e2) {
								foundEdges <- i
								break
							}
						}
					}
				}

			}
			wg.Done()
		}(cpu)
	}

	go func() {
		wg.Wait()
		close(foundEdges)
	}()

	for edge_i := range foundEdges {
		foundBoot[edge_i]++
	}

	for i, count := range foundBoot {
		if !edges[i].Right().Tip() {
			edges[i].SetSupport(float64(count) / float64(nbtrees))
		}
	}

	f.WriteString(reftree.Newick() + "\n")
	f.Close()

}

var classicalEmpirical bool

// classicalCmd represents the classical command
var classicalCmd = &cobra.Command{
	Use:   "classical",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		classical(supportIntree, supportBoottrees, supportOutFile, supportCpus)
	},
}

func init() {
	supportCmd.AddCommand(classicalCmd)

	computeCmd.PersistentFlags().BoolVarP(&classicalEmpirical, "empirical", "e", false, "If the support is computed with comparison to empirical support classical steps (shuffles of the original tree)")
}
