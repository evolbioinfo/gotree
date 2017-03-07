package cmd

import (
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

// pruneCmd represents the prune command
var pruneCmd = &cobra.Command{
	Use:   "prune",
	Short: "Remove tips of the input tree that are not in the compared tree",
	Long: `This tool removes tips of the input reference tree that 
are not present in the compared tree.

In output, we have a tree containing only tips that are common to both trees.

If several trees are present in the file given by -i, they are all analyzed and 
written in the output.

If -c is not given, this command will take taxa names on command line :
gotree prune -i reftree.nw -o outtree.nw t1 t2 t3 

`,
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		var specificTipNames []string

		f := openWriteFile(outtreefile)
		comptree := readTree(intree2file)

		// Read ref Trees
		for reftree := range readTrees(intreefile) {
			if comptree != nil {
				specificTipNames = specificTips(reftree.Tree, comptree)
				err = reftree.Tree.RemoveTips(revert, specificTipNames...)
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
	pruneCmd.Flags().BoolVarP(&revert, "revert", "r", false, "If true, then revert the behavior: will keep only species given in the command line, or remove the species that are in common with compared tree")
}
