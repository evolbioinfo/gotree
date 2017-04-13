package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/fredericlemoine/gotree/io"
	"github.com/fredericlemoine/gotree/upload"
	"github.com/spf13/cobra"
)

var itoluploadid string
var itolprojectname string
var itoltreename string
var itolannotationfile string

// itolCmd represents the itol command
var itolCmd = &cobra.Command{
	Use:   "itol",
	Short: "Upload a tree to iTOL and display the access url",
	Long: `Upload a tree to iTOL and display the access url.

If --id is given, it uploads the tree to the itol account corresponding to the user upload ID.
The upload id is accessible by enabling "Batch upload" option in iTOL user settings. 

If --id is not given, it uploads the tree without account, and will be automatically deleted after 30 days.

If several trees are included in the input file, it will upload all of them, waiting 1 second between each upload

It is possible to give itol annotation files to the uploader:
gotree upload itol -i tree.tree --name tree --user-id uploadkey --project project annotation*.txt

Urls are written on stdout
Server responses are written on stderr

So:
gotree upload itol -i tree.tree --name tree --user-id uploadkey --project project annotation*.txt > urls

Will store only urls in the output file

`,
	Run: func(cmd *cobra.Command, args []string) {
		// args: All annotation files to add to the upload
		upld := upload.NewItolUploader(itoluploadid, itolprojectname, args...)
		i := 0
		treefile, trees := readTrees(intreefile)
		defer treefile.Close()

		for reftree := range trees {
			if reftree.Err != nil {
				io.ExitWithMessage(reftree.Err)
			}
			url, response, err := upld.Upload(fmt.Sprintf("%s_%03d", itoltreename, i), reftree.Tree)
			if err != nil {
				io.ExitWithMessage(err)
			}
			fmt.Println(url)

			fmt.Fprintf(os.Stderr, "-------------------\n")
			fmt.Fprintf(os.Stderr, "<Server response>\n")
			fmt.Fprintf(os.Stderr, response)
			fmt.Fprintf(os.Stderr, "-------------------\n")
			time.Sleep(1 * time.Second)
			i++
		}
	},
}

func init() {
	uploadCmd.AddCommand(itolCmd)
	itolCmd.Flags().StringVar(&itoluploadid, "user-id", "", "iTOL User upload id")
	itolCmd.Flags().StringVar(&itolprojectname, "project", "", "iTOL project to upload the tree")
	itolCmd.Flags().StringVar(&itoltreename, "name", "", "iTOL tree name prefix")
}
