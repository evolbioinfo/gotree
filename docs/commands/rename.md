# Gotree: toolkit and api for phylogenetic tree manipulation

## Commands

### rename
This command renames tips and/or internal nodes of input trees. Several possibilities:

*  An existing map file is given (`-m`), and must be tab separated with columns:
   1) Current name of the tip
   2) Desired new name of the tip
   
  - If `--revert` is specified then it is the other way
  - If a tip name does not appear in the map file, it will not be renamed. 
  - If a name that does not exist appears in the map file, it will not throw an error.

* The `-a` option is given. In this case, tips and/or internal nodes are renamed using automatically generated identifiers of length 10 (or of length `--length`).
  - Correspondance between old names and new generated names is written in the map file given with `-m`. 
  - In this mode, `--revert` has no effect.
  - `--length`  allows to customize length of generated id. Length is set to 5 if given length is less that 5.
  - If several trees in input have different tip names, it does not matter, a new identifier is still generated for each new tip name.

* The `-e` (`--regexp`) and `-b` (`--replace`) is given, then it will replace matching strings in tip/node names by string given by `-b`. It takes advantages of the golang regexp machinery, i.e. it is possible to specify capturing groups and refering to it in the replacement string, for instance: `gotree rename -i tree.nh --regexp 'Tip(\d+)' --replace 'Leaf$1' -m map.txt`  will replace all matches of `Tip(\d+)` with `Leaf$1`, $1 being the matched string inside the capturing group `()`.


Other informations:
- In default mode, only tips are modified (`--tips=true` by default, to inactivate it you must specify `--tips=false`);
- If `--internal` is specified, then internal nodes are renamed;
- If after rename, several tips/nodes have the same name, subsequent commands may fail.

#### Usage

General command
```
Usage:
  gotree rename [flags]

Flags:
  -a, --auto             Renames automatically tips with auto generated id of length 10.
  -h, --help             help for rename
  -i, --input string     Input tree (default "stdin")
      --internal         Internal nodes are taken into account
  -l, --length int       Length of automatically generated id. Only with --auto (default 10)
  -m, --map string       Tip name map file (default "none")
  -o, --output string    Renamed tree output file (default "stdout")
  -e, --regexp string    Regexp to get matching tip/node names (default "none")
  -b, --replace string   String replacement to the given regexp (default "none")
  -r, --revert           Revert orientation of map file
      --tips             Tips are taken into account (--tips=false to cancel) (default true)

Global Flags:
      --format string   Input tree format (newick, nexus, or phyloxml) (default "newick")
```

#### Examples

* Rename all tips from a random tree

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


* Rename automatically all tips an input tree, using identifiers of length 5:

```
gotree generate yuletree -s 10 | gotree rename -a -m map -l 5 -o outtree.nw
```

Should give the following tree:

```
    + T0001                             
+---|                                   
|   |      +---- T0002                  
|   +------|                            
|          + T0003                      
|                                       
|----- T0004                            
|                                       
|     + T0005                           
|+----|                                 
||    |------ T0006                     
||    |                                 
||    +---- T0007                       
+|                                      
 |        +------------------ T0008     
 |      +-|                             
 +------| +----- T0009                  
        |                               
        +------------ T0010             
```

And the following map file:

```
Tip4    T0001
Tip7    T0002
Tip0    T0004
Tip9    T0006
Tip3    T0007
Tip1    T0010
Tip2    T0003
Tip8    T0005
Tip6    T0008
Tip5    T0009
```
