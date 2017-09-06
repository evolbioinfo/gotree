package cmd

import (
	"math/rand"
	"time"

	"github.com/fredericlemoine/gotree/io"
	"github.com/fredericlemoine/gotree/tree"
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

// pruneCmd represents the prune command
var pruneCmd = &cobra.Command{
	Use:   "prune",
	Short: "Remove tips of the input tree that are not in the compared tree",
	Long: `This tool removes tips of the input reference tree that :

1) Are not present in the compared tree (--comp <other tree>) if any or
2) Are present in the given tip file (--tipfile <file>) if any or 
3) Are randomly sampled (--random <num tips>) or
4) Are given on the command line

If several trees are present in the file given by -i, they are all analyzed and 
written in the output.

If -c and -f are not given, this command will take taxa names on command line, for example:
gotree prune -i reftree.nw -o outtree.nw t1 t2 t3 

By order of priority:
1) -f --tipfile <tip file>
2) -c --comp <other tree>
3) --random <number of tips to randomly sample> 
4) tips given on commandline
5) Nothing is done

If -r is given, behavior is reversed, it keep given tips instead of removing them.
`,
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		var specificTipNames []string
		rand.Seed(seed)

		f := openWriteFile(outtreefile)
		comptree := readTree(intree2file)

		// Read ref Trees
		treefile, trees := readTrees(intreefile)
		defer treefile.Close()

		for reftree := range trees {
			if reftree.Err != nil {
				io.ExitWithMessage(reftree.Err)
			}
			var tips []string
			if tipfile != "none" {
				tips = parseTipsFile(tipfile)
				err = reftree.Tree.RemoveTips(revert, tips...)
			} else if comptree != nil {
				specificTipNames = specificTips(reftree.Tree, comptree)
				err = reftree.Tree.RemoveTips(revert, specificTipNames...)
			} else if randomtips > 0 {
				sampled := randomTips(reftree.Tree, randomtips)
				err = reftree.Tree.RemoveTips(revert, sampled...)
			} else {
				err = reftree.Tree.RemoveTips(revert, args...)
			}
			if err != nil {
				io.ExitWithMessage(err)
			}
			f.WriteString(reftree.Tree.Newick() + "\n")
		}
		f.Close()
	},
}

func init() {
	RootCmd.AddCommand(pruneCmd)
	pruneCmd.Flags().StringVarP(&intreefile, "ref", "i", "stdin", "Input reference tree")
	pruneCmd.Flags().StringVarP(&intree2file, "comp", "c", "none", "Input compared tree ")
	pruneCmd.Flags().StringVarP(&outtreefile, "output", "o", "stdout", "Output tree")
	pruneCmd.Flags().StringVarP(&tipfile, "tipfile", "f", "none", "Tip file")
	pruneCmd.Flags().BoolVarP(&revert, "revert", "r", false, "If true, then revert the behavior: will keep only species given in the command line, or keep only the species that are specific to the input tree, or keep only randomly selected taxa")
	pruneCmd.PersistentFlags().Int64VarP(&seed, "seed", "s", time.Now().UTC().UnixNano(), "Initial Random Seed")
	pruneCmd.PersistentFlags().IntVar(&randomtips, "random", 0, "Number of tips to randomly sample")
}
