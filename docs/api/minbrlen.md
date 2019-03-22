# Gotree: toolkit and api for phylogenetic tree manipulation

## API

### minbrlen

Setting a minimun branch length to an input tree
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

	for _, e := range t.Edges() {
		if e.Length() < 0.01 {
			e.SetLength(0.01)
		}
	}
	fmt.Println(t.Newick())
}
```
