# Gotree: toolkit and api for phylogenetic tree manipulation

## API

### collapse

Collapse short branches
```go
package main

import (
	"fmt"
	"strings"

	"github.com/evolbioinfo/gotree/io/newick"
	"github.com/evolbioinfo/gotree/tree"
)

func main() {
	var treeString string
	var t *tree.Tree
	var err error
	treeString = "(Tip4:0.1,Tip0:0.1,(Tip3:0.1,(Tip2:0.2,Tip1:0.2):0.001):0.4);"
	t, err = newick.NewParser(strings.NewReader(treeString)).Parse()
	if err != nil {
		panic(err)
	}
	t.CollapseShortBranches(0.01)
	fmt.Println(t.Newick())
	// Should print (Tip4:0.1,Tip0:0.1,(Tip3:0.1,Tip2:0.2,Tip1:0.2):0.4);
}
```

Collapse lowly supported branches
```go
package main

import (
	"fmt"
	"strings"

	"github.com/evolbioinfo/gotree/io/newick"
	"github.com/evolbioinfo/gotree/tree"
)

func main() {
	var treeString string
	var t *tree.Tree
	var err error
	treeString = "(Tip4,Tip0,(Tip3,(Tip2,Tip1)0.2)0.9);"
	t, err = newick.NewParser(strings.NewReader(treeString)).Parse()
	if err != nil {
		panic(err)
	}
	t.CollapseLowSupport(0.7)
	fmt.Println(t.Newick())
	// Should print (Tip4,Tip0,(Tip3,Tip2,Tip1)0.9);
}
```

Remove external branches, having depth between 2 (cherries) and 3
```go
package main

import (
	"fmt"
	"strings"

	"github.com/evolbioinfo/gotree/io/newick"
	"github.com/evolbioinfo/gotree/tree"
)

func main() {
	var treeString string
	var t *tree.Tree
	var err error
	treeString = "(Tip4,Tip0,(Tip3,(Tip2,Tip1)0.2)0.9);"
	t, err = newick.NewParser(strings.NewReader(treeString)).Parse()
	if err != nil {
		panic(err)
	}
	t.CollapseTopoDepth(2, 3)
	fmt.Println(t.Newick())
	// Should print (Tip4,Tip0,Tip3,Tip2,Tip1);
}
```

Collapse a clade

```go
package main

import (
	"fmt"
	"strings"

	"github.com/evolbioinfo/gotree/io/newick"
	"github.com/evolbioinfo/gotree/tree"
)

func main() {
	var treeString string
	var t *tree.Tree
	var err error
	treeString = "((l4:0.1,((l5:0.2,l7:0.3):0.3,l6:0.3):0.1):0.1,(l2:0.2,l3:0.3):0.1);"
	t, err = newick.NewParser(strings.NewReader(treeString)).Parse()
	if err != nil {
		panic(err)
	}
	t.CollapseClade(false, "l1", "l4", "6")
	fmt.Println(t.Newick())
}
```
