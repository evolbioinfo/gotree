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
	Short: "Removes tips of the input tree that are not in the compared tree",
	Long: `This tool removes tips of the input reference tree that 
are not present in the compared tree.

In output, we have a tree containing only tips that are common to both trees.
`,
	Run: func(cmd *cobra.Command, args []string) {
		// Read ref Tree
		var reftree, comptree *tree.Tree
		var err error
		var specificTipNames []string
		reftree, err = utils.ReadRefTree(prunereftree)
		if err != nil {
			io.ExitWithMessage(err)
		}
		// Read comp Tree
		comptree, err = utils.ReadRefTree(prunecomptree)
		if err != nil {
			io.ExitWithMessage(err)
		}

		specificTipNames = specificTips(reftree, comptree)
		err = reftree.RemoveTips(specificTipNames...)
		if err != nil {
			io.ExitWithMessage(err)
		}
		var f *os.File
		if pruneouttree != "stdout" {
			f, err = os.Create(pruneouttree)
		} else {
			f = os.Stdout
		}
		if err != nil {
			io.ExitWithMessage(err)
		}

		f.WriteString(reftree.Newick() + "\n")
		f.Close()
	},
}

func init() {
	RootCmd.AddCommand(pruneCmd)
	pruneCmd.Flags().StringVarP(&prunereftree, "ref", "i", "stdin", "Input reference tree")
	pruneCmd.Flags().StringVarP(&prunecomptree, "comp", "c", "none", "Input compared tree ")
	pruneCmd.Flags().StringVarP(&pruneouttree, "output", "o", "stdout", "Output tree")
}
