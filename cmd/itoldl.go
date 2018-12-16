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
https://itol.embl.de/help.cgi#bExOpt
`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		var b []byte

		if dloutput == "" {
			err = errors.New("Output file must be specified")
			io.LogError(err)
			return
		}
		if dltreeid == "" {
			err = errors.New("Tree id must be specified")
			io.LogError(err)
			return
		}
		format := download.Format(dlformat)
		if format == download.FORMAT_UNKNOWN {
			err = errors.New("Unkown format: " + dlformat)
			io.LogError(err)
			return
		}
		var config map[string]string
		if dlconfig != "" {
			if config, err = readMapFile(dlconfig, false); err != nil {
				io.LogError(err)
				return
			}
		} else {
			config = make(map[string]string)
		}

		dl := download.NewItolImageDownloader(config)
		if b, err = dl.Download(dltreeid, format); err != nil {
			io.LogError(err)
			return
		}
		ioutil.WriteFile(dloutput, b, 0644)
		return
	},
}

func init() {
	downloadCmd.AddCommand(dlitolCmd)
	dlitolCmd.PersistentFlags().StringVarP(&dlconfig, "config", "c", "", "Itol image config file")
	dlitolCmd.PersistentFlags().StringVarP(&dltreeid, "tree-id", "i", "", "Tree id to download")
	dlitolCmd.PersistentFlags().StringVarP(&dlformat, "format", "f", "pdf", "Image format (png, pdf, eps, svg, newick, nexus, phyloxml)")
	dlitolCmd.PersistentFlags().StringVarP(&dloutput, "output", "o", "", "Image output file")
}
