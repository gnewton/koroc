package main

import (
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/gnewton/pubmedSqlStructs"
	//"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func articleAdder(articleChannel chan ArticlesEnvelope, dbc *DBConnector, db *gorm.DB, commitSize int, wg *sync.WaitGroup) {
	defer wg.Done()
	var err error
	log.Println("Start articleAdder")
	kwjt, err := NewJoinTable("article_keyword", "article_id", "keyword_id", kwjtCreateSql, kwjtInsert)
	if err != nil {
		log.Fatal(err)
	}

	db, err = dbc.Open()
	if err != nil {
		log.Fatal(err)
	}

	tx := db.Begin()
	tx.Debug()

	sdb, err := tx.DB()
	if err != nil {
		log.Fatal(err)
	}
	_, err = sdb.Exec(kwjt.CreateSql())
	if err != nil {
		log.Fatal(err)
	}

	if err = tx.Commit().Error; err != nil {
		log.Fatal(err)
	}
	tx = dbc.DB().Begin()

	//tx = db.Debug()
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

		tmp := articleArray
		//for i := 0; i < len(tmp); i++ {
		for i := 0; i < env.n; i++ {
			article := tmp[i]
			if article == nil {
				log.Println(i, " ******** Article is nil")
				continue
			}

			counter = counter + 1
			totalCount = totalCount + 1
			closeOpenCount = closeOpenCount + 1
			if counter == commitSize {
				tc0 := time.Now()
				if err = tx.Commit().Error; err != nil {
					log.Println(tx.Error)
					log.Fatal("m")
				}
				log.Println("transaction")
				log.Println(tx)

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

			if version, ok := articleIdsInDBCache[article.ID]; !ok {
				log.Println("NEW ]]]]]]]]]]]]]]]]]]]]]]]]]]]", article.ID)
				// New record
				// FIXXX need in update above: remove previous joins and resave
				err = addKeywords(kwjt, tx, article.ID, article.Keywords, article.SourceXMLFilename)
				if err != nil {
					log.Fatal(err)
				}

				//log.Println(".........")
				// if len(article.Keywords) > 0 {
				// 	//for _, kw := range article.Keywords {
				// 	//var zkw pubmedSqlStructs.Keyword
				// 	//result := tx.Where(&pubmedSqlStructs.Keyword{Name: kw.Name}).Find(&zkw)
				// 	//log.Println(".........", zkw, result.RowsAffected, kw.Name)
				// 	// if result.RowsAffected != 0 {
				// 	// 	db.Preload("Keywords").Where(article).Find(&article).Association("Keywords").Append(&zkw)
				// 	// }
				// 	//}
				// }
				articleIdsInDBCache[article.ID] = article.Version
				if doNotWriteToDbFlag {
					continue
				}

				if err := tx.Create(article).Error; err != nil {
					//if err := tx.Save(article).Error; err != nil {
					log.Fatal(err)
				}
			} else {
				// Update record

				// the article version is not more recent than the one already stored
				if article.Version <= version {
					//log.Println("NOT Updating article:", article.ID, article.Version, version)
				} else {
					if doNotWriteToDbFlag {
						continue
					}
					log.Println("Updating article:", article.ID, "old version:", article.Version, "new version:", version)

					var oldArticle pubmedSqlStructs.Article
					if err := tx.Where("ID = ?", article.ID).First(&oldArticle).Error; err != nil {
						log.Fatal(err)
					}
					//if err := tx.Unscoped().Delete(oldArticle).Error; err != nil {
					//if err := tx.Delete(oldArticle).Error; err != nil {
					//	log.Fatal(err)
					//}

					if err := tx.Save(article).Error; err != nil {
						log.Fatal(err)
					}
				}

			}

			if err != nil {
				//log.Println("transaction")
				//log.Println(tx)
				log.Println(err)
				log.Printf("PMID=%d", article.ID)
				if err := tx.Rollback().Error; err != nil {
					log.Println(err)
				}
				log.Fatal(err)
			}

		}
		log.Println("-- END chunk ", chunkCount)
	}
	if !doNotWriteToDbFlag {
		log.Println("Final commit")
		if err := tx.Commit().Error; err != nil {
			log.Fatal(err)
		}

	}

	log.Println("++ END articleAdder")
}

func addKeywords(kwjt *JoinTable, tx *gorm.DB, articleId uint32, keywords []*pubmedSqlStructs.Keyword, sourceFile string) error {
	dups := make(map[string]struct{})
	for i, kw := range keywords {
		log.Println(i, "keyword", kw.ID, kw.MajorTopic, kw.Name)
		//rightId, newItem, joinSql, err := kwjt.AddJoinItem(article.ID, kw.Name)
		rightId, _, joinSql, err := kwjt.AddJoinItem(articleId, kw)
		if err != nil {
			return err
		}
		dupKey := strconv.FormatUint(uint64(articleId), 10) + "_" + strconv.FormatUint(uint64(rightId), 10)
		if _, ok := dups[dupKey]; ok {
			log.Println("------------------------------------------------Duplicate entry:", sourceFile, articleId, rightId, "["+kw.Name+"]", dups)
			continue
		} else {
			dups[dupKey] = empty
		}

		if err = tx.Exec(joinSql).Error; err != nil {
			log.Fatal("articleid=", articleId)
			log.Fatal(joinSql)
			return err
		}
	}
	return nil
}
