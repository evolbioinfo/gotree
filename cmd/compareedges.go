package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"os"

	"github.com/fredericlemoine/gotree/io"
	"github.com/fredericlemoine/gotree/support"
	"github.com/fredericlemoine/gotree/tree"
	"github.com/spf13/cobra"
)

// compareedgesCmd represents the compareedges command
var compareedgesCmd = &cobra.Command{
	Use:   "edges",
	Short: "Compare edges of a reference tree with another tree",
	Long: `Compare edges of a reference tree with another tree

If the compared tree file contains several trees, it will take the first one only
`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Fprintf(os.Stderr, "Reference : %s\n", intreefile)
		fmt.Fprintf(os.Stderr, "Compared  : %s\n", intree2file)

		refTree := readTree(intreefile)
		refTree.ReinitIndexes()
		names := refTree.SortedTips()

		edges1 := refTree.Edges()
		fmt.Printf("tree\tbrid\tlength\tsupport\tterminal\tdepth\ttopodepth\trightname\tfound")
		if transferdist {
			fmt.Printf("\ttransfer\ttaxatomove\tcomparednodename")
		} else {
			fmt.Printf("\tcomparednodename")
		}
		fmt.Printf("\n")
		treefile, treechan := readTrees(intree2file)
		defer treefile.Close()
		for t2 := range treechan {
			if t2.Err != nil {
				io.ExitWithMessage(t2.Err)
			}
			t2.Tree.ReinitIndexes()

			edges2 := t2.Tree.Edges()
			var min_dist []uint16
			var min_dist_edges []int
			if transferdist {
				tips := refTree.Tips()
				min_dist = make([]uint16, len(edges1))
				min_dist_edges = make([]int, len(edges1))
				var i_matrix [][]uint16 = make([][]uint16, len(edges1))
				var c_matrix [][]uint16 = make([][]uint16, len(edges1))
				var hamming [][]uint16 = make([][]uint16, len(edges1))

				for i, e := range edges1 {
					e.SetId(i)
					min_dist[i] = uint16(len(tips))
					i_matrix[i] = make([]uint16, len(edges2))
					c_matrix[i] = make([]uint16, len(edges2))
					hamming[i] = make([]uint16, len(edges2))
				}
				for i, e := range edges2 {
					e.SetId(i)
				}
				support.Update_all_i_c_post_order_ref_tree(refTree, &edges1, t2.Tree, &edges2, &i_matrix, &c_matrix)
				support.Update_all_i_c_post_order_boot_tree(refTree, uint(len(tips)), &edges1, t2.Tree, &edges2, &i_matrix, &c_matrix, &hamming, &min_dist, &min_dist_edges)
			}

			for i, e1 := range edges1 {
				var nodename string = "-"
				found := false
				for _, e2 := range edges2 {
					if e1.SameBipartition(e2) {
						nodename = e2.Name(t2.Tree.Rooted())
						found = true
						break
					}
				}
				fmt.Printf("%d\t%d\t%s\t%t", t2.Id, i, e1.ToStatsString(false), found)

				if transferdist {
					var movedtaxabuf bytes.Buffer
					if movedtaxa {
						be := edges2[min_dist_edges[i]]
						plus, minus := speciesToMove(e1, be, int(min_dist[i]))
						for k, sp := range plus {
							if k > 0 {
								movedtaxabuf.WriteRune(',')
							}
							movedtaxabuf.WriteRune('+')
							movedtaxabuf.WriteString(names[sp])
						}
						for k, sp := range minus {
							if k > 0 || (k == 0 && len(plus) > 0) {
								movedtaxabuf.WriteRune(',')
							}
							movedtaxabuf.WriteRune('-')
							movedtaxabuf.WriteString(names[sp])
						}
						nodename = be.Name(t2.Tree.Rooted())
					} else {
						movedtaxabuf.WriteRune('-')
					}

					fmt.Printf("\t%d\t%s\t%s", min_dist[e1.Id()], movedtaxabuf.String(), nodename)
				} else {
					fmt.Printf("\t%s", nodename)
				}
				fmt.Printf("\n")
			}
		}
	},
}

func init() {
	compareCmd.AddCommand(compareedgesCmd)
	compareedgesCmd.PersistentFlags().BoolVarP(&transferdist, "transfer-dist", "m", false, "If transfer dist must be computed for each edge")
	compareedgesCmd.PersistentFlags().BoolVar(&movedtaxa, "moved-taxa", false, "only if --transfer-dist is given: Then display, for each branch, taxa that must be moved")
}

// Returns the list of species to move to go from one branch to the other
// Its length should correspond to given dist
// If not, exit with an error
func speciesToMove(e, be *tree.Edge, dist int) ([]uint, []uint) {
	var i uint
	ndiff := 0
	neq := 0
	diffplus := make([]uint, 0, 100)
	diffminus := make([]uint, 0, 100)
	equplus := make([]uint, 0, 100)
	equminus := make([]uint, 0, 100)

	for i = 0; i < e.Bitset().Len(); i++ {
		t1 := e.Bitset().Test(i)
		t2 := be.Bitset().Test(i)
		if t1 != t2 {
			ndiff++
			if t1 {
				diffminus = append(diffminus, i)
			} else {
				diffplus = append(diffplus, i)
			}
		} else {
			neq++
			if t1 {
				equminus = append(equminus, i)
			} else {
				equplus = append(equplus, i)
			}
		}
	}
	if ndiff < neq {
		if ndiff != dist {
			io.ExitWithMessage(errors.New(fmt.Sprintf("Length of moved species array (%d) is not equal to the minimum distance found (%d)", ndiff, dist)))
		}
		return diffplus, diffminus
	}
	if neq != dist {
		io.ExitWithMessage(errors.New(fmt.Sprintf("Length of moved species array (%d) is not equal to the minimum distance found (%d)", neq, dist)))
	}
	return equplus, equminus
}
