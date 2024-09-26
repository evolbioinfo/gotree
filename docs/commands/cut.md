# Gotree: toolkit and api for phylogenetic tree manipulation

## Commands

### cut
Cut the input tree by keeping only parts in date window.

It extracts parts of the tree corresponding to >= min-date and <= max-date.

If min-date falls on an internal branch, it will create a new root node and will extract a tree starting at this node.
If max-date is specified (>0) : It removes all tips that are > maxdate

This command considers the input tree as rooted.

Dates are taken from the field [&date=] of the Nexus format.

#### Usage

```
Usage:
  gotree cut date [flags]

Flags:
  -h, --help             help for date
  -i, --input string     Input tree(s) file (default "stdin")
      --max-date float   Maximum date to cut the tree (0=no max date)
      --min-date float   Minimum date to cut the tree
  -o, --output string    Forest output file (default "stdout")

Global Flags:
      --format string   Input tree format (newick, nexus, phyloxml, or nextstrain) (default "newick")
```

#### Example


Initial tree:
```
+-------------------- A[&date="2000"]
[&date="1990"]
|                        +------------ B[&date="2008"]
+------------------------|[&date="2002"]
                         |          +----- C[&date="2010"]
                         +----------|[&date="2007"]
                                    +------- D[&date="2011"]
```

If we cut it between 2003 and 2009:
```
echo '(A[&date="2000"]:10,(B[&date="2008"]:6,(C[&date="2010"]:3,D[&date="2011"]:4)[&date="2007"]:5)[&date="2002"]:12)[&date="1990"];' | ./gotree cut date --min-date 2003 --max-date 2009

(B[&date="2008"]:5)[&date="2003.000000"];
```

- There is only one hanging branch left

If we cut it between 2003 and 2020:

```
echo '(A[&date="2000"]:10,(B[&date="2008"]:6,(C[&date="2010"]:3,D[&date="2011"]:4)[&date="2007"]:5)[&date="2002"]:12)[&date="1990"];' | ./gotree cut date --min-date 2003 --max-date 2020

(B[&date="2008"]:5)[&date="2003.000000"];
((C[&date="2010"]:3,D[&date="2011"]:4)[&date="2007"]:4)[&date="2003.000000"];
```

- There are two trees: 1) only one hanging branch, and 2) two tips left
