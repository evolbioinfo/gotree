# Gotree: toolkit and api for phylogenetic tree manipulation

## API

### generate

Generating random trees with 1000 tips
```go
package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/fredericlemoine/gotree/tree"
)

func main() {
	var t *tree.Tree
	var err error
	var rooted bool = true
	var nbtips int = 1000

	rand.Seed(time.Now().UTC().UnixNano())

	t, err = tree.RandomYuleBinaryTree(nbtips, rooted)
	//t, err = tree.RandomBalancedBinaryTree(depth, rooted)
	//t, err = tree.RandomUniformBinaryTree(nbtips, rooted)
	//t, err = tree.RandomCaterpilarBinaryTree(nbtips, rooted)

	if err != nil {
		panic(err)
	}
	fmt.Println(t.Newick())
}
```
