package main

import (
	"database/sql"
	"log"
	"strconv"
	"strings"

	"github.com/gnewton/pubmedSqlStructs"
	"github.com/jinzhu/gorm"
)

//const createArticlesTable = "CREATE TABLE \"articles\" (\"abstract\" varchar(255),\"day\" integer,\"id\" integer primary key autoincrement,\"issue\" varchar(255),\"journal_id\" bigint,\"keywords_owner\" varchar(255),\"language\" varchar(255),\"month\" varchar(8),\"title\" varchar(255),\"volume\" varchar(255),\"year\" integer,\"date_revised\" bigint );"

const prepInsertArticle = "INSERT INTO articles (abstract,day,id,issue,journal_id,keywords_owner,language,month,title,volume,year,date_revised) VALUES (?,?,?,?,?,?,?,?,?,?,?,?)"

const prepUpdateArticle = "UPDATE articles set abstract=?,day=?,id=?,issue=?,journal_id=?,keywords_owner=?,language=?,month=?,title=?,volume=?,year=?,date_revised=? where id=?"

func articleAdder2(articleChannel chan []*pubmedSqlStructs.Article, db *sql.DB, commitSize int) {
	//func articleAdder2(articleChannel chan []*pubmedSqlStructs.Article, db *sql.DB, commitSize int) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	stmtInsert, err := tx.Prepare(prepInsertArticle)
	if err != nil {
		log.Fatal(err)
	}

	stmtUpdate, err := tx.Prepare(prepUpdateArticle)
	if err != nil {
		log.Fatal(err)
	}

	chunkCount := 0
	var totalCount int64 = 0
	counter := 0
	for articleArray := range articleChannel {
		log.Println("-- Consuming chunk ", chunkCount)
		chunkCount += 1

		for i := 0; i < len(articleArray); i++ {
			a := articleArray[i]

			if a == nil {
				continue
			}

			// Have we already inserted this article (i.e. is this an update?)
			if _, ok := articleIdsInDBCache[a.ID]; ok {
				log.Println("Updating article:", a.ID)
				_, err = stmtUpdate.Exec(a.Abstract, a.Day, a.ID, a.Issue, a.JournalID, a.KeywordsOwner, a.Language, a.Month, a.Title, a.Volume, a.Year, a.DateRevised, a.ID)
			} else {
				articleIdsInDBCache[a.ID] = a.Version
				_, err = stmtInsert.Exec(a.Abstract, a.Day, a.ID, a.Issue, a.JournalID, a.KeywordsOwner, a.Language, a.Month, a.Title, a.Volume, a.Year, a.DateRevised)

			}

			if err != nil {
				log.Println(err)
				if err.Error() == "UNIQUE constraint failed: articles.id" {
					log.Println("*** ", a.ID, "|||", a.Abstract, a.Day, a.Issue, a.JournalID, a.KeywordsOwner, a.Language, a.Month, a.Title, a.Volume, a.Year, a.DateRevised)
					continue
				}
				log.Println(a.ID, "|||", a.Abstract, a.Day, a.Issue, a.JournalID, a.KeywordsOwner, a.Language, a.Month, a.Title, a.Volume, a.Year, a.DateRevised)
				log.Fatal(err)
			}

			counter = counter + 1
			totalCount = totalCount + 1
			if counter == commitSize {
				counter = 0
				log.Println("************ committing", totalCount)
				tx.Commit()
				stmtInsert.Close()
				stmtUpdate.Close()
				tx, err = db.Begin()
				if err != nil {
					log.Fatal(err)
				}
				stmtInsert, err = tx.Prepare(prepInsertArticle)
				if err != nil {
					log.Fatal(err)
				}
				stmtUpdate, err = tx.Prepare(prepUpdateArticle)
				if err != nil {
					log.Fatal(err)
				}
			}
		}
	}
	tx.Commit()
	stmtInsert.Close()
	stmtUpdate.Close()
	db.Close()

}

func updateArticle(article *pubmedSqlStructs.Article) (sql.Result, error) {
	return nil, nil
}

func insertArticle(article *pubmedSqlStructs.Article) (sql.Result, error) {
	return nil, nil
}

func committor(transactionChannel chan *gorm.DB, done chan bool) {
	for tx := range transactionChannel {
		log.Println("COMMIT starting")
		tx.Commit()
		tx.Close()
		log.Println("COMMIT done")
	}
	done <- true
}

func articleAdder(articleChannel chan []*pubmedSqlStructs.Article, dbc *DBConnector, db *gorm.DB, txChannel chan *gorm.DB, commitSize int, done chan bool) {
	log.Println("Start articleAdder")
	var err error
	db, err = dbc.Open()
	if err != nil {
		log.Fatal(err)
	}
	tx := db.Begin()
	//t0 := time.Now()
	var totalCount int64 = 0
	counter := 0
	chunkCount := 0
	for articleArray := range articleChannel {
		log.Println("-- Consuming chunk ", chunkCount)
		chunkCount += 1

		log.Printf("articleAdder counter=%d", counter)
		log.Printf("TOTAL counter=%d", totalCount)

		log.Println("Commit size=", commitSize)
		if doNotWriteToDbFlag {
			counter = counter + len(articleArray)
			totalCount = totalCount + int64(len(articleArray))
			continue
		}

		tmp := articleArray
		for i := 0; i < len(tmp); i++ {
			article := tmp[i]
			if article == nil {
				//log.Println(i, " ******** Article is nil")
				continue
			}

			counter = counter + 1
			totalCount = totalCount + 1
			closeOpenCount = closeOpenCount + 1
			if counter == commitSize {
				//tc0 := time.Now()
				//tx.Commit()
				log.Printf("Transaction channel length=%d", len(txChannel))
				txChannel <- tx
				//log.Println("transaction")
				//log.Println(tx)
				var err error
				tx, err = dbc.Open()
				if err != nil {
					log.Fatal(err)
				}
				//t1 := time.Now()
				//log.Printf("The commit took %v to run.\n", t1.Sub(tc0))
				//log.Printf("The call took %v to run.\n", t1.Sub(t0))
				//t0 = time.Now()
				counter = 0
				tx = tx.Begin()
				//log.Println("transaction")
				//log.Println(tx)
			}
			var err error
			if version, ok := articleIdsInDBCache[article.ID]; ok {
				// the article version is not more recent than the one already stored
				if article.Version < version {
					log.Println("NOT Updating article:", article.ID, article.Version, version)
				} else {
					log.Println("Updating article:", article.ID, "old version:", article.Version, "new version:", version)

					var oldArticle pubmedSqlStructs.Article
					tx.Where("ID = ?", article.ID).First(&oldArticle)
					tx.Delete(oldArticle)

					//err = tx.Update(article).Error
					err = tx.Create(article).Error
				}

			} else {
				err = tx.Create(article).Error
				articleIdsInDBCache[article.ID] = article.Version
			}

			if err != nil {
				//log.Println("transaction")
				//log.Println(tx)
				log.Println(err)
				tx.Rollback()
				log.Println("\\\\\\\\\\\\\\\\")
				log.Println("[", err, "]")
				log.Printf("PMID=%d", article.ID)
				//if !strings.HasSuffix(err.Error(), "PRIMARY KEY must be unique") {
				//continue
				//}
				//log.Println("Returning from articleAdder")
				//log.Fatal(" Fatal\\\\\\\\\\\\\\\\")
				//return
				tx = db.Begin()
			}

		}
		log.Println("-- END chunk ", chunkCount)
	}
	if !doNotWriteToDbFlag {
		tx.Commit()
		var err error
		tx, err = dbc.Open()
		if err != nil {
			log.Fatal(err)
		}
		makeIndexes(tx)
	}
	close(txChannel)
	db.Close()
	log.Println("-- END articleAdder")
	done <- true
	log.Println("++ END articleAdder")
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
