package main

import (
	"bufio"
	"compress/bzip2"
	"compress/gzip"
	"io"
	"os"
	"strings"
)

func genericReader(filename string) (io.ReadCloser, *os.File, error) {
	if filename == "" {
		reader, err := NewReaderToReaderCloserWrapper(bufio.NewReader(os.Stdin))
		if err != nil {
			return nil, nil, err
		}
		return reader, nil, nil
	}
	file, err := os.Open(filename)
	if err != nil {
		return nil, nil, err
	}
	if strings.HasSuffix(filename, "bz2") {
		//return bufio.NewReader(bzip2.NewReader(bufio.NewReader(file))), file, err
		reader, err := NewReaderToReaderCloserWrapper(bzip2.NewReader(bufio.NewReader(file)))
		if err != nil {
			return nil, nil, err
		}
		return reader, file, err
	}

	if strings.HasSuffix(filename, "gz") {
		reader, err := gzip.NewReader(bufio.NewReader(file))
		if err != nil {
			return nil, nil, err
		}
		return reader, file, err
	}

	reader, err := NewReaderToReaderCloserWrapper(bufio.NewReader(file))
	return reader, file, err
}
