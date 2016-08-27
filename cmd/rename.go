package cmd

import (
	"bufio"
	"compress/gzip"
	"errors"
	"fmt"
	"github.com/fredericlemoine/gotree/io"
	"github.com/fredericlemoine/gotree/io/utils"
	"github.com/fredericlemoine/gotree/tree"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

var renameintree string
var renameouttree string
var renamemapfile string
var renamerevert bool

func readMapFile(file string, revert bool) (map[string]string, error) {
	outmap := make(map[string]string, 0)
	var mapfile *os.File
	var err error
	var reader *bufio.Reader

	if mapfile, err = os.Open(file); err != nil {
		return outmap, err
	}

	if strings.HasSuffix(file, ".gz") {
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
		cols := strings.Split(line, "\t")
		if len(cols) != 2 {
			return outmap, errors.New("Map file does not have 2 fields at line: " + fmt.Sprintf("%d", nl))
		}
		if revert {
			outmap[cols[1]] = cols[0]
		} else {
			outmap[cols[0]] = cols[1]
		}
		line, e = utils.Readln(reader)
		nl++
	}

	if err = mapfile.Close(); err != nil {
		return outmap, err
	}

	return outmap, nil
}

// renameCmd represents the rename command
var renameCmd = &cobra.Command{
	Use:   "rename",
	Short: "Renames tips of the input tree, given a map file",
	Long: `Renames tips of the input tree, given a map file.

Map file must be tab separated with columns:
1) Current name of the tip
2) Desired new name of the tip
(if --revert then it is the other way)

If a tip name does not appear in the map file, it will not be renamed. 
If a name that does not exist appears in the map file, it will not throw an error.
`,
	Run: func(cmd *cobra.Command, args []string) {

		if renamemapfile == "none" {
			io.ExitWithMessage(errors.New("map file is not given"))
		}

		// Read map file
		namemap, err := readMapFile(renamemapfile, renamerevert)
		if err != nil {
			io.ExitWithMessage(err)
		}

		// Read Tree
		var tree *tree.Tree
		tree, err = utils.ReadRefTree(renameintree)
		if err != nil {
			io.ExitWithMessage(err)
		}
		var f *os.File
		if renameouttree != "stdout" {
			f, err = os.Create(renameouttree)
		} else {
			f = os.Stdout
		}
		if err != nil {
			io.ExitWithMessage(err)
		}

		err = tree.Rename(namemap)
		if err != nil {
			io.ExitWithMessage(err)
		}

		f.WriteString(tree.Newick() + "\n")
		f.Close()
	},
}

func init() {
	RootCmd.AddCommand(renameCmd)
	renameCmd.Flags().StringVarP(&renameouttree, "output", "o", "stdout", "Renamed tree output file")
	renameCmd.Flags().StringVarP(&renameintree, "input", "i", "stdin", "Input tree")
	renameCmd.Flags().StringVarP(&renamemapfile, "map", "m", "none", "Tip name map file")
	renameCmd.Flags().BoolVarP(&renamerevert, "revert", "r", false, "Revert orientation of map file")

}
