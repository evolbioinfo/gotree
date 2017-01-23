package cmd

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/fredericlemoine/gotree/support"
	"github.com/spf13/cobra"
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
		writeLogMast()
		rand.Seed(mastSeed)
		t := support.MastLike(supportIntree, supportBoottrees, supportLog, mastEmpirical, rootCpus)
		supportOut.WriteString(t.Newick() + "\n")
		supportLog.WriteString(fmt.Sprintf("End         : %s\n", time.Now().Format(time.RFC822)))
	},
}

func init() {
	supportCmd.AddCommand(mastlikeCmd)

	mastlikeCmd.PersistentFlags().BoolVarP(&mastEmpirical, "empirical", "e", false, "If the support is computed with comparison to empirical support classical steps (shuffles of the original tree)")
	mastlikeCmd.PersistentFlags().Int64VarP(&mastSeed, "seed", "s", time.Now().UTC().UnixNano(), "Initial Random Seed if empirical is ON")

}

func writeLogMast() {
	supportLog.WriteString("Mast Support\n")
	supportLog.WriteString(fmt.Sprintf("Date        : %s\n", time.Now().Format(time.RFC822)))
	supportLog.WriteString(fmt.Sprintf("Input tree  : %s\n", supportIntree))
	supportLog.WriteString(fmt.Sprintf("Boot trees  : %s\n", supportBoottrees))
	supportLog.WriteString(fmt.Sprintf("Output tree : %s\n", supportOutFile))
	supportLog.WriteString(fmt.Sprintf("Theor norm  : %t\n", !mastEmpirical))
	supportLog.WriteString(fmt.Sprintf("Seed        : %d\n", mastSeed))
	supportLog.WriteString(fmt.Sprintf("CPUs        : %d\n", rootCpus))
}
