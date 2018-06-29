package main

import (
	"database/sql"
	"github.com/gnewton/pubmedSqlStructs"
	"log"
	"strconv"
	"strings"
	"sync"
)

type Article struct {
	pubmedSqlStructs.Article
	pmidString string
}

func articleAdder3(articleChannel chan []*Article, db *sql.DB, commitSize int, wg *sync.WaitGroup) {
	defer wg.Done()

	ap, err := NewArticlePersist(db, commitSize)
	if err != nil {
		log.Fatal(err)
	}

	var totalCount int64 = 0
	for articleArray := range articleChannel {
		log.Println("X", len(articleArray), totalCount, ap.commitCounter)
		var nilCount int64 = 0
		for i := 0; i < len(articleArray); i++ {
			article := articleArray[i]
			if article == nil {
				nilCount++
			} else {
				totalCount++
				err := ap.Save(article)
				if err != nil {
					log.Fatal(err)
				}
			}
		}
		log.Println("nilCount=", nilCount)
	}

	if ap.commitCounter > 0 {
		log.Println("Starting commit")
		err = ap.tx.Commit()
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Done commit")
		ap.tx, ap.insertPrep, ap.updatePrep, err = newTx(ap.db, prepInsertArticle, prepUpdateArticle)
		ap.commitCounter = 0

	}
}

var dateSeasonYear = []string{"Summer", "Winter", "Spring", "Fall"}

func medlineDate2Year(md string) int {
	// Other cases:
	//   Fall 2017; 8/15/12; Spring 2017; Summer 2017; Fall 2017;

	// case <MedlineDate>1952 Mar-Apr</MedlineDate>
	var year int
	var err error

	// case 2000-2001
	//log.Println(md)

	for i, _ := range dateSeasonYear {
		if strings.HasPrefix(md, dateSeasonYear[i]) {
			return seasonYear(md)
		}
	}

	if len(md) == 5 {
		year, err = strconv.Atoi(md)
		if err != nil {
			log.Println("error!! ", err)
			year = 0
		}
		return year
	}
	if len(md) >= 5 && string(md[4]) == string('-') {
		yearStrings := strings.Split(md, "-")
		//case 1999-00
		if len(yearStrings[1]) != 4 {
			year, err = strconv.Atoi(yearStrings[0])
		} else {
			year, err = strconv.Atoi(yearStrings[1])
		}
		if err != nil {
			log.Println("error!! ", err)
		}
	} else {
		// case 1999 June 6
		yearString := strings.TrimSpace(strings.Split(md, " ")[0])
		yearString = yearString[0:4]
		//year, err = strconv.Atoi(strings.TrimSpace(strings.Split(md, " ")[0]))
		year, err = strconv.Atoi(yearString)
		if err != nil {
			log.Println("error!! yearString=[", yearString, "]", err)
		}
	}
	if year == 0 {
		log.Println("medlineDate2Year [", md, "] [", strings.TrimSpace(string(md[4])), "]")
	}
	return year

}

func seasonYear(md string) int {
	parts := strings.Split(md, " ")
	year, err := strconv.Atoi(parts[1])
	if err != nil {
		log.Fatal(err)
	}
	return year
}
