package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/fredericlemoine/gotree/io/newick"
	gotree "github.com/fredericlemoine/gotree/lib"
	"os"
	"strings"
	"sync"
)

// Type for channel of tree stats
type stats struct {
	id     int
	tree1  int
	common int
	tree2  int
}

// Type for channel of trees
type trees struct {
	tree *gotree.Tree
	id   int
}

// Readln returns a single line (without the ending \n)
// from the input buffered reader.
// An error is returned iff there is an error with the
// buffered reader.
func readln(r *bufio.Reader) (string, error) {
	var (
		isPrefix bool  = true
		err      error = nil
		line, ln []byte
	)
	for isPrefix && err == nil {
		line, isPrefix, err = r.ReadLine()
		ln = append(ln, line...)
	}
	return string(ln), err
}

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
func readCompTrees(inputfile string, compTrees chan<- trees) error {
	var compTreeFile *os.File
	var compTree *gotree.Tree
	var err error
	var reader *bufio.Reader

	if compTreeFile, err = os.Open(inputfile); err != nil {
		return err
	}

	reader = bufio.NewReader(compTreeFile)
	id := 0
	line, e := readln(reader)
	for e == nil {
		parser := newick.NewParser(strings.NewReader(line))
		if compTree, err = parser.Parse(); err != nil {
			return err
		}
		compTrees <- trees{
			compTree,
			id,
		}
		id++
		line, e = readln(reader)
	}
	close(compTrees)
	if err = compTreeFile.Close(); err != nil {
		return err
	}
	return nil
}

func main() {
	tree1 := flag.String("i", "stdin", "Input File containing the reference tree")
	tree2 := flag.String("b", "none", "Input File containing trees to compare to reference tree (mandatory)")
	tips := flag.Bool("tips", false, "If false, does not count tip edges")

	var refTree *gotree.Tree
	var err error
	var nbtips int
	var edges []*gotree.Edge

	compTreesChannel := make(chan trees, 100)
	statsChannel := make(chan stats, 100)

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
	edges = refTree.Edges()
	if nbtips, err = refTree.NbTips(); err != nil {
		panic(err)
	}

	go func() {
		if err = readCompTrees(*tree2, compTreesChannel); err != nil {
			panic(err)
		}
	}()

	var wg sync.WaitGroup
	for cpu := 0; cpu < 10; cpu++ {
		wg.Add(1)
		go func(cpu int) {
			for treeV := range compTreesChannel {
				fmt.Fprintf(os.Stderr, "Thread %d - Received new tree\n", cpu)
				var tree1, common, tree2 int
				var err error
				edges2 := treeV.tree.Edges()

				if tree1, common, tree2, err = gotree.CommonEdges(edges, edges2); err != nil {
					panic(err)
				}
				if !*tips {
					common -= nbtips
				}
				statsChannel <- stats{
					treeV.id,
					tree1,
					common,
					tree2,
				}
				fmt.Fprintf(os.Stderr, "Thread %d - Finished comparison\n", cpu)
			}
			wg.Done()
		}(cpu)
	}

	go func() {
		wg.Wait()
		close(statsChannel)
	}()

	fmt.Printf("Tree\tspecref\tcommon\tspeccomp\n")
	for stats := range statsChannel {
		fmt.Printf("%d\t%d\t%d\t%d\n", stats.id, stats.tree1, stats.common, stats.tree2)
	}
}
