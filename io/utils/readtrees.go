package utils

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/fredericlemoine/gotree/io/newick"
	"github.com/fredericlemoine/gotree/io/nexus"
	"github.com/fredericlemoine/gotree/tree"
)

const (
	FORMAT_NEWICK = iota
	FORMAT_NEXUS
)

func ReadTree(inputfile string, format int) (*tree.Tree, error) {
	if f, r, err := GetReader(inputfile); err != nil {
		return nil, err
	} else {
		if t, err2 := ReadTreeReader(r, format); err2 != nil {
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
// May take several formats: newick or nexus
// In both case, takes the first tree in the file.
func ReadTreeReader(reader *bufio.Reader, format int) (*tree.Tree, error) {
	var reftree *tree.Tree
	var err error

	switch format {
	case FORMAT_NEWICK:
		if reftree, err = newick.NewParser(reader).Parse(); err != nil {
			return nil, err
		}
	case FORMAT_NEXUS:
		if n, err2 := nexus.NewParser(reader).Parse(); err2 != nil {
			return nil, err2
		} else {
			if n.HasTrees {
				reftree = n.FirstTree()
			} else {
				return nil, fmt.Errorf("No tree in the input Nexus file")
			}
		}
	default:
		return nil, fmt.Errorf("Unsupported tree format: %q", format)
	}
	return reftree, nil
}

// Read a bunch of trees from the input reader and send each of them to the output channel
// this function does not close the reader, but closes the channel at the end of the reading.
// It returns almost immediately because parsing is performed in a go routine. Iterating over
// the tree channel will synchronize computations.
// If an error occures while parsing, it stops parsing and sends a nil tree with the error in
// the channel
// Different parsing formats: utils.FORMAT_NEWICK or utils.FORMAT_NEXUS
func ReadMultiTrees(reader *bufio.Reader, format int) <-chan tree.Trees {
	var compTrees chan tree.Trees = make(chan tree.Trees, 10)

	go func() {
		var err error
		var id int = 0
		var compTree *tree.Tree

		switch format {
		case FORMAT_NEWICK:
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
		case FORMAT_NEXUS:
			if n, err := nexus.NewParser(reader).Parse(); err != nil {
				compTrees <- tree.Trees{
					nil,
					id,
					err,
				}
			} else {
				n.IterateTrees(func(name string, t *tree.Tree) {
					compTrees <- tree.Trees{
						t,
						id,
						nil,
					}
					id++
				})
			}
		default:
			compTrees <- tree.Trees{
				nil,
				id,
				fmt.Errorf("Unsupported tree format: %q", format),
			}
		}
		close(compTrees)
	}()
	return compTrees
}
