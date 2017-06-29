package cmd

import (
	"github.com/spf13/cobra"
)

// downloadCmd represents the download command
var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download trees or images from servers",
	Long:  `Download trees or images from different servers (itol, ncbi taxonomy)`,
}

func init() {
	RootCmd.AddCommand(downloadCmd)
}
