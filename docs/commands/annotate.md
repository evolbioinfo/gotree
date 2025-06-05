# Gotree: toolkit and api for phylogenetic tree manipulation

## Commands

### annotate
This command annotates internal nodes/branches of a tree with given information.

Annotations may be (in order of priority):
- A tree with labels on internal nodes (-c). in that case, it will label each branch of 
   the input tree with label of the closest branch of the given compared tree (-c) in terms
   of transfer distance. The labels are of the form: "`label_distance_depth`"; Only internal branches
   are annotated, and no internal branch is annotated with a terminal branch.
- A file with one line per internal node to annotate (-m), and with the following format:
   `<name of internal branch/node n1>:<name of taxon n2>,<name of taxon n3>,...,<name of taxon ni>`
	=> If 0 name is given after ':' an error is returned
	=> If 1 name 'n2' is given after ':' : we search for n2 in the tree (tip or internal node)
       and rename it as n1
    => If > 1 names '[n2,...,ni]' are given after ':' : We find the LCA of every tips whose name 
	   is in '[n2,...,ni]' and rename it as n1.

- If `--comment` is specified, then we do not change the names, but the comments of the given nodes.
- Otherwise output tree won't have bootstrap support at the branches anymore

- If `--subtrees` is given (and `-m`): for each annotation line, not only the given internal node is annotated, but all its descending internal nodes as well (usefull for some branch tests, e.g. hyphy, etc.)


If neither -c nor -m are given, gotree annotate will wait for data on stdin

This command annotates internal branches of a set of trees with given data.

It takes a map file with one line per internal node to annotate:

```
<name of internal node>:<name of taxon 1>,<name of taxon2>,...,<name of taxon n>
```

And will retrieve the last common ancestor of taxa and annotate it with the given name.

The output tree will not have bootstrap support at that branches anymore.


#### Usage

```
Usage:
  gotree annotate [flags]

Flags:
  -c, --compared string   Compared tree file (default "stdin")
  -i, --input string      Input tree(s) file (default "stdin")
  -m, --map-file string   Name map input file (default "none")
  -o, --output string     Resolved tree(s) output file (default "stdout")
```

#### Example

* Using a mapfile
mapfile.txt
```
internal1:Tip6,Tip5,Tip1
```

commands:
```
gotree generate yuletree --seed 10 | gotree annotate -m mapfile.txt | gotree draw svg -w 800 -H 800  -c --with-node-labels > commands/annotate_1.svg
```

Should give:

![Tree image](annotate_1.svg)

* Using a reference tree: See example in documentation for [download](download.md) command. 
