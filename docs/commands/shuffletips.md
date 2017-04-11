# Gotree: toolkit and api for phylogenetic tree manipulation

## Commands

### shuffletips
This commands assigns tip names randomly on the tree, keeping the same topology.

#### Usage

General command
```
Usage:
  gotree shuffletips [flags]

Flags:
  -i, --input string    Input tree (default "stdin")
  -o, --output string   Shuffled tree output file (default "stdout")
  -s, --seed int        Initial Random Seed
```

#### Examples

* Shuffle tips of a random tree

```
gotree generate yuletree -s 10 -o outtree1.nw
gotree shuffletips -i outtree1.nw -o outtree2.nw -s 12
gotree draw svg -w 200 -H 200  -i outtree1.nw -o commands/shuffletips_1.svg
gotree draw svg -w 200 -H 200  -i outtree2.nw -o commands/shuffletips_2.svg
```

Initial random Tree                 | Shuffled Tree
------------------------------------|---------------------------------------
![Random Tree 1](shuffletips_1.svg) | ![Random Supports](shuffletips_2.svg) 
