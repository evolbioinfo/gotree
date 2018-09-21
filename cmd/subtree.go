package cmd

import (
	"fmt"
	goio "io"
	"log"
	"os"

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
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		var f *os.File
		var treefile goio.Closer
		var treechan <-chan tree.Trees

		var nodes []*tree.Node

		if f, err = openWriteFile(outtreefile); err != nil {
			io.LogError(err)
			return
		}
		defer closeWriteFile(f, outtreefile)

		i := 0
		if treefile, treechan, err = readTrees(intreefile); err != nil {
			io.LogError(err)
			return
		}
		defer treefile.Close()
		for t := range treechan {
			if t.Err != nil {
				io.LogError(t.Err)
				return t.Err
			}

			nodes, err = t.Tree.SelectNodes(inputname)
			if err != nil {
				io.LogError(err)
				return
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
		return
	},
}

func init() {
	RootCmd.AddCommand(subtreeCmd)
	subtreeCmd.PersistentFlags().StringVarP(&inputname, "name", "n", "none", "Name of the node to select as the root of the subtree (maybe a regex)")
	subtreeCmd.PersistentFlags().StringVarP(&intreefile, "input", "i", "stdin", "Input tree")
	subtreeCmd.PersistentFlags().StringVarP(&outtreefile, "output", "o", "stdout", "Output tree file")
}
