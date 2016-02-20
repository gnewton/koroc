package main

import (
	"database/sql"
)

type Article struct {
	Abstract  string
	Authors   []Author   `gorm:"many2many:Article_Author;"`
	Chemicals []Chemical `gorm:"many2many:Article_Chemical;"`
	Citations []Citation `gorm:"many2many:Article_Citation;"`
	Day       int
	Genes     []Gene `gorm:"many2many:Article_Gene;"`
	ID        int64  `gorm:"primary_key"` // PMID
	Issue     string
	Journal   Journal
	JournalID sql.NullInt64
	Language  string
	MeshTerms []MeshTerm
	Month     string
	Title     string
	Volume    string
	Year      int
}

type Journal struct {
	ID       int `gorm:"primary_key"`
	Title    string
	Issn     string
	Articles []Article
}

type Author struct {
	ID          int `gorm:"primary_key"`
	LastName    string
	FirstName   string
	MiddleName  string
	Affiliation string
}

type MeshTerm struct {
	ID         int   `gorm:"primary_key"`
	ArticleID  int64 `sql:"index"`
	Descriptor string
	Qualifier  string
}

type Gene struct {
	ID   int `gorm:"primary_key"`
	Name string
}

type Chemical struct {
	ID       int `gorm:"primary_key"`
	Name     string
	Registry string
}

type Citation struct {
	ID        int64 `gorm:"primary_key"`
	RefSource string
	Pmid      int64
}
