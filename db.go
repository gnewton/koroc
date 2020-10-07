package main

import (
	"database/sql"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	//"github.com/jinzhu/gorm"
	//_ "github.com/mattn/go-sqlite3"
	//_ "github.com/mxk/go-sqlite"
	//_ "github.com/go-sql-driver/mysql"
	"log"
)

func dbOpen(filename string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(filename), &gorm.Config{})
	if err != nil {
		log.Println(err)
		return nil, err
	}
	log.Println("Opening db file: ", filename)

	sdb, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}
	err = sdb.Ping()
	if err != nil {
		log.Fatal(err)
	}
	sdb.SetMaxIdleConns(10)
	sdb.SetMaxOpenConns(100)

	sqlite3Config(sdb)
	return db, nil
}

func sqlite3Config(db *sql.DB) {

	db.Exec("PRAGMA cache_size=10000;")
	db.Exec("PRAGMA cache_spill = ON;")
	db.Exec("PRAGMA count_changes = OFF;")
	db.Exec("PRAGMA encoding = \"UTF-8\";")
	db.Exec("PRAGMA journal_mode = delete;")
	db.Exec("PRAGMA locking_mode = EXCLUSIVE;")
	//db.Exec("PRAGMA main.synchronous=NORMAL;")
	db.Exec("PRAGMA page_size = 4096;")
	db.Exec("PRAGMA shrink_memory;")
	db.Exec("PRAGMA synchronous = off;")
	db.Exec("PRAGMA temp_store = memory;")
	db.Exec("legacy_file_format=OFF;")

	//db.Exec("PRAGMA mmap_size=1099511627776;")

	//db.Exec("PRAGMA synchronous=OFF;")
	//db.Exec("PRAGMA synchronous = NORMAL;")
	//db.Exec("PRAGMA temp_store = MEMORY;")
	//db.Exec("PRAGMA threads = 5;")
	//db.Exec("PRAGMA wal_autocheckpoint = 1638400;")
}
