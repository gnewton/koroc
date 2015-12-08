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
	db, err := gorm.Open("sqlite3", "/run/media/gnewton/b2b3f4a1-59af-4860-a100-305ecec24f03/sqlite3/gorm3.db")
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
	db.Raw("PRAGMA page_size = 8192;")

	// // PRAGMA main.cache_size=10000;
	// // PRAGMA main.locking_mode=EXCLUSIVE;
	// // PRAGMA main.synchronous=NORMAL;
	// // PRAGMA main.journal_mode=WAL;
	// // PRAGMA main.cache_size=5000;

	// //db.Raw("PRAGMA automatic_index = ON;")
	// db.Raw("PRAGMA cache_size = 32768;")
	// //db.Raw("PRAGMA cache_spill = OFF;")
	// //db.Raw("PRAGMA foreign_keys = ON;")
	// //db.Raw("PRAGMA journal_mode = WAL;")
	// db.Raw("PRAGMA journal_mode = OFF;")
	// db.Raw("PRAGMA journal_size_limit = 67110000;")
	// db.Raw("PRAGMA locking_mode = EXCLUSIVE;")
	// db.Raw("PRAGMA page_size = 4096;")
	// //db.Raw("PRAGMA recursive_triggers = ON;")
	// db.Raw("PRAGMA secure_delete = ON;")
	// //db.Raw("PRAGMA synchronous = NORMAL;")
	// db.Raw("PRAGMA temp_store = MEMORY;")
	// db.Raw("PRAGMA wal_autocheckpoint = 16384;")
	// db.Raw("PRAGMA mmap_size = 1000000;")
	// db.Raw("PRAGMA soft_heap_limit = 1000000;")
	// db.Raw("`PRAGMA threads = 4;")

	db.Exec("PRAGMA cache_size = 1800000;").Exec("PRAGMA synchronous = OFF;").Exec("PRAGMA journal_mode = OFF;")
	db.Exec("PRAGMA locking_mode = EXCLUSIVE;")
	db.Exec("PRAGMA count_changes = OFF;")
	db.Exec("PRAGMA temp_store = MEMORY;")
	db.Exec("PRAGMA auto_vacuum = NONE;")
	//db.Exec("PRAGMA cache_size=128;")
	//db.Exec("PRAGMA page_size = 65536;")
	db.Exec("PRAGMA threads = 5;")

	db.Table("Article_Author").AddUniqueIndex("articleAuthor", "article_id", "author_id")
	db.Table("Article_Chemical").AddUniqueIndex("articleChemical", "article_id", "chemical_id")
	db.Table("Article_Citation").AddUniqueIndex("articleCitation", "article_id", "citation_id")
	db.Table("Article_Gene").AddUniqueIndex("articleGene", "article_id", "gene_id")
	db.Table("Article_MeshTerm").AddUniqueIndex("articleMeshTerm", "article_id", "mesh_term_id")

	return &db, nil
}

//func foo(db gorm.DB) chan *Article {

//}
