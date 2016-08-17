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

// tipsCmd represents the tips command
var tipsCmd = &cobra.Command{
	Use:   "tips",
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
		var f *os.File
		if statsoutfile != "stdout" {
			f, err = os.Create(statsoutfile)
		} else {
			f = os.Stdout
		}
		if err != nil {
			panic(err)
		}
		f.WriteString("id\tnneigh\tname\n")
		for i, n := range tree.Nodes() {
			if n.Nneigh() == 1 {
				f.WriteString(fmt.Sprintf("%d\t%d\t%s\n", i, n.Nneigh(), n.Name()))
			}
		}
		f.Close()

	},
}

func init() {
	statsCmd.AddCommand(tipsCmd)
}
