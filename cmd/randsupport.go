package cmd

import (
	"github.com/fredericlemoine/gostats"
	"github.com/fredericlemoine/gotree/io"
	"github.com/fredericlemoine/gotree/io/utils"
	"github.com/fredericlemoine/gotree/tree"
	"github.com/spf13/cobra"
	"math/rand"
	"os"
	"time"
)

var setsupportinput string
var setsupportintrees chan tree.Trees
var setsupportout string
var setsupportseed int64
var setsupportmean float64
var setsupportoutfile *os.File

// randsupportCmd represents the randbrlen command
var randsupportCmd = &cobra.Command{
	Use:   "randsupport",
	Short: "Assign a random support to edges of input trees",
	Long: `Assign a random support to edges of input trees.

Support follows a uniform distribution in [0,1].

`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		var err error
		rand.Seed(setsupportseed)
		setsupportintrees = make(chan tree.Trees, 15)
		/* Read ref tree(s) */
		go func() {
			if _, err = utils.ReadCompTrees(setsupportinput, setsupportintrees); err != nil {
				io.ExitWithMessage(err)
			}
		}()
		setsupportoutfile = openWriteFile(setsupportout)
	},
	Run: func(cmd *cobra.Command, args []string) {
		for tr := range setsupportintrees {
			for _, e := range tr.Tree.Edges() {
				if !e.Right().Tip() {
					e.SetSupport(gostats.Float64Range(0, 1))
				}
			}
			setsupportoutfile.WriteString(tr.Tree.Newick() + "\n")
		}
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		setsupportoutfile.Close()
	},
}

func init() {
	RootCmd.AddCommand(randsupportCmd)

	randsupportCmd.PersistentFlags().StringVarP(&setsupportinput, "input", "i", "stdin", "Input tree")
	randsupportCmd.PersistentFlags().StringVarP(&setsupportout, "output", "o", "stdout", "Output file")
	randsupportCmd.Flags().Int64VarP(&setsupportseed, "seed", "s", time.Now().UTC().UnixNano(), "Initial Random Seed")
}
