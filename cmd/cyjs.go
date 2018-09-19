package cmd

import (
	"bufio"
	"fmt"
	goio "io"
	"os"
	"path/filepath"

	"github.com/fredericlemoine/gotree/draw"
	"github.com/fredericlemoine/gotree/io"
	"github.com/fredericlemoine/gotree/tree"
	"github.com/spf13/cobra"
)

// pngCmd represents the png command
var cyjsCmd = &cobra.Command{
	Use:   "cyjs",
	Short: "Draw trees in html file using cytoscape js",
	Long:  `Draw trees in html file using cytoscape js.`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		var l draw.TreeLayout
		var treefile goio.Closer
		var treechan <-chan tree.Trees
		var f *os.File

		ntree := 0

		if treefile, treechan, err = readTrees(intreefile); err != nil {
			io.LogError(err)
			return
		}
		defer treefile.Close()

		for t := range treechan {
			if t.Err != nil {
				io.LogError(t.Err)
				return t.Err
			}
			fname := outtreefile
			if ntree > 0 {
				extension := filepath.Ext(fname)
				if extension == ".html" {
					fname = fname[0 : len(fname)-len(extension)]
				}
				fname = fmt.Sprintf(fname+"_%03d.html", ntree)
			}
			if f, err = openWriteFile(fname); err != nil {
				io.LogError(err)
				return
			}
			w := bufio.NewWriter(f)
			l = draw.NewCytoscapeLayout(w, drawSupport)
			l.SetSupportCutoff(drawSupportCutoff)
			l.DrawTree(t.Tree)
			w.Flush()
			closeWriteFile(f, fname)
			ntree++
		}
		return
	},
}

func init() {
	drawCmd.AddCommand(cyjsCmd)
}
