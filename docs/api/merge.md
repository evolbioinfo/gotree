# Gotree: toolkit and api for phylogenetic tree manipulation

## API

### merge

Merging two trees

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
	var t2 *tree.Tree
	var f *os.File
	var f2 *os.File
	var err error

	// Parsing First tree
	if f, err = os.Open("t1.nw"); err != nil {
		panic(err)
	}
	defer f.Close()
	if t, err = newick.NewParser(f).Parse(); err != nil {
		panic(err)
	}
	// Parsing second tree
	if f2, err = os.Open("t2.nw"); err != nil {
		panic(err)
	}
	defer f2.Close()
	if t2, err = newick.NewParser(f2).Parse(); err != nil {
		panic(err)
	}
	// Initializing tip indexes
	t.UpdateTipIndex()
	t2.UpdateTipIndex()
	// Merging both trees
	err = t.Merge(t2)
	if err != nil {
		panic(err)
	}
	fmt.Println(t.Newick())
}
```
