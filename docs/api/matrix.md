# Gotree: toolkit and api for phylogenetic tree manipulation

## API

### matrix

Generating distance matrix from an input tree
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
	
	rand.Seed(time.Now().UTC().UnixNano())

	t, err = tree.RandomYuleBinaryTree(nbtips, rooted)

	mat = t.ToDistanceMatrix()
	for i, t := range tips {
		f.WriteString(t.Name())
		for j, _ := range tips {
			f.WriteString("\t" + fmt.Sprintf("%.12f", mat[i][j]))
		}
		f.WriteString("\n")
	}
}
```
