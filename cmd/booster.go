package cmd

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/fredericlemoine/gotree/io"
	"github.com/fredericlemoine/gotree/support"
	"github.com/spf13/cobra"
)

// boosterCmd represents the booster command
var boosterCmd = &cobra.Command{
	Use:   "booster",
	Short: "Compute BOOSTER supports",
	Long: `Compute BOOtstrap Support by TransfER
`,
	Run: func(cmd *cobra.Command, args []string) {
		writeLogBooster()
		rand.Seed(seed)
		refTree := readTree(supportIntree)
		boottreefile, boottreechan := readTrees(supportBoottrees)
		defer boottreefile.Close()

		err := support.Booster(refTree, boottreechan, supportLog, supportSilent, movedtaxa, cutoff, rootCpus)
		if err != nil {
			io.ExitWithMessage(err)
		}
		supportOut.WriteString(refTree.Newick() + "\n")
		supportLog.WriteString(fmt.Sprintf("End         : %s\n", time.Now().Format(time.RFC822)))
	},
}

func init() {
	supportCmd.AddCommand(boosterCmd)
	boosterCmd.PersistentFlags().BoolVar(&movedtaxa, "moved-taxa", false, "If true, will print in log file (-l) taxa that move the most around branches")
	boosterCmd.PersistentFlags().Float64Var(&cutoff, "dist-cutoff", 0.05, "If --moved-taxa, then this is the distance cutoff to consider a branch for moving taxa computation. It is the normalized distance to the current bootstrap tree (e.g. 0.05). Must be between 0 and 1, otherwise set to 0")
	boosterCmd.PersistentFlags().Int64VarP(&seed, "seed", "s", time.Now().UTC().UnixNano(), "Initial Random Seed if empirical is ON")
}

func writeLogBooster() {
	supportLog.WriteString("BOOSTER Support\n")
	supportLog.WriteString(fmt.Sprintf("Date        : %s\n", time.Now().Format(time.RFC822)))
	supportLog.WriteString(fmt.Sprintf("Input tree  : %s\n", supportIntree))
	supportLog.WriteString(fmt.Sprintf("Boot trees  : %s\n", supportBoottrees))
	supportLog.WriteString(fmt.Sprintf("Output tree : %s\n", supportOutFile))
	supportLog.WriteString(fmt.Sprintf("Seed        : %d\n", seed))
	supportLog.WriteString(fmt.Sprintf("CPUs        : %d\n", rootCpus))
}
