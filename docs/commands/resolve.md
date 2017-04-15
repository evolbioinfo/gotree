# Gotree: toolkit and api for phylogenetic tree manipulation

## Commands

### resolve
This command randomly resolves multifurcations by adding 0 length branches. If any node has more than 3 neighbors: it randomly adds 0 length branches until it has 3 neighbors.

#### Usage

General command
```
Usage:
  gotree resolve [flags]

Flags:
  -i, --input string    Input tree(s) file (default "stdin")
  -o, --output string   Resolved tree(s) output file (default "stdout")
  -s, --seed int        Initial Random Seed
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

