# Gotree
[![build](https://travis-ci.org/fredericlemoine/gotree.svg?branch=master)](https://travis-ci.org/fredericlemoine/gotree) [![Anaconda-Server Badge](https://anaconda.org/bioconda/gotree/badges/installer/conda.svg)](https://anaconda.org/bioconda/gotree)  [![Docker hub](https://img.shields.io/docker/build/evolbioinfo/gotree.svg)](https://hub.docker.com/r/evolbioinfo/gotree/builds/)

GoTree is a set of command line tools to manipulate phylogenetic trees. It is implemented in [Go](https://golang.org/) language.

The goal is to handle phylogenetic trees in [Newick](https://en.wikipedia.org/wiki/Newick_format), Nexus and PhyloXML formats, through several basic commands. Each command may print result (a tree for example) in the standard output, and thus can be piped to the standard input of the next gotree command.

Input files may be local or remote files:

- If file name is of the form `http://<URL>`, the file is download from the given URL.
- If file name is of the form `itol://<ID>`, the tree having the given ID is downloaded from [iTOL](http://itol.embl.de/) using the iTOL api.
- If file name is of the form `treebase://<ID>`, the tree having the given ID is downloaded from [TreeBase](https://treebase.org).
- Otherwise, the file is considered local.

Gzipped input files (`.gz` extension) are supported.


**Note**:

To manipulate multiple alignments, See also [Goalign](https://github.com/fredericlemoine/goalign).

**Examples:**

```[bash]
$ echo "(1,(2,(3,4,5,6)polytomy)internal)root;" | gotree draw text --with-node-labels -w 50
+--------------- 1                                          
|                                                           
root            +---------------- 2                         
|               |                                           
+---------------|internal        +--------------- 3         
                |                |                          
                |                |--------------- 4         
                +----------------|polytomy                  
                                 |--------------- 5         
                                 |                          
                                 +--------------- 6         

```

```[bash]
$ gotree generate uniformtree -l 100 -n 10 | gotree stats

|tree  |  nodes  |  tips  |  edges  |  meanbrlen   |  sumbrlen     |  meansupport  |  mediansupport  |  rooted    |
|------|---------|--------|---------|--------------|---------------|---------------|-----------------|------------|
|0     |  198    |  100   |  197    |  0.09029828  |  17.78876078  |  NaN          |  NaN            |  unrooted  |
|1     |  198    |  100   |  197    |  0.08391711  |  16.53167037  |  NaN          |  NaN            |  unrooted  |
|2     |  198    |  100   |  197    |  0.08369861  |  16.48862662  |  NaN          |  NaN            |  unrooted  |
|3     |  198    |  100   |  197    |  0.08652623  |  17.04566698  |  NaN          |  NaN            |  unrooted  |
|4     |  198    |  100   |  197    |  0.07970206  |  15.70130625  |  NaN          |  NaN            |  unrooted  |
|5     |  198    |  100   |  197    |  0.09145831  |  18.01728772  |  NaN          |  NaN            |  unrooted  |
|6     |  198    |  100   |  197    |  0.08482117  |  16.70977068  |  NaN          |  NaN            |  unrooted  |
|7     |  198    |  100   |  197    |  0.08470308  |  16.68650662  |  NaN          |  NaN            |  unrooted  |
|8     |  198    |  100   |  197    |  0.08646811  |  17.03421732  |  NaN          |  NaN            |  unrooted  |
|9     |  198    |  100   |  197    |  0.07088132  |  13.96362091  |  NaN          |  NaN            |  unrooted  |

```
This will generate 10 random unrooted uniform binary trees, each having 100 tips, and print statistics about them.

## Installation
### Easy way: Binaries
You can download ready to run binaries for the latest release in the [release](https://github.com/fredericlemoine/gotree/releases) section.
Binaries are available for MacOS, Linux, and Windows (32 and 64 bits).

Once downloaded, you can just run the executable without any other downloads.

### Docker
Gotree Docker image is accessible from [docker hub](https://hub.docker.com/r/evolbioinfo/gotree/). You may use it as following:

```[bash]
# Display gotree help
docker run -v $PWD:$PWD -w $PWD -i -t evolbioinfo/gotree:v0.2.8b -h
```

### Singularity
Gotree [Docker image](https://hub.docker.com/r/evolbioinfo/gotree/) is usable from singularity . You may use it as following:

```[bash]
# Pull image from docker hub
singularity pull docker://evolbioinfo/gotree:v0.2.8b
# Display gotree help
./gotree-v0.2.8b.simg -h
```
### Conda
Gotree is also available on [bioconda](https://anaconda.org/bioconda/gotree). Just type:

```
conda install -c bioconda gotree
```

### From sources
In order to compile gotree, you must first [download](https://golang.org/dl/) and [install](https://golang.org/doc/install) Go on your system.

Then you just have to type :
```
go get github.com/fredericlemoine/gotree/
go get -u github.com/golang/dep/cmd/dep
```
This will download GoTree sources from github, and all its dependencies.

You can then build it with:
```
cd $GOPATH/src/github.com/fredericlemoine/gotree/
dep ensure
make
```

The `gotree` executable should be located in the `$GOPATH/bin` folder.

## Auto completion

### Bash
* Install bash-completion:
```
# MacOS
brew install bash-completion
# Linux
yum install bash-completion -y
apt-get install bash-completion
```

* Activate gotree bash completion
```
# Once
source <(gotree completion bash)
# Permanently
mkdir ~/.gotree
gotree completion bash > ~/.gotree/completion.bash.inc
printf "
# gotree shell completion
source '$HOME/.gotree/completion.bash.inc'
" >> $HOME/.bashrc
```

### Zsh (not tested)

```
# Once
source <(kubectl completion zsh)
# Permanently
gotree completion zsh > "${fpath[1]}/_gotree"
```

## Usage
gotree implements several tree manipulation commands. 

You may go to the [doc](docs/index.md) for a more detailed documentation of the commands.

### List of commands
*  annotate:    Annotate internal nodes of a tree with given data
*  brlen:       Modify branch lengths
    * clear       Clear lengths from input trees
    * scale       Scale lengths from input trees by a given factor
	* setmin      Set a min branch length to all branches with length < cutoff
	* setrand     Assign a random length to edges of input trees
*  collapse:    Collapse branches of input trees
    * depth
    * length
    * support
*  comment:     Modify branch/node comments
    * clear:    Remove node/tip comments
*  compare:     Compare full trees, edges, or tips
    * edges: Individually compare edges of the reference tree to a compared tree
    * tips: Compare the set of tips of the reference tree to a compared tree
    * trees: Compare 2 trees in terms of common and specific branches
*  compute:     Computations such as consensus and supports
    * bipartitiontree: Builds one tree with only one given bipartition
    * consensus: Compute the consensus from a set of input trees
    * edgetrees: Write one output tree per branch of the input tree, with only one branch
    * support: Compute bootstrap supports
      * classical ([Felsenstein Bootstrap](https://www.jstor.org/stable/2408678))
      * booster ([Transfer Bootstrap](http://biorxiv.org/content/early/2017/06/23/154542))
*  divide:      Divide an input tree file into several tree files
*  download:     Download a tree image from a server
    * itol: download a tree image from iTOL, with given image options
    * ncbitax: Download the full ncbi taxonomy in newick format
*  draw: Draw tree(s) with different layouts
    * text: Display tree(s) in ASCII text format
    * png : Draw tree(s) in png format, with normal, radial/unrooted or circular layout
    * svg : Draw tree(s) in svg format, with normal, radial/unrooted or circular layout
	* cyjs: Draw tree(s) in a html file, using cytoscape js
*  generate:    Generate random trees, branch lengths are simply drawn from an expontential(1) law
    * balancedtree
    * caterpillartree
	* topologies: all possible topologies
    * uniformtree
    * yuletree
*  matrix:      Print (patristic) distance matrix associated to the input tree
*  merge:       Merges two rooted trees
*  prune:       Remove tips of the input tree that are not in the compared tree, or that are given on the command line
*  reformat: Convert input file between nexus and newick formats
    * newick
    * nexus
*  rename:      Rename tips of the input tree, given a map file, or a regexp, or automatically
*  reroot:      Reroot trees using an outgroup or at midpoint
    * midpoint
    * outgroup
* rotate: Reorders neighbors of internal nodes. Does not change the topology, but just traversal order
	* rand: Randomly reorders neighbors of internal nodes 
	* sort: Sort neighbors of internal nodes by ascending number of tips
*  resolve:     Resolve multifurcations by adding 0 length branches
*  sample:      Takes a sample (with or without replacement) from the set of input trees
*  shuffletips: Shuffle tip names of an input tree
*  subtree: extract a subtree
*  support: Modify branch supports
    * clear       Clear supports from input trees
    * setrand     Assign a random support to edges of input trees
    * scale       Scale branch supports from input trees by a given factor
*  stats:       Print statistics about the tree, its edges, its nodes, if it is rooted, and its tips
    * edges
    * nodes
    * rooted
    * tips
    * splits
*  unroot:      Unroot input tree
*  upload:      Upload a tree to a given server
    * itol : Upload a tree to itol, with given annotations
*  version:     Display version of gotree

### Gotree commandline examples

* Generate 10 random unrooted uniform binary trees
```[bash]
$ gotree generate uniformtree -l 100 -n 10 | gotree stats
```

* Generate 1 Yule-Harding tree with 50 tips, and display it on the terminal (width 100)
```[bash]
$ gotree generate yuletree -l 50 | gotree draw text -w 100
```

* Generate 1 tree with 50 tips, and draw it on a SVG image
```[bash]
$ gotree generate yuletree -l 50 | gotree draw svg -w 1000 -H 1000 -o tree.svg
$ gotree generate yuletree -l 50 | gotree draw svg -w 1000 -H 1000 -r -o tree_radial.svg
```

* Reformating 10 input random trees into Nexus format:
```[bash]
$ gotree generate yuletree -n 4 -l 8 -s 10 | gotree clear lengths | gotree reformat nexus
```
Will output:
```
#NEXUS
BEGIN TAXA;
 TAXLABELS Tip4 Tip7 Tip2 Tip0 Tip3 Tip6 Tip5 Tip1;
END;
BEGIN TREES;
  TREE tree0 = ((Tip4,(Tip7,Tip2)),Tip0,(Tip3,((Tip6,Tip5),Tip1)));
  TREE tree1 = (Tip5,Tip0,((Tip6,Tip4),((Tip3,Tip2),(Tip7,Tip1))));
  TREE tree2 = (((Tip7,Tip3),(Tip4,Tip2)),Tip0,((Tip6,Tip5),Tip1));
  TREE tree3 = (Tip4,Tip0,((Tip5,Tip2),(Tip3,(Tip6,(Tip7,Tip1)))));
END;
```

* Unrooting a tree
```[bash]
$ gotree unroot -i tree.tre -o unrooted.tre
```

* Collapsing short branches
```[bash]
$ gotree collapse length -i tree.tre -l 0.001 -o collapsed.tre
```

* Collapsing lowly supported branches
```[bash]
$ gotree collapse support -i tree.tre -s 0.8 -o collapsed.tre
```

* Removing length information
```[bash]
$ gotree clear lengths -i tree.nw -o nolength.nw
```

* Removing support information
```[bash]
$ gotree clear supports -i tree.nw -o nosupport.nw
```
Note that you can pipe the two previous commands:

```[bash]
$ gotree clear supports -i tree.nw | gotree clear lengths -o nosupport.nw
```

* Printing tree statistics
```[bash]
$ gotree stats -i tree.tre
```

* Printing edge statistics
```[bash]
$ gotree stats edges -i tree.tre
```

Example of result:

tree  |  brid  |  length    |  support  |  terminal  |  depth  |  topodepth  |  rightname
------|--------|------------|-----------|------------|---------|-------------|-------------
0     |  0     |  0.107614  |  N/A      |  false     |  1      |  6          |  
0     |  1     |  0.149560  |  N/A      |  true      |  0      |  1          |  Tip51
0     |  2     |  0.051126  |  N/A      |  false     |  1      |  5          |  
0     |  3     |  0.003992  |  N/A      |  false     |  1      |  4          |  
0     |  4     |  0.030974  |  N/A      |  false     |  1      |  3          |  
0     |  5     |  0.270017  |  N/A      |  true      |  0      |  1          |  Tip84
0     |  6     |  0.029931  |  N/A      |  false     |  1      |  2          |  
0     |  7     |  0.001136  |  N/A      |  true      |  0      |  1          |  Tip70
0     |  8     |  0.011658  |  N/A      |  true      |  0      |  1          |  Tip45
0     |  9     |  0.104188  |  N/A      |  true      |  0      |  1          |  Tip34
0     |  10    |  0.003361  |  N/A      |  true      |  0      |  1          |  Tip16
0     |  11    |  0.021988  |  N/A      |  true      |  0      |  1          |  Node0

* Printing tips
```[bash]
$ gotree stats tips -i tree.tre
```
Example of result:

|tree  |  id  |  nneigh  |  name   |
|------|------|----------|---------|
|0     |  1   |  1       |  Tip8   |
|0     |  2   |  1       |  Node0  |
|0     |  5   |  1       |  Tip4   |
|0     |  8   |  1       |  Tip9   |
|0     |  9   |  1       |  Tip7   |
|0     |  11  |  1       |  Tip6   |
|0     |  13  |  1       |  Tip5   |
|0     |  14  |  1       |  Tip3   |
|0     |  16  |  1       |  Tip2   |
|0     |  17  |  1       |  Tip1   |

* Comparing tips of two trees
```[bash]
$ gotree compare tips -i tree.tre -c tree2.tre
```
This will compare the two sets of tips.

Example:
```
$ gotree compare tips -i <(gotree generate uniformtree -l 10 -n 1) \
                      -c <(gotree generate uniformtree -l 11 -n 1)
> Tip10
= 10
```
10 tips are equal, and "Tip10" is present only in the second tree.

* Removing tips that are absent from another tree
```[bash]
$ gotree prune -i tree.tre -c other.tre -o pruned.tre
```

You can test with
```[bash]
$ gotree prune -i <(gotree generate uniformtree -l 1000 -n 1) \
               -c <(gotree generate uniformtree -l 100 -n 1) \
               | gotree stats
```
It should print 100 tips.

* Comparing bipartitions
Count the number of common/specific bipartitions between two trees.

```[bash]
$ gotree compare trees -i tree.tre -c other.tre
```

You can test with random trees (there should be very few common bipartitions)

```[bash]
$ gotree compare trees -i <(gotree generate uniformtree -l 100 -n 1) \
                       -c <(gotree generate uniformtree -l 100 -n 1)
```
Tree  | reference | common | compared
------|-----------|--------|---------
   0  |     97    |    0   |    97   

* Renaming tips of the tree
If you have a file containing the mapping between current names and new names of the tips, you can rename the tips:

```[bash]
$ gotree rename -i tree.tre -m mapfile.txt -o newtree.tre
```

You can try by doing:
```[bash]
$ gotree generate uniformtree -l 100 -n 1 -o tree.tre
$ gotree stats tips -i tree.tre | awk '{if(NR>1){print $4 "\tNEWNAME" $4}}' > mapfile.txt
$ gotree rename -i tree.tre -m mapfile.txt | gotree stats tips
```

### Gotree api usage examples

* Parsing a newick string
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
	treeString = "(Tip2,Tip0,(Tip3,(Tip4,Tip1)));"
	t, err = newick.NewParser(strings.NewReader(treeString)).Parse()
	if err != nil {
		panic(err)
	}
	fmt.Println(t.Newick())
}
```

* Parsing a newick file
```go
package main

import (
	"fmt"
	"os"

	"github.com/fredericlemoine/gotree/io/newick"
	"github.com/fredericlemoine/gotree/tree"
)

func main() {

	var t *tree.Tree
	var err error
	var f *os.File
	if f, err = os.Open("t.nw"); err != nil {
		panic(err)
	}
	t, err = newick.NewParser(f).Parse()
	if err != nil {
		panic(err)
	}
	fmt.Println(t.Newick())
}
```

* Helper functions to parse multi tree newick file
```go
package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/fredericlemoine/gotree/io/utils"
	"github.com/fredericlemoine/gotree/tree"
)

func main() {
	var t tree.Trees
	var err error
	var ntrees int = 0
	var trees <-chan tree.Trees
	var treefile *os.File
	var treereader *bufio.Reader

	/* File reader (plain text or gzip) */
	if treefile, treereader, err = utils.GetReader("trees.nw"); err != nil {
		panic(err)
	}
	defer treefile.Close()
	trees = utils.ReadMultiTrees(treereader)

	for t = range trees {
		if t.Err != nil {
			panic(t.Err)
		}
		ntrees++
		fmt.Println(t.Tree.Newick())
	}
	fmt.Printf("Number of trees: %d\n", ntrees)
}
```

* Tree functions
```go
// Getting edges
var edges []*tree.Edge = t.Edges()
// Internal edges only
var iedges []*tree.Edge = t.InternalEdges()
// Tip edges only
var tedges []*tree.Edge = t.TipEdges()
// Getting Nodes
var nodes []*tree.Node = t.Nodes()
// Tips only
var tips []*tree.Node = t.Tips()
// Getting tips
var tips []*tree.Node = t.Tips()
// Getting tip names
var tipnames []string = t.AllTipNames()
// Root/Pseudoroot node
var root *tree.Node = t.Root()
// If the tree is rooted or not
var rooted bool = t.Rooted()
```

* Branch functions
```go
// Branch length
var l float64 = e.Length()
// Branch support
var s float64 = e.Support()
// Node on the "right"
var n1 *tree.Node = e.Right()
// Node on the "left"
var n2 *tree.Node = e.Left()
// Number of leaves under this edge
var nt uint = e.NumTips()
```

* Node functions
```go
// Node name
var n string = n.Name()
// Number of neighbors
var nn int = n.Nneigh()
// List of neighbors (including "parent")
var neighb []*tree.Node = n.Neigh()
// If a node is a tip or not
var tip bool = n.Tip()
// List of branches going from this node (including "parent")
var edges []*tree.Edge = n.Edges()
```

* Removing tips
```go
if err = t.RemoveTips(false, "Tip0"); err != nil {
	panic(err)
}
fmt.Println(t.Newick())
```

* Knowning if a tip exists in the tree
```go
var exists bool
var err error
exists,err = t.ExistsTip("Tip0")
```

* Shuffling tip names of the tree
```go
t.ShuffleTips()
```

* Removing branches
```go
// Short branches
t.CollapseShortBranches(0.01)
// Lowly supported branches
t.CollapseLowSupport(0.7)
// Branches with "depth" <=10 && >= 1
t.CollapseTopoDepth(1,10)
```

* Randomly resolving multifurcations
```go
t.Resolve()
```

* Removing branch informations
```go
// Branch lengths
t.ClearLengths()
// Branch supports
t.ClearSupports()
```

* Unrooting the tree
```go
t.Unroot()
```

* Cloning the tree
```go
t.Clone()
```

* Rerooting at midpoint
```go
t.RerootMidPoint()
```

* Generating random trees
```go
var ntips int = 100
var rooted bool = true
// Uniform tree
t,err = tree.RandomUniformBinaryTree(ntips, rooted)
// Balanced tree
t,err = tree.RandomBalancedBinaryTree(ntips, rooted)
// Yule-Harding tree
t,err = tree.RandomYuleBinaryTree(ntips, rooted)
```

* Computing bootstrap supports from tree files
```go
import "github.com/fredericlemoine/gotree/support"
...
var cpus int = 1
boottree,err := support.ClassicalFile("referencetreefile", "bootstraptreesfile", cpus)
```

* SVG Tree drawing 
```go
import "github.com/fredericlemoine/gotree/draw"
...
f, err := os.Create("image.svg")
d = draw.NewSvgTreeDrawer(f, 800, 800, 30, 30, 30, 30)
l = draw.NewRadialLayout(d, false, false, false, false)
// or l = draw.NewCircularLayout(d, false, false, false, false)
// or l = draw.NewNormalLayout(d, false, false, false, false)
l.DrawTree(t)
f.Close()
```

* PNG Tree drawing 
```go
import "github.com/fredericlemoine/gotree/draw"
...
f, err := os.Create("image.svg")
d = draw.NewPngTreeDrawer(f, 800, 800, 30, 30, 30, 30)
l = draw.NewRadialLayout(d, false, false, false, false)
// or l = draw.NewCircularLayout(d, false, false, false, false)
// or l = draw.NewNormalLayout(d, false, false, false, false)
l.DrawTree(t)
f.Close()
```
