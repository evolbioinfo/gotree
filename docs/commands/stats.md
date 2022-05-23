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
   7. Names of parent node (for unrooted tree, depend on the way the tree is written)
   8. Names of children nodes (for unrooted tree, depend on the way the tree is written) 

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
   5. Length of the external branch leading to the tip
   6. Sum of branch lengths from the root to the tip

* `gotree stats monophyletic` : Tells wether a set of tips form a monophyletic clade in the given trees. Output is in tab delimited format, with columns:
   1. Tree id (input file order)
   2. Monophyletic (true/false)

#### Usage

General command
```
Usage:
  gotree stats [flags]
  gotree stats [command]

Available Commands:
  edges        Displays statistics on edges of input tree
  monophyletic Tells wether input tips form a monophyletic group in each of the input trees
  nodes        Displays statistics on nodes of input tree
  rooted       Tells wether the tree is rooted or unrooted
  splits       Prints all the splits from an input tree
  tips         Displays statistics on tips of input tree

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

|tree  |  brid  |  length                |  support  |  terminal  |  depth  |  topodepth  |  rootdepth  |  rightname  |
|------|--------|------------------------|-----------|------------|---------|-------------|-------------|-------------|
|0     |  0     |  0.0912341925030609    |  N/A      |  false     |  1      |  3          |  -1         |             |
|0     |  1     |  0.020616211789029896  |  N/A      |  true      |  0      |  1          |  -1         |  Tip4       |
|0     |  2     |  0.12939642466438622   |  N/A      |  false     |  1      |  2          |  -1         |             |
|0     |  3     |  0.09740195047110385   |  N/A      |  true      |  0      |  1          |  -1         |  Tip7       |
|0     |  4     |  0.015450672710905129  |  N/A      |  true      |  0      |  1          |  -1         |  Tip2       |
|0     |  5     |  0.12959932895259058   |  N/A      |  true      |  0      |  1          |  -1         |  Tip0       |
|0     |  6     |  0.022969404523534506  |  N/A      |  false     |  1      |  4          |  -1         |             |
|0     |  7     |  0.09604804621401375   |  N/A      |  false     |  1      |  3          |  -1         |             |
|0     |  8     |  0.027845992087631298  |  N/A      |  true      |  0      |  1          |  -1         |  Tip8       | 
|0     |  9     |  0.005132906169455565  |  N/A      |  false     |  1      |  2          |  -1         |             |
|0     |  10    |  0.13492605122032592   |  N/A      |  true      |  0      |  1          |  -1         |  Tip9       |
|0     |  11    |  0.10309294031874587   |  N/A      |  true      |  0      |  1          |  -1         |  Tip3       |
|0     |  12    |  0.15075207292513051   |  N/A      |  false     |  1      |  3          |  -1         |             |
|0     |  13    |  0.029087690784364996  |  N/A      |  false     |  1      |  2          |  -1         |             |
|0     |  14    |  0.3779897840448691    |  N/A      |  true      |  0      |  1          |  -1         |  Tip6       |
|0     |  15    |  0.1120177846434196    |  N/A      |  true      |  0      |  1          |  -1         |  Tip5       |
|0     |  16    |  0.239082088939295     |  N/A      |  true      |  0      |  1          |  -1         |  Tip1       |

As the tree is unrooted, `rootdepth` is set to -1. However, with a rooted tree:

```
gotree generate yuletree -r --seed 10 | gotree stats edges
```

|tree | brid | length                | support | terminal | depth | topodepth | rootdepth | rightname | comments | leftname | rightcomment | leftcomment |
|-----|------|-----------------------|---------|----------|-------|-----------|-----------|-----------|----------|----------|--------------|-------------|
|0    | 0    | 0.054743875470795914  | N/A     | false    | 2     | 2         | 1         |           | []       |          |    []        |   []        |
|0    | 1    | 0.13492605122032592   | N/A     | false    | 1     | 2         | 2         |           | []       |          |    []        |   []        |
|0    | 2    | 0.10309294031874587   | N/A     | true     | 0     | 1         | 3         | Tip9      | []       |          |    []        |   []        |
|0    | 3    | 0.03707446096764306   | N/A     | true     | 0     | 1         | 3         | Tip2      | []       |          |    []        |   []        |
|0    | 4    | 0.13604994737755394   | N/A     | false    | 1     | 4         | 2         |           | []       |          |    []        |   []        |
|0    | 5    | 0.19852695409349608   | N/A     | true     | 0     | 1         | 3         | Tip3      | []       |          |    []        |   []        |
|0    | 6    | 0.020616211789029896  | N/A     | false    | 1     | 5         | 3         |           | []       |          |    []        |   []        |
|0    | 7    | 0.08184535681853511   | N/A     | false    | 1     | 4         | 4         |           | []       |          |    []        |   []        |
|0    | 8    | 0.3779897840448691    | N/A     | false    | 1     | 3         | 5         |           | []       |          |    []        |   []        |
|0    | 9    | 0.027845992087631298  | N/A     | false    | 1     | 2         | 6         |           | []       |          |    []        |   []        |
|0    | 10   | 0.0440885662122905    | N/A     | true     | 0     | 1         | 7         | Tip8      | []       |          |    []        |   []        |
|0    | 11   | 0.14809735366802398   | N/A     | true     | 0     | 1         | 7         | Tip6      | []       |          |    []        |   []        |
|0    | 12   | 0.18347097513974125   | N/A     | true     | 0     | 1         | 6         | Tip5      | []       |          |    []        |   []        |
|0    | 13   | 0.03199874235185574   | N/A     | true     | 0     | 1         | 5         | Tip4      | []       |          |    []        |   []        |
|0    | 14   | 0.10033210749794116   | N/A     | true     | 0     | 1         | 4         | Tip1      | []       |          |    []        |   []        |
|0    | 15   | 0.09740195047110385   | N/A     | false    | 1     | 2         | 1         |           | []       |          |    []        |   []        |
|0    | 16   | 0.015450672710905129  | N/A     | true     | 0     | 1         | 2         | Tip7      | []       |          |    []        |   []        |
|0    | 17   | 0.17182241382980687   | N/A     | true     | 0     | 1         | 2         | Tip0      | []       |          |    []        |   []        |


Here `rootdepth` gives the number of branches from the current branch (included) to the root.

* Check wether a set of tips form a monophyletic clade

tips.txt:
```
T1
T2
T3
```

tree.nw
```
((T1,T2),(T3,T4),(T5,T6,(T7,T8)));
```

Then
```
> gotree stats monophyletic -i tree.nw -l tips.txt
Tree	Monophyletic
0	false
```

Or
```
> echo "((T1,T2),(T3,T4),(T5,T6,(T7,T8)));" | gotree stats monophyletic T5 T6 T7 T8
Tree	Monophyletic
0	true
```
