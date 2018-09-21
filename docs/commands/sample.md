# Gotree: toolkit and api for phylogenetic tree manipulation

## Commands

### sample
This command takes a sample of the set of trees from the input file.

It can be with or without replacement depending on the presence of the `--replace` option.

If the number of desired trees is > number of input trees: 
  - with `--replace`: Will take `-n` trees with replacement;
  - without `--replace`: Will take all input trees.

#### Usage

General command
```
Usage:
  gotree sample [flags]

Flags:
  -i, --input string    Input reference trees (default "stdin")
  -n, --nbtrees int     Number of trees to sample from input file (default 1)
  -o, --output string   Output trees (default "stdout")
      --replace         If given, samples with replacement
      --seed int        Random Seed: -1 = nano seconds since 1970/01/01 00:00:00 (default -1)
```

#### Examples

* We generate 10 random trees, take a sample of size 10 with replacement, and count different sampled trees:

```
gotree generate yuletree --seed 10 -l 4 -n 10 | gotree sample -n 10 --replace --seed 10 | sort | uniq -c
```

Should give:
```
      1 (Tip2:0.08565127428804534,Tip0:0.021705093532846203,(Tip3:0.07735928108468929,Tip1:0.08967099662302665):0.002022779393968061);
      3 (Tip2:0.18347097513974125,Tip0:0.0560088923217098,(Tip3:0.17182241382980687,Tip1:0.1920960924280275):0.007725336355452564);
      3 (Tip3:0.004413852122575761,Tip0:0.0015988844434968424,(Tip2:0.20136539775567347,Tip1:0.0330285678892811):0.05565924680439801);
      1 ((Tip3:0.043398414425198095,Tip2:0.11878365225567722):0.010523628571460197,Tip0:0.04494388758695279,Tip1:0.1683671542243689);
      1 ((Tip3:0.2704987376078806,Tip2:0.004939317945689752):0.02131497189184287,Tip0:0.06748946801240394,Tip1:0.32968964999864436);
      1 ((Tip3:0.28668879505468636,Tip2:0.21552905482949153):0.01853723048382153,Tip0:0.07404867683401199,Tip1:0.13492605122032592);
```
