# Gotree: toolkit and api for phylogenetic tree manipulation

## API

### reformat

Parsing nexus trees and output them as newick

```go
package main

import (
	"fmt"
	"os"

	"github.com/fredericlemoine/gotree/io/nexus"
	"github.com/fredericlemoine/gotree/tree"
)

func main() {
	var f *os.File
	var err error
	if f, err = os.Open("input.nexus"); err != nil {
		panic(err)
	}
	n, err := nexus.NewParser(f).Parse()

	if err != nil {
		panic(err)
	} else {
		n.IterateTrees(func(name string, tree *tree.Tree) {
			fmt.Println(tree.Newick())
		})
	}
}
```

Parsing phyloxml trees and output them as newick

```go
package main

import (
	"fmt"
	"os"

	"github.com/fredericlemoine/gotree/io/phyloxml"
	"github.com/fredericlemoine/gotree/tree"
)

func main() {
	var f *os.File
	var err error
	var xml *phyloxml.PhyloXML
	if f, err = os.Open("bcl_2.xml"); err != nil {
		panic(err)
	}
	defer f.Close()
	p := phyloxml.NewParser(f)
	xml, err = p.Parse()
	if err != nil {
		panic(err)
	}
	xml.IterateTrees(func(t *tree.Tree, err error) {
		if err != nil {
			panic(err)
		}
		fmt.Println(t.Newick())
	})
}
```
