package cmd

import (
	"os"

	"github.com/evolbioinfo/gotree/download"
	"github.com/evolbioinfo/gotree/io"
	"github.com/evolbioinfo/gotree/tree"
	"github.com/spf13/cobra"
)

var ncbioutput string
var ncbitiptaxid bool
var ncbinodetaxid bool
var ncbitaxidtoname string

// dlncbiCmd represents the dlncbi command
var dlncbiCmd = &cobra.Command{
	Use:   "ncbitax",
	Short: "Downloads the full ncbi taxonomy in newick format",
	Long:  `Downloads the full ncbi taxonomy in newick format`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		var t *tree.Tree
		var f *os.File

		dl := download.NewNcbiTreeDownloader()
		dl.SetInternalNodesTaxId(ncbinodetaxid)
		dl.SetTipsTaxId(ncbitiptaxid)
		if ncbitaxidtoname != "none" {
			dl.SetMapFileOutputPath(ncbitaxidtoname)
		}
		if t, err = dl.Download(""); err != nil {
			io.LogError(err)
			return
		}
		if f, err = openWriteFile(ncbioutput); err != nil {
			io.LogError(err)
			return
		}
		f.WriteString(t.Newick() + "\n")
		closeWriteFile(f, ncbioutput)
		return
	},
}

func init() {
	downloadCmd.AddCommand(dlncbiCmd)
	dlncbiCmd.PersistentFlags().StringVarP(&ncbioutput, "output", "o", "stdout", "NCBI newick output file")
	dlncbiCmd.PersistentFlags().BoolVar(&ncbitiptaxid, "tips-taxid", false, "Keeps tax id as tip names")
	dlncbiCmd.PersistentFlags().BoolVar(&ncbinodetaxid, "nodes-taxid", false, "Keeps tax id as internal nodes identifiers")
	dlncbiCmd.PersistentFlags().StringVar(&ncbitaxidtoname, "map", "none", "Output mapping file between taxid and species name (tab separated)")
}
