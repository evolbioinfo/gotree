# Gotree: toolkit and api for phylogenetic tree manipulation

## API

### repopulate

* Adding a set of identical tips into a tree:

```go
package main

import (
	"fmt"
	"strings"

	"github.com/evolbioinfo/gotree/io/newick"
	"github.com/evolbioinfo/gotree/tree"
)

func main() {
	var t *tree.Tree
	var err error

	groups := [][]string{[]string{"Tip2", "Tip5", "Tip6", "Tip7"}, []string{"Tip4", "Tip8"}}

	// Parsing single tree from newick string
	if t, err = newick.NewParser(strings.NewReader("(Tip4:0.1,Tip0:0.1,(Tip3:0.1,(Tip2:0.2,Tip1:0.2)0.8:0.3)0.9:0.4);")).Parse(); err != nil {
		panic(err)
	}

	t.UpdateTipIndex()
	if err = t.InsertIdenticalTips(groups); err != nil {
		panic(err)
	}

	fmt.Println(t.Newick())
}
```
