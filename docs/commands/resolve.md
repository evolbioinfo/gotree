# Gotree: toolkit and api for phylogenetic tree manipulation

## Commands

### resolve
This command randomly resolves multifurcations by adding 0 length branches. If any node has more than 3 neighbors: it randomly adds 0 length branches until it has 3 neighbors.

If `--rooted` is specified, then the root is randomly resolved as well, to produce a rooted output tree.

#### Usage

General command
```
Usage:
  gotree resolve [flags]

Flags:
  -i, --input string    Input tree(s) file (default "stdin")
  -o, --output string   Resolved tree(s) output file (default "stdout")
      --rooted          Considers the tree as rooted (will randomly resolve the root also if needed)
      --seed int        Random Seed: -1 = nano seconds since 1970/01/01 00:00:00 (default -1)
```

#### Examples

* We generate a random tree, collapse branches with length < 0.5 and resolve randomly the multifurcations

```
gotree generate yuletree -s 10 -o outtree1.nw
gotree collapse length -i outtree1.nw -l 0.05 -o outtree2.nw
gotree resolve -i outtree2.nw -o outtree3.nw
gotree draw svg -w 200 -H 200  -i outtree1.nw -o commands/resolve_1.svg
gotree draw svg -w 200 -H 200  -i outtree2.nw -o commands/resolve_2.svg
gotree draw svg -w 200 -H 200  -i outtree3.nw --no-branch-lengths --with-branch-support --support-cutoff 0.5 -o commands/resolve_3.svg
```

Initial random Tree             | Collapsed Tree                     | Resolved Tree
--------------------------------|------------------------------------|---------------------------------
![Random Tree 1](resolve_1.svg) | ![Collapsed tree](resolve_2.svg)   | ![Resolved Supports](resolve_3.svg)


### resolve named
This command resolves internal named nodes as new tips with 0 length branches

#### Usage

General command
```
Usage:
  gotree resolve named [flags]

Flags:
  -h, --help   help for named

Global Flags:
      --format string   Input tree format (newick, nexus, phyloxml, or nextstrain) (default "newick")
  -i, --input string    Input tree(s) file (default "stdin")
  -o, --output string   Resolved tree(s) output file (default "stdout")
```

#### Examples

* resolving internal named nodes on a simple tree:

```
> echo "(T1:1,T2:1,T3:1)N1;" | gotree resolve named
(T1:1,T2:1,T3:1,N1:0)N1;
```

```
 	      -------T1             -------T1
	      |                     |
	T3----*N1        =>   T3----*N1---N1
	      |                     |
	      -------T2             -------T2 
```
