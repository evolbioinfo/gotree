package cmd

import (
	"github.com/fredericlemoine/gotree/io"
	"github.com/fredericlemoine/gotree/io/utils"
	"github.com/fredericlemoine/gotree/tree"
	"github.com/spf13/cobra"
	"os"
)

var minbrlenCutoff float64

// minbrlenCmd represents the minbrlen command
var minbrlenCmd = &cobra.Command{
	Use:   "minbrlen",
	Short: "This will set a min branch length to all branches with length < cutoff",
	Long: `This will set a min branch length to all branches with length < cutoff

`,
	Run: func(cmd *cobra.Command, args []string) {
		// Read Tree
		var t *tree.Tree
		var err error
		t, err = utils.ReadRefTree(transformInputTree)
		if err != nil {
			io.ExitWithMessage(err)
		}
		var f *os.File
		if transformOutputTree != "stdout" {
			f, err = os.Create(transformOutputTree)
		} else {
			f = os.Stdout
		}
		if err != nil {
			io.ExitWithMessage(err)
		}

		for _, e := range t.Edges() {
			if e.Length() < minbrlenCutoff {
				e.SetLength(minbrlenCutoff)
			}
		}
		if err != nil {
			io.ExitWithMessage(err)
		}

		f.WriteString(t.Newick() + "\n")
		f.Close()

	},
}

func init() {
	transformCmd.AddCommand(minbrlenCmd)
	minbrlenCmd.Flags().Float64VarP(&minbrlenCutoff, "length", "l", 0.0, "Min Length cutoff")

}
