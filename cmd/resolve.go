package cmd

import (
	"math/rand"
	"os"
	"time"

	"github.com/fredericlemoine/gotree/io"
	"github.com/fredericlemoine/gotree/io/utils"
	"github.com/fredericlemoine/gotree/tree"
	"github.com/spf13/cobra"
)

var resolveInputTree string
var resolveOutputTree string
var resolveIntrees chan tree.Trees
var resolveOutTrees *os.File
var resolveSeed int64

// resolveCmd represents the resolve command
var resolveCmd = &cobra.Command{
	Use:   "resolve",
	Short: "Resolve multifurcations by adding 0 length branches",
	Long: `Resolve multifurcations by adding 0 length branches.

* If any node has more than 3 neighbors :
   Resolve neighbors randomly by adding 0 length 
   branches until it has 3 neighbors
`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		var err error
		var nbtrees int = 0
		resolveIntrees = make(chan tree.Trees, 15)
		/* Read ref tree(s) */
		go func() {
			if nbtrees, err = utils.ReadCompTrees(resolveInputTree, resolveIntrees); err != nil {
				io.ExitWithMessage(err)
			}
		}()
		resolveOutTrees = openWriteFile(resolveOutputTree)
		rand.Seed(resolveSeed)
	},
	Run: func(cmd *cobra.Command, args []string) {
		for t := range resolveIntrees {
			t.Tree.Resolve()
			resolveOutTrees.WriteString(t.Tree.Newick() + "\n")
		}

	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		resolveOutTrees.Close()
	},
}

func init() {
	RootCmd.AddCommand(resolveCmd)
	resolveCmd.PersistentFlags().StringVarP(&resolveInputTree, "input", "i", "stdin", "Input tree(s) file")
	resolveCmd.PersistentFlags().StringVarP(&resolveOutputTree, "output", "o", "stdout", "Resolved tree(s) output file")
	resolveCmd.PersistentFlags().Int64VarP(&resolveSeed, "seed", "s", time.Now().UTC().UnixNano(), "Initial Random Seed")
}
