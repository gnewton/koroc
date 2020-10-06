package main

import (
	"database/sql"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gnewton/pubmedSqlStructs"
	"github.com/jinzhu/gorm"
)

func updateArticle(article *pubmedSqlStructs.Article) (sql.Result, error) {
	return nil, nil
}

func insertArticle(article *pubmedSqlStructs.Article) (sql.Result, error) {
	return nil, nil
}

//func articleAdder(articleChannel chan []*pubmedSqlStructs.Article, dbc *DBConnector, db *gorm.DB, commitSize int, wg *sync.WaitGroup) {
func articleAdder(articleChannel chan ArticlesEnvelope, dbc *DBConnector, db *gorm.DB, commitSize int, wg *sync.WaitGroup) {
	defer wg.Done()

	log.Println("Start articleAdder")
	var err error
	db, err = dbc.Open()
	if err != nil {
		log.Fatal(err)
	}
	db.LogMode(true)
	tx := db.Begin()
	t0 := time.Now()
	var totalCount int64 = 0
	counter := 0
	chunkCount := 0
	for env := range articleChannel {
		articleArray := env.articles
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
				tc0 := time.Now()
				tx.Commit()
				if tx.Error != nil {
					log.Println(tx.Error)
					handleDbErrors(tx.GetErrors())
				}
				log.Println("transaction")
				log.Println(tx)

				// var err error
				// tx, err = dbc.Open()
				// if err != nil {
				// 	log.Fatal(err)
				// }
				t1 := time.Now()
				log.Printf("The commit took %v to run.\n", t1.Sub(tc0))
				log.Printf("The call took %v to run.\n", t1.Sub(t0))
				t0 = time.Now()
				counter = 0
				tx = dbc.DB().Begin()
				log.Println("transaction")
				log.Println(tx)
			}
			var err error
			// New Version
			if version, ok := articleIdsInDBCache[article.ID]; ok {
				log.Println("$$$$$$$$$$$$$$$$$$$$$$$$$")
				// the article version is not more recent than the one already stored
				if article.Version <= version {
					log.Println("NOT Updating article:", article.ID, article.Version, version)
				} else {
					log.Println("Updating article:", article.ID, "old version:", article.Version, "new version:", version)

					var oldArticle pubmedSqlStructs.Article
					if err := tx.Where("ID = ?", article.ID).First(&oldArticle).Error; err != nil {
						log.Fatal(err)
					}
					if err := tx.Unscoped().Delete(oldArticle).Error; err != nil {
						//if err := tx.Delete(oldArticle).Error; err != nil {
						log.Fatal(err)
					}

					if err := tx.Save(article).Error; err != nil {
						// if err := tx.Update(article).Error; err != nil {
						log.Fatal(err)
					}

					//if err := tx.Create(article).Error; err != nil {
					//log.Fatal(err)
					//}
				}

			} else {
				log.Println(".........")
				if len(article.Keywords) > 0 {
					for _, kw := range article.Keywords {
						var zkw pubmedSqlStructs.Keyword
						result := tx.Where(&pubmedSqlStructs.Keyword{Name: kw.Name}).Find(&zkw)
						log.Println(".........", zkw, result.RowsAffected, kw.Name)
						if result.RowsAffected != 0 {
							db.Preload("Keywords").Where(article).Find(&article).Association("Keywords").Append(&zkw)
						}
					}
				}
				if err := tx.Create(article).Error; err != nil {
					//if err := tx.Save(article).Error; err != nil {
					log.Fatal(err)
				}
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
				log.Fatal(err)
			}

		}
		log.Println("-- END chunk ", chunkCount)
	}
	if !doNotWriteToDbFlag {
		log.Println("Final commit")
		db := tx.Commit()
		if db.Error != nil {
			log.Fatal(db.Error)
		}
		if tx.Error != nil {
			log.Fatal(tx.Error)
		}

		db, err = dbc.Open()
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Making indexes")
		makeIndexes(db)
	}
	log.Println("-- Close DB")
	db.Close()

	log.Println("++ END articleAdder")
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

func handleDbErrors(errors []error) {
	for i, e := range errors {
		log.Println(i, e)
	}
	log.Fatal("Errors")
}
