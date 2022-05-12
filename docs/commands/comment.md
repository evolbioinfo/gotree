# Gotree: toolkit and api for phylogenetic tree manipulation

## Commands

### comment
This command modifies branch/node comments of input trees.

Please note that comments may be associated to nodes or to edges depending on where they are located in the newick representation. If the comment is located after branch length (i.e. `:0.0011[edge comment]`, it will be associated to the branch. Otherwise, it will be associated to the node, i.e. `(t1,t2)[node comment]:0.001[edge comment]`).

If the tree has no branch lengths, it is not possible to differentiate them, thus all comments are associated to nodes.

#### Usage

Modify branch/node comments

Version:

Usage:
  gotree comment [command]

Available Commands:
  clear       Removes node/tip comments
  transfer    Transfers node names to comments

Flags:
  -h, --help            help for comment
  -i, --input string    Input tree (default "stdin")
  -o, --output string   Cleared tree output file (default "stdout")

Global Flags:
      --format string   Input tree format (newick, nexus, or phyloxml) (default "newick")
      --seed int        Random Seed: -1 = nano seconds since 1970/01/01 00:00:00 (default -1)
  -t, --threads int     Number of threads (Max=4) (default 1)
  ```

transfer subcommand
```
Usage:
  gotree comment clear [flags]

Flags:
      --edges-only   Clear comments on edges only
  -h, --help         help for clear
      --nodes-only   Clear comments on nodes only

Global Flags:
  -i, --input string    Input tree (default "stdin")
  -o, --output string   Cleared tree output file (default "stdout")
```

clear subcommand
```
Usage:
  gotree comment transfer [flags]

Flags:
  -h, --help   help for transfer
  --reverse    

Global Flags:
  -i, --input string    Input tree (default "stdin")
  -o, --output string   Cleared tree output file (default "stdout")
```


#### Examples

* Removing node comments from an input tree
```
echo "(t1[n1]:1[e1],t2[n2]:1[e2],(t3[n3]:1[e3],t4[n4]:1[e4])[n5]:1[e5]);" | gotree comment clear --nodes-only
```

Should print:
```
(t1:1[e1],t2:1[e2],(t3:1[e3],t4:1[e4]):1[e5]);
```

* Removing edge comments from an input tree
```
echo "(t1[n1]:1[e1],t2[n2]:1[e2],(t3[n3]:1[e3],t4[n4]:1[e4])[n5]:1[e5]);" | gotree comment clear --edges-only
```

Should print:
```
(t1[n1]:1,t2[n2]:1,(t3[n3]:1,t4[n4]:1)[n5]:1);
```

* Removing all comments from an input tree
```
echo "(t1[n1]:1[e1],t2[n2]:1[e2],(t3[n3]:1[e3],t4[n4]:1[e4])[n5]:1[e5]);" | gotree comment clear 
```

Should print:
```
(t1:1,t2:1,(t3:1,t4:1):1);
```

* Transfering node names to node comments:
```
echo "(t1:1,t2:1,(t3:1,t4:1)n5:1);" | gotree comment transfer
```

Should print:
```
(t1:1,t2:1,(t3:1,t4:1)[n5]:1);
```

* Transfering node comment to node name:
```
echo "(t1:1,t2:1,(t3:1,t4:1)[n5]:1);" | gotree comment transfer --reverse
```

Should print:
```
(t1:1,t2:1,(t3:1,t4:1)n5:1);
```

