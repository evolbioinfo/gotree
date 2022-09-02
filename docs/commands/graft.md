# Gotree: toolkit and api for phylogenetic tree manipulation

## Commands

### graft
Graft a tree t2 on a tree t1, at the position of a given tip.

The root of t2 will replace the given tip of t2.

#### Usage

General command
```
Usage:
  gotree graft [flags]

Flags:
  -c, --graft string     Tree to graft (default "none")
  -h, --help             help for graft
  -o, --output string    Output tree (default "stdout")
  -i, --reftree string   Reference tree input file (default "stdin")
  -l, --tip string       Name of the tip to graft the second tree at (default "none")

Global Flags:
      --format string   Input tree format (newick, nexus, phyloxml, or nextstrain) (default "newick")
```

#### Examples

* grafting t2 on t1, at tip l1

```
	t1:      t2:
	/--- l1  /---l4
	|----l2  |---l5
	\---l3   \---l6
```

```
gotree graft -i t1.nw -c t2.nw -l l1
```

Result:
```
	     /---l4
	/--- |---l5
	|    \---l6
	|---l2
	\---l3
```
