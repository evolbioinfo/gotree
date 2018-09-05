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

Round branch lengths
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
	t.RoundLengths(3)
	fmt.Println(t.Newick())
	// Should print (Tip4,Tip0,(Tip3,(Tip2,Tip1)));
}
```

Cut long branches and list the connected components

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
	var bags []*tree.TipBag

	treeString = "(((1:0.1,2:0.1):0.5,((3:0.1,4:0.1):0.2,5:0.1):0.5):0.6,(6:0.1,7:0.1):0.5,(8:0.1,9:0.1):0.5);"
	t, err = newick.NewParser(strings.NewReader(treeString)).Parse()
	if err != nil {
		panic(err)
	}
	bags, err = t.CutEdgesMaxLength(0.2)
	if err != nil {
		panic(err)
	}
	for _, b := range bags {
		for i, tip := range b.Tips() {
			if i > 0 {
				fmt.Printf(",")
			}
			fmt.Printf("%s", tip.Name())
		}
		fmt.Printf("\n")
	}

}
```
