# Gotree: toolkit and api for phylogenetic tree manipulation

## API

### matrix

Generating distance matrix from an random tree
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
	var mat [][]float64
	var tips []*tree.Node

	rand.Seed(time.Now().UTC().UnixNano())

	t, err = tree.RandomYuleBinaryTree(nbtips, rooted)
	if err != nil {
		panic(err)
	}

	mat = t.ToDistanceMatrix()
	tips = t.Tips()

	for i, tip := range tips {
		fmt.Print(tip.Name())
		for j, _ := range tips {
			fmt.Print("\t" + fmt.Sprintf("%.12f", mat[i][j]))
		}
		fmt.Println()
	}
}
```
