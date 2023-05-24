package cmd

import (
	"errors"
	"fmt"
	goio "io"
	"math"
	"runtime"

	"github.com/spf13/cobra"

	"github.com/evolbioinfo/gotree/io"
	"github.com/evolbioinfo/gotree/tree"
)

var comparetreeidentical bool
var comparetreerf bool
var comparetreeweighted bool

// compareCmd represents the compare command
var compareTreesCmd = &cobra.Command{
	Use:   "trees",
	Short: "Compare a reference tree with a set of trees",
	Long: `Compare a reference tree with a set of trees.

If --binary is given:
For each trees in the compared tree file, it will print tab separated values with:
1) The index of the compared tree in the file
2) "true" if the tree is identical, 
   "false" otherwise

Otherwise:
For each trees in the compared tree file, it will print tab separated values with:
1) The index of the compared tree in the file
2) The number of branches that are specific to the reference tree
3) The number of branches that are common to both trees
4) The number of branches that are specific to the compared tree

If --rf is given, it only computes the Robinson-Foulds distance, as the sum of 
reference + compared specific branches.

If --weighted is given:
For each trees in the compared tree file, it will print tab separated values with:
1) The index of the compared tree in the file
2) The weighted Robinson-Foulds distance (Robinson & Foulds, 1979)
3) The Khuner-Felsenstein branch score (Khuner & Felsenstein, 1994)

If --weighted and --binary are given:
For each trees in the compared tree file, it will print tab separated values with:
1) The index of the compared tree in the file
2) "true" if the tree is identical, both in topology and branch lengths, 
   "false" otherwise
`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		var treefile goio.Closer
		var treechan <-chan tree.Trees
		var refTree *tree.Tree
		var stats <-chan tree.BipartitionStats

		if intree2file == "none" {
			err = errors.New("You must provide a file containing compared trees")
			io.LogError(err)
			return
		}

		maxcpus := runtime.NumCPU()
		if rootCpus > maxcpus {
			rootCpus = maxcpus
		}
		if refTree, err = readTree(intreefile); err != nil {
			io.LogError(err)
			return
		}

		if err = refTree.ReinitIndexes(); err != nil {
			io.LogError(err)
		}

		if treefile, treechan, err = readTrees(intree2file); err != nil {
			io.LogError(err)
			return
		}
		defer treefile.Close()

		if comparetreeweighted {
			if stats2, err := tree.CompareWeighted(refTree, treechan, compareTips, comparetreeidentical, rootCpus); err != nil {
				io.LogError(err)
				return err
			} else {
				if comparetreeidentical {
					fmt.Printf("tree\tidentical\n")
				} else {
					fmt.Printf("tree\tweighted_RF\tKF\n")
				}
				for st := range stats2 {
					if st.Err != nil {
						/* We empty the channel if needed*/
						for range stats {
						}
						io.LogError(st.Err)
						return st.Err
					}

					if comparetreeidentical {
						fmt.Printf("%d\t%v\n", st.Id, st.Sametree)
					} else {
						// Computing Weighted Robinson-Foulds and Khuner-Felsenstein distances.
						wrf := 0.0
						kf := 0.0

						for _, diff := range st.Common {
							wrf += math.Abs(diff)
							kf += math.Pow(diff, 2.0)
						}

						for _, container := range [][]float64{st.Tree1, st.Tree2} {
							for _, length := range container {
								wrf += length
								kf += math.Pow(length, 2.0)
							}
						}

						fmt.Printf("%d\t%E\t%E\n", st.Id, wrf, math.Sqrt(kf))
					}
				}
			}

			return
		}

		if stats, err = tree.Compare(refTree, treechan, compareTips, comparetreeidentical, rootCpus); err != nil {
			io.LogError(err)
			return
		}

		if comparetreeidentical {
			fmt.Printf("tree\tidentical\n")
			for st := range stats {
				if st.Err != nil {
					/* We empty the channel if needed*/
					for range stats {
					}
					io.LogError(st.Err)
					return st.Err
				}
				fmt.Printf("%d\t%v\n", st.Id, st.Sametree)
			}
		} else if comparetreerf {
			for st := range stats {
				if st.Err != nil {
					/* We empty the channel if needed*/
					for range stats {
					}
					io.LogError(st.Err)
					return st.Err
				}
				fmt.Printf("%d\n", st.Tree1+st.Tree2)
			}
		} else {
			fmt.Printf("tree\treference\tcommon\tcompared\n")
			for st := range stats {
				if st.Err != nil {
					/* We empty the channel if needed*/
					for range stats {
					}
					io.LogError(st.Err)
					return st.Err
				}
				fmt.Printf("%d\t%d\t%d\t%d\n", st.Id, st.Tree1, st.Common, st.Tree2)
			}
		}
		return
	},
}

func init() {
	compareCmd.AddCommand(compareTreesCmd)
	compareTreesCmd.Flags().BoolVarP(&compareTips, "tips", "l", false, "Include tips in the comparison")
	compareTreesCmd.Flags().BoolVar(&comparetreeidentical, "binary", false, "If true, then just print true (identical tree) or false (different tree) for each compared tree")
	compareTreesCmd.Flags().BoolVar(&comparetreerf, "rf", false, "If true, outputs Robinson-Foulds distance, as the sum of reference + compared specific branches")
	compareTreesCmd.Flags().BoolVar(&comparetreeweighted, "weighted", false, "If true, outputs comparison metrics including branch lengths")
}
