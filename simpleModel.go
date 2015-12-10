package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	//_ "github.com/go-sql-driver/mysql"
	"log"
)

func dbInit() (*gorm.DB, error) {
	//db, err := gorm.Open("sqlite3", "/tmp/gorm_15000-x.db")

	//db, err := gorm.Open("mysql", "gnewton:@/pubmed?charset=utf8&parseTime=True&loc=Local")
	db, err := gorm.Open("sqlite3", "/run/media/gnewton/b2b3f4a1-59af-4860-a100-305ecec24f03/sqlite3/gorm4-test.db")
	//db, err := gorm.Open("sqlite3", "/tmp/gorm4.db")
	if err != nil {
		log.Println(err)
		return nil, err
	}
	//db.LogMode(true)

	db.DB()
	db.DB().Ping()
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)

	db.CreateTable(&Article{})
	db.CreateTable(&Author{})
	db.CreateTable(&Chemical{})
	db.CreateTable(&Citation{})
	db.CreateTable(&Gene{})
	db.CreateTable(&Journal{})
	db.CreateTable(&MeshTerm{})

	db.Exec("PRAGMA cache_size = 1800000;").Exec("PRAGMA synchronous = OFF;")
	//db.Exec("PRAGMA journal_mode = OFF;")
	db.Exec("PRAGMA locking_mode = EXCLUSIVE;")
	db.Exec("PRAGMA count_changes = OFF;")
	db.Exec("PRAGMA temp_store = MEMORY;")
	//db.Exec("PRAGMA auto_vacuum = NONE;")
	db.Exec("PRAGMA cache_size=16;")
	db.Exec("PRAGMA page_size = 65536;")
	db.Exec("PRAGMA threads = 5;")

	//db.Exec("PRAGMA mmap_size=12884901888;")
	db.Exec("PRAGMA mmap_size=0;")

	db.Table("Article_Author").AddUniqueIndex("articleAuthor", "article_id", "author_id")
	db.Table("Article_Chemical").AddUniqueIndex("articleChemical", "article_id", "chemical_id")
	db.Table("Article_Citation").AddUniqueIndex("articleCitation", "article_id", "citation_id")
	db.Table("Article_Gene").AddUniqueIndex("articleGene", "article_id", "gene_id")
	db.Table("Article_MeshTerm").AddUniqueIndex("articleMeshTerm", "article_id", "mesh_term_id")

	//return &db, nil
	return &db, nil
}

//func foo(db gorm.DB) chan *Article {

//}
