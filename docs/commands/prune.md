# Gotree: toolkit and api for phylogenetic tree manipulation

## Commands

### prune
This command removes (or retain with `-r`) a given set of tips from input trees. Several possibilities, in order of priorities :
1. Giving a tip file (`-f`): This file contains one tip name per line. In this case, it will remove (or retain with `-r`) only tips given in the file; 
2. Giving a compared tree (`-c`): In this case, tips that are specific to the input tree are removed (or retained if `-r`) from the input tree;
3. Giving a number of random tips (`--random`), tips that are sampled are removed (or retained if `-r`) from the input tree or
4. Giving tip names on the commandline: In this case, it will remove (or retain with `-r`) only the tips given on the command line. 

If  2 branches need to be merged after a tip removal, length of these branches are added, and the bootstrap support of the new branch is the maximum of the bootstrap supports of the two branches.

#### Usage

```
Usage:
  gotree prune [flags]

Flags:
  -c, --comp string      Input compared tree  (default "none")
  -o, --output string    Output tree (default "stdout")
      --random int       Number of tips to randomly sample
  -i, --ref string       Input reference tree (default "stdin")
  -r, --revert           If true, then revert the behavior: will keep only species given in the command line, or remove the species that are in common with compared tree (no effect with --random)
      --seed int         Random Seed: -1 = nano seconds since 1970/01/01 00:00:00 (default -1)
  -f, --tipfile string   Tip file (default "none")
```

#### Example

* Removing two tips from the tree

```
gotree generate yuletree --seed 10 -o outtree.nw
gotree prune -i outtree.nw -o pruned.nw Tip1 Tip2
gotree draw svg -w 200 -H 200  -i outtree.nw -o commands/prune_1.svg
gotree draw svg -w 200 -H 200  -i pruned.nw -o commands/prune_2.svg
```
Random Tree                          | Pruned Tree
-------------------------------------|-----------------------------------
![Random Tree](prune_1.svg)          | ![Pruned Tree](prune_2.svg) 


* Removing tips that are not common between two trees.
```
gotree generate yuletree --seed 10 -l 20 -o outtree1.nw
gotree generate yuletree --seed 12 -l 10 -o outtree2.nw
gotree prune -i outtree1.nw -c outtree2.nw -o pruned.nw
gotree draw svg -w 200 -H 200  -i outtree1.nw -o commands/prune_3.svg
gotree draw svg -w 200 -H 200  -i outtree2.nw -o commands/prune_4.svg
gotree draw svg -w 200 -H 200  -i pruned.nw -o commands/prune_5.svg
```

Random Tree 20 Tips           | Random Tree 10 Tips          | Pruned tree
------------------------------|------------------------------|---------------------------------
![Random Tree 1](prune_3.svg) | ![Random Tree 2](prune_4.svg)| ![Pruned Tree 2](prune_5.svg) 

* Removing 4 tips randomly:
```
gotree generate yuletree --seed 10 -l 20 -o outtree1.nw
gotree prune -i outtree1.nw --random 4 --seed 10 -o pruned.nw
gotree draw svg -w 200 -H 200  -i outtree1.nw -o commands/prune_6.svg
gotree draw svg -w 200 -H 200  -i pruned.nw -o commands/prune_7.svg
```

Random Tree                   | Pruned Tree                  
------------------------------|------------------------------
![Random Tree](prune_6.svg) | ![Pruned Tree](prune_7.svg)
