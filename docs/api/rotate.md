# Gotree: toolkit and api for phylogenetic tree manipulation

## API

### rotate

Sort internal node neighbors by number of tips

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
	t.SortNeighborsByTips()

	fmt.Println(t.Newick())
}
```

Randomly rotate neighbors of internal nodes

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
	t.RotateInternalNodes()

	fmt.Println(t.Newick())
}
```
