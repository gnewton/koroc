package main

import (
	"bufio"
	"compress/bzip2"
	"compress/gzip"
	"errors"
	"io"
	"log"
	"os"
	"strings"
)

type FileReader struct {
	reader io.ReadCloser
	file   *os.File
}

func NewFileReader(filename string) (io.Reader, error) {
	return NewFileReaderSize(filename, 4096)
}

func NewFileReaderSize(filename string, size int) (io.Reader, error) {
	if size < 1 {
		return nil, errors.New("Size cannot be <1")
	}

	r := new(FileReader)
	var err error
	r.file, err = os.Open(filename)
	if err != nil {
		return nil, err
	}
	bufioReader := bufio.NewReader(r.file)
	if strings.HasSuffix(filename, "bz2") {
		r.reader, err = NewReaderToReaderCloserWrapper(bzip2.NewReader(bufioReader))
		if err != nil {
			return nil, err
		}
	}

	if strings.HasSuffix(filename, "gz") {
		r.reader, err = gzip.NewReader(bufioReader)
		if err != nil {
			return nil, err
		}
	}
	return r, nil
}

func (r FileReader) Close() (err error) {
	if r.reader != nil {
		err = r.reader.Close()
		if err != nil {
			log.Println(err)
		}
	}

	if r.file != nil {
		err = r.file.Close()
		if err != nil {
			log.Println(err)
		}
	}
	return err
}

func (r FileReader) Read(p []byte) (n int, err error) {
	return 0, nil
}

type ReaderToReaderCloserWrapper struct {
	reader io.Reader
}

func NewReaderToReaderCloserWrapper(r io.Reader) (io.ReadCloser, error) {
	if r == nil {
		return nil, errors.New("Reader cannot be nil")
	}
	rw := new(ReaderToReaderCloserWrapper)
	rw.reader = r
	return rw, nil
}

func (rw *ReaderToReaderCloserWrapper) Close() error {
	return nil
}

func (rw *ReaderToReaderCloserWrapper) Read(p []byte) (n int, err error) {
	if rw.reader == nil {
		return -1, errors.New("Underlying reader is nil")
	}
	return rw.reader.Read(p)
}
