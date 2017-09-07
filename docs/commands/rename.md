# Gotree: toolkit and api for phylogenetic tree manipulation

## Commands

### rename
This command renames tips of input trees. Several possibilities:

* An existing map file is given (-m). The map file must be tab separated with the following columns:
  1. Current name of the tip
  2. Desired new name of the tip

if `--revert` is given, then it is the other way.
If a tip name does not appear in the map file, it will not be renamed. If a name that does not exist appears in the map file, it will not throw an error.

* The `-a` option is given. In this case, tips are renamed using automatically generated identifiers of length 10 (or of length `--length`).
  * Correspondance between old names and new generated names is written in the map file given with `-m`. 
  * In this mode, `--revert` has no effect.
  * `--length`  allows to customize length of generated id. Length is set to 5 if given length is less that 5.
  * If several trees in input have different tip names, it does not matter, a new identifier is still generated for each new tip name.

#### Usage

General command
```
  gotree rename [flags]

Flags:
  -a, --auto            Renames automatically tips with auto generated id of length 10.
  -i, --input string    Input tree (default "stdin")
  -l, --length int      Length of automatically generated id. (default 10)
  -m, --map string      Tip name map file (default "none")
  -o, --output string   Renamed tree output file (default "stdout")
  -r, --revert          Revert orientation of map file
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
