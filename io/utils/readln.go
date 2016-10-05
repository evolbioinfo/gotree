package utils

import (
	"bufio"
)

// Readln returns a single line (without the ending \n)
// from the input buffered reader.
// An error is returned iff there is an error with the
// buffered reader.
func Readln(r *bufio.Reader) (string, error) {
	var (
		isPrefix bool  = true
		err      error = nil
		line, ln []byte
	)
	for isPrefix && err == nil {
		line, isPrefix, err = r.ReadLine()
		ln = append(ln, line...)
	}
	return string(ln), err
}

// ReadUntilSemiColon returns a string (without the ending \n)
// from the input buffered reader, ending at ';' or at end of file
// It allows to read a newick tree on several lines
// An error is returned iff there is an error with the
// buffered reader.
func ReadUntilSemiColon(r *bufio.Reader) (string, error) {
	var (
		isPrefix bool  = true
		err      error = nil
		lastChar byte  = '0'
		line, ln []byte
	)
	for err == nil && (isPrefix || lastChar != ';') {
		line, isPrefix, err = r.ReadLine()
		ln = append(ln, line...)
		if len(ln) > 0 {
			lastChar = ln[len(ln)-1]
		}
	}
	return string(ln), err
}
