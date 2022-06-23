package cmd

import (
	"bufio"
	"compress/gzip"
	"fmt"
	goio "io"
	"os"
	"strings"

	"github.com/evolbioinfo/gotree/io"
	"github.com/evolbioinfo/gotree/io/fileutils"
	"github.com/evolbioinfo/gotree/mutils"
	"github.com/evolbioinfo/gotree/support"
	"github.com/evolbioinfo/gotree/tree"

	"github.com/spf13/cobra"
)

var annotateComment bool

// annotateCmd represents the annotate command
var annotateCmd = &cobra.Command{
	Use:   "annotate",
	Short: "Annotates internal branches of a tree with given data",
	Long: `Annotates internal branches of a tree with given data.

Annotations may be (in order of priority):
- A tree with labels on internal nodes (-c). in that case, it will label each branch of 
   the input tree with label of the closest branch of the given compared tree (-c) in terms
   of transfer distance. The labels are of the form: "label_distance_depth"; Only internal branches
   are annotated, and no internal branch is annotated with a terminal branch.
- A file with one line per internal node to annotate (-m), and with the following format:
   <name of internal branch/node n1>:<name of taxon n2>,<name of taxon n3>,...,<name of taxon ni>
	=> If 0 name is given after ':' an error is returned
	=> If 1 name 'n2' is given after ':' : we search for n2 in the tree (tip or internal node)
       and rename it as n1
    => If > 1 names '[n2,...,ni]' are given after ':' : We find the LCA of every tips whose name 
	   is in '[n2,...,ni]' and rename it as n1.


If --comment is specified, then we do not change the names, but the comments of the given nodes.
Otherwise output tree won't have bootstrap support at the branches anymore

If neither -c nor -m are given, gotree annotate will wait for a reference tree on stdin
`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		var f *os.File
		var treefile goio.Closer
		var treechan <-chan tree.Trees
		var compTree *tree.Tree
		var annotateNames [][]string

		if f, err = openWriteFile(outtreefile); err != nil {
			io.LogError(err)
			return
		}
		defer closeWriteFile(f, outtreefile)

		if treefile, treechan, err = readTrees(intreefile); err != nil {
			io.LogError(err)
			return
		}
		defer treefile.Close()

		if mapfile != "none" {
			annotateNames, err = readAnnotateNameMap(mapfile)
			if err != nil {
				io.LogError(err)
				return
			}

			for t := range treechan {
				if t.Err != nil {
					io.LogError(t.Err)
					return t.Err
				}
				t.Tree.Annotate(annotateNames, annotateComment)
				f.WriteString(t.Tree.Newick() + "\n")
			}
		} else {
			if intree2file == "none" {
				intree2file = "stdin"
			}
			// We will annotate branches using labels of closest branches in
			// the closest tree
			if compTree, err = readTree(intree2file); err != nil {
				io.LogError(err)
				return
			}
			if err = compTree.ReinitIndexes(); err != nil {
				io.LogError(err)
				return
			}

			edges2 := compTree.Edges()
			for i, e := range edges2 {
				e.SetId(i)
				e.SetSupport(tree.NIL_SUPPORT)
			}

			for t := range treechan {
				if t.Err != nil {
					io.LogError(t.Err)
					return t.Err
				}

				if err = t.Tree.ReinitIndexes(); err != nil {
					io.LogError(err)
					return
				}

				edges1 := t.Tree.Edges()
				tips := t.Tree.Tips()

				for i, e := range edges1 {
					e.SetId(i)
				}

				for _, e1 := range edges1 {
					if !e1.Right().Tip() {

						dist, minedges, _, _ := support.MinTransferDist(e1, t.Tree, compTree, len(tips), edges2, false)
						depth := e1.NumTipsRight()
						if len(minedges) > 0 {
							// We take the first min dist edge
							e2 := minedges[0]
							// If root edge and rooted tree, we take the closest branch
							if e2.Left().Nneigh() == 2 {
								var t3, t2, t1 int

								e3 := e2.Left().Edges()[0]
								if e3 == e2 {
									e3 = e2.Left().Edges()[1]
								}

								t3 = e3.NumTipsRight()
								t2 = e2.NumTipsRight()
								t1 = e1.NumTipsRight()

								if mutils.Abs(t3-t1) < mutils.Abs(t2-t1) {
									e2 = e3
								}
							}

							if !e2.Right().Tip() {
								if annotateComment {
									e1.Right().AddComment(fmt.Sprintf("%s_%d_%d", e2.Name(true), dist, depth))
								} else {
									e1.Right().SetName(fmt.Sprintf("%s_%d_%d", e2.Name(true), dist, depth))
								}
							}
						}
					}
				}
				f.WriteString(t.Tree.Newick() + "\n")
			}
		}
		return
	},
}

func init() {
	RootCmd.AddCommand(annotateCmd)
	annotateCmd.PersistentFlags().StringVarP(&intreefile, "input", "i", "stdin", "Input tree(s) file")
	annotateCmd.PersistentFlags().StringVarP(&intree2file, "compared", "c", "stdin", "Compared tree file")
	annotateCmd.PersistentFlags().StringVarP(&mapfile, "map-file", "m", "none", "Name map input file")
	annotateCmd.PersistentFlags().BoolVar(&annotateComment, "comment", false, "Annotations are stored in Newick comment fields")
	annotateCmd.PersistentFlags().StringVarP(&outtreefile, "output", "o", "stdout", "Resolved tree(s) output file")
}

// returns the list of annotations
// Each element of the slice is a list of strings
// The first element of each list is the annotation
// the other elements are the tip names
func readAnnotateNameMap(annotateInputMap string) ([][]string, error) {
	outmap := make([][]string, 0)
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
	line, e := fileutils.Readln(reader)
	nl := 1
	for e == nil {
		cols := strings.Split(line, ":")
		if len(cols) != 2 {
			return outmap, fmt.Errorf("Map file does not have 2 fields separated by \":\" at line: %d", nl)
		}
		if len(cols[0]) == 0 {
			return outmap, fmt.Errorf("Internal node name must have length > 0 at line : %d", nl)
		}

		cols2 := strings.Split(cols[1], ",")
		if len(cols2) < 1 {
			return outmap, fmt.Errorf("More than one taxon must be given for an ancestral node: node %s at line: %d", cols[0], nl)
		}
		cols2 = append([]string{cols[0]}, cols2...)
		outmap = append(outmap, cols2)

		line, e = fileutils.Readln(reader)
		nl++
	}

	if err = mapfile.Close(); err != nil {
		return outmap, err
	}

	return outmap, nil

}
