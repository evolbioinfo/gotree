# Gotree: toolkit and api for phylogenetic tree manipulation

## Commands

### matrix
This command prints the distance matrix associated to the input tree.


#### Usage

```
Usage:
  gotree matrix [flags]

Flags:
  -i, --input string    Input tree (default "stdin")
  -o, --output string   Matrix output file (default "stdout")
```

#### Example

We generate a random tree and print its associated distance matrix, and infer a tree from the distance matrix using FastME. Finally we display both trees

```
gotree generate yuletree --seed 10 | gotree matrix -o matrix.txt
fastme-2.1.5-osx -i matrix.txt
gotree reroot midpoint -i matrix.txt_fastme_tree.nwk | gotree draw svg -w 200 -H 200 --no-tip-labels -r  -o commands/matrix_1.svg
gotree generate yuletree --seed 10 | gotree reroot midpoint | gotree draw svg -w 200 -H 200 -r --no-tip-labels -o commands/matrix_2.svg
```

Random Tree                          | Inferred Tree
-------------------------------------|--------------------------------
![Random tree](matrix_2.svg)         | ![Inferred tree](matrix_1.svg) 
