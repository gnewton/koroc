package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	//_ "github.com/go-sql-driver/mysql"
	"log"
)

func dbOpen() (*gorm.DB, error) {
	//db, err := gorm.Open("sqlite3", "/tmp/gorm_15000-x.db")

	//db, err := gorm.Open("mysql", "gnewton:@/pubmed?charset=utf8&parseTime=True&loc=Local")
	//db, err := gorm.Open("sqlite3", "/run/media/gnewton/b2b3f4a1-59af-4860-a100-305ecec24f03/sqlite3/gorm4-test.db")
	//db, err := gorm.Open("sqlite3", "/tmp/gorm4_tmp.db")
	//db, err := gorm.Open("sqlite3", "/run/media/gnewton/f34c5c5b-48de-4ae1-8ef2-28e95139cb06/tmp/gorm_all2.db")
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
	db.Exec("PRAGMA cache_size = 1800000;").Exec("PRAGMA synchronous = OFF;")
	db.Exec("PRAGMA journal_mode = OFF;")

	db.Exec("PRAGMA auto_vacuum = 0;")
	db.Exec("PRAGMA encoding = \"UTF-8\";")
	db.Exec("PRAGMA quick_check;")
	db.Exec("PRAGMA shrink_memory")
	db.Exec("PRAGMA synchronous = 0")

	db.Exec("PRAGMA count_changes = OFF;")
	db.Exec("PRAGMA temp_store = MEMORY;")
	//db.Exec("PRAGMA auto_vacuum = NONE;")
	db.Exec("PRAGMA cache_size=10000;")
	db.Exec("PRAGMA page_size = 32768;")
	db.Exec("PRAGMA threads = 5;")
	//db.Exec("PRAGMA mmap_size=12884901888;")
	//db.Exec("PRAGMA mmap_size=1099511627776;")
	db.Exec("PRAGMA mmap_size=0;")
	return &db, nil
}

func dbInit() (*gorm.DB, error) {
	db, err := dbOpen()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	db.CreateTable(&Article{})
	db.CreateTable(&Author{})
	db.CreateTable(&Chemical{})
	db.CreateTable(&Citation{})
	db.CreateTable(&Gene{})
	db.CreateTable(&Journal{})
	db.CreateTable(&MeshTerm{})

	//db.Exec("CREATE VIRTUAL TABLE pages USING fts4(title, body);")

	return db, nil
}

func makeIndexes(db *gorm.DB) {
	log.Println("makeing indexes START")
	db.Table("Article_Author").AddUniqueIndex("articleAuthor", "article_id", "author_id")
	db.Table("Article_Chemical").AddUniqueIndex("articleChemical", "article_id", "chemical_id")
	db.Table("Article_Citation").AddUniqueIndex("articleCitation", "article_id", "citation_id")
	db.Table("Article_Gene").AddUniqueIndex("articleGene", "article_id", "gene_id")
	db.Table("Article_MeshTerm").AddUniqueIndex("articleMeshTerm", "article_id", "mesh_term_id")
	db.Table("Article_MeshTerm").AddUniqueIndex("articleMeshTerm", "article_id", "mesh_term_id")
	db.Table("articles").AddIndex("articlesYear", "year")
	db.Table("articles").AddIndex("articlesJournalId", "journal_id")
	log.Println("makeing indexes END")
}

func dbCloseOpen(prevDb *gorm.DB) (*gorm.DB, error) {
	prevDb.Close()
	return dbOpen()
}

//func foo(db gorm.DB) chan *Article {

//}
