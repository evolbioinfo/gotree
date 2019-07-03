package cmd

import (
	"os"

	"github.com/evolbioinfo/gotree/io"
	"github.com/evolbioinfo/gotree/tree"
	"github.com/fredericlemoine/gostats"
	"github.com/spf13/cobra"
)

func starTree(nbtrees int, nbtips int, output string) error {
	var f *os.File
	var err error
	var t *tree.Tree

	if output != "stdout" && output != "-" {
		f, err = os.Create(output)
		defer f.Close()
	} else {
		f = os.Stdout
	}
	if err != nil {
		return err
	}

	for i := 0; i < nbtrees; i++ {
		if t, err = tree.StarTree(nbtips); err != nil {
			return err
		}
		for _, e := range t.Edges() {
			e.SetLength(gostats.Exp(1.0 / setlengthmean))
		}

		f.WriteString(t.Newick() + "\n")
	}

	return nil
}

// startreeCmd represents the binarytree command
var startreeCmd = &cobra.Command{
	Use:   "startree",
	Short: "Generates a star tree",
	Long: `Generates a star tree.

--rooted option is not functional here.
`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := starTree(generateNbTrees, generateNbTips, generateOutputfile); err != nil {
			io.LogError(err)
			return
		}
	},
}

func init() {
	generateCmd.AddCommand(startreeCmd)
	startreeCmd.PersistentFlags().IntVarP(&generateNbTips, "nbtips", "l", 10, "Number of tips/leaves of the tree to generate")
}
