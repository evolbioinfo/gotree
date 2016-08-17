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
		var f *os.File
		if statsoutfile != "stdout" {
			f, err = os.Create(statsoutfile)
		} else {
			f = os.Stdout
		}
		if err != nil {
			panic(err)
		}
		f.WriteString("id\tlength\tsupport\tterminal\trightname\n")
		for i, e := range tree.Edges() {
			var length = "N/A"
			if e.Length() != -1 {
				length = fmt.Sprintf("%f", e.Length())
			}
			var support = "N/A"
			if e.Support() != -1 {
				support = fmt.Sprintf("%f", e.Support())
			}

			f.WriteString(fmt.Sprintf("%d\t%s\t%s\t%t\t%s\n", i, length, support, e.Right().Tip(), e.Right().Name()))
		}
		f.Close()
	},
}

func init() {
	statsCmd.AddCommand(edgesCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// edgesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// edgesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
