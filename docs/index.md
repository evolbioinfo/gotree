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
[brlen](commands/brlen.md) ([api](api/brlen.md))                   |                   | Modifies branch lengths
--                                                                 | clear             | Clear lengths from input trees
--                                                                 | scale             | Scales branch lengths from input trees by a given factor
--                                                                 | setmin            | Sets a min branch length to all branches with length < cutoff
--                                                                 | setrand           | Assigns a random length to edges of input trees
[collapse](commands/collapse.md) ([api](api/collapse.md))          |                   | Collapses/Removes branches of input trees
--                                                                 | depth             | Collapses/Removes branches of input trees having a given depth
--                                                                 | length            | Collapses/Removes short branches of input trees
--                                                                 | single            | Collapses/Removes branches that connect single internal nodes (linear paths)
--                                                                 | support           | Collapses/Removes lowly supported branches of input trees 
[comment](commands/comment.md) ([api](api/comment.md))             |                   | Modifies branch/node comments
--                                                                 | clear             | Clears branch/node comments from input trees
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
--                                                                 | cyjs              | Draws tree(s) in a html file, using cytoscape js
[generate](commands/generate.md) ([api](api/generate.md))          |                   | Generates random trees, branch lengths are simply drawn from an expontential(0.1) law
--                                                                 | balancedtree      | Randomly generates perfectly balanced trees
--                                                                 | caterpillartree   | Randomly generates perfectly caterpillar trees
--                                                                 | topologies        | Generates all possible tree topologies
--                                                                 | uniformtree       | Randomly generates uniform trees
--                                                                 | yuletree          | Randomly generates Yule-Harding trees
[matrix](commands/matrix.md) ([api](api/matrix.md))                |                   | Prints distance matrix associated to the input tree
[merge](commands/merge.md) ([api](api/merge.md))                   |                   | Merges two rooted trees
[prune](commands/prune.md) ([api](api/prune.md))                   |                   | Removes tips of input trees
[reformat](commands/reformat.md) ([api](api/reformat.md))          |                   | Reformats input file
--                                                                 | newick            | Reformats input file (nexus, newick, phyloxml) into newick
--                                                                 | nexus             | Reformats input file (nexus, newick, phyloxml) into nexus
--                                                                 | phyloxml          | Reformats input file (nexus, newick, phyloxml) into phyloxml
[rename](commands/rename.md) ([api](api/rename.md))                |                   | Renames tips of the input tree
[reroot](commands/reroot.md) ([api](api/reroot.md))                |                   | Reroots trees using an outgroup or at midpoint
--                                                                 | midpoint          | Reroots trees at midpoint position
--                                                                 | outgroup          | Reroots trees using a given outgroup
[rotate](commands/rotate.md) ([api](api/rotate.md))                |                   | Reorders neighbors of internal nodes. Does not change the topology, but just traversal order.
--                                                                 | sort              | Sort neighbors of internal nodes by ascending number of tips
--                                                                 | rand              | Randomly reorders neighbors of internal nodes 
[resolve](commands/resolve.md) ([api](api/resolve.md))             |                   | Resolves multifurcations by adding 0 length branches
[sample](commands/sample.md)                                       |                   | Samples trees from a set of input trees
[shuffletips](commands/shuffletips.md) ([api](api/shuffletips.md)) |                   | Shuffles tip names of an input tree
[subtree](commands/subtree.md) ([api](api/subtree.md))             |                   | Extracts a subtree starting at a given node
[support](commands/support.md) ([api](api/support.md))             |                   | Modifies branch supports
--                                                                 | clear             | Clears branch supports from input trees
--                                                                 | scale             | Scales branch supports from input trees by a given factor
--                                                                 | setrand           | Assigns a random support to edges of input trees
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
