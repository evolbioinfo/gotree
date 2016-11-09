package cmd

import (
	"github.com/fredericlemoine/gotree/io"
	"github.com/fredericlemoine/gotree/support"
	"github.com/spf13/cobra"
	"math/rand"
	"os"
	"time"
)

var mastEmpirical bool
var mastSeed int64

// mastlikeCmd represents the mastlike command
var mastlikeCmd = &cobra.Command{
	Use:   "mastlike",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples

`,
	Run: func(cmd *cobra.Command, args []string) {
		var f *os.File
		var err error
		rand.Seed(mastSeed)

		if supportOutFile != "stdout" {
			f, err = os.Create(supportOutFile)
		} else {
			f = os.Stdout
		}
		if err != nil {
			io.ExitWithMessage(err)
		}
		t := support.MastLike(supportIntree, supportBoottrees, mastEmpirical, rootCpus)
		f.WriteString(t.Newick() + "\n")
		f.Close()
	},
}

func init() {
	supportCmd.AddCommand(mastlikeCmd)

	mastlikeCmd.PersistentFlags().BoolVarP(&mastEmpirical, "empirical", "e", false, "If the support is computed with comparison to empirical support classical steps (shuffles of the original tree)")
	mastlikeCmd.PersistentFlags().Int64VarP(&mastSeed, "seed", "s", time.Now().UTC().UnixNano(), "Initial Random Seed if empirical is ON")

}
