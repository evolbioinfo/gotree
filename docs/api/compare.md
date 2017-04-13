# Gotree: toolkit and api for phylogenetic tree manipulation

## API

### compare

Comparing a reference tree to a set of compared trees
```go
package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/fredericlemoine/gotree/io/newick"
	"github.com/fredericlemoine/gotree/io/utils"
	"github.com/fredericlemoine/gotree/tree"
)

func main() {
	var reftree *tree.Tree
	var f, treefile *os.File
	var treereader *bufio.Reader
	var err error
	var trees <-chan tree.Trees
	var stats <-chan tree.BipartitionStats

	// Parsing multi tree newick (compared trees
	if treefile, treereader, err = utils.GetReader("trees.nw"); err != nil {
		panic(err)
	}
	defer treefile.Close()
	trees = utils.ReadMultiTrees(treereader)
	// Parsing single tree newick file
	if f, err = os.Open("ref.nw"); err != nil {
		panic(err)
	}
	reftree, err = newick.NewParser(f).Parse()
	if err != nil {
		panic(err)
	}
	f.Close()
	// Comparing reftree with all comp trees
	stats, err = tree.Compare(reftree, trees, false, false, 1)
	// Iterating over statistic channel
	fmt.Printf("tree\treference\tcommon\tcompared\n")
	for stats := range stats {
		if stats.Err != nil {
			panic(err)
		}
		fmt.Printf("%d\t%d\t%d\t%d\n", stats.Id, stats.Tree1, stats.Common, stats.Tree2)
	}
}
```

Comparing reference tree edges to a set of compared trees
```go
package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/fredericlemoine/gotree/io/newick"
	"github.com/fredericlemoine/gotree/io/utils"
	"github.com/fredericlemoine/gotree/tree"
)

func main() {
	var reftree *tree.Tree
	var f, treefile *os.File
	var treereader *bufio.Reader
	var err error
	var trees <-chan tree.Trees
	var refEdges []*tree.Edge
	var edgeIndex *tree.EdgeIndex
	// Parsing multi tree newick (compared trees
	if treefile, treereader, err = utils.GetReader("trees.nw"); err != nil {
		panic(err)
	}
	defer treefile.Close()
	trees = utils.ReadMultiTrees(treereader)
	// Parsing single tree newick file
	if f, err = os.Open("ref.nw"); err != nil {
		panic(err)
	}
	reftree, err = newick.NewParser(f).Parse()
	if err != nil {
		panic(err)
	}
	f.Close()
	// Building reference edge index
	refEdges = reftree.Edges()
	edgeIndex = tree.NewEdgeIndex(int64(len(refEdges)*2), 0.75)
	for _, e := range refEdges {
		edgeIndex.PutEdgeValue(e, -1, -1)
	}
	// All trees
	for t2 := range trees {
		if t2.Err != nil {
			panic(t2.Err)
		}
		// All edges
		for j, e := range t2.Tree.Edges() {
			// We check if the edge is present in the index
			_, ok := edgeIndex.Value(e)
			fmt.Printf("tree %d | branch %d | #tips %v | found %v\n", t2.Id, j, e.NumTips(), ok)
		}
	}
}
```
