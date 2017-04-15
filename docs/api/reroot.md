# Gotree: toolkit and api for phylogenetic tree manipulation

## API

### reroot

Rerooting a tree at midpoint position
```go
package main

import (
	"fmt"
	"os"

	"github.com/fredericlemoine/gotree/io/newick"
	"github.com/fredericlemoine/gotree/tree"
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
	t.RerootMidPoint()

	fmt.Println(t.Newick())
}
```

Rerooting a tree using an outgroup

```go
package main

import (
	"fmt"
	"os"

	"github.com/fredericlemoine/gotree/io/newick"
	"github.com/fredericlemoine/gotree/tree"
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

	err = t.RerootOutGroup("Tip1", "Tip2", "Tip3")
	if err != nil {
		panic(err)
	}

	fmt.Println(t.Newick())
}
```
