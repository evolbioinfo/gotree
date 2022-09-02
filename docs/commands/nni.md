# Gotree: toolkit and api for phylogenetic tree manipulation

## Commands

### nni
This command generates all NNI neighbors from a given tree.

#### Usage

```
Usage:
  gotree nni [flags]

Flags:
  -h, --help            help for nni
  -i, --input string    Input Tree (default "stdin")
  -o, --output string   NNI output tree file (default "stdout")

Global Flags:
      --format string   Input tree format (newick, nexus, phyloxml, or nextstrain) (default "newick")
```

#### Example

* Generates NNI neighbors

```
echo "(n1_1,n1_2,(n2_1,n2_2)n2)n1;" | gotree nni
```

It should give the following trees:
```
(n1_1,n2_2,(n2_1,n1_2)n2)n1;
(n1_1,n2_1,(n1_2,n2_2)n2)n1;
```


Tree:
```
n1_1 ---------------+                +--------------- n2_1
                    |----------------|                            
n1_2 ---------------+                +--------------- n2_2
```

NNIs:
```
n1_1 ---------------+                +--------------- n2_1
                    |----------------|                            
n2_2 ---------------+                +--------------- n1_2

n1_1 ---------------+                +--------------- n1_2
                    |----------------|                            
n2_1 ---------------+                +--------------- n2_2
```
