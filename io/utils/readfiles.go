package utils

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/fredericlemoine/gotree/download"
)

func OpenFile(inputfile string) (*os.File, error) {
	var infile *os.File
	var err error
	if inputfile == "" || inputfile == "stdin" || inputfile == "-" {
		infile = os.Stdin
	} else {
		infile, err = os.Open(inputfile)
		if err != nil {
			return nil, err
		}
	}
	return infile, nil
}

// Returns the opened file and a buffered reader (gzip or not) for the file
//
// The file may be a remote file:
//     * http://
//     * itol://<itol id>
//
// Or a local file
func GetReader(inputfile string) (io.Closer, *bufio.Reader, error) {
	var reader *bufio.Reader

	var err error
	var f io.ReadCloser

	if isHttpFile(inputfile) {
		var res *http.Response
		if res, err = http.Get(inputfile); err != nil {
			return nil, nil, err
		}
		f = res.Body
	} else if isItol(inputfile) {
		var b []byte
		dl := download.NewItolImageDownloader(make(map[string]string))
		itolid := strings.TrimPrefix(inputfile, "itol://")
		if b, err = dl.Download(itolid, download.TXTFORMAT_NEWICK); err != nil {
			return nil, nil, err
		}
		f = ioutil.NopCloser(bytes.NewReader(b))
	} else {
		if f, err = OpenFile(inputfile); err != nil {
			return nil, nil, err
		}
	}

	if strings.HasSuffix(inputfile, ".gz") {
		if gr, err := gzip.NewReader(f); err != nil {
			return nil, nil, err
		} else {
			reader = bufio.NewReader(gr)
		}
	} else {
		reader = bufio.NewReader(f)
	}
	return f, reader, nil
}

func isHttpFile(file string) bool {
	return strings.HasPrefix(file, "http://") ||
		strings.HasPrefix(file, "https://")
}

func isItol(file string) bool {
	return strings.HasPrefix(file, "itol://")
}
