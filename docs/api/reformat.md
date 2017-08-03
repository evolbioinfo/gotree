# Gotree: toolkit and api for phylogenetic tree manipulation

## API

### prune

Parsing nexus trees and output them as newick
```go
package main

import (
	"fmt"
	"os"

	"github.com/fredericlemoine/gotree/io/nexus"
	"github.com/fredericlemoine/gotree/tree"
)

func main() {
	var f *os.File
	var err error
	if f, err = os.Open("input.nexus"); err != nil {
		panic(err)
	}
	n, err := nexus.NewParser(f).Parse()

	if err != nil {
		panic(err)
	} else {
		n.IterateTrees(func(name string, tree *tree.Tree) {
			fmt.Println(tree.Newick())
		})
	}
}
```
