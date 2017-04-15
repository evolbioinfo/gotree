# Gotree: toolkit and api for phylogenetic tree manipulation

## Commands

### annotate
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
  -i, --input string      Input tree(s) file (default "stdin")
  -m, --map-file string   Name map input file (default "none")
  -o, --output string     Resolved tree(s) output file (default "stdout")

Global Flags:
  -t, --threads int   Number of threads (Max=12) (default 1)
```

#### Example

mapfile.txt
```
internal1:Tip6,Tip5,Tip1
```

commands:
```
gotree generate yuletree -s 10 | gotree annotate -m mapfile.txt | gotree draw svg -w 800 -H 800  -c --with-node-labels > commands/annotate_1.svg
```

Should give:

![Tree image](annotate_1.svg)
