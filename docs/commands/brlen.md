# Gotree: toolkit and api for phylogenetic tree manipulation

## Commands

### brlen
This command modifies branch lengths of input trees.

#### Usage

```
Usage:
  gotree brlen [command]

Available Commands:
  clear       Clear lengths from input trees
  multiply    Multiply lengths from input trees by a given factor
  setmin      Set a min branch length to all branches with length < cutoff
  setrand     Assign a random length to edges of input trees

Flags:
  -i, --input string    Input tree (default "stdin")
  -o, --output string   Cleared tree output file (default "stdout")
```

clear subcommand
```
Usage:
  gotree brlen clear [flags]

Flags:
  -h, --help   help for clear

Global Flags:
  -i, --input string    Input tree (default "stdin")
  -o, --output string   Cleared tree output file (default "stdout")
```

multiply subcommand
```
Usage:
  gotree brlen multiply [flags]

Flags:
  -f, --factor float   Branch length multiplication factor (default 1)
  -h, --help           help for multiply

Global Flags:
  -i, --input string    Input tree (default "stdin")
  -o, --output string   Cleared tree output file (default "stdout")
```

setmin subcommand
```
Usage:
  gotree brlen setmin [flags]

Flags:
  -h, --help           help for setmin
  -l, --length float   Min Length cutoff

Global Flags:
  -i, --input string    Input tree (default "stdin")
  -o, --output string   Cleared tree output file (default "stdout")
```

setrand subcommand
```
Usage:
  gotree brlen setrand [flags]

Flags:
  -h, --help         help for setrand
  -m, --mean float   Mean of the exponential distribution of branch lengths (default 0.1)
  -s, --seed int     Initial Random Seed (default random)

Global Flags:
  -i, --input string    Input tree (default "stdin")
  -o, --output string   Cleared tree output file (default "stdout")
```

#### Examples

1. Removing branch lengths from a set of 10 trees
```
gotree generate yuletree -s 10 -n 10
```

should output:
```
((Tip4:0.020616211789029896,(Tip7:0.09740195047110385,Tip2:0.015450672710905129):0.12939642466438622):0.0912341925030609,Tip0:0.12959932895259058,((Tip8:0.027845992087631298,(Tip9:0.13492605122032592,Tip3:0.10309294031874587):0.005132906169455565):0.09604804621401375,((Tip6:0.3779897840448691,Tip5:0.1120177846434196):0.029087690784364996,Tip1:0.239082088939295):0.15075207292513051):0.022969404523534506);
(Tip5:0.08565127428804534,Tip0:0.021705093532846203,((Tip6:0.07735928108468929,(Tip7:0.13497893602480787,Tip4:0.0867277155865104):0.0006245685056451973):0.002022779393968061,(Tip2:0.06180495431033288,((Tip8:0.04262994378368574,(Tip9:0.026931366387984084,Tip3:0.05072862497546402):0.00862321921222685):0.11773314378908921,Tip1:0.0059751058197291896):0.003071914457094241):0.16885535253756265):0.5612486967794823);
(Tip6:0.03760877531440096,Tip0:0.08031752649721294,((((Tip5:0.08222558679467758,Tip4:0.010497094428057133):0.027829623402199005,((Tip9:0.0296932046009012,Tip8:0.02260725650758202):0.06699766638404583,Tip3:0.16069947681135582):0.01698933419943515):0.01651428394464055,(Tip7:0.001981710342918917,Tip2:0.127330604774822):0.03371804481869709):0.049385742433923574,Tip1:0.03570249031638065):0.005818760903794294);
(Tip3:0.08432656980817163,Tip0:0.0392134640479104,(((Tip6:0.08005469426954873,Tip5:0.11152589538874742):0.02335921802464195,Tip2:0.049021107917367204):0.08037917732988885,(((Tip9:0.0718913673258644,Tip8:0.06678163204148892):0.009795677738096756,Tip4:0.006547827146823905):0.010394091957264661,(Tip7:0.021074675690276046,Tip1:0.5144747430477511):0.07802407523731346):0.04908580033651718):0.14303241054024363);
(((Tip4:0.00655127976296051,(Tip8:0.24219725619786042,Tip3:0.24838794524945645):0.011720903684812777):0.005446012106551122,((Tip6:0.0539059450392461,Tip5:0.10942616446609373):0.005289510347927662,Tip2:0.18081164518105278):0.0574949822058484):0.05738813637845891,Tip0:0.01808985295613734,(Tip7:0.21455700994877352,(Tip9:0.10644397906022163,Tip1:0.22560831422003544):0.041375531864334826):0.09839681300839312);
(Tip4:0.0031460432390235603,Tip0:0.029576757818637184,(Tip3:0.05598929193748072,((((Tip8:0.012003522473523527,Tip7:0.160964296184698):0.07470084358335415,Tip5:0.18536691292780375):0.05514219780759042,((Tip9:0.05851819289780007,Tip6:0.031143921272508944):0.001732306024297812,Tip2:0.09851551626665338):0.00045107592939929787):0.06843299046011575,Tip1:0.034819467974243366):0.07128205374566994):0.04020059938668683);
(Tip5:0.14713471646659498,Tip0:0.14944717675975588,((Tip8:0.05200350559587298,Tip2:0.09186060066201236):0.019272872692755696,((Tip7:0.017013901822259973,Tip3:0.00722538128881932):0.07059101683797457,((Tip9:0.0827846288269361,Tip4:0.20495216708389177):0.02442099920627207,(Tip6:0.04442441212934027,Tip1:0.19354745472945992):0.010477259056996566):0.0008920078286351969):0.09109441262917872):0.12239266854289019);
((Tip9:0.0704495070914909,Tip6:0.1496980728758482):0.05103728351371519,Tip0:0.0017628717782335215,((((Tip8:0.03738034531983412,Tip4:0.14371286017305163):0.012961268183538757,Tip3:0.13687316515489634):0.009803779593708666,Tip2:0.10132090953011344):0.01600828297222675,((Tip7:0.012628182891023559,Tip5:0.01872508396298026):0.09142355974341702,Tip1:0.21222302845284774):0.001979896506191992):0.08419144019242206);
(Tip9:0.33396108494477056,Tip0:0.01725960043191933,(Tip8:0.15822189216152865,((Tip7:0.2709048460428041,Tip4:0.08286703341169069):0.07754455231612335,((Tip6:0.26229989177453267,Tip2:0.2048675593588131):0.011776444737892341,(Tip3:0.06758944244464968,(Tip5:0.18542831796386594,Tip1:0.0894819395274044):0.062007175935054296):0.01712637494567372):0.024056939891891758):0.12500239546255096):0.08584948451822802);
((((Tip7:0.0027189616811839684,(Tip8:0.14853094179740908,(Tip9:0.01638924236300552,Tip6:0.10624453018746868):0.029000513563448127):0.045781336959022444):0.007633361283361719,Tip5:0.08158495263005563):0.04025350062495641,Tip4:0.05971791841919005):0.09621210503909404,Tip0:0.042406511103404765,(Tip2:0.14262660284294432,(Tip3:0.07913840818550329,Tip1:0.27122555921925007):0.11635871034580876):0.02022686032583556);
```

```
gotree generate yuletree -s 10 -n 10 | gotree brlen clear
```

should output :
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

2. Assign random lengths to a random tree
```
gotree generate yuletree -s 10 -o outtree1.nw
gotree brlen setrand -i outtree.nw -s 13 -o outtree2.nw
gotree draw svg -w 200 -H 200  -i outtree1.nw -o commands/randbrlen_1.svg
gotree draw svg -w 200 -H 200  -i outtree2.nw -o commands/randbrlen_2.svg
```

Initial random Tree               | Random lengths
----------------------------------|-----------------------------------
![Random Tree 1](randbrlen_1.svg) | ![Random Lengths](randbrlen_2.svg)


3. Setting min branch length to 0.1.

```
gotree generate yuletree -s 10 -l 100 -o outtree.nw
gotree draw svg -r -w 200 -H 200 --no-tip-labels -i outtree.nw -o commands/minbrlen_1.svg
gotree brlen setmin -i outtree.nw -l 0.1 | gotree draw svg -r -w 200 -H 200 --no-tip-labels -o commands/minbrlen_2.svg
```

Random Tree                          | Min brlen tree
-------------------------------------|-----------------------------------
![Random Tree](minbrlen_1.svg)       | ![Min brlen tree](minbrlen_2.svg) 


4. Multiplying branch lengths by 3.0

```
gotree generate yuletree -s 10 -l 100 -o outtree.nw
gotree brlen multiply -f 3.0 -i outtree.nw
```
