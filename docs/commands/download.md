# Gotree: toolkit and api for phylogenetic tree manipulation

## Commands

### download
This command downloads trees or tree images from a given source. Two subcommands so far:
* `gotree download itol`, which downloads a tree file/image from [iTOL](http://itol.embl.de/), given a tree id (`-i`) and a configuration file (`-c`). Formats may be "png", "eps", "svg", "pdf", "newick", "nexus", "phyloxml". The configuration file (used only with image formats) is a tab separated key/value file corresponding to the iTOL [api optional parameters](http://itol.embl.de/help.cgi#bExOpt).
* `gotree download ncbitax`, which downloads the ncbi taxonomy from NCBI ftp server and converts it in Newick format. Internal and tip node names are NCBI names given by the file "names.dmp". Please not that to conform to Newick format, following character are replaced by `_` : `()[]:, ;`. Moreover, the NCBI taxononomy may have species (~tips) with children (ex: [taxid:9606](https://www.ncbi.nlm.nih.gov/Taxonomy/Browser/wwwtax.cgi?mode=Tree&id=9606)). These cases are resolved by Gotree by adding a new corresponding tip.

#### Usage

* General command
```
Usage:
  gotree download [command]

Available Commands:
  itol        Download a tree image from iTOL
  ncbitax     Downloads the full ncbi taxonomy in newick format
```


* itol subcommand
```
Usage:
  gotree download itol [flags]

Flags:
  -c, --config string   Itol image config file

Global Flags:
  -f, --format string    Image format (png, pdf, eps, svg, newick, nexus, phyloxml) (default "pdf")
  -o, --output string    Image output file
  -i, --tree-id string   Tree id to download
```

* ncbitax subcommand
```
Usage:
  gotree download ncbitax [flags]

Flags:
      --map string      Output mapping file between taxid and species name (tab separated) (default "none")
      --nodes-taxid     Keeps tax id as internal nodes identifiers
  -o, --output string   NCBI newick output file (default "stdout")
      --tips-taxid      Keeps tax id as tip names
```

#### Example

* We generate a tree that we upload to iTOL and get the tree ID
```
gotree generate yuletree --seed 10 | gotree upload itol > url
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
gotree download itol -i $ID -f svg -c config.txt -o commands/download_1.svg
```

![Image from iTOL](download_1.svg)

* We download NCBI taxonomy, prune it to get the same taxa than the tree to test (test.nw) and we compare internal branches to the NCBI topology:

test.nw:
```
(((((Hylobates_pileatus:0.23988592,(Pongo_pygmaeus_abelii:0.11809071,(Gorilla_gorilla_gorilla:0.13596645,(Homo_sapiens:0.11344407,Pan_troglodytes:0.11665038)0.62:0.02364476)0.78:0.04257513)0.93:0.15711475)0.56:0.03966791,(Macaca_sylvanus:0.06332916,(Macaca_fascicularis_fascicularis:0.07605049,(Macaca_mulatta:0.06998962,Macaca_fuscata:0)0.98:0.08492791)0.47:0.02236558)0.89:0.11208218)0.43:0.0477543,Saimiri_sciureus:0.25824985)0.71:0.14311537,(Tarsius_tarsier:0.62272677,Lemur_sp.:0.40249393)0.35:0)0.62:0.077084225,(Mus_musculus:0.4057381,Bos_taurus:0.65776307)0.62:0.077084225);
```

```bash
gotree download ncbitax -o ncbi.nw
gotree prune -i ncbi.nw -c test.nw -o ncbi_prune.nw
gotree annotate -i test.nw -c ncbi_prune.nw -o test_annotated.nw
gotree draw text -i test_annotated.nw -w 100
```

It should give a tree like that:
```
                                         +-------------------------------- Hylobates_pileatus                 
                                    +----|Hominoidea_0_5                                                      
                                    |    |                     +--------------- Pongo_pygmaeus_abelii         
                                    |    +---------------------|Pongidae_0_4                                  
                                    |                          |     +----------------- Gorilla_gorilla_gorill
                                    |                          +-----|Homo/Pan/Gorilla_group_0_3              
                             +------|Catarrhini_0_5                  |  +-------------- Homo_sapiens          
                             |      |                                +--|Pan_troglodytes_1_2                  
                             |      |                                   +--------------- Pan_troglodytes      
                             |      |                                                                         
                             |      |              +-------- Macaca_sylvanus                                  
                             |      +--------------|Macaca_0_4                                                
          +------------------|Simiiformes_0_4      |  +---------- Macaca_fascicularis_fascicularis            
          |                  |                     +--|Macaca_1_3                                             
          |                  |                        |           +-------- Macaca_mulatta                    
          |                  |                        +-----------|Macaca_fuscata_1_2                         
+---------|Primates_0_12     |                                     Macaca_fuscata                             
|         |                  |                                                                                
|         |                  +----------------------------------- Saimiri_sciureus                            
|         |                                                                                                   
|         |------------------------------------------------------------------------------------ Tarsius_tarsie
|         |Tarsius_tarsier_1_2                                                                                
|         +------------------------------------------------------ Lemur_sp.                                   
|                                                                                                             
|         +------------------------------------------------------ Mus_musculus                                
+---------|                                                                                                   
          +---------------------------------------------------------------------------------------- Bos_taurus
```
