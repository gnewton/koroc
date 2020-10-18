package main

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type DBConnector struct {
	dbFilename string
	gdb        *gorm.DB
	tx         *gorm.DB
}

func (dbc *DBConnector) Open() (*gorm.DB, error) {
	var err error
	dbc.gdb, err = gorm.Open(sqlite.Open(dbc.dbFilename), &gorm.Config{})
	return dbc.gdb, err
}

func (dbc *DBConnector) DB() *gorm.DB {
	return dbc.gdb
}
