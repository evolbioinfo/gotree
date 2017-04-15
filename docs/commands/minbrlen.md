# Gotree: toolkit and api for phylogenetic tree manipulation

## Commands

### minbrlen
This command sets a minimum length to all branches of the input tree.


#### Usage

```
Usage:
  gotree minbrlen [flags]

Flags:
  -i, --input string    Input tree (default "stdin")
  -l, --length float    Min Length cutoff
  -o, --output string   Length corrected tree output file (default "stdout")
```

#### Example

We generate a random tree and sets min branch length to 0.1.

```
gotree generate yuletree -s 10 -l 100 -o outtree.nw
gotree draw svg -r -w 200 -H 200 --no-tip-labels -i outtree.nw -o commands/minbrlen_1.svg
gotree minbrlen -i outtree.nw -l 0.1 | gotree draw svg -r -w 200 -H 200 --no-tip-labels -o commands/minbrlen_2.svg
```

Random Tree                          | Min brlen tree
-------------------------------------|-----------------------------------
![Random Tree](minbrlen_1.svg)       | ![Min brlen tree](minbrlen_2.svg) 
