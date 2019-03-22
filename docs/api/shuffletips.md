# Gotree: toolkit and api for phylogenetic tree manipulation

## API

### shuffletips

Shuffing tip names of an input tree
```go
package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/evolbioinfo/gotree/io/newick"
	"github.com/evolbioinfo/gotree/tree"
)

func main() {
	var t *tree.Tree
	var f *os.File
	var err error

	rand.Seed(time.Now().UTC().UnixNano())

	// Parsing single tree newick file
	if f, err = os.Open("ref.nw"); err != nil {
		panic(err)
	}
	defer f.Close()

	t, err = newick.NewParser(f).Parse()
	if err != nil {
		panic(err)
	}

	t.ShuffleTips()

	fmt.Println(t.Newick())
}
```
