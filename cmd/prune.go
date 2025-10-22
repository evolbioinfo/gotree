package cmd

import (
	goio "io"
	"math/rand"
	"os"

	"github.com/evolbioinfo/gotree/io"
	"github.com/evolbioinfo/gotree/tree"
	"github.com/spf13/cobra"
)

func specificTips(ref *tree.Tree, comp *tree.Tree) []string {
	compmap := make(map[string]*tree.Node)
	spectips := make([]string, 0)
	for _, n := range comp.Nodes() {
		if n.Nneigh() == 1 {
			compmap[n.Name()] = n
		}
	}

	for _, n := range ref.Nodes() {
		if n.Nneigh() == 1 {
			_, ok := compmap[n.Name()]
			if !ok {
				spectips = append(spectips, n.Name())
			}
		}
	}
	return spectips
}

// returns a random list of n tips from the tree
// uses reservoir sampling to select random tips
// if n>nb tips: returns the tips
func randomTips(tr *tree.Tree, n int) (sampled []string) {
	sampled = make([]string, n)
	total := 0
	for i, tip := range tr.Tips() {
		if i < n {
			sampled[i] = tip.Name()
		} else {
			j := rand.Intn(i)
			if j < n {
				sampled[j] = tip.Name()
			}
		}
		total++
	}
	if total < n {
		sampled = sampled[:total]
	}
	return
}

var randomtips int
var diversity bool

// pruneCmd represents the prune command
var pruneCmd = &cobra.Command{
	Use:   "prune",
	Short: "Remove tips of the input tree that are not in the compared tree",
	Long: `This tool removes tips of the input reference tree that :

1) Are not present in the compared tree (--comp <other tree>) if any or
2) Are present in the given tip file (--tipfile <file>) if any or 
3) Are randomly sampled (--random <num tips>), accounting for diversity (--diversity) or not, or
4) Are given on the command line

If several trees are present in the file given by -i, they are all analyzed and 
written in the output.

If -c and -f are not given, this command will take taxa names on command line, for example:
gotree prune -i reftree.nw -o outtree.nw t1 t2 t3 

By order of priority:
1) -f --tipfile <tip file>
2) -c --comp <other tree>
3) --random <number of tips to randomly sample>  (with or without --diversity)
4) tips given on commandline
5) Nothing is done

If -r is given, behavior is reversed, it keep given tips instead of removing them.

If --random and --diversity are given: Tips to be removed are selected in order to keep the highest diversity in the tree.
To do so, until the desired number of tips is reached, the closest tips pairs are selected, and one of the tip is chosen 
(randomly) to be deleted. In case of equality, one random pair is selected. The process stops when
the number of desired number of tips to remove is reached (--random <int>). If revert is true (-r --revert), then --random <i> indicates 
the number of tips to keep (as opposed to the number of tips to remove).
`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		var f *os.File
		var comptree *tree.Tree
		var treefile goio.Closer
		var treechan <-chan tree.Trees

		var specificTipNames []string

		if f, err = openWriteFile(outtreefile); err != nil {
			io.LogError(err)
			return
		}
		defer closeWriteFile(f, outtreefile)

		if intree2file != "none" {
			if comptree, err = readTree(intree2file); err != nil {
				io.LogError(err)
				return
			}
		}

		// Read ref Trees
		if treefile, treechan, err = readTrees(intreefile); err != nil {
			io.LogError(err)
			return
		}
		defer treefile.Close()

		var tips []string
		if tipfile != "none" {
			if tips, err = parseTipsFile(tipfile); err != nil {
				io.LogError(err)
				return
			}
		}

		var ntips int
		for reftree := range treechan {
			ntips = len(reftree.Tree.Tips())

			if reftree.Err != nil {
				io.LogError(reftree.Err)
				return reftree.Err
			}
			if tipfile != "none" {
				err = reftree.Tree.RemoveTips(revert, tips...)
			} else if comptree != nil {
				specificTipNames = specificTips(reftree.Tree, comptree)
				err = reftree.Tree.RemoveTips(revert, specificTipNames...)
			} else if randomtips > 0 {
				if diversity {
					if !revert {
						randomtips = ntips - randomtips
					}
					sampled := reftree.Tree.SubSampleDiversity(randomtips)
					err = reftree.Tree.RemoveTips(false, sampled...)
				} else {
					sampled := randomTips(reftree.Tree, randomtips)
					err = reftree.Tree.RemoveTips(revert, sampled...)
				}
			} else {
				err = reftree.Tree.RemoveTips(revert, args...)
			}
			if err != nil {
				io.LogError(err)
				return
			}
			f.WriteString(reftree.Tree.Newick() + "\n")
		}
		return
	},
}

func init() {
	RootCmd.AddCommand(pruneCmd)
	pruneCmd.Flags().StringVarP(&intreefile, "ref", "i", "stdin", "Input reference tree")
	pruneCmd.Flags().StringVarP(&intree2file, "comp", "c", "none", "Input compared tree ")
	pruneCmd.Flags().StringVarP(&outtreefile, "output", "o", "stdout", "Output tree")
	pruneCmd.Flags().StringVarP(&tipfile, "tipfile", "f", "none", "Tip file")
	pruneCmd.Flags().BoolVarP(&revert, "revert", "r", false, "If true, then revert the behavior: will keep only species given in the command line, or keep only the species that are specific to the input tree, or keep only randomly selected taxa")
	pruneCmd.Flags().IntVar(&randomtips, "random", 0, "Number of tips to randomly sample")
	pruneCmd.Flags().BoolVar(&diversity, "diversity", false, "If the random pruning takes into account diversity (only with --random)")
}
