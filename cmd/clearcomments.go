package cmd

import (
	"github.com/fredericlemoine/gotree/io"
	"github.com/spf13/cobra"
)

var edgecomments, nodecomments bool

// commentsCmd represents the comments command
var clearcommentsCmd = &cobra.Command{
	Use:   "clear",
	Short: "Removes node/tip comments",
	Long: `Removes node/tip/edges comments from all nodes/tips/edges of the tree

Example:
t.nw : (t1[c1],t2[c2],(t3[c3],t4[c4])[c5]);

gotree clear comments -i t.nw :
(t1,t2,(t3,t4));

If --edges-only is given: will only remove edge comments
If --nodes-only is given: will only remove nodes comments
If both or none are given, will remove every comments.
`,
	Run: func(cmd *cobra.Command, args []string) {
		if !edgecomments && !nodecomments {
			edgecomments = true
			nodecomments = true
		}

		f := openWriteFile(outtreefile)
		treefile, treechan := readTrees(intreefile)
		defer treefile.Close()
		for t := range treechan {
			if t.Err != nil {
				io.ExitWithMessage(t.Err)
			}
			if edgecomments {
				t.Tree.ClearNodeComments()
			}
			if nodecomments {
				t.Tree.ClearEdgeComments()
			}
			f.WriteString(t.Tree.Newick() + "\n")
		}
		f.Close()

	},
}

func init() {
	commentCmd.AddCommand(clearcommentsCmd)
	clearcommentsCmd.PersistentFlags().BoolVar(&edgecomments, "edges-only", false, "Clear comments on edges only")
	clearcommentsCmd.PersistentFlags().BoolVar(&nodecomments, "nodes-only", false, "Clear comments on nodes only")
}
