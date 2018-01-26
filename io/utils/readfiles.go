package utils

import (
	"bufio"
	"compress/gzip"
	"io"
	"net/http"
	"os"
	"strings"
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

/* Returns the opened file and a buffered reader (gzip or not) for the file */
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
