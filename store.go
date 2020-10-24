package main

import (
	//"database/sql"
	//	"errors"
	"log"
	//	"strconv"
	"sync"
	"time"

	//"github.com/gnewton/
	//"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func articleAdder(articleChannel chan ArticlesEnvelope, dbc *DBConnector, db *gorm.DB, commitSize int, wg *sync.WaitGroup) {
	defer wg.Done()
	var err error
	log.Println("Start articleAdder")
	//kwjt, err := NewJoinTable("article_keyword", "article_id", "keyword_id", kwjtCreateSql, kwjtInsert, kwtSavePreparedSql)
	if err != nil {
		log.Fatal(err)
	}

	gormdb, err := dbc.Open()
	if err != nil {
		log.Fatal(err)
	}

	gormtx := gormdb.Begin()
	gormtx.Debug()

	// sdb, err := gormtx.DB()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// create, index := kwjt.CreateSql()
	// _, err = sdb.Exec(create)
	// if err != nil {
	// 	log.Println(create)
	// 	log.Fatal(err)
	// }
	// _, err = sdb.Exec(index)
	// if err != nil {
	// 	log.Println(index)
	// 	log.Fatal(err)
	// }

	if err = gormtx.Commit().Error; err != nil {
		log.Fatal(err)
	}
	xdb, err := dbc.DB().DB()
	if err != nil {
		log.Fatal(err)
	}

	tx, err := xdb.Begin()
	if err != nil {
		log.Fatal(err)
	}

	// kwjtDeletePrepared, err := kwjt.DeletePreparedStatement(tx)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// kwjtSavePrepared, err := kwjt.SavePreparedStatement(tx)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// kwPrepared, err := Keyword{}.MakePreparedStatement(tx)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	articleInsertPrepared, articleDeletePrepared, err := Article{}.MakePreparedStatements(tx)
	if err != nil {
		log.Fatal(err)
	}
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
		for i := 0; i < len(articleArray); i++ {
			//for i := 0; i < env.n; i++ {
			article := tmp[i]

			counter = counter + 1
			totalCount = totalCount + 1
			closeOpenCount = closeOpenCount + 1
			if counter == commitSize {
				if err := articleInsertPrepared.Close(); err != nil {
					log.Fatal(err)
				}
				if err := articleDeletePrepared.Close(); err != nil {
					log.Fatal(err)
				}

				// if err := kwjtSavePrepared.Close(); err != nil {
				// 	log.Fatal(err)
				// }
				// if err := kwjtDeletePrepared.Close(); err != nil {
				// 	log.Fatal(err)
				// }

				log.Println("------------------------------------------start commit")
				tc0 := time.Now()
				if err = tx.Commit(); err != nil {
					//log.Println(tx.Error)
					log.Fatal(err)
				}
				log.Println("------------------------------------------end commit")
				log.Println(tx)

				t1 := time.Now()
				log.Printf("The commit took %v to run.\n", t1.Sub(tc0))
				log.Printf("The call took %v to run.\n", t1.Sub(t0))
				t0 = time.Now()
				counter = 0

				tx, err = xdb.Begin()
				if err != nil {
					log.Fatal(err)
				}
				log.Println("transaction")
				log.Println(tx)
				articleInsertPrepared, articleDeletePrepared, err = Article{}.MakePreparedStatements(tx)
				if err != nil {
					log.Fatal(err)
				}
				// kwPrepared, err = Keyword{}.MakePreparedStatement(tx)
				// if err != nil {
				// 	log.Fatal(err)
				// }
				// kwjtDeletePrepared, err = kwjt.DeletePreparedStatement(tx)
				// if err != nil {
				// 	log.Fatal(err)
				// }
				// kwjtSavePrepared, err = kwjt.SavePreparedStatement(tx)
				// if err != nil {
				// 	log.Fatal(err)
				// }
			}
			var err error

			if version, ok := articleIdsInDBCache[article.ID]; !ok {
				// New record
				// err = addKeywords(kwPrepared, kwjt, tx, article.ID, article.Keywords, article.SourceXMLFilename)
				// if err != nil {
				// 	log.Fatal(err)
				// }
				articleIdsInDBCache[article.ID] = article.Version
				if doNotWriteToDbFlag {
					continue
				}

				err = article.Save(articleInsertPrepared)
				if err != nil {
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

					// Delete previous version of article
					err = article.Delete(articleDeletePrepared)
					if err != nil {
						log.Fatal(err)
					}
					// err = kwjt.Delete(article.ID, kwjtDeletePrepared)
					// if err != nil {
					// 	log.Fatal(err)
					// }

					// err = addKeywords(kwPrepared, kwjt, tx, article.ID, article.Keywords, article.SourceXMLFilename)
					// if err != nil {
					// 	log.Fatal(err)
					// }
					if err = article.Save(articleInsertPrepared); err != nil {
						log.Fatal(err)
					}
				}

			}

			if err != nil {
				//log.Println("transaction")
				//log.Println(tx)
				log.Println(err)
				log.Printf("PMID=%d", article.ID)
				if err := tx.Rollback(); err != nil {
					log.Println(err)
				}
				log.Fatal(err)
			}

		}
		log.Println("-- END chunk ", chunkCount)
	}
	if !doNotWriteToDbFlag {
		log.Println("Final commit")
		if err := tx.Commit(); err != nil {
			log.Fatal(err)
		}

	}

	log.Println("++ END articleAdder")
}

// func addKeywords(kwPrepared *sql.Stmt, kwjt *JoinTable, tx *sql.Tx, articleId uint32, keywords []*Keyword, sourceFile string) error {
// 	dups := make(map[string]struct{})
// 	for _, kw := range keywords {

// 		//rightId, newItem, joinSql, err := kwjt.AddJoinItem(article.ID, kw.Name)
// 		rightId, newKeyword, joinSql, err := kwjt.AddJoinItem(articleId, kw)
// 		if err != nil {
// 			return err
// 		}
// 		if newKeyword {
// 			kw.ID = rightId
// 			kw.Save(kwPrepared)
// 		}
// 		dupKey := strconv.FormatUint(uint64(articleId), 10) + "_" + strconv.FormatUint(uint64(rightId), 10)
// 		if _, ok := dups[dupKey]; ok {
// 			log.Println("------------------------------------------------Duplicate entry:", sourceFile, articleId, rightId, "["+kw.Name+"]", dups)
// 			continue
// 		} else {
// 			dups[dupKey] = empty
// 		}

// 		//if err = tx.Exec(joinSql).Error; err != nil {
// 		r, err := tx.Exec(joinSql)
// 		if err != nil {
// 			log.Println(err)
// 			log.Println("articleid=", articleId)
// 			log.Fatal("Fatal ", joinSql)
// 		}
// 		n, err := r.RowsAffected()
// 		if err != nil {
// 			log.Println("articleid=", articleId)
// 			log.Fatal(joinSql)
// 		}
// 		if n != 1 {
// 			log.Println("articleid=", articleId)
// 			log.Fatal(errors.New("Join insert affected >1 record"))
// 		}

// 	}
// 	return nil
// }
