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
      --seed int        Random Seed: -1 = nano seconds since 1970/01/01 00:00:00 (default -1)
```

#### Examples

* Shuffle tips of a random tree

```
gotree generate yuletree --seed 10 -o outtree1.nw
gotree shuffletips -i outtree1.nw -o outtree2.nw --seed 12
gotree draw svg -w 200 -H 200  -i outtree1.nw -o commands/shuffletips_1.svg
gotree draw svg -w 200 -H 200  -i outtree2.nw -o commands/shuffletips_2.svg
```

Initial random Tree                 | Shuffled Tree
------------------------------------|---------------------------------------
![Random Tree 1](shuffletips_1.svg) | ![Shuffled tree](shuffletips_2.svg) 
