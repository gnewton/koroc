package main

import (
	//"github.com/gnewton/
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

func createTables(db *gorm.DB) {
	log.Printf("%v\n", db)

	//db = db.Debug()
	if err := db.AutoMigrate(&Article{}); err != nil {
		log.Fatal(err)
	}
	err := db.AutoMigrate(&Author{})
	if err != nil {
		log.Fatal(err)
	}

	err = db.AutoMigrate(&Chemical{})
	if err != nil {
		log.Fatal(err)
	}

	err = db.AutoMigrate(&Citation{})
	if err != nil {
		log.Fatal(err)
	}

	err = db.AutoMigrate(&Gene{})
	if err != nil {
		log.Fatal(err)
	}

	err = db.AutoMigrate(&Journal{})
	if err != nil {
		log.Fatal(err)
	}

	err = db.AutoMigrate(&Keyword{})
	if err != nil {
		log.Fatal(err)
	}

	err = db.AutoMigrate(&MeshDescriptor{})
	if err != nil {
		log.Fatal(err)
	}

	err = db.AutoMigrate(&MeshQualifier{})
	if err != nil {
		log.Fatal(err)
	}

	err = db.AutoMigrate(&Other{})
	if err != nil {
		log.Fatal(err)
	}

	err = db.AutoMigrate(&DataBank{})
	if err != nil {
		log.Fatal(err)
	}

	err = db.AutoMigrate(&AccessionNumber{})
	if err != nil {
		log.Fatal(err)
	}

	err = db.AutoMigrate(&ArticleID{})
	if err != nil {
		log.Fatal(err)
	}

	err = db.AutoMigrate(&PublicationType{})
	if err != nil {
		log.Fatal(err)
	}

	log.Println("dbInit DONE ---------------------")
	//db.Exec("CREATE VIRTUAL TABLE pages USING fts4(title, body);")
}
