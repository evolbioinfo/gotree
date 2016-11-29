package cmd

import (
	"fmt"
	"github.com/fredericlemoine/gotree/io"
	"github.com/fredericlemoine/gotree/io/utils"
	"github.com/fredericlemoine/gotree/tree"
	"github.com/spf13/cobra"
	"sync"
)

var roccurveIntree string
var roccurveTruetree string
var roccurveMinThr float64
var roccurveMaxThr float64
var roccurveStepThr float64
var roccurveOutFile string
var roccurvePvalue float64
var roccurveMinBrLen float64
var roccurveMaxBrLen float64

type roccurveOutStats struct {
	found   bool
	length  float64
	support float64
	pvalue  float64
}

// roccurveCmd represents the roccurve command
var roccurveCmd = &cobra.Command{
	Use:   "roccurve",
	Short: "Computes true positives and false positives at different thresholds",
	Long: `Computes true positives and false positives at different thresholds

At a given threshold t, the true positives (TP) are the branches that have a support >= t
and that are found in the true tree

At a given threshold t, the false positives (FP) are the branches that have a support >= t
and that are not found in the true tree

You need to provide:
-i           : Input tree, the tree to test
-r           : Reference tree, the true tree
-m           : min threshold 
-M           : max threshold
-s           : step
--length-leq : keep only branches with length <= value
--length-geq : keep only branches with length >= value

As output, a tab delimited file with columns:
1) threshold
2) number of TP
3) number of FP
4) number of TP with branch length filter
5) number of FP with branch length filter

`,
	Run: func(cmd *cobra.Command, args []string) {
		nbsteps := int((float64(roccurveMaxThr)-float64(roccurveMinThr))/float64(roccurveStepThr)) + 1
		inputTreeEdges := make(chan *tree.Edge, 100)
		statResults := make(chan roccurveOutStats, 100)
		tp := make([]int, nbsteps)
		fp := make([]int, nbsteps)
		tplen := make([]int, nbsteps)
		fplen := make([]int, nbsteps)

		var intree, truetree *tree.Tree
		var err error

		f := openWriteFile(roccurveOutFile)
		if intree, err = utils.ReadRefTree(roccurveIntree); err != nil {
			io.ExitWithMessage(err)
		}
		if truetree, err = utils.ReadRefTree(roccurveTruetree); err != nil {
			io.ExitWithMessage(err)
		}

		/* We fill the edges channel */
		go func() {
			for _, e := range intree.Edges() {
				inputTreeEdges <- e
			}
			close(inputTreeEdges)
		}()

		/* Now we compute numbers (multithreaded) */
		var wg sync.WaitGroup
		for cpu := 0; cpu < rootCpus; cpu++ {
			wg.Add(1)
			go func() {
				for e := range inputTreeEdges {
					found := false
					length := e.Length()
					support := e.Support()
					pvalue := e.PValue()
					if !e.Right().Tip() {
						for _, e2 := range truetree.Edges() {
							if !e2.Right().Tip() && e.SameBipartition(e2) {
								found = true
								break
							}
						}
						statResults <- roccurveOutStats{
							found,
							length,
							support,
							pvalue,
						}
					}

				}
				wg.Done()
			}()
		}

		/*
			wait the end of computation in a thread that will close
			the statResults channel
		*/
		go func() {
			wg.Wait()
			close(statResults)
		}()

		for result := range statResults {
			i := 0
			for thr := float64(roccurveMinThr); thr <= float64(roccurveMaxThr); thr += roccurveStepThr {
				if result.support >= thr && result.pvalue <= roccurvePvalue {
					if result.found {
						tp[i]++
					} else {
						fp[i]++
					}
					// Branch length filter
					if result.length >= roccurveMinBrLen || roccurveMinBrLen == -1 {
						if result.length <= roccurveMaxBrLen || roccurveMaxBrLen == -1 {
							if result.found {
								tplen[i]++
							} else {
								fplen[i]++
							}
						}
					}
				}
				i++
			}
		}

		i := 0
		fmt.Fprintf(f, "thr\tTP\tFP\tTPLen\tFPLen\n")
		for thr := float64(roccurveMinThr); thr <= float64(roccurveMaxThr); thr += roccurveStepThr {
			fmt.Fprintf(f, "%f\t%d\t%d\t%d\t%d\n", thr, tp[i], fp[i], tplen[i], fplen[i])
			i++
		}

		f.Close()
	},
}

func init() {
	computeCmd.AddCommand(roccurveCmd)

	roccurveCmd.PersistentFlags().StringVarP(&roccurveIntree, "intree", "i", "stdin", "Input tree file")
	roccurveCmd.PersistentFlags().StringVarP(&roccurveTruetree, "truetree", "r", "none", "True tree file")
	roccurveCmd.PersistentFlags().Float64VarP(&roccurveMinThr, "min", "m", 0, "Min threshold")
	roccurveCmd.PersistentFlags().Float64VarP(&roccurveMaxThr, "max", "M", 1, "Max threshold")
	roccurveCmd.PersistentFlags().Float64VarP(&roccurveStepThr, "step", "s", 0.1, "Step between each threshold")
	roccurveCmd.PersistentFlags().StringVarP(&roccurveOutFile, "out", "o", "stdout", "Output tree file, with supports")
	roccurveCmd.PersistentFlags().Float64VarP(&roccurvePvalue, "pvalue", "p", 1.0, "Keep only branches that have a pvalue <=  value")
	roccurveCmd.PersistentFlags().Float64Var(&roccurveMaxBrLen, "length-leq", -1.0, "Keep only branches that are <= value (-1=No filter) ")
	roccurveCmd.PersistentFlags().Float64Var(&roccurveMinBrLen, "length-geq", -1.0, "Keep only branches that are >= value (-1=No filter) ")
}
