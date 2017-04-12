# Gotree: toolkit and api for phylogenetic tree manipulation

## API

### compare

Compare a reference tree to a set of compared trees
```go
package main

import (
	"fmt"
	"os"

	"github.com/fredericlemoine/gotree/io/newick"
	"github.com/fredericlemoine/gotree/io/utils"
	"github.com/fredericlemoine/gotree/tree"
)

func main() {
	var reftree *tree.Tree
	var f *os.File
	var err error
	var comptrees chan tree.Trees
	stats := make(chan tree.BipartitionStats)
	comptrees = make(chan tree.Trees)
	// Parsing multi tree newick (compared trees
	go func() {
		if _, err = utils.ReadMultiTreeFile("comp.nw", comptrees); err != nil {
			panic(err)
		}
		close(comptrees)
	}()
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
	err = tree.Compare(reftree, comptrees, false, false, 1, stats)
	// Iterating over statistic channel
	fmt.Printf("tree\treference\tcommon\tcompared\n")
	for stats := range stats {
		fmt.Printf("%d\t%d\t%d\t%d\n", stats.Id, stats.Tree1, stats.Common, stats.Tree2)
	}
}
```
