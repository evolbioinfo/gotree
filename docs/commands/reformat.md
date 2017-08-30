# Gotree: toolkit and api for phylogenetic tree manipulation

## Commands

### reformat
This command reformats an input tree file into different formats.

So far, formats can be :
- Input formats: Newick, Nexus, PhyloXML
- Output formats: Newick, Nexus, PhyloXML.


#### Usage

General command
```
Usage:
  gotree reformat [command]

Available Commands:
  newick      Reformats an input tree file into Newick format
  nexus       Reformats an input tree file into Nexus format
  phyloxml    Reformats an input tree file into PhyloXML format

Flags:
  -f, --format string   Input format (newick, nexus, phyloxml) (default "newick")
  -h, --help            help for reformat
  -i, --input string    Input tree (default "stdin")
  -o, --output string   Output file (default "stdout")
```


#### Examples

* Reformat input nexus format into newick
```
gotree reformat newick -i input.nexus -f nexus -o output.nw
```

* Reformat input newick format into nexus
```
gotree reformat nexus -i input.nw -f newick -o output.nexus
```

* Reformat input phyloxml format into newick
```
gotree reformat newick -i input.xml -f phyloxml -o output.nexus
```
