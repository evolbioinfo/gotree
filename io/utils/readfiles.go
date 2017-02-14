package utils

import (
	"bufio"
	"compress/gzip"
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
func GetReader(inputfile string) (*os.File, *bufio.Reader, error) {
	var reader *bufio.Reader
	if f, err := OpenFile(inputfile); err != nil {
		return nil, nil, err
	} else {

		if strings.HasSuffix(f.Name(), ".gz") {
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
}
