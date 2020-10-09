package main

import (
	"bufio"
	"compress/bzip2"
	"compress/gzip"
	"io"
	"log"
	"os"
	"strconv"
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

var dateSeasonYear = []string{"Summer", "Winter", "Spring", "Fall"}

func medlineDate2Year(md string) uint16 {
	// Other cases:
	//   Fall 2017; 8/15/12; Spring 2017; Summer 2017; Fall 2017;

	// case <MedlineDate>1952 Mar-Apr</MedlineDate>
	var year uint16
	var err error

	// case 2000-2001
	//log.Println(md)

	for i, _ := range dateSeasonYear {
		if strings.HasPrefix(md, dateSeasonYear[i]) {
			return seasonYear(md)
		}
	}
	var tmp uint64
	if len(md) == 5 {
		//year, err = strconv.Atoi(md)
		tmp, err = strconv.ParseUint(md, 10, 16)
		if err != nil {
			log.Println("error!! ", err)
			tmp = 0
		}
		return uint16(tmp)
	}

	if len(md) >= 5 && string(md[4]) == string('-') {
		yearStrings := strings.Split(md, "-")
		//case 1999-00
		if len(yearStrings[1]) != 4 {
			//year, err = strconv.Atoi(yearStrings[0])
			tmp, err = strconv.ParseUint(yearStrings[0], 10, 16)
		} else {
			//year, err = strconv.Atoi(yearStrings[1])
			tmp, err = strconv.ParseUint(yearStrings[1], 10, 16)
		}
		if err != nil {
			log.Println("error!! ", err)
		}
		year = uint16(tmp)
	} else {
		// case 1999 June 6
		yearString := strings.TrimSpace(strings.Split(md, " ")[0])
		yearString = yearString[0:4]
		//year, err = strconv.Atoi(strings.TrimSpace(strings.Split(md, " ")[0]))
		//year, err = strconv.Atoi(yearString)
		tmp, err = strconv.ParseUint(yearString, 10, 16)
		if err != nil {
			log.Println("error!! yearString=[", yearString, "]", err)
		}
		year = uint16(tmp)
	}
	if year == 0 {
		log.Println("medlineDate2Year [", md, "] [", strings.TrimSpace(string(md[4])), "]")
	}
	return year

}

func seasonYear(md string) uint16 {
	parts := strings.Split(md, " ")
	tmp, err := strconv.ParseUint(parts[1], 10, 16)
	if err != nil {
		log.Fatal(err)
	}

	return uint16(tmp)
}
