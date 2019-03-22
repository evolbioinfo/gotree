# Gotree: toolkit and api for phylogenetic tree manipulation

## API

### prune

Removing a set of tips from an input tree
```go
package main

import (
	"fmt"
	"os"

	"github.com/evolbioinfo/gotree/io/newick"
	"github.com/evolbioinfo/gotree/tree"
)

func main() {
	var t *tree.Tree
	var f *os.File
	var err error

	// Parsing single tree newick file
	if f, err = os.Open("ref.nw"); err != nil {
		panic(err)
	}
	defer f.Close()

	t, err = newick.NewParser(f).Parse()
	if err != nil {
		panic(err)
	}

	err = r.RemoveTips(false, "Tip1","Tip2","Tip3")
	if err != nil {
		panic(err)
	}
	
	fmt.Println(t.Newick())
}
```
