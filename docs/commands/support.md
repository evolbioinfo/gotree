# Gotree: toolkit and api for phylogenetic tree manipulation

## Commands

### support
This command modifies branch lengths of input trees.

#### Usage

```
Usage:
  gotree support [command]

Available Commands:
  clear       Clear supports from input trees
  setrand     Assign a random support to edges of input trees

Flags:
  -h, --help            help for support
  -i, --input string    Input tree (default "stdin")
  -o, --output string   Cleared tree output file (default "stdout")
```

clear subcommand
```
Usage:
  gotree support clear [flags]

Flags:
  -h, --help   help for clear

Global Flags:
  -i, --input string    Input tree (default "stdin")
  -o, --output string   Cleared tree output file (default "stdout")
```

setrand subcommand
```
Usage:
  gotree support setrand [flags]

Flags:
  -h, --help       help for setrand
  -s, --seed int   Initial Random Seed (default 1513954241492262668)

Global Flags:
  -i, --input string    Input tree (default "stdin")
  -o, --output string   Cleared tree output file (default "stdout")
```

#### Examples

1. Removing branch supports from a set of 10 trees
```
gotree generate yuletree -s 10 -n 10 | gotree brlen clear | gotree support setrand -s 10
```

Should give: 
```
((Tip4,(Tip7,Tip2)0.41765200380165207)0.5660920659323543,Tip0,((Tip8,(Tip9,Tip3)0.48924257454472525)0.42157058562840155,((Tip6,Tip5)0.36832994780378014,Tip1)0.9167074899036827)0.925128845219594);
(Tip5,Tip0,((Tip6,(Tip7,Tip4)0.027115744632158528)0.8626564105921588,(Tip2,((Tip8,(Tip9,Tip3)0.6333402875097828)0.7125611076588594,Tip1)0.18629884992437365)0.838731378864137)0.7434673863639695);
(Tip6,Tip0,((((Tip5,Tip4)0.2738418304689928,((Tip9,Tip8)0.977174976899892,Tip3)0.44108273138920473)0.9084455030937538,(Tip7,Tip2)0.6737782279152807)0.5588867210642816,Tip1)0.9509561978377322);
(Tip3,Tip0,(((Tip6,Tip5)0.6224376297409722,Tip2)0.9248243943261302,(((Tip9,Tip8)0.8535338483942572,Tip4)0.8206155728480076,(Tip7,Tip1)0.24305139502127918)0.14316227191812703)0.8403401639121292);
(((Tip4,(Tip8,Tip3)0.09756455768006228)0.7725838174161286,((Tip6,Tip5)0.6433246886961355,Tip2)0.7405679634457042)0.3565337250650908,Tip0,(Tip7,(Tip9,Tip1)0.9431243487113095)0.3097798053899953);
(Tip4,Tip0,(Tip3,((((Tip8,Tip7)0.059589022482781605,Tip5)0.592973812044721,((Tip9,Tip6)0.05800099267108681,Tip2)0.8143090492413795)0.4252740874863805,Tip1)0.46100530105520376)0.884130479969827);
(Tip5,Tip0,((Tip8,Tip2)0.3520774712997434,((Tip7,Tip3)0.8152133800706765,((Tip9,Tip4)0.35215374439996194,(Tip6,Tip1)0.5753596666753715)0.7226776216380938)0.6951198476027183)0.1897987224883182);
((Tip9,Tip6)0.996348024291363,Tip0,((((Tip8,Tip4)0.5920905095134684,Tip3)0.5386474652850751,Tip2)0.03964818517503376,((Tip7,Tip5)0.012413676784426743,Tip1)0.10316065974432471)0.5095342920507389);
(Tip9,Tip0,(Tip8,((Tip7,Tip4)0.9630021878803063,((Tip6,Tip2)0.34707919634675627,(Tip3,(Tip5,Tip1)0.04819317477170198)0.9331288319823794)0.9050744971731132)0.5799061565448809)0.740705127298385);
((((Tip7,(Tip8,(Tip9,Tip6)0.954400690745701)0.5819457063280987)0.3978726130049126,Tip5)0.2360963909039234,Tip4)0.15841173792024738,Tip0,(Tip2,(Tip3,Tip1)0.5647012033384201)0.6749510929466774);
```

```
gotree generate yuletree -s 10 -n 10 | gotree brlen clear | gotree support setrand -s 10 | gotree support clear
```

Should produce:
```
((Tip4,(Tip7,Tip2)),Tip0,((Tip8,(Tip9,Tip3)),((Tip6,Tip5),Tip1)));
(Tip5,Tip0,((Tip6,(Tip7,Tip4)),(Tip2,((Tip8,(Tip9,Tip3)),Tip1))));
(Tip6,Tip0,((((Tip5,Tip4),((Tip9,Tip8),Tip3)),(Tip7,Tip2)),Tip1));
(Tip3,Tip0,(((Tip6,Tip5),Tip2),(((Tip9,Tip8),Tip4),(Tip7,Tip1))));
(((Tip4,(Tip8,Tip3)),((Tip6,Tip5),Tip2)),Tip0,(Tip7,(Tip9,Tip1)));
(Tip4,Tip0,(Tip3,((((Tip8,Tip7),Tip5),((Tip9,Tip6),Tip2)),Tip1)));
(Tip5,Tip0,((Tip8,Tip2),((Tip7,Tip3),((Tip9,Tip4),(Tip6,Tip1)))));
((Tip9,Tip6),Tip0,((((Tip8,Tip4),Tip3),Tip2),((Tip7,Tip5),Tip1)));
(Tip9,Tip0,(Tip8,((Tip7,Tip4),((Tip6,Tip2),(Tip3,(Tip5,Tip1))))));
((((Tip7,(Tip8,(Tip9,Tip6))),Tip5),Tip4),Tip0,(Tip2,(Tip3,Tip1)));

```

2. Assigning random supports to a random tree and highlight branches with support > 0.5
```
gotree generate yuletree -s 10 -o outtree1.nw
gotree support setrand -i outtree.nw -s 12 -o outtree2.nw
gotree draw svg -w 200 -H 200  -i outtree1.nw -o commands/randsupport_1.svg
gotree draw svg -w 200 -H 200  -i outtree2.nw --with-branch-support --support-cutoff 0.5 -o commands/randsupport_2.svg
```

Initial random Tree                 | Random Supports
------------------------------------|---------------------------------------
![Random Tree 1](randsupport_1.svg) | ![Random Supports](randsupport_2.svg)
