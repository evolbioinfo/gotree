# Gotree: toolkit and api for phylogenetic tree manipulation

## API

### upload

#### Upload a tree to [iTOL](https://itol.embl.de/)

	
```go
package main

import (
	"fmt"
	"os"

	"github.com/evolbioinfo/gotree/io/utils"
	"github.com/evolbioinfo/gotree/tree"
	"github.com/evolbioinfo/gotree/upload"
)

func main() {
	var reftree *tree.Tree
	var err error

	upld := upload.NewItolUploader("", "")

	if reftree, err = utils.ReadTree("tree.nhx", utils.FORMAT_NEWICK); err != nil {
		panic(err)
	}

	url, response, err := upld.Upload("Tree_name", reftree)
	if err != nil {
		panic(err)
	}

	// URL to iTOL visualization
	fmt.Println(url)
	// Response from iTOL server
	fmt.Fprintf(os.Stderr, "-------------------\n")
	fmt.Fprintf(os.Stderr, "<Server response>\n")
	fmt.Fprintf(os.Stderr, response)
	fmt.Fprintf(os.Stderr, "-------------------\n")
}
```
