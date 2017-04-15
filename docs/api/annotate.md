# Gotree: toolkit and api for phylogenetic tree manipulation

## API

### annotate

Annotate function may be used to set a name to an internal node

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
	treeString = "(Tip4,Tip0,(Tip3,(Tip2,Tip1)));"
	t, err = newick.NewParser(strings.NewReader(treeString)).Parse()
	if err != nil {
		panic(err)
	}
	namemap := make(map[string][]string)
	namemap["internalnode"] = []string{"Tip1", "Tip2", "Tip3"}
	t.Annotate(namemap)
	fmt.Println(t.Newick())
	// Should print (Tip4,Tip0,(Tip3,(Tip2,Tip1))internalnode);
}
```

Also, we can use LeastCommonAncestor function
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
	treeString = "(Tip4,Tip0,(Tip3,(Tip2,Tip1)));"
	t, err = newick.NewParser(strings.NewReader(treeString)).Parse()
	if err != nil {
		panic(err)
	}

	tipnames := []string{"Tip1","Tip2","Tip3"}
	nodeindex := tree.NewNodeIndex(t)
	n, _, _, err := t.LeastCommonAncestorUnrooted(nodeindex, tipnames...)
	if err != nil {
		panic(err)
	}
	n.SetName("internalnode")
	fmt.Println(t.Newick())
	// Should print (Tip4,Tip0,(Tip3,(Tip2,Tip1))internalnode);
}
```
