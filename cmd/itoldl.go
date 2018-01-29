package cmd

import (
	"errors"
	"io/ioutil"

	"github.com/fredericlemoine/gotree/download"
	"github.com/fredericlemoine/gotree/io"
	"github.com/spf13/cobra"
)

var dlconfig string
var dltreeid string
var dlformat string
var dloutput string

// dlitolCmd represents the dlitol command
var dlitolCmd = &cobra.Command{
	Use:   "itol",
	Short: "Download a tree image/file from iTOL",
	Long: `Download a tree image/file from iTOL

Option -c allows to give a configuration file having tab separated key value pairs, 
as defined here:
http://itol.embl.de/help.cgi#bExOpt
`,
	Run: func(cmd *cobra.Command, args []string) {
		if dloutput == "" {
			io.ExitWithMessage(errors.New("Output file must be specified"))
		}
		if dltreeid == "" {
			io.ExitWithMessage(errors.New("Tree id must be specified"))
		}
		format := download.Format(dlformat)
		if format == download.FORMAT_UNKNOWN {
			io.ExitWithMessage(errors.New("Unkown format: " + dlformat))
		}
		var config map[string]string
		if dlconfig != "" {
			var err error
			config, err = readMapFile(dlconfig, false)
			if err != nil {
				io.ExitWithMessage(err)
			}
		} else {
			config = make(map[string]string)
		}

		dl := download.NewItolImageDownloader(config)
		b, err := dl.Download(dltreeid, format)
		if err != nil {
			io.ExitWithMessage(err)
		}
		ioutil.WriteFile(dloutput, b, 0644)
	},
}

func init() {
	downloadCmd.AddCommand(dlitolCmd)
	dlitolCmd.PersistentFlags().StringVarP(&dlconfig, "config", "c", "", "Itol image config file")
	dlitolCmd.PersistentFlags().StringVarP(&dltreeid, "tree-id", "i", "", "Tree id to download")
	dlitolCmd.PersistentFlags().StringVarP(&dlformat, "format", "f", "pdf", "Image format (png, pdf, eps, svg, newick, nexus, phyloxml)")
	dlitolCmd.PersistentFlags().StringVarP(&dloutput, "output", "o", "", "Image output file")
}
