package cmd

import (
	"github.com/fredericlemoine/gotree/io"
	"github.com/fredericlemoine/gotree/support"
	"github.com/spf13/cobra"
	"math/rand"
	"os"
	"time"
)

var parsimonyEmpirical bool
var parsimonySeed int64

// parsimonyCmd represents the parsimony command
var parsimonyCmd = &cobra.Command{
	Use:   "parsimony",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		var f *os.File
		var err error
		rand.Seed(parsimonySeed)

		if supportOutFile != "stdout" {
			f, err = os.Create(supportOutFile)
		} else {
			f = os.Stdout
		}
		if err != nil {
			io.ExitWithMessage(err)
		}
		t := support.Parsimony(supportIntree, supportBoottrees, parsimonyEmpirical, rootCpus)
		f.WriteString(t.Newick() + "\n")
		f.Close()
	},
}

func init() {
	supportCmd.AddCommand(parsimonyCmd)
	parsimonyCmd.PersistentFlags().BoolVarP(&parsimonyEmpirical, "empirical", "e", false, "If the support is computed with comparison to empirical support classical steps (shuffles of the original tree)")
	parsimonyCmd.PersistentFlags().Int64VarP(&parsimonySeed, "seed", "s", time.Now().UTC().UnixNano(), "Initial Random Seed if empirical is ON")

}
