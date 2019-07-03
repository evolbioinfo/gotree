# Gotree: toolkit and api for phylogenetic tree manipulation

## Commands

### generate
This command generates random trees according to different models:
* `gotree generate balancedtree` : perfectly balanced binary tree
* `gotree generate caterpillartree`: caterpillar tree
* `gotree generate topologies`: all topologies
* `gotree generate uniform tree` : uniform tree (edges are added randomly in the middle of any previous edge)
* `gotree generate yuletree`: Yule-Harding model (edges are added randomly in the middle of any external edge). If `-r` is not specified, the tree is unrooted.

All commands take a number of taxa/leaves (`-l`) as option except the balancedtree commands that takes a depth (`-d`).

#### Usage

General command
```
Usage:
  gotree generate [command]

Available Commands:
  balancedtree    Generates a random balanced binary tree
  caterpillartree Generates a random caterpilar binary tree
  startree        Generates a star tree (no internal branch)
  topologies      Generates all possible tree topologies
  uniformtree     Generates a random uniform binary tree
  yuletree        Generates a random yule binary tree

Flags:
  -n, --nbtrees int     Number of trees to generate (default 1)
  -o, --output string   Number of tips of the tree to generate (default "stdout")
  -r, --rooted          Generate rooted trees
      --seed int        Random Seed: -1 = nano seconds since 1970/01/01 00:00:00 (default -1)
```

#### Examples

* Generate Yule-Harding tree with 1000 taxa
```
gotree generate yuletree --seed 10 -l 1000 | gotree draw svg -r -w 200 -H 200 --no-tip-labels -o commands/generate_1.svg
```

![yule](generate_1.svg)

* Generate caterpillar tree with 1000 taxa
```
gotree generate caterpillartree --seed 10 -l 1000 | gotree draw svg -r -w 200 -H 200 --no-tip-labels -o commands/generate_2.svg
```

![caterpillar](generate_2.svg)

* Generate Balanced tree with depth 10
```
gotree generate balancedtree --seed 10 -d 10 | gotree draw svg -r -w 200 -H 200 --no-tip-labels -o commands/generate_3.svg
```

![balanced](generate_3.svg)

* Generate star tree
```
gotree generate startree -l 100 | gotree draw svg -r -w 200 -H 200 --no-tip-labels -o commands/generate_5.svg
```

![star](generate_5.svg)


* Generate uniform tree
```
gotree generate uniformtree --seed 10 -l 1000 | gotree draw svg -r -w 200 -H 200 --no-tip-labels -o commands/generate_4.svg
```

![uniform](generate_4.svg)

* Generate all 5 tips unrooted trees
```
gotree generate topologies -l 5
```

```
((Tip5,(Tip4,Tip1)),Tip2,Tip3);
(((Tip5,Tip4),Tip1),Tip2,Tip3);
((Tip4,(Tip5,Tip1)),Tip2,Tip3);
((Tip4,Tip1),(Tip5,Tip2),Tip3);
((Tip4,Tip1),Tip2,(Tip5,Tip3));
((Tip5,Tip1),(Tip4,Tip2),Tip3);
(Tip1,(Tip5,(Tip4,Tip2)),Tip3);
(Tip1,((Tip5,Tip4),Tip2),Tip3);
(Tip1,(Tip4,(Tip5,Tip2)),Tip3);
(Tip1,(Tip4,Tip2),(Tip5,Tip3));
((Tip5,Tip1),Tip2,(Tip4,Tip3));
(Tip1,(Tip5,Tip2),(Tip4,Tip3));
(Tip1,Tip2,(Tip5,(Tip4,Tip3)));
(Tip1,Tip2,((Tip5,Tip4),Tip3));
(Tip1,Tip2,(Tip4,(Tip5,Tip3)));
```

* Generate all 5 tips unrooted trees with tip names taken from another tree

input.nw
```
(A,(B,D),(C,E));
```

```
gotree generate topologies -i input.nw
```

```
((E,(C,A)),B,D);
(((E,C),A),B,D);
((C,(E,A)),B,D);
((C,A),(E,B),D);
((C,A),B,(E,D));
((E,A),(C,B),D);
(A,(E,(C,B)),D);
(A,((E,C),B),D);
(A,(C,(E,B)),D);
(A,(C,B),(E,D));
((E,A),B,(C,D));
(A,(E,B),(C,D));
(A,B,(E,(C,D)));
(A,B,((E,C),D));
(A,B,(C,(E,D)));
```
