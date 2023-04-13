# Gotree: toolkit and api for phylogenetic tree manipulation

## Commands

### matrix
Prints distance matrix associated to the input tree.

The distance matrix can be computed in several ways, depending on the "metric" option:
* --metric brlen : distances correspond to the sum of branch lengths between the tips (patristic distance). If there is no length for a given branch, 0.0 is the default.
* --metric boot : distances correspond to the sum of supports of the internal branches separating the tips. If there is no support for a given branch (e.g. for a tip), 1.0 is the default. If branch supports range from 0 to 100, you may consider to use gotree support scale -f 0.01 first.
* --metric none : distances correspond to the sum of the branches separating the tips, but each individual branch is counted as having a length of 1 (topological distance)

#### Usage

```

Usage:
  gotree matrix [flags]

Flags:
  -i, --input string    Input tree (default "stdin")
  -m, --metric string   Distance metric (brlen|boot|none) (default "brlen")
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
