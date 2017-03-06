package cmd

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/fredericlemoine/gotree/io"
	"github.com/fredericlemoine/gotree/support"
	"github.com/spf13/cobra"
)

var boosterSeed int64
var boosterPrintMovedTaxa bool

// boosterCmd represents the booster command
var boosterCmd = &cobra.Command{
	Use:   "booster",
	Short: "Compute BOOSTER supports",
	Long: `Compute BOOtstrap Support by TransfER
`,
	Run: func(cmd *cobra.Command, args []string) {
		writeLogBooster()
		rand.Seed(boosterSeed)
		t, err := support.Booster(supportIntree, supportBoottrees, supportLog, supportSilent, boosterPrintMovedTaxa, rootCpus)
		if err != nil {
			io.ExitWithMessage(err)
		}
		supportOut.WriteString(t.Newick() + "\n")
		supportLog.WriteString(fmt.Sprintf("End         : %s\n", time.Now().Format(time.RFC822)))
	},
}

func init() {
	supportCmd.AddCommand(boosterCmd)
	boosterCmd.PersistentFlags().BoolVar(&boosterPrintMovedTaxa, "moved-taxa", false, "If true, will print in log file (-l) taxa that move the most around branches")
	boosterCmd.PersistentFlags().Int64VarP(&boosterSeed, "seed", "s", time.Now().UTC().UnixNano(), "Initial Random Seed if empirical is ON")
}

func writeLogBooster() {
	supportLog.WriteString("BOOSTER Support\n")
	supportLog.WriteString(fmt.Sprintf("Date        : %s\n", time.Now().Format(time.RFC822)))
	supportLog.WriteString(fmt.Sprintf("Input tree  : %s\n", supportIntree))
	supportLog.WriteString(fmt.Sprintf("Boot trees  : %s\n", supportBoottrees))
	supportLog.WriteString(fmt.Sprintf("Output tree : %s\n", supportOutFile))
	supportLog.WriteString(fmt.Sprintf("Seed        : %d\n", boosterSeed))
	supportLog.WriteString(fmt.Sprintf("CPUs        : %d\n", rootCpus))
}
