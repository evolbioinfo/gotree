# Gotree: toolkit and api for phylogenetic tree manipulation

## Commands

### divide
This command divides a multi tree input file into several output one tree files.

#### Usage

```
gotree divide -i trees.nw -o prefix_

Usage:
  gotree divide [flags]

Flags:
  -i, --input string    Input tree(s) file (default "stdin")
  -o, --output string   Divided trees output file prefix (default "prefix")
```

#### Example

We generate 10 random trees and put one of them per file
```
gotree generate yuletree -n 10 | gotree divide -o tree
```
Should poduce the files
```
tree_000.nw
tree_001.nw
tree_002.nw
tree_003.nw
tree_004.nw
tree_005.nw
tree_006.nw
tree_007.nw
tree_008.nw
tree_009.nw
```
