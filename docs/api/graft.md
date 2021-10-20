# Gotree: toolkit and api for phylogenetic tree manipulation

## API

### graft

Grafting a tree t2 on a tree t1 at tip "tip1":

```go
package api

import (
	"fmt"
	"os"

	"github.com/evolbioinfo/gotree/io"
	"github.com/evolbioinfo/gotree/io/newick"
	"github.com/evolbioinfo/gotree/tree"
)

func main() {
	var refTree, graftTree *tree.Tree
	var f, f2 *os.File
	var err error

	// Parsing ref tree newick file
	if f, err = os.Open("ref.nw"); err != nil {
		panic(err)
	}
	defer f.Close()

	refTree, err = newick.NewParser(f).Parse()
	if err != nil {
		panic(err)
	}

	// Parsing graft tree newick file
	if f2, err = os.Open("graft.nw"); err != nil {
		panic(err)
	}
	defer f2.Close()

	graftTree, err = newick.NewParser(f2).Parse()
	if err != nil {
		panic(err)
	}

	if err = refTree.UpdateTipIndex(); err != nil {
		io.LogError(err)
	}

	refTree.GraftTreeOnTip("tip1", graftTree)

	fmt.Println(refTree.Newick())

}
```
