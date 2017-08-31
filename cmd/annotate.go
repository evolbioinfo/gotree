package cmd

import (
	"bufio"
	"compress/gzip"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/fredericlemoine/gotree/io"
	"github.com/fredericlemoine/gotree/io/utils"
	"github.com/fredericlemoine/gotree/support"
	"github.com/fredericlemoine/gotree/tree"

	"github.com/spf13/cobra"
)

// annotateCmd represents the annotate command
var annotateCmd = &cobra.Command{
	Use:   "annotate",
	Short: "Annotates internal branches of a tree with given data",
	Long: `Annotates internal branches of a tree with given data.

Data for annotation may be (in order of priority):
- A file with one line per internal node to annotate (-m), and with the following format:
   <name of internal branch>:<name of taxon 1>,<name of taxon2>,...,<name of taxon n>
   => It will take the lca of taxa and annotate it with the given name
   => Output tree won't have bootstrap support at the branches anymore
- A tree with labels on internal nodes (-c). in that case, it will label each branch of 
   the input tree with label of the closest branch of the given compared tree (-c) in terms
   of transfer distance. The labels are of the form: "label_distance_depth)";
If neither -c nor -m are given, gotree annotate will wait for a reference tree on stdin
`,
	Run: func(cmd *cobra.Command, args []string) {
		f := openWriteFile(outtreefile)
		defer f.Close()
		treefile, treechan := readTrees(intreefile)
		defer treefile.Close()

		if mapfile != "none" {
			annotateNames, err := readAnnotateNameMap(mapfile)
			if err != nil {
				io.ExitWithMessage(err)
			}

			for t := range treechan {
				if t.Err != nil {
					io.ExitWithMessage(t.Err)
				}
				t.Tree.Annotate(annotateNames)
				f.WriteString(t.Tree.Newick() + "\n")
			}
		} else {
			if intree2file == "none" {
				intree2file = "stdin"
			}
			// We will annotate branches using labels of closest branches in
			// the closest tree
			compTree := readTree(intree2file)
			compTree.ReinitIndexes()
			edges2 := compTree.Edges()
			for i, e := range edges2 {
				e.SetId(i)
				e.SetSupport(tree.NIL_SUPPORT)
			}

			for t := range treechan {
				if t.Err != nil {
					io.ExitWithMessage(t.Err)
				}
				t.Tree.ReinitIndexes()
				edges1 := t.Tree.Edges()
				var min_dist []uint16
				var min_dist_edges []int
				tips := t.Tree.Tips()
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
				support.Update_all_i_c_post_order_ref_tree(t.Tree, &edges1, compTree, &edges2, &i_matrix, &c_matrix)
				support.Update_all_i_c_post_order_boot_tree(t.Tree, uint(len(tips)), &edges1, compTree, &edges2, &i_matrix, &c_matrix, &hamming, &min_dist, &min_dist_edges)
				for _, e1 := range edges1 {
					if !e1.Right().Tip() {
						e2 := edges2[min_dist_edges[e1.Id()]]
						dist := min_dist[e1.Id()]
						depth, _ := e1.NumTipsRight()
						if dist > uint16(len(tips))/2 {
							dist = uint16(len(tips)) - dist
						}
						e1.Right().SetName(fmt.Sprintf("%s_%d_%d", e2.Name(true), dist, depth))
					}
				}
				f.WriteString(t.Tree.Newick() + "\n")
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(annotateCmd)
	annotateCmd.PersistentFlags().StringVarP(&intreefile, "input", "i", "stdin", "Input tree(s) file")
	annotateCmd.PersistentFlags().StringVarP(&intree2file, "compared", "c", "stdin", "Compared tree file")
	annotateCmd.PersistentFlags().StringVarP(&mapfile, "map-file", "m", "none", "Name map input file")
	annotateCmd.PersistentFlags().StringVarP(&outtreefile, "output", "o", "stdout", "Resolved tree(s) output file")
}

func readAnnotateNameMap(annotateInputMap string) (map[string][]string, error) {
	outmap := make(map[string][]string, 0)
	var mapfile *os.File
	var err error
	var reader *bufio.Reader

	if mapfile, err = os.Open(annotateInputMap); err != nil {
		return outmap, err
	}

	if strings.HasSuffix(annotateInputMap, ".gz") {
		if gr, err2 := gzip.NewReader(mapfile); err2 != nil {
			return outmap, err2
		} else {
			reader = bufio.NewReader(gr)
		}
	} else {
		reader = bufio.NewReader(mapfile)
	}
	line, e := utils.Readln(reader)
	nl := 1
	for e == nil {
		cols := strings.Split(line, ":")
		if len(cols) != 2 {
			return outmap, errors.New(fmt.Sprintf("Map file does not have 2 fields separated by \":\" at line: %d", nl))
		}
		if len(cols[0]) == 0 {
			return outmap, errors.New(fmt.Sprintf("Internal node name must have length > 0 at line : %d", nl))
		}

		cols2 := strings.Split(cols[1], ",")
		if len(cols2) <= 1 {
			return outmap, errors.New(fmt.Sprintf("More than one taxon must be given for an ancestral node: node %s at line: %d", cols[0], nl))
		}

		if _, ok := outmap[cols[0]]; ok {
			return outmap, errors.New(fmt.Sprintf("Internal node name already given: %s at line %d", cols[0], nl))
		}

		outmap[cols[0]] = cols2

		line, e = utils.Readln(reader)
		nl++
	}

	if err = mapfile.Close(); err != nil {
		return outmap, err
	}

	return outmap, nil

}
