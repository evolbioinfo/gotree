# Gotree: toolkit and api for phylogenetic tree manipulation

## Commands

### reformat
This command reformats an input tree file into different formats.

So far, formats can be :
- Input formats: Newick, Nexus
- Output formats: Newick, Nexus.


#### Usage

General command
```
Usage:
  gotree reformat [command]

Available Commands:
  newick      Reformats an input tree file into Newick format
  nexus       Reformats an input tree file into Nexus format

Flags:
  -f, --format string   Input format (newick, nexus) (default "newick")
  -h, --help            help for reformat
  -i, --input string    Input tree (default "stdin")
  -o, --output string   Output file (default "stdout")
```


#### Examples

* Reformat input nexus format into newick
```
gotree reformat newick -i input.nexus -f nexus -o output.nw
```

* Reformat input nexick format into nexus
```
gotree reformat nexus -i input.nw -f newick -o output.nexus
```
