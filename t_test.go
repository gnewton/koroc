package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
	"testing"
)

var db *sql.DB

func TestMain(m *testing.M) {
	var err error
	db, err = dbOpen2("test1.sqlite3")
	defer db.Close()
	if err != nil {
		log.Fatal(err)
	}
	// call flag.Parse() here if TestMain uses flags
	os.Exit(m.Run())
}

func TestCreateTable(t *testing.T) {

	_, err := NewTable("journals")
	if err != nil {
		t.Error(err)
	}

}

func TestCreateTableZeroName(t *testing.T) {

	_, err := NewTable("")
	if err == nil {
		t.Error(err)
	}

}

func TestAddField(t *testing.T) {

	journals, err := NewTable("journals")
	if err != nil {
		t.Error(err)
	}
	_, err = journals.NewFieldString("issn")
	if err != nil {
		log.Fatal(err)
	}

}

func TestAddFieldZeroName(t *testing.T) {

	journals, err := NewTable("journals")
	if err != nil {
		t.Error(err)
	}
	_, err = journals.NewFieldString("")
	if err == nil {
		log.Fatal(err)
	}

}

func TestFoo2(t *testing.T) {
	var err error
	if err != nil {
		t.Error(err)
	}
}
