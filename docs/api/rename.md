# Gotree: toolkit and api for phylogenetic tree manipulation

## API

### rename

* Renaming a set of tips from an input tree using a map:

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
	var namemap map[string]string

	// Parsing single tree newick file
	if f, err = os.Open("ref.nw"); err != nil {
		panic(err)
	}
	defer f.Close()

	t, err = newick.NewParser(f).Parse()
	if err != nil {
		panic(err)
	}

	namemap = make(map[string]string)
	namemap["Tip1"] = "NewTip1"
	namemap["Tip2"] = "NewTip2"
	namemap["Tip3"] = "NewTip3"

	err = t.Rename(namemap)
	if err != nil {
		panic(err)
	}

	fmt.Println(t.Newick())
}
```

* Renaming tips from an input tree with automatically generated names of length 10:

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
	var namemap map[string]string

	// Parsing single tree newick file
	if f, err = os.Open("ref.nw"); err != nil {
		panic(err)
	}
	defer f.Close()

	t, err = newick.NewParser(f).Parse()
	if err != nil {
		panic(err)
	}

	namemap = make(map[string]string)
	id := 1
	err = t.RenameAuto(false, true, 5, &id, namemap)
	if err != nil {
		panic(err)
	}

	fmt.Println(t.Newick())
}
```

* Renaming tips from an input tree using regexp:

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
	var namemap map[string]string

	// Parsing single tree newick file
	if f, err = os.Open("ref.nw"); err != nil {
		panic(err)
	}
	defer f.Close()

	t, err = newick.NewParser(f).Parse()
	if err != nil {
		panic(err)
	}

	namemap = make(map[string]string)
	err = t.RenameRegexp(false, true, "Tip(\\d+)", "Leaf$1", namemap)
	if err != nil {
		panic(err)
	}

	fmt.Println(t.Newick())
}
```
