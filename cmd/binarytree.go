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
	"github.com/fredericlemoine/gotree/tree"
	"github.com/spf13/cobra"
	"math/rand"
	"os"
	"time"
)

var binarytreeNbTips int
var binarytreeNbTrees int
var binarytreeOutputfile string
var binarytreeSeed int64

func binarytree(nbtrees int, nbtips int, output string, binarytreeSeed int64) error {
	var f *os.File
	var err error
	var t *tree.Tree

	rand.Seed(binarytreeSeed)

	if output != "stdout" {
		f, err = os.Create(output)
	} else {
		f = os.Stdout
	}
	if err != nil {
		return err
	}

	for i := 0; i < nbtrees; i++ {
		t, err = tree.RandomBinaryTree(nbtips)
		if err != nil {
			return err
		}
		f.WriteString(t.Newick() + "\n")
	}
	f.Close()
	return nil
}

// binarytreeCmd represents the binarytree command
var binarytreeCmd = &cobra.Command{
	Use:   "binarytree",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := binarytree(binarytreeNbTrees, binarytreeNbTips, binarytreeOutputfile, binarytreeSeed); err != nil {
			io.ExitWithMessage(err)
		}
	},
}

func init() {
	generateCmd.AddCommand(binarytreeCmd)
	binarytreeCmd.Flags().IntVarP(&binarytreeNbTips, "nbtips", "t", 10, "Number of tips of the tree to generate")
	binarytreeCmd.Flags().IntVarP(&binarytreeNbTrees, "nbtrees", "n", 1, "Number of trees to generate")
	binarytreeCmd.Flags().Int64VarP(&binarytreeSeed, "seed", "s", time.Now().UTC().UnixNano(), "Initial Random Seed")
	binarytreeCmd.Flags().StringVarP(&binarytreeOutputfile, "output", "o", "stdout", "Number of tips of the tree to generate")
}
