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
	"github.com/spf13/cobra"
	"os"
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// edgesCmd represents the edges command
var edgesCmd = &cobra.Command{
	Use:   "edges",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		tree, err := utils.ReadRefTree(statsintree)
		if err != nil {
			panic(err)
		}
		tree.ComputeDepths()
		var f *os.File
		if statsoutfile != "stdout" {
			f, err = os.Create(statsoutfile)
		} else {
			f = os.Stdout
		}
		if err != nil {
			panic(err)
		}
		f.WriteString("id\tlength\tsupport\tterminal\tdepth\ttopodepth\trightname\n")
		for i, e := range tree.Edges() {
			var length = "N/A"
			if e.Length() != -1 {
				length = fmt.Sprintf("%f", e.Length())
			}
			var support = "N/A"
			if e.Support() != -1 {
				support = fmt.Sprintf("%f", e.Support())
			}
			var depth, leftdepth, rightdepth int

			if leftdepth, err = e.Left().Depth(); err != nil {
				panic(err)
			}
			if rightdepth, err = e.Right().Depth(); err != nil {
				panic(err)
			}
			depth = min(leftdepth, rightdepth)
			var topodepth int
			topodepth, err = e.TopoDepth()
			if err != nil {
				panic(err)
			}

			f.WriteString(fmt.Sprintf("%d\t%s\t%s\t%t\t%d\t%d\t%s\n", i, length, support, e.Right().Tip(), depth, topodepth, e.Right().Name()))
		}
		f.Close()
	},
}

func init() {
	statsCmd.AddCommand(edgesCmd)
}
