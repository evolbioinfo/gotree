# Gotree: toolkit and api for phylogenetic tree manipulation

## API

### brlen

Clear branch lengths
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
	treeString = "(Tip4:0.1,Tip0:0.1,(Tip3:0.1,(Tip2:0.2,Tip1:0.2):0.3):0.4);"
	t, err = newick.NewParser(strings.NewReader(treeString)).Parse()
	if err != nil {
		panic(err)
	}
	t.ClearLengths()
	fmt.Println(t.Newick())
	// Should print (Tip4,Tip0,(Tip3,(Tip2,Tip1)));
}
```

