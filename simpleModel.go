package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

func dbInit() (*gorm.DB, error) {
	db, err := gorm.Open("sqlite3", "/tmp/gorm_15000-x.db")
	if err != nil {
		log.Println(err)
		return nil, err
	}
	db.LogMode(true)
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
	db.Raw("PRAGMA page_size = 4096;")

	// PRAGMA main.cache_size=10000;
	// PRAGMA main.locking_mode=EXCLUSIVE;
	// PRAGMA main.synchronous=NORMAL;
	// PRAGMA main.journal_mode=WAL;
	// PRAGMA main.cache_size=5000;

	//db.Raw("PRAGMA automatic_index = ON;")
	db.Raw("PRAGMA cache_size = 32768;")
	//db.Raw("PRAGMA cache_spill = OFF;")
	//db.Raw("PRAGMA foreign_keys = ON;")
	//db.Raw("PRAGMA journal_mode = WAL;")
	db.Raw("PRAGMA journal_mode = OFF;")
	db.Raw("PRAGMA journal_size_limit = 67110000;")
	db.Raw("PRAGMA locking_mode = EXCLUSIVE;")
	db.Raw("PRAGMA page_size = 4096;")
	//db.Raw("PRAGMA recursive_triggers = ON;")
	db.Raw("PRAGMA secure_delete = ON;")
	//db.Raw("PRAGMA synchronous = NORMAL;")
	db.Raw("PRAGMA temp_store = MEMORY;")
	db.Raw("PRAGMA wal_autocheckpoint = 16384;")
	db.Raw("PRAGMA mmap_size = 1000000;")
	db.Raw("PRAGMA soft_heap_limit = 1000000;")
	db.Raw("PRAGMA threads = 4;")

	db2 := db.Table("Article_Author")
	db2.AddUniqueIndex("a7", "article_id", "author_id")

	return &db, nil
}

//func foo(db gorm.DB) chan *Article {

//}
