# Gotree: toolkit and api for phylogenetic tree manipulation

## Commands

### comment
This command modifies branch/node comments of input trees.

#### Usage

```
Usage:
  gotree comment [command]

Available Commands:
  clear       Removes node/tip comments

Flags:
  -h, --help            help for comment
  -i, --input string    Input tree (default "stdin")
  -o, --output string   Cleared tree output file (default "stdout")
```

clear subcommand
```
Usage:
  gotree comment clear [flags]

Flags:
  -h, --help   help for clear

Global Flags:
  -i, --input string    Input tree (default "stdin")
  -o, --output string   Cleared tree output file (default "stdout")
```

#### Examples

* Removing node comments from an input tree
```
echo "(t1[c1],t2[c2],(t3[c3],t4[c4])[c5]);" | gotree comment clear
```

Should print:
```
(t1,t2,(t3,t4));
```
