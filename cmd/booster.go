package cmd

import (
	"fmt"
	goio "io"
	"time"

	"github.com/evolbioinfo/gotree/io"
	"github.com/evolbioinfo/gotree/support"
	"github.com/evolbioinfo/gotree/tree"
	"github.com/spf13/cobra"
)

var boosterdistcutoff float64

// boosterCmd represents the booster command
// Just to keep the alias
var boosterCmd = &cobra.Command{
	Hidden: true,
	Use:    "booster",
	Short:  "Compute TBE supports",
	Long: `Compute BOOtstrap Support by TransfER

	For more information, See:
	Lemoine, F. and Domelevo Entfellner, J.-B. and Wilkinson, E. and Correia, D. and Dávila Felipe, M. and De Oliveira, T. and Gascuel, O.
	Renewing Felsenstein’s phylogenetic bootstrap in the era of big data. Nature, 556:452–456
`,
	RunE: booster,
}

// boosterCmd represents the booster command
var tbeCmd = &cobra.Command{
	Use:   "tbe",
	Short: "Compute TBE supports",
	Long: `Compute BOOtstrap Support by TransfER

	For more information, See:
	Lemoine, F. and Domelevo Entfellner, J.-B. and Wilkinson, E. and Correia, D. and Dávila Felipe, M. and De Oliveira, T. and Gascuel, O.
	Renewing Felsenstein’s phylogenetic bootstrap in the era of big data. Nature, 556:452–456
`,
	RunE: booster,
}

func booster(cmd *cobra.Command, args []string) (err error) {
	var refTree *tree.Tree
	var rawtree *tree.Tree
	var boottreefile goio.Closer
	var boottreechan <-chan tree.Trees
	//var f *os.File

	//f, err = os.Create("cpuprof")
	//if err != nil {
	//	log.Fatal("could not create CPU profile: ", err)
	//}
	//defer f.Close()
	//if err := pprof.StartCPUProfile(f); err != nil {
	//	log.Fatal("could not start CPU profile: ", err)
	//}
	//defer pprof.StopCPUProfile()

	writeLogBooster()
	if refTree, err = readTree(supportIntree); err != nil {
		io.LogError(err)
		return
	}
	if boottreefile, boottreechan, err = readTrees(supportBoottrees); err != nil {
		io.LogError(err)
		return
	}
	defer boottreefile.Close()

	if err = refTree.ReinitIndexes(); err != nil {
		io.LogError(err)
		return
	}

	// Compute average supports (non normalized, e.g normalizedByExpected=false)
	if rawtree, err = support.TBE(refTree, boottreechan, rootCpus, rawSupportOutputFile != "none", movedtaxa, taxperbranches, boosterdistcutoff, supportLog, nil); err != nil {
		io.LogError(err)
		return
	}
	// If rawSupportOutputFile is set, then we print the raw support tree first
	if rawSupportOutputFile != "none" {
		rawSupportOut.WriteString(rawtree.Newick() + "\n")
	}
	supportOut.WriteString(refTree.Newick() + "\n")
	supportLog.WriteString(fmt.Sprintf("End         : %s\n", time.Now().Format(time.RFC822)))

	return
}

func addTBEFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().BoolVar(&movedtaxa, "moved-taxa", false, "If true, will print in log file (-l) taxa that move the most around branches")
	cmd.PersistentFlags().BoolVar(&taxperbranches, "per-branches", false, "If true, will print in log file (-l) average taxa transfers for all taxa per banches of the reference tree")
	//boosterCmd.PersistentFlags().BoolVar(&hightaxperbranches, "highest-per-branches", false, "If true, will print in log file (-l) average taxa transfers for highly transfered taxa per banches of the reference tree (i.e. the x most transfered, with x~ average distance)")
	cmd.PersistentFlags().StringVarP(&rawSupportOutputFile, "out-raw", "r", "none", "If given, then prints the same tree with non normalized supports (average transfer distance) as branch names, in the form branch_id|avg_distance|branch_depth")
	cmd.PersistentFlags().Float64Var(&boosterdistcutoff, "dist-cutoff", 0.3, "If --moved-taxa, then this is the distance cutoff to consider a branch for moving taxa computation. It is the normalized distance to the current bootstrap tree (e.g. 0.05). Must be between 0 and 1, otherwise set to 0")
}

func init() {
	computesupportCmd.AddCommand(boosterCmd)
	computesupportCmd.AddCommand(tbeCmd)

	addTBEFlags(boosterCmd)
	addTBEFlags(tbeCmd)
}

func writeLogBooster() {
	supportLog.WriteString("BOOSTER Support\n")
	supportLog.WriteString(fmt.Sprintf("Date        : %s\n", time.Now().Format(time.RFC822)))
	supportLog.WriteString(fmt.Sprintf("Input tree  : %s\n", supportIntree))
	supportLog.WriteString(fmt.Sprintf("Boot trees  : %s\n", supportBoottrees))
	supportLog.WriteString(fmt.Sprintf("Output tree : %s\n", supportOutFile))
	supportLog.WriteString(fmt.Sprintf("CPUs        : %d\n", rootCpus))
}
