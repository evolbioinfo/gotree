// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
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
	"log"

	"github.com/fredericlemoine/gotree/io"
	"github.com/fredericlemoine/gotree/tree"
	"github.com/spf13/cobra"
)

// subtreeCmd represents the subtree command
var subtreeCmd = &cobra.Command{
	Use:   "subtree",
	Short: "Select a subtree from the input tree whose root has the given name",
	Long: `Select a subtree from the input tree whose root has the given name.

The name may be a regexp, for example :
gotree subtree -i tree.nhx -n "^Mammal.*"

If several nodes match the given name/regexp, do nothing, and print the name of matching nodes.

`,
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		var nodes []*tree.Node

		f := openWriteFile(outtreefile)
		i := 0
		for reftree := range readTrees(intreefile) {
			nodes, err = reftree.Tree.SelectNodes(inputname)
			if err != nil {
				io.ExitWithMessage(err)
			}
			switch len(nodes) {
			case 1:
				subtree := reftree.Tree.SubTree(nodes[0])
				f.WriteString(subtree.Newick() + "\n")
			case 0:
				log.Print(fmt.Sprintf("Tree %d: No node matches input name", i))
			default:
				log.Print(fmt.Sprintf("Tree %d: Two many nodes match input name (%d)", i, len(nodes)))
				for _, n := range nodes {
					log.Print(fmt.Sprintf("Node: %s", n.Name()))
				}
			}
			i++
		}
		f.Close()
	},
}

func init() {
	RootCmd.AddCommand(subtreeCmd)
	subtreeCmd.PersistentFlags().StringVarP(&inputname, "name", "n", "none", "Name of the node to select as the root of the subtree (maybe a regex)")
	subtreeCmd.PersistentFlags().StringVarP(&intreefile, "input", "i", "stdin", "Input tree")
	subtreeCmd.PersistentFlags().StringVarP(&outtreefile, "output", "o", "stdout", "Output tree file")
}
