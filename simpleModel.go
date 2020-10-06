package main

import (
	"github.com/gnewton/pubmedSqlStructs"
	//_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
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
	var errs []error
	db.AutoMigrate(&pubmedSqlStructs.Article{}, &pubmedSqlStructs.Author{}, &pubmedSqlStructs.Chemical{}, &pubmedSqlStructs.Citation{}, &pubmedSqlStructs.Gene{}, &pubmedSqlStructs.Journal{}, &pubmedSqlStructs.Keyword{}, &pubmedSqlStructs.MeshDescriptor{}, &pubmedSqlStructs.MeshQualifier{}, &pubmedSqlStructs.Other{}, &pubmedSqlStructs.DataBank{}, &pubmedSqlStructs.AccessionNumber{}, &pubmedSqlStructs.ArticleID{}, &pubmedSqlStructs.PublicationType{})
	handleErrors(errs)

	// Relations
	// var meshDescriptor pubmedSqlStructs.MeshDescriptor
	//var article pubmedSqlStructs.Article
	// errs = db.Model(&meshDescriptor).Related(&article).GetErrors()
	// handleErrors(errs)
	//var dataBank pubmedSqlStructs.DataBank
	//errs = db.Model(&dataBank).Related(&article).GetErrors()
	//handleErrors(errs)

	// var accessionNumber pubmedSqlStructs.AccessionNumber
	// errs = db.Model(&accessionNumber).Related(&dataBank).GetErrors()
	// handleErrors(errs)

	// var articleID pubmedSqlStructs.ArticleID
	// errs = db.Model(&articleID).Related(&article).GetErrors()
	// handleErrors(errs)

	//var author pubmedSqlStructs.Author
	//errs = db.Model(&article).Related(&author).GetErrors()
	//handleErrors(errs)

	//	Var others pubmedSqlStructs.Other
	//errs = db.Model(&others).Related(&article).GetErrors()
	//handleErrors(errs)

	//var publicationType pubmedSqlStructs.PublicationType
	//db.Model(&publicationType).Related(&article)

	//db.Exec("CREATE VIRTUAL TABLE pages USING fts4(title, body);")
}

func makeIndexes(db *gorm.DB) {
	if false {
		log.Println("makeing indexes START")

		db.Table("articles").AddIndex("articlesYear", "year")
		db.Table("articles").AddIndex("articlesJournalId", "journal_id")

		db.Table("mesh_descriptors").AddIndex("mesh_descriptor_article", "article_id")
		db.Table("mesh_qualifiers").AddIndex("mesh_qualifier_descriptor", "mesh_descriptor_id")

		db.Table("accession_numbers").AddIndex("data_bank", "data_bank_id")
		db.Table("data_banks").AddIndex("databank_article", "article_id")
		db.Table("article_ids").AddIndex("ids_article", "article_id")

		log.Println("makeing indexes END")
	}
}
