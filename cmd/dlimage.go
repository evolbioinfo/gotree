package cmd

import (
	"github.com/spf13/cobra"
)

var dltreeid string
var dlformat string
var dloutput string

// dlimageCmd represents the download command
var dlimageCmd = &cobra.Command{
	Use:   "dlimage",
	Short: "Download a tree image from from a server",
	Long:  `Download a tree image from a server`,
}

func init() {
	RootCmd.AddCommand(dlimageCmd)
	dlimageCmd.PersistentFlags().StringVarP(&dltreeid, "tree-id", "i", "", "Tree id to download")
	dlimageCmd.PersistentFlags().StringVarP(&dlformat, "format", "f", "pdf", "Image format (png, pdf, eps, svg)")
	dlimageCmd.PersistentFlags().StringVarP(&dloutput, "output", "o", "", "Image output file")
}
