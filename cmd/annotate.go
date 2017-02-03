package cmd

import (
	"bufio"
	"compress/gzip"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/fredericlemoine/gotree/io"
	"github.com/fredericlemoine/gotree/io/utils"
	"github.com/fredericlemoine/gotree/tree"
	"github.com/spf13/cobra"
)

var annotateInputTree string
var annotateOutputTree string
var annotateIntrees chan tree.Trees
var annotateOutTrees *os.File
var annotateNames map[string][]string
var annotateInputMap string

// annotateCmd represents the annotate command
var annotateCmd = &cobra.Command{
	Use:   "annotate",
	Short: "Annotates internal branches of a tree with given data",
	Long: `Annotates internal branches of a tree with given data.

Takes a map file with one line per internal node to annotate:
<name of internal branch>:<name of taxon 1>,<name of taxon2>,...,<name of taxon n>

=> It will take the lca of taxa and annotate it with the given name
=> Output tree won't have bootstrap support at the branches anymore
`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if annotateInputMap == "none" {
			io.ExitWithMessage(errors.New("You should give a map file for node names"))
		}
		var err error
		var nbtrees int = 0
		annotateIntrees = make(chan tree.Trees, 15)
		/* Read ref tree(s) */
		go func() {
			if nbtrees, err = utils.ReadCompTrees(annotateInputTree, annotateIntrees); err != nil {
				io.ExitWithMessage(err)
			}
		}()
		annotateOutTrees = openWriteFile(annotateOutputTree)

		if annotateNames, err = readAnnotateNameMap(annotateInputMap); err != nil {
			io.ExitWithMessage(err)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		for t := range annotateIntrees {
			t.Tree.Annotate(annotateNames)
			annotateOutTrees.WriteString(t.Tree.Newick() + "\n")
		}

	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		annotateOutTrees.Close()
	},
}

func init() {
	RootCmd.AddCommand(annotateCmd)
	annotateCmd.PersistentFlags().StringVarP(&annotateInputTree, "input", "i", "stdin", "Input tree(s) file")
	annotateCmd.PersistentFlags().StringVarP(&annotateInputMap, "map-file", "m", "none", "Name map input file")
	annotateCmd.PersistentFlags().StringVarP(&annotateOutputTree, "output", "o", "stdout", "Resolved tree(s) output file")
}

func readAnnotateNameMap(annotateInputMap string) (map[string][]string, error) {
	outmap := make(map[string][]string, 0)
	var mapfile *os.File
	var err error
	var reader *bufio.Reader

	if mapfile, err = os.Open(annotateInputMap); err != nil {
		return outmap, err
	}

	if strings.HasSuffix(annotateInputMap, ".gz") {
		if gr, err2 := gzip.NewReader(mapfile); err2 != nil {
			return outmap, err2
		} else {
			reader = bufio.NewReader(gr)
		}
	} else {
		reader = bufio.NewReader(mapfile)
	}
	line, e := utils.Readln(reader)
	nl := 1
	for e == nil {
		cols := strings.Split(line, ":")
		if len(cols) != 2 {
			return outmap, errors.New(fmt.Sprintf("Map file does not have 2 fields separated by \":\" at line: %d", nl))
		}
		if len(cols[0]) == 0 {
			return outmap, errors.New(fmt.Sprintf("Internal node name must have length > 0 at line : %d", nl))
		}

		cols2 := strings.Split(cols[1], ",")
		if len(cols2) <= 1 {
			return outmap, errors.New(fmt.Sprintf("More than one taxon must be given for an ancestral node: node %s at line: %d", cols[0], nl))
		}

		if _, ok := outmap[cols[0]]; ok {
			return outmap, errors.New(fmt.Sprintf("Internal node name already given: %s at line %d", cols[0], nl))
		}

		outmap[cols[0]] = cols2

		line, e = utils.Readln(reader)
		nl++
	}

	if err = mapfile.Close(); err != nil {
		return outmap, err
	}

	return outmap, nil

}
