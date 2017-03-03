package main

import (
	"database/sql"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	//_ "github.com/mxk/go-sqlite"
	//_ "github.com/go-sql-driver/mysql"
	"log"
)

func dbOpen(filename string) (*gorm.DB, error) {
	db, err := gorm.Open("sqlite3", filename)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	log.Println("Opening db file: ", filename)
	if sqliteLogFlag {
		db.LogMode(true)
	}

	db.DB()
	db.DB().Ping()
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)

	sqlite3Config(db.DB())
	return db, nil
}

func sqlite3Config(db *sql.DB) {
	//db.Exec("PRAGMA auto_vacuum = 0;")
	//db.Exec("PRAGMA cache_size=32768;")
	db.Exec("PRAGMA cache_size=65536;")
	db.Exec("PRAGMA count_changes = OFF;")
	db.Exec("PRAGMA cache_spill = ON;")
	//db.Exec("PRAGMA journal_size_limit = 67110000;")
	//db.Exec("PRAGMA locking_mode = EXCLUSIVE;")
	db.Exec("PRAGMA locking_mode = OFF;")
	db.Exec("PRAGMA encoding = \"UTF-8\";")
	db.Exec("PRAGMA journal_mode = WAL;")

	//db.Exec("busy_timeout=0;")
	db.Exec("legacy_file_format=OFF;")

	//db.Exec("PRAGMA mmap_size=1099511627776;")
	db.Exec("PRAGMA page_size = 40960;")

	db.Exec("PRAGMA shrink_memory;")
	//db.Exec("PRAGMA synchronous=OFF;")
	//db.Exec("PRAGMA synchronous = NORMAL;")
	//db.Exec("PRAGMA temp_store = MEMORY;")
	//db.Exec("PRAGMA threads = 5;")
	//db.Exec("PRAGMA wal_autocheckpoint = 1638400;")
}
