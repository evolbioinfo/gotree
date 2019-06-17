# Gotree: toolkit and api for phylogenetic tree manipulation

## Commands

### repopulate
Re populate the tree with identical tips (having the same sequences).

When a tree is inferred, some tools first remove identical sequences.

However, it may be useful to keep them in the tree. To do so, this command takes:

1. A input tree
2. A file containing a list of tips that are identical, in the following format:

```
Tip1,Tip2
Tip3,Tip4
```

This means that Tip1 is identical to Tip2, and Tip3 is identical to Tip4.

"repopulate" command then adds Tip2 next to Tip1 if Tip1 is present in the tree, or 
Tip1 next to Tip2 if Tip2 is present in the tree. To do so, it adds two 0.0 length
 branches. 

Example with Tip1,Tip2 :

           Before                  |          After (if `l>0.0`)       |         After (if `l=0.0`)
-----------------------------------+-----------------------------------+----------------------------------
![Repopulate 1](repopulate_1.png)  | ![Repopulate 2](repopulate_2.png) | ![Repopulate 3](repopulate_3.png)


Each identical group must contain exactly 1 already present tip, otherwise it returns
 an error.

If a new tip is present in several groups, then returns and error.

The tree after "repopulate" command may contain polytomies.


#### Usage

General command
```
Usage:
  gotree repopulate [flags]

Flags:
  -h, --help               help for repopulate
  -g, --id-groups string   File with groups of identical tips (default "none")
  -i, --input string       Input tree (default "stdin")
  -o, --output string      Output tree file (default "stdout")

Global Flags:
      --format string   Input tree format (newick, nexus, or phyloxml) (default "newick")
```

#### Examples

* Add several identical tips to an input tree

Tip file `tips.txt`:

```
Tip2,Tip5,Tip6,Tip7
Tip4,Tip8
```

```
$ echo "(Tip4:0.1,Tip0:0.1,(Tip3:0.1,(Tip2:0.2,Tip1:0.2)0.8:0.3)0.9:0.4);" | gotree repopulate -g tips.txt
((Tip8:0,Tip4:0):0.1,Tip0:0.1,(Tip3:0.1,((Tip5:0,Tip2:0,Tip6:0,Tip7:0):0.2,Tip1:0.2)0.8:0.3)0.9:0.4);
```


           Before                  |          After
-----------------------------------+-----------------------------------
![Repopulate 4](repopulate_4.png)  | ![Repopulate 5](repopulate_5.png) 

