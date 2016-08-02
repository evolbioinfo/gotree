package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/fredericlemoine/gotree/io/newick"
	gotree "github.com/fredericlemoine/gotree/lib"
	"os"
	"strings"
)

func readRefTree(inputfile string) (*gotree.Tree, error) {
	var refTreeFile *os.File
	var reftree *gotree.Tree
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
func readCompTrees(inputfile string) ([]*gotree.Tree, error) {
	var compTreeFile *os.File
	var compTrees []*gotree.Tree
	var compTree *gotree.Tree
	var err error
	var scanner *bufio.Scanner

	if compTreeFile, err = os.Open(inputfile); err != nil {
		return nil, err
	}

	scanner = bufio.NewScanner(compTreeFile)
	for scanner.Scan() {
		parser := newick.NewParser(strings.NewReader(scanner.Text()))
		if compTree, err = parser.Parse(); err != nil {
			return nil, err
		}
		compTrees = append(compTrees, compTree)
	}

	if err = compTreeFile.Close(); err != nil {
		return nil, err
	}

	return compTrees, nil
}

func main() {
	tree1 := flag.String("i", "stdin", "Input File containing the reference tree")
	tree2 := flag.String("b", "none", "Input File containing trees to compare to reference tree (mandatory)")
	tips := flag.Bool("tips", false, "If false, does not count tip edges")

	var refTree *gotree.Tree
	var err error
	var compTrees []*gotree.Tree
	var commonEdges []int
	var tree1Edges []int
	var tree2Edges []int

	flag.Parse()

	if *tree2 == "none" {
		panic("You must provide a file containing compared trees")
	}

	fmt.Fprintf(os.Stderr, "Reference : %s\n", *tree1)
	fmt.Fprintf(os.Stderr, "Compared  : %s\n", *tree2)
	fmt.Fprintf(os.Stderr, "With tips : %b\n", *tips)

	if refTree, err = readRefTree(*tree1); err != nil {
		panic(err)
	}

	if compTrees, err = readCompTrees(*tree2); err != nil {
		panic(err)
	}

	tree1Edges, commonEdges, tree2Edges, _ = refTree.CompareEdges(compTrees, *tips)

	fmt.Printf("Tree\tspecref\tcommon\tspeccomp\n")
	for i := 0; i < len(commonEdges); i++ {
		fmt.Printf("%d\t%d\t%d\t%d\n", i, tree1Edges[i], commonEdges[i], tree2Edges[i])
	}
}
