# Gotree: toolkit and api for phylogenetic tree manipulation

## API

### clear

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

Clear branch supports
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
	treeString = "(Tip4:0.1,Tip0:0.1,(Tip3:0.1,(Tip2:0.2,Tip1:0.2)0.8:0.3)0.9:0.4);"
	t, err = newick.NewParser(strings.NewReader(treeString)).Parse()
	if err != nil {
		panic(err)
	}
	t.ClearSupports()
	fmt.Println(t.Newick())
	// Should print (Tip4:0.1,Tip0:0.1,(Tip3:0.1,(Tip2:0.2,Tip1:0.2):0.3):0.4)
}
```

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
