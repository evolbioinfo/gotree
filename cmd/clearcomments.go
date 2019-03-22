package cmd

import (
	goio "io"
	"os"

	"github.com/evolbioinfo/gotree/io"
	"github.com/evolbioinfo/gotree/tree"
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
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		var f *os.File
		var treefile goio.Closer
		var treechan <-chan tree.Trees

		if !edgecomments && !nodecomments {
			edgecomments = true
			nodecomments = true
		}

		if f, err = openWriteFile(outtreefile); err != nil {
			io.LogError(err)
			return
		}
		defer closeWriteFile(f, outtreefile)

		treefile, treechan, err = readTrees(intreefile)
		defer treefile.Close()
		for t := range treechan {
			if t.Err != nil {
				io.LogError(t.Err)
				return t.Err
			}
			if nodecomments {
				t.Tree.ClearNodeComments()
			}
			if edgecomments {
				t.Tree.ClearEdgeComments()
			}
			f.WriteString(t.Tree.Newick() + "\n")
		}
		return
	},
}

func init() {
	commentCmd.AddCommand(clearcommentsCmd)
	clearcommentsCmd.PersistentFlags().BoolVar(&edgecomments, "edges-only", false, "Clear comments on edges only")
	clearcommentsCmd.PersistentFlags().BoolVar(&nodecomments, "nodes-only", false, "Clear comments on nodes only")
}
