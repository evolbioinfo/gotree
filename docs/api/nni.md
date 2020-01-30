# Gotree: toolkit and api for phylogenetic tree manipulation

## API

### nni

Merging two trees

```go
package main

import (
	"fmt"
	"os"

	"github.com/evolbioinfo/gotree/io/newick"
	"github.com/evolbioinfo/gotree/tree"
)

func main() {
	var t *tree.Tree
	var f *os.File
	var err error

	if f, err = os.Open("t1.nw"); err != nil {
		panic(err)
	}
	defer f.Close()
	if t, err = newick.NewParser(f).Parse(); err != nil {
		panic(err)
	}

	r := &tree.NNIRearranger{}

	r.Rearrange(t, func(re tree.Rearrangement) bool {
		if err = re.Apply(); err != nil {
			return false
		}
		if err = t.CheckTreePostOrder(); err != nil {
			return false
		}
		fmt.Println(t.Newick())
		if err = re.Undo(); err != nil {
			return false
		}
		if err = t.CheckTreePostOrder(); err != nil {
			return false
		}
		return true
	})
}
```
