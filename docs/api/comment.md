# Gotree: toolkit and api for phylogenetic tree manipulation

## API

### comment

Clear node comments
```go
package main

import (
	"fmt"
	"strings"

	"github.com/fredericlemoine/gotree/io/newick"
	"github.com/fredericlemoine/gotree/tree"
)

func main() {
	var treeString string
	var t *tree.Tree
	var err error
	treeString = "(t1[c1],t2[c2],(t3[c3],t4[c4])[c5]);"
	t, err = newick.NewParser(strings.NewReader(treeString)).Parse()
	if err != nil {
		panic(err)
	}
	t.ClearComments()
	fmt.Println(t.Newick())
	// Should print (t1,t2,(t3,t4));
}
```
