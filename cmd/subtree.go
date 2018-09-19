package cmd

import (
	"fmt"
	"log"

	"github.com/fredericlemoine/gotree/io"
	"github.com/fredericlemoine/gotree/tree"
	"github.com/spf13/cobra"
)

// subtreeCmd represents the subtree command
var subtreeCmd = &cobra.Command{
	Use:   "subtree",
	Short: "Select a subtree from the input tree whose root has the given name",
	Long: `Select a subtree from the input tree whose root has the given name.

The name may be a regexp, for example :
gotree subtree -i tree.nhx -n "^Mammal.*"

If several nodes match the given name/regexp, do nothing, and print the name of matching nodes.

The only matching node must be an internal node, otherwise, it will do nothing and print the tip.

`,
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		var nodes []*tree.Node

		f := openWriteFile(outtreefile)
		defer closeWriteFile(f, outtreefile)

		i := 0
		treefile, treechan := readTrees(intreefile)
		defer treefile.Close()
		for t := range treechan {
			if t.Err != nil {
				io.ExitWithMessage(t.Err)
			}

			nodes, err = t.Tree.SelectNodes(inputname)
			if err != nil {
				io.ExitWithMessage(err)
			}
			switch len(nodes) {
			case 1:
				n := nodes[0]
				if n.Tip() {
					log.Print(fmt.Sprintf("Tree %d: Node %s is a tip", i, n.Name()))
				} else {
					subtree := t.Tree.SubTree(n)
					f.WriteString(subtree.Newick() + "\n")
				}
			case 0:
				log.Print(fmt.Sprintf("Tree %d: No node matches input name", i))
			default:
				log.Print(fmt.Sprintf("Tree %d: Two many nodes match input name (%d)", i, len(nodes)))
				for _, n := range nodes {
					log.Print(fmt.Sprintf("Node: %s", n.Name()))
				}
			}
			i++
		}
	},
}

func init() {
	RootCmd.AddCommand(subtreeCmd)
	subtreeCmd.PersistentFlags().StringVarP(&inputname, "name", "n", "none", "Name of the node to select as the root of the subtree (maybe a regex)")
	subtreeCmd.PersistentFlags().StringVarP(&intreefile, "input", "i", "stdin", "Input tree")
	subtreeCmd.PersistentFlags().StringVarP(&outtreefile, "output", "o", "stdout", "Output tree file")
}
