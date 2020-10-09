package main

import (
	"github.com/gnewton/pubmedSqlStructs"
	//_ "github.com/go-sql-driver/mysql"
	//"github.com/jinzhu/gorm"
	//_ "github.com/mattn/go-sqlite3"
	"gorm.io/gorm"
	"log"
)

func handleErrors(errs []error) {
	if len(errs) > 0 {
		for i, _ := range errs {
			log.Println(i, errs[i])
		}

		log.Fatal("====================")
	}
}

func dbInit(db *gorm.DB) {
	log.Printf("%v\n", db)

	//db = db.Debug()
	err := db.AutoMigrate(&pubmedSqlStructs.Article{})
	if err != nil {
		log.Fatal(err)
	}
	err = db.AutoMigrate(&pubmedSqlStructs.Author{})
	if err != nil {
		log.Fatal(err)
	}

	err = db.AutoMigrate(&pubmedSqlStructs.Chemical{})
	if err != nil {
		log.Fatal(err)
	}

	err = db.AutoMigrate(&pubmedSqlStructs.Citation{})
	if err != nil {
		log.Fatal(err)
	}

	err = db.AutoMigrate(&pubmedSqlStructs.Gene{})
	if err != nil {
		log.Fatal(err)
	}

	err = db.AutoMigrate(&pubmedSqlStructs.Journal{})
	if err != nil {
		log.Fatal(err)
	}

	err = db.AutoMigrate(&pubmedSqlStructs.Keyword{})
	if err != nil {
		log.Fatal(err)
	}

	err = db.AutoMigrate(&pubmedSqlStructs.MeshDescriptor{})
	if err != nil {
		log.Fatal(err)
	}

	err = db.AutoMigrate(&pubmedSqlStructs.MeshQualifier{})
	if err != nil {
		log.Fatal(err)
	}

	err = db.AutoMigrate(&pubmedSqlStructs.Other{})
	if err != nil {
		log.Fatal(err)
	}

	err = db.AutoMigrate(&pubmedSqlStructs.DataBank{})
	if err != nil {
		log.Fatal(err)
	}

	err = db.AutoMigrate(&pubmedSqlStructs.AccessionNumber{})
	if err != nil {
		log.Fatal(err)
	}

	err = db.AutoMigrate(&pubmedSqlStructs.ArticleID{})
	if err != nil {
		log.Fatal(err)
	}

	err = db.AutoMigrate(&pubmedSqlStructs.PublicationType{})
	if err != nil {
		log.Fatal(err)
	}

	log.Println("dbInit DONE ---------------------")
	//db.Exec("CREATE VIRTUAL TABLE pages USING fts4(title, body);")
}
