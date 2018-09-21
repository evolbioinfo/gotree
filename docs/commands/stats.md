# Gotree: toolkit and api for phylogenetic tree manipulation

## Commands

### stats
This commands display informations about trees, edges, nodes or tips. Several subcommands:
* `gotree stats` Without subcommand: Display informations about input trees, in tab delimited format, with columns:
   1. Tree id (input file order)
   2. Number of nodes (including tips)
   3. Number of tips
   4. Number of edges
   5. Average branch length
   6. Sum of all branch lengths
   7. Average bootstrap support
   8. Median bootstrap support
   9. Rooted: true/false
* `gotree stats edges` : Display informations about edges of input trees, in tab delimited format, with columns:
   1. Tree id (input file order)
   2. Branch id (newick parsing order)
   3. Branch length
   4. Branch support if any
   5. Terminal branch or not : true/false
   6. Depth 1: length of the shortest path to a tip
   7. Depth 2: number of tips on the lightest side of the branch
   8. Name of the node on the right (tip name if it is a terminal branch)
   
* `gotree stats nodes` : Display informations about nodes of input trees, in tab delimited format, with columns:
   1. Tree id (input file order)
   2. Node id (newick parsing order)
   3. Number of neighbors of the node (3 if internal node without multifurcation, 1 if a tip)
   4. Name of the node (Tip name if tip, Internal node name if any)
   5. Depth of the node: length of the shortest path to a tip
   6. Comments associated to tree nodes (in the form `(1,2,3)[comment]` in Newick format)
   
* `gotree stats rooted` : Tells if each input tree is rooted or not, in tab delimited format, with columns:
   1. Tree id (input file order)
   2. rooted : true/false

* `gotree stats splits` : Displays each branches of input trees in binary vector format, in tab delimited format, with columns:
   1. Tree id (input file order)
   2. binary vector

* `gotree stats tips` : Displays informations about tips of input trees, in tab delimited format, with columns:
   1. Tree id (input file order)
   2. tip id (same as its node id in `gotree stats nodes`)
   3. Number of neighbors (always 1...)
   4. Name of the tip

#### Usage

General command
```
Usage:
  gotree stats [flags]
  gotree stats [command]

Available Commands:
  edges       Displays statistics on edges of input tree
  nodes       Displays statistics on nodes of input tree
  rooted      Tells wether the tree is rooted or unrooted
  splits      Prints all the splits from an input tree
  tips        Displays statistics on tips of input tree

Flags:
  -i, --input string    Input tree (default "stdin")
  -o, --output string   Output file (default "stdout")
```

#### Examples

* Generate a random tree and display informations about it

```
gotree generate yuletree --seed 10 | gotree stats
gotree generate yuletree --seed 10 | gotree stats edges
```

Should give

|tree  |  nodes  |  tips  |  edges  |  meanbrlen   |  sumbrlen    |  meansupport  |  mediansupport  |  rooted    |
|------|---------|--------|---------|--------------|--------------|---------------|-----------------|------------|
|0     |  18     |  10    |  17     |  0.10486138  |  1.78264354  |  N/A          |  N/A            |  unrooted  |

and

|tree  |  brid  |  length                |  support  |  terminal  |  depth  |  topodepth  |  rightname  |
|------|--------|------------------------|-----------|------------|---------|-------------|-------------|
|0     |  0     |  0.0912341925030609    |  N/A      |  false     |  1      |  3          |             |
|0     |  1     |  0.020616211789029896  |  N/A      |  true      |  0      |  1          |  Tip4       |
|0     |  2     |  0.12939642466438622   |  N/A      |  false     |  1      |  2          |             |
|0     |  3     |  0.09740195047110385   |  N/A      |  true      |  0      |  1          |  Tip7       |
|0     |  4     |  0.015450672710905129  |  N/A      |  true      |  0      |  1          |  Tip2       |
|0     |  5     |  0.12959932895259058   |  N/A      |  true      |  0      |  1          |  Tip0       |
|0     |  6     |  0.022969404523534506  |  N/A      |  false     |  1      |  4          |             |
|0     |  7     |  0.09604804621401375   |  N/A      |  false     |  1      |  3          |             |
|0     |  8     |  0.027845992087631298  |  N/A      |  true      |  0      |  1          |  Tip8       | 
|0     |  9     |  0.005132906169455565  |  N/A      |  false     |  1      |  2          |             |
|0     |  10    |  0.13492605122032592   |  N/A      |  true      |  0      |  1          |  Tip9       |
|0     |  11    |  0.10309294031874587   |  N/A      |  true      |  0      |  1          |  Tip3       |
|0     |  12    |  0.15075207292513051   |  N/A      |  false     |  1      |  3          |             |
|0     |  13    |  0.029087690784364996  |  N/A      |  false     |  1      |  2          |             |
|0     |  14    |  0.3779897840448691    |  N/A      |  true      |  0      |  1          |  Tip6       |
|0     |  15    |  0.1120177846434196    |  N/A      |  true      |  0      |  1          |  Tip5       |
|0     |  16    |  0.239082088939295     |  N/A      |  true      |  0      |  1          |  Tip1       |

