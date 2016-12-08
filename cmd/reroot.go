package cmd

import (
	"bufio"
	"compress/gzip"
	"github.com/fredericlemoine/gotree/io"
	"github.com/fredericlemoine/gotree/io/utils"
	"github.com/fredericlemoine/gotree/tree"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

var reroottipfile string
var rerootinputfile string
var rerootoutputfile string

// rerootCmd represents the reroot command
var rerootCmd = &cobra.Command{
	Use:   "reroot",
	Short: "Reroot the tree using an outgroup",
	Long: `Reroot the tree using an outgroup given in argument or in stdin.

Example:


`,
	Run: func(cmd *cobra.Command, args []string) {
		tips := parseTipsFile(reroottipfile)

		var err error
		var nbtrees int

		compareChannel := make(chan tree.Trees, 15)

		go func() {
			if nbtrees, err = utils.ReadCompTrees(rerootinputfile, compareChannel); err != nil {
				io.ExitWithMessage(err)
			}
		}()

		var f *os.File
		if rerootoutputfile != "stdout" {
			f, err = os.Create(rerootoutputfile)
		} else {
			f = os.Stdout
		}
		if err != nil {
			io.ExitWithMessage(err)
		}

		for t2 := range compareChannel {
			err = t2.Tree.RerootOutGroup(tips...)
			if err != nil {
				io.ExitWithMessage(err)
			}

			f.WriteString(t2.Tree.Newick() + "\n")
		}

		f.Close()
	},
}

func init() {
	RootCmd.AddCommand(rerootCmd)
	rerootCmd.PersistentFlags().StringVarP(&reroottipfile, "tip-file", "l", "stdin", "File containing names of tips of the outgroup")
	rerootCmd.PersistentFlags().StringVarP(&rerootinputfile, "input", "i", "stdin", "Input Tree")
	rerootCmd.PersistentFlags().StringVarP(&rerootoutputfile, "output", "o", "stdout", "Rerooted output tree file")
}

func parseTipsFile(file string) []string {
	var f *os.File
	var r *bufio.Reader
	tips := make([]string, 0, 100)
	var err error
	if file == "stdin" || file == "-" {
		f = os.Stdin
	} else {
		f, err = os.Open(file)
		if err != nil {
			io.ExitWithMessage(err)
		}
	}

	if strings.HasSuffix(file, ".gz") {
		if gr, err := gzip.NewReader(f); err != nil {
			io.ExitWithMessage(err)
		} else {
			r = bufio.NewReader(gr)
		}
	} else {
		r = bufio.NewReader(f)
	}

	l, e := Readln(r)
	for e == nil {
		for _, name := range strings.Split(l, ",") {
			tips = append(tips, name)
		}
		l, e = Readln(r)
	}
	return tips
}
