package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	//_ "github.com/go-sql-driver/mysql"
	"log"
)

func dbOpen() (*gorm.DB, error) {
	db, err := gorm.Open("sqlite3", dbFileName)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	log.Println("Opening db file: ", dbFileName)
	if sqliteLogFlag {
		db.LogMode(true)
	}

	db.DB()
	db.DB().Ping()
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)

	db.DB()
	db.DB().Ping()
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)

	db.Exec("PRAGMA auto_vacuum = 0;")
	db.Exec("PRAGMA cache_size=10000;")
	db.Exec("PRAGMA count_changes = OFF;")
	db.Exec("PRAGMA encoding = \"UTF-8\";")
	db.Exec("PRAGMA journal_mode = WAL;")
	db.Exec("PRAGMA mmap_size=1099511627776;")
	db.Exec("PRAGMA page_size = 4096;")
	db.Exec("PRAGMA quick_check;")
	db.Exec("PRAGMA shrink_memory")
	db.Exec("PRAGMA synchronous = NORMAL;")
	db.Exec("PRAGMA temp_store = MEMORY;")
	db.Exec("PRAGMA threads = 5;")
	return &db, nil
}

func dbInit() (*gorm.DB, error) {
	db, err := dbOpen()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	log.Printf("%v\n", *db)

	db.CreateTable(&Article{})
	db.CreateTable(&Author{})
	db.CreateTable(&Chemical{})
	db.CreateTable(&Citation{})
	db.CreateTable(&Gene{})
	db.CreateTable(&Journal{})
	db.CreateTable(&Keyword{})
	db.CreateTable(&MeshDescriptor{})
	db.CreateTable(&MeshHeading{})
	db.CreateTable(&MeshQualifier{})
	db.CreateTable(&OtherID{})

	db.Table("Article_Author").AddUniqueIndex("articleAuthor", "article_id", "author_id")
	db.Table("Article_Chemical").AddUniqueIndex("articleChemical", "article_id", "chemical_id")
	db.Table("Article_Citation").AddUniqueIndex("articleCitation", "article_id", "citation_id")
	db.Table("Article_Gene").AddUniqueIndex("articleGene", "article_id", "gene_id")
	db.Table("Article_Keyword").AddUniqueIndex("articleKeyword", "article_id", "keyword_id")
	db.Table("Article_MeshHeading").AddUniqueIndex("articleMeshHeading", "article_id", "mesh_heading_id")

	//db.Exec("CREATE VIRTUAL TABLE pages USING fts4(title, body);")

	return db, nil
}

func makeIndexes(db *gorm.DB) {
	log.Println("makeing indexes START")

	db.Table("articles").AddIndex("articlesYear", "year")
	db.Table("articles").AddIndex("articlesJournalId", "journal_id")
	log.Println("makeing indexes END")
}

func dbCloseOpen(prevDb *gorm.DB) (*gorm.DB, error) {
	prevDb.Close()
	return dbOpen()
}
