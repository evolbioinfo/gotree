package cmd

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/fredericlemoine/gotree/io"
	"github.com/fredericlemoine/gotree/support"
	"github.com/fredericlemoine/gotree/tree"
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

		// Compute average supports (non normalized, e.g normalizedByExpected=false)
		err := support.Booster(refTree, boottreechan, supportLog, supportSilent, movedtaxa, taxperbranches, hightaxperbranches, cutoff, false, rootCpus)
		if err != nil {
			io.ExitWithMessage(err)
		}
		// If rawSupportOutputFile is set, then we print the raw support tree first
		if rawSupportOutputFile != "none" {
			reformated := refTree.Clone()
			reformatAvgDistance(reformated)
			rawSupportOut.WriteString(reformated.Newick() + "\n")
		}
		// We normalize the supports and print them
		normalizeTransferDistancesByDepth(refTree)
		supportOut.WriteString(refTree.Newick() + "\n")
		supportLog.WriteString(fmt.Sprintf("End         : %s\n", time.Now().Format(time.RFC822)))
	},
}

func init() {
	supportCmd.AddCommand(boosterCmd)
	boosterCmd.PersistentFlags().BoolVar(&movedtaxa, "moved-taxa", false, "If true, will print in log file (-l) taxa that move the most around branches")
	boosterCmd.PersistentFlags().BoolVar(&taxperbranches, "per-branches", false, "If true, will print in log file (-l) average taxa transfers for all taxa per banches of the reference tree")
	boosterCmd.PersistentFlags().BoolVar(&hightaxperbranches, "highest-per-branches", false, "If true, will print in log file (-l) average taxa transfers for highly transfered taxa per banches of the reference tree (i.e. the x most transfered, with x~ average distance)")
	boosterCmd.PersistentFlags().StringVarP(&rawSupportOutputFile, "out-raw", "r", "none", "If given, then prints the same tree with non normalized supports (average transfer distance) as branch names, in the form branch_id|avg_distance|branch_depth")
	boosterCmd.PersistentFlags().Float64Var(&cutoff, "dist-cutoff", 0.3, "If --moved-taxa, then this is the distance cutoff to consider a branch for moving taxa computation. It is the normalized distance to the current bootstrap tree (e.g. 0.05). Must be between 0 and 1, otherwise set to 0")
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

/*
This finction writes on the child node name the string: "branch_id|avg_dist|depth"
and removes support information from each branch
*/
func reformatAvgDistance(t *tree.Tree) {
	for i, e := range t.Edges() {
		if e.Support() != tree.NIL_SUPPORT {
			td, _ := e.TopoDepth()
			e.Right().SetName(fmt.Sprintf("%d|%s|%d", i, e.SupportString(), td))
			e.SetSupport(tree.NIL_SUPPORT)
		}
	}
}

/*
This function takes all branch support values (that are considered as average
transfer distances over bootstrap trees), normalizes them by the depth and
convert them to similarity, i.e:
    1-avg_dist/(depth-1)
*/
func normalizeTransferDistancesByDepth(t *tree.Tree) {
	for _, e := range t.Edges() {
		avgdist := e.Support()
		if avgdist != tree.NIL_SUPPORT {
			td, _ := e.TopoDepth()
			e.SetSupport(float64(1) - avgdist/float64(td-1))
		}
	}
}
