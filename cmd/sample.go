package cmd

import (
	//"fmt"
	"math/rand"
	"time"

	"github.com/fredericlemoine/gotree/io"
	"github.com/fredericlemoine/gotree/tree"
	"github.com/spf13/cobra"
)

var numtrees int

// sampleCmd represents the sample command
var sampleCmd = &cobra.Command{
	Use:   "sample",
	Short: "Takes a subsample of the set of trees from the input file",
	Long: `Takes a subsample of the set of trees from the input file.

It can be with or without replacement depending on the presence of the --replace option

If the number of desired trees is > number of input trees: 
  - with --replace: Will take -n trees
  - without --replace: Will take all trees.
`,
	Run: func(cmd *cobra.Command, args []string) {
		totaltrees := 0
		outtrees := make([]*tree.Tree, numtrees)
		rand.Seed(seed)

		treefile, treechan := readTrees(intreefile)
		defer treefile.Close()

		// Standard reservoir sampling
		if !replace {
			for t := range treechan {
				if t.Err != nil {
					io.ExitWithMessage(t.Err)
				}
				if totaltrees < numtrees {
					outtrees[totaltrees] = t.Tree
				} else {
					j := rand.Intn(totaltrees)
					if j < numtrees {
						outtrees[j] = t.Tree
					}
				}
				totaltrees++
			}
			if totaltrees < numtrees {
				outtrees = outtrees[:totaltrees]
			}
		} else {
			// Naive Reservoir sampling with replacement
			// See http://www.sciencedirect.com/science/article/pii/S0167947307001089
			// https://doi.org/10.1016/j.csda.2007.03.010
			for t := range treechan {
				if t.Err != nil {
					io.ExitWithMessage(t.Err)
				}
				totaltrees++
				for j := 0; j < numtrees; j++ {
					// One chance over current number of trees to replace the tree at j
					r := rand.Intn(totaltrees)
					if r == 0 {
						outtrees[j] = t.Tree
					}
				}
			}
		}

		f := openWriteFile(outtreefile)
		defer closeWriteFile(f, outtreefile)

		for _, t := range outtrees {
			f.WriteString(t.Newick() + "\n")
		}
	},
}

func init() {
	RootCmd.AddCommand(sampleCmd)
	sampleCmd.Flags().StringVarP(&intreefile, "input", "i", "stdin", "Input reference trees")
	sampleCmd.Flags().StringVarP(&outtreefile, "output", "o", "stdout", "Output trees")
	sampleCmd.PersistentFlags().IntVarP(&numtrees, "nbtrees", "n", 1, "Number of trees to sample from input file")
	sampleCmd.PersistentFlags().Int64VarP(&seed, "seed", "s", time.Now().UTC().UnixNano(), "Initial Random Seed")
	sampleCmd.PersistentFlags().BoolVar(&replace, "replace", false, "If given, samples with replacement")
}
