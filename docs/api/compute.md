# Gotree: toolkit and api for phylogenetic tree manipulation

## API

### compute

Building one bipartition tree
```go
package main

import (
	"fmt"

	"github.com/evolbioinfo/gotree/tree"
)

func main() {
	var bipartitionTree *tree.Tree
	var err error

	rightTips := []string{"T1", "T2", "T3", "T4"}
	leftTips := []string{"T5", "T6", "T7"}

	bipartitionTree, err = tree.BipartitionTree(leftTips, rightTips)
	if err != nil {
		panic(err)
	}

	fmt.Println(bipartitionTree.Newick())
}
```


Computing consensus tree
```go
package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/evolbioinfo/gotree/io/utils"
	"github.com/evolbioinfo/gotree/tree"
)

func main() {
	var consensus *tree.Tree
	var treefile *os.File
	var treereader *bufio.Reader
	var err error
	var trees <-chan tree.Trees

	// Parsing multi tree newick
	if treefile, treereader, err = utils.GetReader("trees.nw"); err != nil {
		panic(err)
	}
	defer treefile.Close()
	trees = utils.ReadMultiTrees(treereader)

	// Computing majority consensus
	consensus, err = tree.Consensus(trees, 0.5)
	if err != nil {
		panic(err)
	}
	fmt.Println(consensus.Newick())
}
```

Computing standard bootstrap support (fbp)
```go
package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/evolbioinfo/gotree/io/newick"
	"github.com/evolbioinfo/gotree/io/utils"
	"github.com/evolbioinfo/gotree/support"
	"github.com/evolbioinfo/gotree/tree"
)

func main() {
	var reftree *tree.Tree
	var f, treefile *os.File
	var treereader *bufio.Reader
	var err error
	var trees <-chan tree.Trees

	// Parsing multi tree newick (compared trees
	if treefile, treereader, err = utils.GetReader("trees.nw"); err != nil {
		panic(err)
	}
	defer treefile.Close()
	trees = utils.ReadMultiTrees(treereader)

	// Parsing single tree newick file
	if f, err = os.Open("ref.nw"); err != nil {
		panic(err)
	}
	defer f.Close()
	reftree, err = newick.NewParser(f).Parse()
	if err != nil {
		panic(err)
	}

	// Computing fbp
	err = support.Classical(reftree, trees, 4)
	if err != nil {
		panic(err)
	}
	fmt.Println(reftree.Newick())
}
```
Computing booster support (tbe)

```go
package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/evolbioinfo/gotree/io/newick"
	"github.com/evolbioinfo/gotree/io/utils"
	"github.com/evolbioinfo/gotree/support"
	"github.com/evolbioinfo/gotree/tree"
)

func main() {
	var reftree *tree.Tree
	var f, treefile *os.File
	var treereader *bufio.Reader
	var err error
	var trees <-chan tree.Trees

	// Parsing multi tree newick (compared trees
	if treefile, treereader, err = utils.GetReader("trees.nw"); err != nil {
		panic(err)
	}
	defer treefile.Close()
	trees = utils.ReadMultiTrees(treereader)

	// Parsing single tree newick file
	if f, err = os.Open("ref.nw"); err != nil {
		panic(err)
	}
	defer f.Close()
	reftree, err = newick.NewParser(f).Parse()
	if err != nil {
		panic(err)
	}

	// Computing tbe
	err = support.Booster(reftree, trees, nil, false, false, 0, 2)
	if err != nil {
		panic(err)
	}
	fmt.Println(reftree.Newick())
}
```
