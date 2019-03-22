# Gotree: toolkit and api for phylogenetic tree manipulation

## API

### draw

Locally drawing a tree image
```go
package main

import (
	"os"

	"github.com/evolbioinfo/gotree/draw"
	"github.com/evolbioinfo/gotree/io/newick"
	"github.com/evolbioinfo/gotree/tree"
)

func main() {
	var d draw.TreeDrawer
	var l draw.TreeLayout
	var reftree *tree.Tree
	var f, outfile *os.File
	var err error

	// Parsing single tree newick file
	if f, err = os.Open("ref.nw"); err != nil {
		panic(err)
	}
	defer f.Close()

	reftree, err = newick.NewParser(f).Parse()
	if err != nil {
		panic(err)
	}

	outfile, err = os.Create("image.svg")
	defer outfile.Close()
	d = draw.NewSvgTreeDrawer(outfile, 800, 800, 30, 30, 30, 30)
	l = draw.NewRadialLayout(d, true, true, false, false)
	//l = draw.NewCircularLayout(d, true, true, false, false)
	//l = draw.NewNormalLayout(d, true, true, false, false)
	l.DrawTree(reftree)
}
```

Rendering a tree using cytoscape js
```go
package main

import (
	"bufio"
	"os"

	"github.com/evolbioinfo/gotree/draw"
	"github.com/evolbioinfo/gotree/io/newick"
	"github.com/evolbioinfo/gotree/tree"
)

func main() {
	var l draw.TreeLayout
	var reftree *tree.Tree
	var f, outfile *os.File
	var err error

	// Parsing single tree newick file
	if f, err = os.Open("ref.nw"); err != nil {
		panic(err)
	}
	defer f.Close()

	reftree, err = newick.NewParser(f).Parse()
	if err != nil {
		panic(err)
	}

	outfile, err = os.Create("tree.html")
	defer outfile.Close()
	w := bufio.NewWriter(outfile)
	l = draw.NewCytoscapeLayout(w, true)
	l.DrawTree(reftree)
	w.Flush()
}
```
