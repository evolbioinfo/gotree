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
  -h, --help       help for labels
      --internal   Internal node labels are listed
      --tips       Tip labels are listed (--tips=false to cancel) (default true)

Global Flags:
      --format string   Input tree format (newick, nexus, phyloxml, or nextstrain) (default "newick")
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

* List all tips and internal nodes labels 

```
echo "(1,(2,(3,4,5,6)polytomy)internal)root;" | gotree labels --internal
root
1
internal
2
polytomy
3
4
5
6
```


* List only internal node labels 

```
echo "(1,(2,(3,4,5,6)polytomy)internal)root;" | gotree labels --tips=false --internal
root
internal
polytomy
```
