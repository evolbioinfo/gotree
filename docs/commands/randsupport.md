# Gotree: toolkit and api for phylogenetic tree manipulation

## Commands

### randsupport

This command assigns a random bootstrap support to edges of input trees. The support follows a uniform distribution between 0.0 and 1.0.

#### Usage

General command
```
Usage:
  gotree randsupport [flags]

Flags:
  -i, --input string    Input tree (default "stdin")
  -o, --output string   Output file (default "stdout")
  -s, --seed int        Initial Random Seed
```

#### Examples

* We assign random supports to a random tree and highlight branches with support > 0.5
```
gotree generate yuletree -s 10 -o outtree1.nw
gotree randsupport -i outtree.nw -s 12 -o outtree2.nw
gotree draw svg -w 200 -H 200  -i outtree1.nw -o commands/randsupport_1.svg
gotree draw svg -w 200 -H 200  -i outtree2.nw --with-branch-support --support-cutoff 0.5 -o commands/randsupport_2.svg
```

Initial random Tree                 | Random Supports
------------------------------------|---------------------------------------
![Random Tree 1](randsupport_1.svg) | ![Random Supports](randsupport_2.svg)
