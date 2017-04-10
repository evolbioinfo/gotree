# Gotree: toolkit and api for phylogenetic tree manipulation

## Commands

### dlimage
This command downloads a tree image from a server. So far the only sub-commands is
* `gotree dlimage itol`, which downloads a tree image from [iTOL](http://itol.embl.de/), given a tree id (`-i`) and a configuration file (`-c`). The configuration file is a tab separated key/value file corresponding to the iTOL [api optionql parameters](http://itol.embl.de/help.cgi#bExOpt).

#### Usage

```
Usage:
  gotree dlimage itol [flags]

Flags:
  -c, --config string   Itol image config file

Global Flags:
  -f, --format string    Image format (png, pdf, eps, svg) (default "pdf")
  -o, --output string    Image output file
  -i, --tree-id string   Tree id to download
```

#### Example

* We generate a tree that we upload to iTOL and get the tree ID
```
gotree generate yuletree -s 10 | gotree upload itol > url
TREEID=$(basename $(cat url))
```

* We write a configuration file `config.txt` for iTOL
```
display_mode	2
label_display	1
align_labels	0
ignore_branch_length	0
bootstrap_display	 1
bootstrap_type	1
bootstrap_symbol	1
bootstrap_slider_min	0.7
bootstrap_slider_max	1
bootstrap_symbol_min	20
bootstrap_symbol_max	20
bootstrap_symbol_color	#c8c7fc
current_font_size	30
line_width	2
inverted	0
```

* We download the tree from iTOL
```
gotree dlimage itol -i $ID -f svg -c config.txt -o commands/dlimage_1.svg
```

![Image from iTOL](dlimage_1.svg)
