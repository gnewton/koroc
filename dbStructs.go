package main

import (
	"database/sql"
)

type Article struct {
	Abstract  string
	Authors   []Author   `gorm:"many2many:article_author;"`
	Chemicals []Chemical `gorm:"many2many:article_chemical;"`
	Citations []Citation `gorm:"many2many:article_citation;"`
	Day       int
	Genes     []Gene `gorm:"many2many:article_gene;"`
	ID        int64  `gorm:"primary_key"` // PMID
	Issue     string
	Journal   Journal
	JournalID sql.NullInt64
	Language  string
	MeshTerms []MeshTerm `gorm:"many2many:article_meshterm;"`
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
	ID                   int `gorm:"primary_key"`
	Descriptor           MeshDescriptor
	DescriptorMajorTopic bool
	Qualifiers           []MeshQualifier
}

type MeshDescriptor struct {
	ID             int `gorm:"primary_key"`
	DescriptorName string
}

type MeshQualifier struct {
	ID                   int `gorm:"primary_key"`
	QualifiersMajorTopic bool
	MeshQualifierName    MeshQualifierName
}

type MeshQualifierName struct {
	ID   int `gorm:"primary_key"`
	Name string
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
