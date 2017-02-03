package cmd

import (
	"bufio"
	"compress/gzip"
	"errors"
	"os"
	"strings"

	"github.com/fredericlemoine/gotree/io"
	"github.com/fredericlemoine/gotree/io/utils"
	"github.com/fredericlemoine/gotree/tree"
	"github.com/spf13/cobra"
)

// outgroupCmd represents the outgroup command
var outgroupCmd = &cobra.Command{
	Use:   "outgroup",
	Short: "Reroot trees using an outgroup",
	Long: `Reroot the tree using an outgroup given in argument or in stdin.

Example:

Reroot on 1 tip named "Tip10" using stdin:
echo "Tip10" | gotree reroot outgroup -i tree.nw -l - > reroot.nw

Reroot using an outgroup defined by 3 tips using stdin:
echo "Tip1,Tip2,Tip10" | gotree reroot outgroup -i tree.nw -l - > reroot.nw

Reroot using an outgroup defined by 3 tips using command args:

gotree reroot outgroup -i tree.nw Tip1 Tip2 Tip3 > reroot.nw

`,
	Run: func(cmd *cobra.Command, args []string) {
		var tips []string
		if reroottipfile != "none" {
			tips = parseTipsFile(reroottipfile)
		} else if len(args) > 0 {
			tips = args
		} else {
			io.ExitWithMessage(errors.New("Not group given"))
		}

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
	rerootCmd.AddCommand(outgroupCmd)
	outgroupCmd.PersistentFlags().StringVarP(&reroottipfile, "tip-file", "l", "none", "File containing names of tips of the outgroup")
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
