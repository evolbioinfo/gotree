package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/evolbioinfo/gotree/io"
	"github.com/evolbioinfo/gotree/tree"
	"github.com/spf13/cobra"
)

// bipartitiontreeCmd represents the bipartitiontree command
var bipartitiontreeCmd = &cobra.Command{
	Use:   "bipartitiontree",
	Short: "Builds a tree with only one branch/bipartition",
	Long: `Builds a tree with only one branch/bipartition.

To do so, it takes an input tree, and one set of tip/leave names.

It will output a tree with one branch separating the given tips from the others of the input tree.

If a given tip does not exist in the input tree, it will not be taken into account (with a warning).

If not tips remain, it will give an error.

Tips may be given using a file with --tipfile (-f) or as last arguments of the command line:

   gotree compute bipartitiontree -i tree.nw -f tipfile -o outtree.nw
or gotree compute bipartitiontree -i tree.nw -o outtree.nw tip1 tip2 tip3
`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		var tipNames []string
		var leftTipsMap map[string]bool = make(map[string]bool)
		var leftTips []string = make([]string, 0, 10)
		var rightTips []string = make([]string, 0, 10)
		var existTip bool
		var outtree *tree.Tree
		var f *os.File
		var tr *tree.Tree

		if f, err = openWriteFile(outtreefile); err != nil {
			io.LogError(err)
			return
		}
		defer closeWriteFile(f, outtreefile)

		if tr, err = readTree(intreefile); err != nil {
			io.LogError(err)
			return
		}

		if err = tr.UpdateTipIndex(); err != nil {
			io.LogError(err)
			return
		}

		if tipfile != "none" {
			if tipNames, err = parseTipsFile(tipfile); err != nil {
				io.LogError(err)
				return
			}
		} else {
			tipNames = args
		}

		// We take the tips that are present in the tree => left side of the bipartition
		for _, t := range tipNames {
			existTip, err = tr.ExistsTip(t)
			if err != nil {
				io.LogError(err)
				return
			}
			if !existTip {
				io.LogWarning(fmt.Errorf("Tip %s does not exist in the tree", t))
			} else {
				leftTipsMap[t] = true
				leftTips = append(leftTips, t)
			}
		}

		if len(leftTips) == 0 {
			io.LogError(errors.New("No given tips exist in the input tree"))
			return
		}

		//We take the tips of the input tree that are not in the map => right side of the bipartition
		for _, t := range tr.Tips() {
			if _, existTip = leftTipsMap[t.Name()]; !existTip {
				rightTips = append(rightTips, t.Name())
			}
		}
		if len(rightTips) == 0 {
			io.LogError(errors.New("No tips left on the right side of the bipartition"))
			return
		}

		outtree, err = tree.BipartitionTree(leftTips, rightTips)

		if err != nil {
			io.LogError(err)
			return
		}

		f.WriteString(outtree.Newick() + "\n")
		return
	},
}

func init() {
	computeCmd.AddCommand(bipartitiontreeCmd)
	bipartitiontreeCmd.Flags().StringVarP(&intreefile, "input", "i", "stdin", "Input tree")
	bipartitiontreeCmd.Flags().StringVarP(&outtreefile, "output", "o", "stdout", "Output tree")
	bipartitiontreeCmd.Flags().StringVarP(&tipfile, "tipfile", "f", "none", "Tip file")
}
