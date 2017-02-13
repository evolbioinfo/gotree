package utils

import (
	"bufio"
	"compress/gzip"
	"github.com/fredericlemoine/gotree/io"
	"github.com/fredericlemoine/gotree/io/newick"
	"github.com/fredericlemoine/gotree/tree"
	"os"
	"strings"
)

func ReadRefTree(inputfile string) (*tree.Tree, error) {
	if r, err := GetReader(inputfile); err != nil {
		return nil, err
	} else {
		if t, err2 := ReadRefTreeFile(r); err2 != nil {
			return nil, err2
		} else {
			if err = r.Close(); err != nil {
				return nil, err
			}
			return t, nil
		}
	}

}

// Reads one tree from the input file
func ReadRefTreeFile(reader *io.Reader) (*tree.Tree, error) {
	var reftree *tree.Tree
	var err error

	if reftree, err = newick.NewParser(reader).Parse(); err != nil {
		return nil, err
	}

	return reftree, nil
}

func ReadCompTrees(inputfile string, compTrees chan<- tree.Trees) (int, error) {
	if r, err := GetReader(inputfile); err != nil {
		return 0, err
	} else {
		if i, err2 := ReadCompTreesFile(r, compTrees); err2 != nil {
		} else {
			if err = compTreeFile.Close(); err != nil {
				return id, err
			}
			return i, nil
		}
	}
}

// Read a bunch of trees from the input file. One line must define One tree.
// One tree per line
func ReadCompTreesFile(reader *io.Reader, compTrees chan<- tree.Trees) (int, error) {
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
	close(compTrees)
	return id, nil
}
