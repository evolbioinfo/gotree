# Gotree: toolkit and api for phylogenetic tree manipulation

## Commands

### unroot
This command unroots a set of input trees. If a tree is already unrooted it does nothing. Otherwise it places the new pseudo root on a trifurcated node and removes the old root. When removing the old root, two branches have to be merged. Length of the new branch will be the sum of the two merged branches, and its support will be the maximum. 

```
             ------C         
             |z	         
    ---------*	                       ------C 
    |x       |t	                 x+y   |z	   
ROOT*        ------B   =>    A---------*ROOT   
    |y		                       |t	   		 
    ---*A                              ------B 
```

#### Usage

General command
```
Usage:
  gotree unroot [flags]

Flags:
  -i, --input string    Input tree (default "stdin")
  -o, --output string   Collapsed tree output file (default "stdout")

```

#### Examples

* Generate a rooted random tree, and unroots it

```
gotree generate yuletree -r -s 10 -o outtree1.nw
gotree unroot -i outtree1.nw -o outtree2.nw
gotree draw svg -w 200 -H 200 -i outtree1.nw -o commands/unroot_1.svg
gotree draw svg -w 200 -H 200 -i outtree2.nw -o commands/unroot_2.svg
```

Initial random Tree                 | Unrooted tree
------------------------------------|---------------------------------------
![Random Tree](unroot_1.svg)       | ![Subtree](unroot_2.svg) 
