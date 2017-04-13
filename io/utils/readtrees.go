package utils

import (
	"bufio"
	"strings"

	"github.com/fredericlemoine/gotree/io/newick"
	"github.com/fredericlemoine/gotree/tree"
)

func ReadTree(inputfile string) (*tree.Tree, error) {
	if f, r, err := GetReader(inputfile); err != nil {
		return nil, err
	} else {
		if t, err2 := ReadTreeReader(r); err2 != nil {
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
func ReadTreeReader(reader *bufio.Reader) (*tree.Tree, error) {
	var reftree *tree.Tree
	var err error

	if reftree, err = newick.NewParser(reader).Parse(); err != nil {
		return nil, err
	}

	return reftree, nil
}

// Read a bunch of trees from the input reader and send each of them to the output channel
// this function does not close the reader, but closes the channel at the end of the reading.
// It returns almost immediately because parsing is performed in a go routine. Iterating over
// the tree channel will synchronize computations.
// If an error occures while parsing, it stops parsing and sends a nil tree with the error in
// the channel
func ReadMultiTrees(reader *bufio.Reader) <-chan tree.Trees {
	var compTrees chan tree.Trees = make(chan tree.Trees, 10)

	go func() {
		var err error
		var id int = 0
		var compTree *tree.Tree

		line, e := ReadUntilSemiColon(reader)
		for e == nil {
			parser := newick.NewParser(strings.NewReader(line))
			if compTree, err = parser.Parse(); err != nil {
				compTrees <- tree.Trees{
					nil,
					id,
					err,
				}
				break
			} else {
				compTrees <- tree.Trees{
					compTree,
					id,
					nil,
				}
			}
			id++
			line, e = ReadUntilSemiColon(reader)
		}
		close(compTrees)
	}()
	return compTrees
}
