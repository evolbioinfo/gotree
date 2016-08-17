package utils

import (
	"bufio"
	"github.com/fredericlemoine/gotree/io/newick"
	"github.com/fredericlemoine/gotree/tree"
	"os"
	"strings"
)

// Type for channel of trees
type Trees struct {
	Tree *tree.Tree
	Id   int
}

// Reads one tree from the input file
func ReadRefTree(inputfile string) (*tree.Tree, error) {
	var refTreeFile *os.File
	var reftree *tree.Tree
	var err error
	var reader *bufio.Reader

	if inputfile == "" || inputfile == "stdin" {
		refTreeFile = os.Stdin
	} else {
		refTreeFile, err = os.Open(inputfile)
		if err != nil {
			return nil, err
		}
	}

	reader = bufio.NewReader(refTreeFile)
	if reftree, err = newick.NewParser(reader).Parse(); err != nil {
		return nil, err
	}
	if err = refTreeFile.Close(); err != nil {
		return nil, err
	}

	return reftree, nil
}

// Read a bunch of trees from the input file. One line must define One tree.
// One tree per line
func ReadCompTrees(inputfile string, compTrees chan<- Trees) error {
	var compTreeFile *os.File
	var compTree *tree.Tree
	var err error
	var reader *bufio.Reader

	if compTreeFile, err = os.Open(inputfile); err != nil {
		return err
	}

	reader = bufio.NewReader(compTreeFile)
	id := 0
	line, e := Readln(reader)
	for e == nil {
		parser := newick.NewParser(strings.NewReader(line))
		if compTree, err = parser.Parse(); err != nil {
			return err
		}
		compTrees <- Trees{
			compTree,
			id,
		}
		id++
		line, e = Readln(reader)
	}
	close(compTrees)
	if err = compTreeFile.Close(); err != nil {
		return err
	}
	return nil
}
