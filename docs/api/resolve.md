# Gotree: toolkit and api for phylogenetic tree manipulation

## API

### resolve

Randomly resolving multifurcations
```go
package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/fredericlemoine/gotree/io/newick"
	"github.com/fredericlemoine/gotree/tree"
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
	t.Resolve()

	fmt.Println(t.Newick())
}
```
