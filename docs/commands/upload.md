# Gotree: toolkit and api for phylogenetic tree manipulation

## Commands

### upload
This command uploads a set of input trees to a given server. So far only iTOL is supported, with the subcommand `gotree upload itol`.

iTOL subcommand: It prints the url(s) at which tree(s) is(are) accessible. If `--id` is given, it uploads the tree to the itol account corresponding to the user upload ID, into a specific project given by `--project`. The upload id is accessible by enabling "Batch upload" option in iTOL user settings. If `--id` is not given, it uploads the tree without account, and it will be automatically deleted after 30 days.

If several trees are included in the input file, it will upload all of them, waiting 1 second between each upload to avoid overload of iTOL server.

It is possible to give itol annotation files (see [iTOL documentation](https://itol.embl.de/help.cgi#annot)) to the uploader at the end of command line:
```
gotree upload itol -i tree.nw --name tree --user-id uploadkey --project project annotation*.txt
```

As output, urls are written on stdout and server responses are written on stderr/

So:
```
gotree upload itol -i tree.tree --name tree --user-id uploadkey --project project annotation*.txt > urls
```

Will store only urls in the output file.

#### Usage

General command

```
Usage:
  gotree upload itol [flags] [<annotation file1> ... <annotation file n>]

Flags:
      --name string      iTOL tree name prefix
      --project string   iTOL project to upload the tree
      --user-id string   iTOL User upload id
```

#### Examples

* Generating a random tree, uploading it to iTOL with specific annotations

annotations.txt
```
TREE_COLORS
SEPARATOR SPACE
DATASET_LABEL Clades
COLOR #ff0000
DATA
Tip0 branch #FF0000 normal 5
Tip1 branch #0C00FF normal 5
Tip2 branch #90FF00 normal 5
Tip3 branch #00FF78 normal 5
Tip4 branch #00EAFF normal 5
Tip5 branch #4477AA normal 5
Tip6 branch #006CFF normal 5
Tip7 branch #F000FF normal 5
Tip8 branch #FF006C normal 5
Tip9 branch #FF9C00 normal 5
```


```
gotree generate yuletree -r --seed 10 -o outtree1.nw
gotree upload itol -i outtree1.nw annotations.txt > url
```
It is also possible to download the image from iTOL with `gotree dlimage itol`:

config.txt
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

```
ID=$(basename $(cat url))
gotree dlimage itol -i $ID -f svg -c config.txt -o commands/upload_1.svg
```

It should give:

![iTOL image](upload_1.svg) 
