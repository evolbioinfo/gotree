package cmd

import (
	"bufio"
	"fmt"
	"path/filepath"

	"github.com/fredericlemoine/gotree/draw"
	"github.com/fredericlemoine/gotree/io"
	"github.com/spf13/cobra"
)

// pngCmd represents the png command
var cyjsCmd = &cobra.Command{
	Use:   "cyjs",
	Short: "Draw trees in html file using cytoscape js",
	Long:  `Draw trees in html file using cytoscape js.`,
	Run: func(cmd *cobra.Command, args []string) {
		var l draw.TreeLayout
		ntree := 0
		treefile, trees := readTrees(intreefile)
		defer treefile.Close()

		for t := range trees {
			if t.Err != nil {
				io.ExitWithMessage(t.Err)
			}
			fname := outtreefile
			if ntree > 0 {
				extension := filepath.Ext(fname)
				if extension == ".html" {
					fname = fname[0 : len(fname)-len(extension)]
				}
				fname = fmt.Sprintf(fname+"_%03d.html", ntree)
			}
			f := openWriteFile(fname)
			w := bufio.NewWriter(f)
			l = draw.NewCytoscapeLayout(w, drawSupport)
			l.SetSupportCutoff(drawSupportCutoff)
			l.DrawTree(t.Tree)
			w.Flush()
			f.Close()
			ntree++
		}
	},
}

func init() {
	drawCmd.AddCommand(cyjsCmd)
}
