package cmd

import (
	"github.com/fredericlemoine/gotree/io"
	"github.com/spf13/cobra"
)

// commentsCmd represents the comments command
var clearcommentsCmd = &cobra.Command{
	Use:   "clear",
	Short: "Removes node/tip comments",
	Long: `Removes node/tip comments from all nodes/tips of the tree

Example:
t.nw : (t1[c1],t2[c2],(t3[c3],t4[c4])[c5]);

gotre clear comments -i t.nw :
(t1,t2,(t3,t4));
`,
	Run: func(cmd *cobra.Command, args []string) {
		f := openWriteFile(outtreefile)
		treefile, treechan := readTrees(intreefile)
		defer treefile.Close()
		for t := range treechan {
			if t.Err != nil {
				io.ExitWithMessage(t.Err)
			}
			t.Tree.ClearComments()
			f.WriteString(t.Tree.Newick() + "\n")
		}
		f.Close()

	},
}

func init() {
	commentCmd.AddCommand(clearcommentsCmd)
}
