package utils

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"github.com/fredericlemoine/gotree/io"
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

	if inputfile == "" || inputfile == "stdin" || inputfile == "-" {
		refTreeFile = os.Stdin
	} else {
		refTreeFile, err = os.Open(inputfile)
		if err != nil {
			return nil, err
		}
	}

	if strings.HasSuffix(inputfile, ".gz") {
		if gr, err := gzip.NewReader(refTreeFile); err != nil {
			io.ExitWithMessage(err)
		} else {
			reader = bufio.NewReader(gr)
		}
	} else {
		reader = bufio.NewReader(refTreeFile)
	}

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
func ReadCompTrees(inputfile string, compTrees chan<- Trees) (int, error) {
	var compTreeFile *os.File
	var compTree *tree.Tree
	var err error
	var reader *bufio.Reader
	id := 0

	if inputfile == "" || inputfile == "stdin" || inputfile == "-" {
		compTreeFile = os.Stdin
	} else {
		if compTreeFile, err = os.Open(inputfile); err != nil {
			return id, err
		}
	}

	if strings.HasSuffix(inputfile, ".gz") {
		if gr, err := gzip.NewReader(compTreeFile); err != nil {
			io.ExitWithMessage(err)
		} else {
			reader = bufio.NewReader(gr)
		}
	} else {
		reader = bufio.NewReader(compTreeFile)
	}

	line, e := ReadUntilSemiColon(reader)
	for e == nil {
		parser := newick.NewParser(strings.NewReader(line))
		if compTree, err = parser.Parse(); err != nil {
			return id, err
		}
		fmt.Println(compTree.Newick())
		compTrees <- Trees{
			compTree,
			id,
		}
		id++
		line, e = ReadUntilSemiColon(reader)
	}
	close(compTrees)
	if err = compTreeFile.Close(); err != nil {
		return id, err
	}
	return id, nil
}
