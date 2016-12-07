# GoTree
GoTree is a set of command line tools to manipulate phylogenetic trees. It is implemented in [Go](https://golang.org/) language.

The goal is to handle phylogenetic trees in [Newick](https://en.wikipedia.org/wiki/Newick_format) format, through several basic commands. Each command may print result (a tree for example) in the standard output, and thus can piped to the standard input of the next gotree command.

**Example:**
```[bash]
$ gotree generate uniformtree -l 100 -n 10 | gotree stats
```
This will generate 10 random unrooted uniform binary trees, each having 100 tips, and print statistics about them, for example:

tree | nodes | tips | edges | meanbrlen | meansupport | mediansupport | rooted
-----|-------|------|-------|-----------|-------------|---------------|----------
0    | 198   | 100  | 197   | 0.0821    | -1.0000     | -1.0000       | unrooted
1    | 198   | 100  | 197   | 0.0898    | -1.0000     | -1.0000       | unrooted
2    | 198   | 100  | 197   | 0.0765    | -1.0000     | -1.0000       | unrooted
3    | 198   | 100  | 197   | 0.0746    | -1.0000     | -1.0000       | unrooted
4    | 198   | 100  | 197   | 0.0846    | -1.0000     | -1.0000       | unrooted
5    | 198   | 100  | 197   | 0.0784    | -1.0000     | -1.0000       | unrooted
6    | 198   | 100  | 197   | 0.0884    | -1.0000     | -1.0000       | unrooted
7    | 198   | 100  | 197   | 0.0943    | -1.0000     | -1.0000       | unrooted
8    | 198   | 100  | 197   | 0.0885    | -1.0000     | -1.0000       | unrooted
9    | 198   | 100  | 197   | 0.0839    | -1.0000     | -1.0000       | unrooted


## Installation
### Binaries
You can download already compiled binaries for the latest release in the [release](https://github.com/fredericlemoine/gotree/releases) section.
Binaries are available for MacOS, Linux, and Windows (32 and 64 bits).

Once downloaded, you can just run the executable without any other downloads.

### From sources
In order to compile gotree, you must first [download](https://golang.org/dl/) and [install](https://golang.org/doc/install) Go on your system.

Then you just have to type :
```
go get github.com/fredericlemoine/gotree/
```
This will download GoTree sources from github, and all its dependencies.

You can then build it with:
```
cd $GOPATH/src/github.com/fredericlemoine/gotree/
make
```
The `gotree` executable should be located in the `$GOPATH/bin` folder.

## Usage
gotree implements several tree manipulation commands. Here are some short examples:

### Generate random unrooted uniform binary trees
```[bash]
$ gotree generate uniformtree -l 100 -n 10 | gotree stats
```

### Unrooting a tree

```[bash]
$ gotree unroot -i tree.tre -o unrooted.tre
```

### Collapsing short branches

```[bash]
$ gotree collapse length -i tree.tre -l 0.001 -o collapsed.tre
```
### Collapsing lowly supported branches

```[bash]
$ gotree collapse support -i tree.tre -s 0.8 -o collapsed.tre
```

### Clearing length information

```[bash]
$ gotree clear lengths -i tree.nw -o nolength.nw
```
### Clearing support information

```[bash]
$ gotree clear supports -i tree.nw -o nosupport.nw
```
Note that you can pipe the two previous commands:

```[bash]
$ gotree clear supports -i tree.nw | gotree clear lengths -o nosupport.nw
```

### Printing tree statistics

```[bash]
$ gotree stats -i tree.tre
```
### Printing edge statistics

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

### Printing tips
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

### Comparing tips of two trees

```[bash]
$ gotree difftips -i tree.tre -c tree2.tre
```
This will compare the two sets of tips.

Example:
```
$ gotree difftips -i <(gotree generate uniformtree -l 10 -n 1) \
               -c <(gotree generate uniformtree -l 11 -n 1)
> Tip10
= 10
```
10 tips are equal, and "Tip10" is present only in the second tree.

### Removing tips that are absent from another tree

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

### Comparing bipartitions
Count the number of common/specific bipartitions between two trees.

```[bash]
$ gotree compare -i tree.tre -c other.tre
```

You can test with random trees (there should be very few common bipartitions)

```[bash]
$ gotree compare -i <(gotree generate uniformtree -l 100 -n 1) \
                -c <(gotree generate uniformtree -l 100 -n 1)
```

Tree  |  specref  |  common
------|-----------|---------
0     |  97       |  0


### Renaming tips of the tree
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
