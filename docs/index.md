# Gotree: toolkit and api for phylogenetic tree manipulation
## Github repository
[Gotree github repository](https://github.com/fredericlemoine/gotree).
## Introduction
GoTree is a set of command line tools to manipulate phylogenetic trees. It is implemented in [Go](https://golang.org/) language.

The goal is to handle phylogenetic trees in [Newick](https://en.wikipedia.org/wiki/Newick_format) format, through several basic commands. Each command may print result (a tree for example) in the standard output, and thus can be piped to the standard input of the next gotree command.

## Installation
### Binaries
You can download already compiled binaries for the latest release in the [release](https://github.com/fredericlemoine/gotree/releases) section.
Binaries are available for MacOS, Linux, and Windows (32 and 64 bits).

Once downloaded, you can just run the executable without any other downloads.

### From sources
In order to compile gotree, you must first [download](https://golang.org/dl/) and [install](https://golang.org/doc/install) Go on your system.

Then you just have to type :
```bash
go get github.com/fredericlemoine/gotree/
```
This will download GoTree sources from github, and all its dependencies.

You can then build it with:
```bash
cd $GOPATH/src/github.com/fredericlemoine/gotree/
make
```
The `gotree` executable should be located in the `$GOPATH/bin` folder.

## Commands

Here is the list of all commands, with the link to the full description, and a link to a snippet that does it in GO.

Command                                                            | Subcommand        |        Description
-------------------------------------------------------------------|-------------------|-------------------------------------------------------------------------------------------------
[annotate](commands/annotate.md) ([api](api/annotate.md))          |                   | Annotates internal nodes of a tree with given data
[clear](commands/clear.md) ([api](api/clear.md))                   |                   | Clears lengths or supports from input trees
--                                                                 | lengths           | Clears lengths from input trees
--                                                                 | supports          | Clears supports from input trees
[collapse](commands/collapse.md) ([api](api/collapse.md))          |                   | Collapses/Removes branches of input trees
--                                                                 | depth             | Collapses/Removes branches of input trees having a given depth
--                                                                 | length            | Collapses/Removes short branches of input trees
--                                                                 | support           | Collapses/Removes lowly supported branches of input trees 
[compare](commands/compare.md) ([api](api/compare.md))             |                   | Compares full trees, edges, or tips
--                                                                 | edges             | Individually compares edges of the reference tree to a compared tree
--                                                                 | tips              | Compares the set of tips of the reference tree to a compared tree
--                                                                 | trees             | Compare 2 trees in terms of common and specific branches
[compute](commands/compute.md) ([api](api/compute.md))             |                   | Computations such as consensus and supports
--                                                                 | bipartitiontree   | Builds one tree with only one given bipartition
--                                                                 | consensus         | Computes the consensus from a set of input trees
--                                                                 | edgetrees         | Writes one output tree per branch of the input tree, with only one branch
--                                                                 | support classical | Computes classical bootstrap supports
--                                                                 | support booster   | Computes booster bootstrap supports
[divide](commands/divide.md)                                       |                   | Divides an input tree file into several tree files
[download](commands/download.md) ([api](api/download.md))          |                   | Downloads trees from a server
--                                                                 | itol              | Downloads a tree image from iTOL, with given image options
--                                                                 | ncbitax           | Downloads the full ncbi taxonomy from NCBI ftp server and cinverts it in Newick
[draw](commands/draw.md) ([api](api/draw.md))                      |                   | Draws tree(s) with different layouts
--                                                                 | text              | Draws tree(s) in text/ascii format
--                                                                 | png               | Draws tree(s) in png format
--                                                                 | svg               | Draws tree(s) in svg format
[generate](commands/generate.md) ([api](api/generate.md))          |                   | Generates random trees, branch lengths are simply drawn from an expontential(0.1) law
--                                                                 | balancedtree      | Randomly generates perfectly balanced trees
--                                                                 | caterpillartree   | Randomly generates perfectly caterpillar trees
--                                                                 | uniformtree       | Randomly generates uniform trees
--                                                                 | yuletree          | Randomly generates Yule-Harding trees
[matrix](commands/matrix.md) ([api](api/matrix.md))                |                   | Prints distance matrix associated to the input tree
[minbrlen](commands/minbrlen.md) ([api](api/minbrlen.md))          |                   | Sets a minimum branch length to all branches with length < cutoff
[prune](commands/prune.md) ([api](api/prune.md))                   |                   | Removes tips of the input tree that are not in the compared tree, or that are given on the command line
[randbrlen](commands/randbrlen.md)                                 |                   | Assigns a random length to edges of input trees
[randsupport](commands/randsupport.md)                             |                   | Assigns a random support to edges of input trees
[rename](commands/rename.md) ([api](api/rename.md))                |                   | Renames tips of the input tree, given a map file
[reroot](commands/reroot.md) ([api](api/reroot.md))                |                   | Reroots trees using an outgroup or at midpoint
--                                                                 | midpoint          | Reroots trees at midpoint position
--                                                                 | outgroup          | Rerootes trees using a given outgroup
[resolve](commands/resolve.md) ([api](api/resolve.md))             |                   | Resolves multifurcations by adding 0 length branches
[shuffletips](commands/shuffletips.md) ([api](api/shuffletips.md)) |                   | Shuffles tip names of an input tree
[subtree](commands/subtree.md) ([api](api/subtree.md))             |                   | Extracts a subtree starting at a given node
[stats](commands/stats.md) ([api](api/stats.md))                   |                   | Prints statistics about the tree, its edges, its nodes, if it is rooted, and its tips
--                                                                 | edges             | Prints informations about all the edges
--                                                                 | nodes             | Prints informations about all the nodes
--                                                                 | rooted            | Tells if the tree is rooted or not
--                                                                 | tips              | Prints informations about all the tips
--                                                                 | splits            | Prints all the splits/bipartitions of the tree  (bit vectors)
[unroot](commands/unroot.md) ([api](api/unroot.md))                |                   | Unroots input tree(s)
[upload](commands/upload.md) ([api](api/upload.md))                |                   | Uploads trees to a given server
--                                                                 | itol              | Uploads trees to itol, with given annotations
version                                                            |                   | Prints gotree version
