# Gotree: toolkit and api for phylogenetic tree manipulation

## API

### download

#### Download from [iTOL](https://itol.embl.de/) (SVG)
Downloading a tree image from iTOL
```go
package main

import (
	"io/ioutil"

	"github.com/evolbioinfo/gotree/download"
	"github.com/evolbioinfo/gotree/io"
)

func main() {
	var config map[string]string
	var dl *download.ItolImageDownloader
	var dltreeid string = "<itol tree id>"
	var dloutput string = "image.svg"
	var b []byte
	var err error

	// See https://itol.embl.de/help.cgi#bExOpt for all config options
	config = make(map[string]string)
	config["display_mode"] = "3"          // Unrooted
	config["ignore_branch_length"] = "0"  // Take branch length into account
	config["bootstrap_display"] = "1"     // Display bootstrap supports
	config["bootstrap_type"] = "1"        // Display Bootstrap as symbols
	config["bootstrap_symbol_min"] = "1"  // Min bootstrap symbol size
	config["bootstrap_symbol_max"] = "20" // Max bootstrap symbol size

	dl = download.NewItolImageDownloader(config)
	b, err = dl.Download(dltreeid, download.IMGFORMAT_SVG)
	if err != nil {
		io.ExitWithMessage(err)
	}
	ioutil.WriteFile(dloutput, b, 0644)
}
```

#### Download from [iTOL](https://itol.embl.de/) (Newick format)
Downloading a tree image from iTOL
```go
package main

import (
	"io/ioutil"

	"github.com/evolbioinfo/gotree/download"
	"github.com/evolbioinfo/gotree/io"
)

func main() {
	var dl *download.ItolImageDownloader
	var dltreeid string = "<itol tree id>"
	var dloutput string = "tree.nhx"
	var b []byte
	var err error

	dl = download.NewItolImageDownloader(make(map[string]string))
	b, err = dl.Download(dltreeid, download.TXTFORMAT_NEWICK)
	if err != nil {
		io.ExitWithMessage(err)
	}
	ioutil.WriteFile(dloutput, b, 0644)
}
```


#### Download and convert NCBI taxonomy

```go
package main

import (
	"github.com/evolbioinfo/gotree/download"
)

func main(){
	dl := download.NewNcbiTreeDownloader()
	t, err := dl.Download("")
	if err != nil {
		panic(err)
	}
	f := openWriteFile("ncbi.nw")
	f.WriteString(t.Newick() + "\n")
	f.Close()
}
```


#### Download a tree from Panther database

```go
package main

import (
	"github.com/evolbioinfo/gotree/download"
)

func main(){

	dl := download.NewPantherTreeDownloader()
	if t, err := dl.Download("PTHR10000"); err != nil {
		panic(err)
	}
	f := openWriteFile("ncbi.nw")
	f.WriteString(t.Newick() + "\n")
	f.Close()
}
```
