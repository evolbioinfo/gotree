# Gotree: toolkit and api for phylogenetic tree manipulation

## Commands

### rotate

This command reorders neighbors of all internal nodes of an input tree, in two ways:
1. `gotree rotate sort` : Sorts internal node neighbors by number of tips;
2. `gotree rotate rand`: Randomly reorders neighbors of internal nodes.

This commands do not change the tree topology, but instead modify the trasversal of the tree, and then the newick output, or the drawing of the tree.

#### Usage

General command
```
Usage:
  gotree rotate [command]

Available Commands:
  rand        Randomly rotates children of internal nodes
  sort        Sorts children of internal nodes by number of tips

Flags:
  -h, --help            help for rotate
  -i, --input string    Input tree (default "stdin")
  -o, --output string   Rotated tree output file (default "stdout")

Global Flags:
      --format string   Input tree format (newick, nexus, or phyloxml) (default "newick")
```

Specificity of `rand` subcommand:

```
-s, --seed int   Initial Random Seed
```

#### Examples

* Sort internal node neighbors by number of tips

```
echo "(((9,10),8),(((1,(2,8)),(3,4)),5),(6,7));" | gotree rotate sort
```

Should print
```
((6,7),(8,(9,10)),(5,((3,4),(1,(2,8)))));
```
