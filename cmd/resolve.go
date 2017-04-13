package cmd

import (
	"math/rand"
	"time"

	"github.com/fredericlemoine/gotree/io"
	"github.com/spf13/cobra"
)

// resolveCmd represents the resolve command
var resolveCmd = &cobra.Command{
	Use:   "resolve",
	Short: "Resolve multifurcations by adding 0 length branches",
	Long: `Resolve multifurcations by adding 0 length branches.

* If any node has more than 3 neighbors :
   Resolve neighbors randomly by adding 0 length 
   branches until it has 3 neighbors
`,
	Run: func(cmd *cobra.Command, args []string) {
		rand.Seed(seed)
		f := openWriteFile(outtreefile)
		treefile, trees := readTrees(intreefile)
		defer treefile.Close()

		for tr := range trees {
			if tr.Err != nil {
				io.ExitWithMessage(tr.Err)
			}
			tr.Tree.Resolve()
			f.WriteString(tr.Tree.Newick() + "\n")
		}
		f.Close()
	},
}

func init() {
	RootCmd.AddCommand(resolveCmd)
	resolveCmd.PersistentFlags().StringVarP(&intreefile, "input", "i", "stdin", "Input tree(s) file")
	resolveCmd.PersistentFlags().StringVarP(&outtreefile, "output", "o", "stdout", "Resolved tree(s) output file")
	resolveCmd.PersistentFlags().Int64VarP(&seed, "seed", "s", time.Now().UTC().UnixNano(), "Initial Random Seed")
}
