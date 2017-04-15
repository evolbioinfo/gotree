# Gotree: toolkit and api for phylogenetic tree manipulation

## Commands

### rename
This command rename tips of the input tree, given a map file.

The map file must be tab separated with columns:
1. Current name of the tip
2. Desired new name of the tip

(if `--revert` is given then it is the other way)

If a tip name does not appear in the map file, it will not be renamed. If a name that does not exist appears in the map file, it will not throw an error.

#### Usage

General command
```
  gotree rename [flags]

Flags:
  -i, --input string    Input trees (default "stdin")
  -m, --map string      Tip name map file (default "none")
  -o, --output string   Renamed tree output file (default "stdout")
  -r, --revert          Revert orientation of map file
```

#### Examples

* We will rename all tips from a random tree

mapfile.txt
```
Tip0	Tax0
Tip1	Tax1
Tip2	Tax2
Tip3	Tax3
Tip4	Tax4
Tip5	Tax5
Tip6	Tax6
Tip7	Tax7
Tip8	Tax8
Tip9	Tax9
```

```
gotree generate yuletree -s 10 -o outtree1.nw
gotree rename -i outtree1.nw -o outtree2.nw -m mapfile.txt
gotree draw svg -w 200 -H 200 -i outtree1.nw -o commands/rename_1.svg
gotree draw svg -w 200 -H 200 -i outtree2.nw --with-branch-support --support-cutoff 0.5 -o commands/rename_2.svg
```

Initial random Tree            | Renamed Tree
-------------------------------|---------------------------------------
![Random Tree 1](rename_1.svg) | ![Renamed tree](rename_2.svg)
