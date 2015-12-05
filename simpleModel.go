package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

func dbInit() (*gorm.DB, error) {
	db, err := gorm.Open("sqlite3", "/tmp/gorm_15000.db")
	if err != nil {
		log.Println(err)
		return nil, err
	}
	//db.LogMode(true)
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

	// chemical1 := Chemical{ID: 1, Name: "acetone"}
	// chemical2 := Chemical{ID: 2, Name: "water"}

	// gene1 := Gene{ID: 1, Name: "ALM3"}
	// gene2 := Gene{ID: 2, Name: "Kenoacetase"}

	// meshTerm1 := MeshTerm{ID: 1, Descriptor: "d1"}
	// meshTerm2 := MeshTerm{ID: 2, Descriptor: "d2"}

	// citation1 := Citation{SourcePmid: 124, CitationPmid: 43}
	// citation2 := Citation{SourcePmid: 124, CitationPmid: 933}

	// author1 := Author{ID: "1", LastName: "Smith", FirstName: "Bill"}
	// author2 := Author{ID: "3", LastName: "Williams", FirstName: "James"}

	// journal := Journal{ID: 22, Title: "art title", Issn: "34454"}
	// article := Article{ID: 124, Title: "a new article title...", Journal: journal, Chemicals: []Chemical{chemical1, chemical2}, Genes: []Gene{gene1, gene2}, MeshTerms: []MeshTerm{meshTerm1, meshTerm2},
	// 	Citations: []Citation{citation1, citation2}, Authors: []Author{author1, author2}}
	// db.Create(&article)

	// citation1 = Citation{SourcePmid: 129, CitationPmid: 43}
	// article2 := Article{ID: 129, Title: "ZZZZZZZZZZZz a new article title...", Journal: journal, Chemicals: []Chemical{chemical2}, Genes: []Gene{gene1}, MeshTerms: []MeshTerm{meshTerm2},
	// 	Citations: []Citation{citation1}, Authors: []Author{author2}}
	// db.Create(&article2)
	return &db, nil
}

//func foo(db gorm.DB) chan *Article {

//}
