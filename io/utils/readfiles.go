package utils

import (
	"os"
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

func GetReader(inputfile string) (*io.Reader, err) {
	var reader *bufio.Reader
	if f, err := OpenFile(inputfile); err != nil {
		return nil, err
	} else {

		if strings.HasSuffix(refTreeFile.Name(), ".gz") {
			if gr, err := gzip.NewReader(refTreeFile); err != nil {
				return nil, err
			} else {
				reader = bufio.NewReader(gr)
			}
		} else {
			reader = bufio.NewReader(refTreeFile)
		}
	}
	return reader
}
