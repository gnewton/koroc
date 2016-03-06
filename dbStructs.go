package main

import (
	"database/sql"
)

type Article struct {
	Abstract        string
	Authors         []Author    `gorm:"many2many:article_author;"`
	Chemicals       []*Chemical `gorm:"many2many:article_chemical;"`
	Citations       []*Citation `gorm:"many2many:article_citation;"`
	Day             int
	Genes           []Gene `gorm:"many2many:article_gene;"`
	ID              int64  `gorm:"primary_key"` // PMID
	Issue           string
	Journal         *Journal
	JournalID       sql.NullInt64
	Keywords        []*Keyword `gorm:"many2many:article_keyword;"`
	KeywordsOwner   string
	Language        string
	MeshDescriptors []*MeshDescriptor
	Month           string
	OtherId         []OtherID
	Title           string
	Volume          string
	Year            int
}

type OtherID struct {
	ID      int `gorm:"primary_key"`
	Source  string
	OtherID string
}

type Journal struct {
	ID              int `gorm:"primary_key"`
	Articles        []Article
	IsoAbbreviation string
	Issn            string
	Title           string
}

type Author struct {
	ID          int `gorm:"primary_key"`
	LastName    string
	FirstName   string
	MiddleName  string
	Affiliation string
}

type Keyword struct {
	ID         int `gorm:"primary_key"`
	MajorTopic bool
	Name       string
}

type MeshDescriptor struct {
	ID         int `gorm:"primary_key"`
	Name       string
	Type       string
	MajorTopic bool
	Qualifiers []*MeshQualifier
	ArticleID  int64
}

type MeshQualifier struct {
	ID               int `gorm:"primary_key"`
	MajorTopic       bool
	Name             string
	MeshDescriptorID int
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
	ID int64 `gorm:"primary_key"`
}
