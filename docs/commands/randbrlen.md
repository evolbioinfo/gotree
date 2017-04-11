# Gotree: toolkit and api for phylogenetic tree manipulation

## Commands

### randbrlen

This command assigns a random length to edges of input trees. Random length follows an exponential distribution of mean 0.1. The mean can be set with `-m`.

#### Usage

General command
```
Usage:
  gotree randbrlen [flags]

Flags:
  -i, --input string    Input tree (default "stdin")
  -m, --mean float      Mean of the exponential distribution of branch lengths (default 0.1)
  -o, --output string   Output file (default "stdout")
  -s, --seed int        Initial Random Seed (default 1491922735075976690)
```

#### Examples

* Assign random lengths to a random tree
```
gotree generate yuletree -s 10 -o outtree1.nw
gotree randbrlen -i outtree.nw -s 13 -o outtree2.nw
gotree draw svg -w 200 -H 200  -i outtree1.nw -o commands/randbrlen_1.svg
gotree draw svg -w 200 -H 200  -i outtree2.nw -o commands/randbrlen_2.svg
```

Initial random Tree               | Random lengths
----------------------------------|-----------------------------------
![Random Tree 1](randbrlen_1.svg) | ![Random Lengths](randbrlen_2.svg)
