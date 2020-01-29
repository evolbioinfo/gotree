package cmd

import (
	"fmt"
	"os"

	"github.com/evolbioinfo/gotree/download"
	"github.com/evolbioinfo/gotree/io"
	"github.com/evolbioinfo/gotree/tree"
	"github.com/spf13/cobra"
)

var pantheroutput string
var pantherfamily string

// dlncbiCmd represents the dlncbi command
var dlpantherCmd = &cobra.Command{
	Use:   "panther",
	Short: "Downloads a panther family tree from panther (http://pantherdb.org/)",
	Long:  `Downloads a panther family tree from panther (http://pantherdb.org/)`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		var t *tree.Tree
		var f *os.File

		dl := download.NewPantherTreeDownloader()

		if pantherfamily == "none" {
			err = fmt.Errorf("Panther Family ID must be provided")
			io.LogError(err)
			return
		}
		if t, err = dl.Download(pantherfamily); err != nil {
			io.LogError(err)
			return
		}
		if f, err = openWriteFile(pantheroutput); err != nil {
			io.LogError(err)
			return
		}
		f.WriteString(t.Newick() + "\n")
		closeWriteFile(f, ncbioutput)
		return
	},
}

func init() {
	downloadCmd.AddCommand(dlpantherCmd)
	dlpantherCmd.PersistentFlags().StringVarP(&pantheroutput, "output", "o", "stdout", "Panther family newick output file")
	dlpantherCmd.PersistentFlags().StringVarP(&pantherfamily, "family-id", "f", "none", "Panther Family ID to download")
}
