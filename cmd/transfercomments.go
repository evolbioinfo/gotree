package cmd

import (
	goio "io"
	"os"
	"strings"

	"github.com/evolbioinfo/gotree/io"
	"github.com/evolbioinfo/gotree/tree"
	"github.com/spf13/cobra"
)

var transferreverse bool

// commentsCmd represents the comments command
var transferCommentsCmd = &cobra.Command{
	Use:   "transfer",
	Short: "Transfers node names to comments",
	Long: `Transfers node names to comments and removes node names
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
			if nodecomments {
				for _, n := range t.Tree.Nodes() {
					if transferreverse {
						commentsjoined := strings.Join(n.Comments(), ",")
						if !n.Tip() && commentsjoined != "" {
							n.SetName(commentsjoined)
							n.ClearComments()
						}
					} else {
						if !n.Tip() && n.Name() != "" {
							n.AddComment(n.Name())
							n.SetName("")
						}
					}
				}
			}
			f.WriteString(t.Tree.Newick() + "\n")
		}
		return
	},
}

func init() {
	commentCmd.AddCommand(transferCommentsCmd)
	transferCommentsCmd.PersistentFlags().BoolVar(&transferreverse, "reverse", false, "Reverses the orientation of the transfer (comment to name)")
}
