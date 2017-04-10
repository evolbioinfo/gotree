# Gotree: toolkit and api for phylogenetic tree manipulation

## Commands

### generate
This command generates random trees according to different models:
* `gotree generate balancedtree` : perfectly balanced binary tree
* `gotree generate caterpillartree`: caterpillar tree
* `gotree generate uniform tree` : uniform tree (edges are added randomly in the middle of any previous edge)
* `gotree generate yuletree`: Yule-Harding model (edges are added randomly in the middle of any external edge)

All commands take a number of taxa/leaves (`-l`) as option except the balancedtree commands that takes a depth (`-d`).

#### Usage

General command
```
Usage:
  gotree generate [command]

Available Commands:
  balancedtree    Generates a random balanced binary tree
  caterpillartree Generates a random caterpilar binary tree
  uniformtree     Generates a random uniform binary tree
  yuletree        Generates a random yule binary tree

Flags:
  -n, --nbtrees int     Number of trees to generate (default 1)
  -o, --output string   Number of tips of the tree to generate (default "stdout")
  -r, --rooted          Generate rooted trees
  -s, --seed int        Initial Random Seed (default 1491857625914835301)
```

#### Examples

* Generate Yule-Harding tree with 1000 taxa
```
gotree generate yuletree -s 10 -l 1000 | gotree draw svg -r -w 200 -H 200 --no-tip-labels -o commands/generate_1.svg
```

![yule](generate_1.svg)

* Generate caterpillar tree with 1000 taxa
```
gotree generate caterpillartree -s 10 -l 1000 | gotree draw svg -r -w 200 -H 200 --no-tip-labels -o commands/generate_2.svg
```

![caterpillar](generate_2.svg)

* Generate Balanced tree with depth 10
```
gotree generate balancedtree -s 10 -d 10 | gotree draw svg -r -w 200 -H 200 --no-tip-labels -o commands/generate_3.svg
```

![balanced](generate_3.svg)

* Generate uniform tree
```
gotree generate uniformtree -s 10 -l 1000 | gotree draw svg -r -w 200 -H 200 --no-tip-labels -o commands/generate_4.svg
```

![uniform](generate_4.svg)
