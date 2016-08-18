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
	"github.com/fredericlemoine/gotree/support"
	"github.com/spf13/cobra"
	"os"
)

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
		var f *os.File
		var err error
		if supportOutFile != "stdout" {
			f, err = os.Create(supportOutFile)
		} else {
			f = os.Stdout
		}
		if err != nil {
			panic(err)
		}
		t := support.Classical(supportIntree, supportBoottrees, supportCpus)
		f.WriteString(t.Newick() + "\n")
		f.Close()
	},
}

func init() {
	supportCmd.AddCommand(classicalCmd)

	computeCmd.PersistentFlags().BoolVarP(&classicalEmpirical, "empirical", "e", false, "If the support is computed with comparison to empirical support classical steps (shuffles of the original tree)")
}
