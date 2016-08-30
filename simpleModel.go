package main

import (
	"github.com/gnewton/pubmedSqlStructs"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

func dbInit() (*gorm.DB, error) {
	db, err := dbOpen()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	log.Printf("%v\n", *db)

	db.CreateTable(&pubmedSqlStructs.Article{})
	db.CreateTable(&pubmedSqlStructs.Author{})
	db.CreateTable(&pubmedSqlStructs.Chemical{})
	db.CreateTable(&pubmedSqlStructs.Citation{})
	db.CreateTable(&pubmedSqlStructs.Gene{})
	db.CreateTable(&pubmedSqlStructs.Journal{})
	db.CreateTable(&pubmedSqlStructs.Keyword{})
	db.CreateTable(&pubmedSqlStructs.MeshDescriptor{})
	db.CreateTable(&pubmedSqlStructs.MeshQualifier{})
	db.CreateTable(&pubmedSqlStructs.OtherID{})

	db.Table("Article_Author").AddUniqueIndex("articleAuthor", "article_id", "author_id")
	db.Table("Article_Chemical").AddUniqueIndex("articleChemical", "article_id", "chemical_id")
	db.Table("Article_Citation").AddUniqueIndex("articleCitation", "article_id", "citation_id")
	db.Table("Article_Gene").AddUniqueIndex("articleGene", "article_id", "gene_id")
	db.Table("Article_Keyword").AddUniqueIndex("articleKeyword", "article_id", "keyword_id")
	db.Table("mesh_descriptors").AddIndex("mesh_descriptor_article", "article_id")
	db.Table("mesh_qualifiers").AddIndex("mesh_qualifier_descriptor", "mesh_descriptor_id")

	//db.Exec("CREATE VIRTUAL TABLE pages USING fts4(title, body);")

	return db, nil
}

func makeIndexes(db *gorm.DB) {
	log.Println("makeing indexes START")

	db.Table("articles").AddIndex("articlesYear", "year")
	db.Table("articles").AddIndex("articlesJournalId", "journal_id")
	log.Println("makeing indexes END")
}
