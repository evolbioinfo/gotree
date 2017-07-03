package cmd

import (
	"github.com/fredericlemoine/gotree/download"
	"github.com/fredericlemoine/gotree/io"
	"github.com/spf13/cobra"
)

var ncbioutput string
var ncbitiptaxid bool
var ncbinodetaxid bool

// dlncbiCmd represents the dlncbi command
var dlncbiCmd = &cobra.Command{
	Use:   "ncbitax",
	Short: "Downloads the full ncbi taxonomy in newick format",
	Long:  `Downloads the full ncbi taxonomy in newick format`,
	Run: func(cmd *cobra.Command, args []string) {
		dl := download.NewNcbiTreeDownloader()
		dl.SetInternalNodesTaxId(ncbinodetaxid)
		dl.SetTipsTaxId(ncbitiptaxid)
		t, err := dl.Download("")
		if err != nil {
			io.ExitWithMessage(err)
		}
		f := openWriteFile(ncbioutput)
		f.WriteString(t.Newick() + "\n")
		f.Close()
	},
}

func init() {
	downloadCmd.AddCommand(dlncbiCmd)
	dlncbiCmd.PersistentFlags().StringVarP(&ncbioutput, "output", "o", "stdout", "NCBI newick output file")
	dlncbiCmd.PersistentFlags().BoolVar(&ncbitiptaxid, "tips-taxid", false, "Keeps tax id as tip names")
	dlncbiCmd.PersistentFlags().BoolVar(&ncbinodetaxid, "nodes-taxid", false, "Keeps tax id as internal nodes identifiers")
}
