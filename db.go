package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	//_ "github.com/mxk/go-sqlite"
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
	db.Exec("PRAGMA cache_size=32768;")
	db.Exec("PRAGMA count_changes = OFF;")
	db.Exec("PRAGMA cache_spill = OFF;")
	db.Exec("PRAGMA journal_size_limit = 67110000;")
	db.Exec("PRAGMA locking_mode = NORMAL;")
	db.Exec("PRAGMA encoding = \"UTF-8\";")
	db.Exec("PRAGMA journal_mode = WAL;")

	db.Exec("PRAGMA mmap_size=1099511627776;")
	db.Exec("PRAGMA page_size = 4096;")
	db.Exec("PRAGMA quick_check;")
	db.Exec("PRAGMA shrink_memory")
	db.Exec("PRAGMA synchronous = NORMAL;")
	db.Exec("PRAGMA temp_store = MEMORY;")
	db.Exec("PRAGMA threads = 5;")
	db.Exec("PRAGMA wal_autocheckpoint = 1638400;")
	return db, nil
}

func dbCloseOpen(prevDb *gorm.DB) (*gorm.DB, error) {
	prevDb.Close()
	return dbOpen()
}
