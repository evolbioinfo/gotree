# Gotree: toolkit and api for phylogenetic tree manipulation

## API

### dlimage

Downloading a tree image from iTOL
```go
package main

import (
	"io/ioutil"

	"github.com/fredericlemoine/gotree/download"
	"github.com/fredericlemoine/gotree/io"
)

func main() {
	var config map[string]string
	var dl *download.ItolImageDownloader
	var dltreeid string = "<itol tree id>"
	var dloutput string = "image.svg"
	var b []byte
	var err error

	// See http://itol.embl.de/help.cgi#bExOpt for all config options
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
