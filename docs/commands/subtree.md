# Gotree: toolkit and api for phylogenetic tree manipulation

## Commands

### subtree
This commands selects an internal node of the input tree, and outputs the subtree starting at that node. The internal node is selected by its name, which can be specified with a regexp.

If several nodes match the given name/regexp, it does nothing, and print the name of matching nodes.

The only matching node must be an internal node, otherwise, it will do nothing and print the tip.

#### Usage

General command
```
Usage:
  gotree subtree [flags]

Flags:
  -i, --input string    Input tree (default "stdin")
  -n, --name string     Name of the node to select as the root of the subtree (maybe a regex) (default "none")
  -o, --output string   Output tree file (default "stdout")
```

#### Examples

* Generate a random tree, annotate an internal node, and take the subtree starting at this internal node

clade.txt
```
clade:Tip2,Tip4,Tip7
```

```
gotree generate yuletree -s 10 -o outtree1.nw
gotree annotate -m clade.txt -i outtree.nw | gotree subtree -n clade -o outtree2.nw
gotree draw svg -w 200 -H 200  -i outtree1.nw -o commands/subtree_1.svg
gotree draw svg -w 200 -H 200  -i outtree2.nw -o commands/subtree_2.svg
```

Initial random Tree                 | Subtree
------------------------------------|---------------------------------------
![Random Tree 1](subtree_1.svg)     | ![Random Supports](subtree_2.svg) 
