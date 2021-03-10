# Gotree: toolkit and api for phylogenetic tree manipulation

## Commands

### labels
Lists labels of all tree tips

If several trees are present in the input file, labels of all trees are listed.

Example of usage:

```
gotree labels -i t.mw
```

#### Usage

General command
```
Usage:
  gotree labels [flags]

Flags:
  -h, --help   help for labels

Global Flags:
      --format string   Input tree format (newick, nexus, or phyloxml) (default "newick")
```

#### Examples

* List all tips labels of an input tree


```
$ echo "((Tip4,(Tip7,Tip2)),Tip0,((Tip8,(Tip9,Tip3)),((Tip6,Tip5),Tip1)));" | gotree labels
Tip4
Tip7
Tip2
Tip0
Tip8
Tip9
Tip3
Tip6
Tip5
Tip1
```

