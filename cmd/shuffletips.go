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
	"math/rand"
	"os"
	"time"
)

var shuffletipsSeed int64

// shuffletipsCmd represents the shuffletips command
var shuffletipsCmd = &cobra.Command{
	Use:   "shuffletips",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Read Tree

		rand.Seed(shuffletipsSeed)

		var t *tree.Tree
		var err error
		t, err = utils.ReadRefTree(transformInputTree)
		if err != nil {
			panic(err)
		}
		var f *os.File
		if renameouttree != "stdout" {
			f, err = os.Create(transformOutputTree)
		} else {
			f = os.Stdout
		}
		if err != nil {
			panic(err)
		}

		t.ShuffleTips()

		if err != nil {
			panic(err)
		}

		f.WriteString(t.Newick() + "\n")
		f.Close()
	},
}

func init() {
	transformCmd.AddCommand(shuffletipsCmd)
	shuffletipsCmd.Flags().Int64VarP(&shuffletipsSeed, "seed", "s", time.Now().UTC().UnixNano(), "Initial Random Seed")
}
