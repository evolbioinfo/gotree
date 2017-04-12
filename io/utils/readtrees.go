package utils

import (
	"bufio"
	"strings"

	"github.com/fredericlemoine/gotree/io/newick"
	"github.com/fredericlemoine/gotree/tree"
)

func ReadRefTree(inputfile string) (*tree.Tree, error) {
	if f, r, err := GetReader(inputfile); err != nil {
		return nil, err
	} else {
		if t, err2 := ReadTreeFile(r); err2 != nil {
			return nil, err2
		} else {
			if err = f.Close(); err != nil {
				return nil, err
			}
			return t, nil
		}
	}

}

// Reads one tree from the input reader
// this function does not close the reader
func ReadTreeFile(reader *bufio.Reader) (*tree.Tree, error) {
	var reftree *tree.Tree
	var err error

	if reftree, err = newick.NewParser(reader).Parse(); err != nil {
		return nil, err
	}

	return reftree, nil
}

// Reads all the trees from the input file and send them to the channel
func ReadMultiTreeFile(inputfile string, compTrees chan<- tree.Trees) (int, error) {
	var i int
	var readerr, closerr error
	if f, r, err := GetReader(inputfile); err != nil {
		return 0, err
	} else {
		if i, readerr = ReadMultiTrees(r, compTrees); readerr != nil {
			return i, readerr
		} else {
			if closerr = f.Close(); closerr != nil {
				return i, closerr
			}
		}
	}
	return i, nil
}

// Read a bunch of trees from the input reader and send each of them to the channel
// this function does not close the reader
func ReadMultiTrees(reader *bufio.Reader, compTrees chan<- tree.Trees) (int, error) {
	var compTree *tree.Tree
	var err error
	id := 0

	line, e := ReadUntilSemiColon(reader)
	for e == nil {
		parser := newick.NewParser(strings.NewReader(line))
		if compTree, err = parser.Parse(); err != nil {
			return id, err
		}
		compTrees <- tree.Trees{
			compTree,
			id,
		}
		id++
		line, e = ReadUntilSemiColon(reader)
	}
	return id, nil
}
