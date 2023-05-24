# Gotree: toolkit and api for phylogenetic tree manipulation

## Commands

### compare
This command compares a reference tree -given with `-i` with a set of compared trees given with `-c`. Three subcommands :
* `gotree compare edges`: Compares each edges/branches of the reference tree to all compared trees, by giving the following informations in a tab-separated format:
 1. Compared tree index;
 2. Reference branch id;
 3. Length of reference branch;
 4. Bootstrap Support of reference branch;
 5. "true" if terminal branch, "false" otherwise;
 6. Depth1: Minimum path to a tip;
 7. Depth2: Number of tips in the lightest side of the branch;
 8. RightName: Name of the node on the right of the branch;
 9. "true" if the branch is present in the compared tree, "false" otherwise;
 10. if `-m` is given : transfer distance between the reference branch and its closest branch of the compared tree;
 11. if `-m` and `--moved-taxa` are given: List of taxa to move from left to right, and from right to left, to go from the reference branch to its closest branch of the compared tree.
 12. Name of the matching node in the compared tree if any (best match if -m is given of exact match otherwise). If the tree is rooted, the node name is the name of the descendent node. Otherwise the node name is the name of the node on the lightest side of the matching  bipartition.

* `gotree compare tips`: Compares the set of tips of the reference tree with the set of tips of all the compared trees, in the manner of unix diff. Output:
  * For each missing tip in the compared tree, will print: `(Tree <id>) < TipName`,
  * For each missing tip in the reference tree, will print: `(Tree <id>) > TipName`,
  * Then print the number of common tips: `(Tree <id>) = <nb common>`,
* `gotree compare trees`: Compares the reference tree with all the compared trees, in terms of common bi-partitions. If `--rf` option is given, only "number of branches specific to reference tree" + "number of branches specific to compared tree" is given. Otherwise, the output is tab separated with the following columns:
 1. Compared tree index;
 2. Number of branches specific to the reference tree;
 3. Number of common branches between reference and compared trees;
 4. Number of branches specific to the compared tree.

  If the `--weighted` option is given the trees will be compared with branch length aware metrics. The output is tab separated with the following columns:
  1. Compared tree index;
  2. Weighted Robinson-Foulds distance [(Robinson & Foulds, 1979)](https://doi.org/10.1007/BFb0102690);
  3. Khuner-Felsenstein distance [(Khuner & Felsenstein, 1994)](https://doi.org/10.1093/oxfordjournals.molbev.a040126);

#### Usage

General command
```
Usage:
  gotree compare [command]

Available Commands:
  edges       Compare edges of a reference tree with another tree
  tips        Print diff between tip names of two trees
  trees       Compare a reference tree with a set of trees

Flags:
  -c, --compared string   Compared trees input file (default "none")
  -i, --reftree string    Reference tree input file (default "stdin")
```

edges sub-command
```
Usage:
  gotree compare edges [flags]

Flags:
      --moved-taxa      only if --transfer-dist is given: Then display, for each branch, taxa that must be moved
  -m, --transfer-dist   If transfer dist must be computed for each edge

Global Flags:
  -c, --compared string   Compared trees input file (default "none")
  -i, --reftree string    Reference tree input file (default "stdin")
```

tips sub-command
```
Usage:
  gotree compare tips [flags]

Global Flags:
  -c, --compared string   Compared trees input file (default "none")
  -i, --reftree string    Reference tree input file (default "stdin")
```

trees sub-command
```
Usage:
  gotree compare trees [flags]

Flags:
      --binary   If true, then just print true (identical tree) or false (different tree) for each compared tree
  -l, --tips     Include tips in the comparison
  --rf           If true, outputs Robinson-Foulds distance, as the sum of reference + compared specific branches

Global Flags:
  -c, --compared string   Compared trees input file (default "none")
  -i, --reftree string    Reference tree input file (default "stdin")
```

#### Examples

1. Comparing edges

```
gotree compare edges -i <(gotree generate yuletree --seed 10) -c <(gotree generate yuletree --seed 12 -n 1) -m --moved-taxa
```

Should give:

tree|brid| length | support |terminal|depth|topodepth|rightname|found|transfer|taxatomove
----|----|--------|---------|--------|-----|---------|---------|-----|--------|------------
0   |0   | 0.0912 |   N/A   |false   |1    |3        |         |false|2       |-Tip2,-Tip7
0   |1   | 0.0206 |   N/A   |true    |0    |1        |Tip4     |true |0       |
0   |2   | 0.1293 |   N/A   |false   |1    |2        |         |false|1       |-Tip7
0   |3   | 0.0974 |   N/A   |true    |0    |1        |Tip7     |true |0       |
0   |4   | 0.0154 |   N/A   |true    |0    |1        |Tip2     |true |0       |
0   |5   | 0.1295 |   N/A   |true    |0    |1        |Tip0     |true |0       |
0   |6   | 0.0229 |   N/A   |false   |1    |4        |         |false|3       |+Tip0,+Tip2,+Tip7
0   |7   | 0.0960 |   N/A   |false   |1    |3        |         |false|1       |-Tip3
0   |8   | 0.0278 |   N/A   |true    |0    |1        |Tip8     |true |0       |
0   |9   | 0.0051 |   N/A   |false   |1    |2        |         |false|1       |-Tip9
0   |10  | 0.1349 |   N/A   |true    |0    |1        |Tip9     |true |0       |
0   |11  | 0.1030 |   N/A   |true    |0    |1        |Tip3     |true |0       |
0   |12  | 0.1507 |   N/A   |false   |1    |3        |         |false|2       |-Tip1,-Tip5
0   |13  | 0.0290 |   N/A   |false   |1    |2        |         |false|1       |-Tip5
0   |14  | 0.3779 |   N/A   |true    |0    |1        |Tip6     |true |0       |
0   |15  | 0.1120 |   N/A   |true    |0    |1        |Tip5     |true |0       |
0   |16  | 0.2390 |   N/A   |true    |0    |1        |Tip1     |true |0       |

2. Comparing tips

```
gotree compare tips -i <(gotree generate yuletree --seed 10) -c <(gotree generate yuletree --seed 12 -n 1 -l 12)
```

Should give:

```
(Tree 0) > Tip11
(Tree 0) > Tip10
(Tree 0) = 10
```

3. Comparing trees

```
gotree compare trees -i <(gotree generate yuletree --seed 10) -c <(gotree generate yuletree --seed 12 -n 1)
```

|tree  |  reference  |  common  |  compared  |
|------|-------------|----------|------------|
|0     |  7          |  0       |  7         |

If we want to take the branch lengths into account when comparing the trees we can
specify the `--weighted` flag:

```
gotree compare trees --weighted -i <(gotree generate yuletree --seed 10) -c <(gotree generate yuletree --seed 12 -n 1)
```

| tree | weighted_RF  | KF           |
|------|--------------|--------------|
|0     | 1.310593E+00 | 5.056856E-01 |