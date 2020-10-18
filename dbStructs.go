package main

import (
	"github.com/jinzhu/gorm"
)

type Article struct {
	//gorm.Model
	FullAbstract                 string
	AbstractSections             []AbstractText `gorm:"-"`
	AbstractCopyrightInformation string
	//ArticleIDs       []*ArticleID `gorm:"many2many:article_articleid;"` // https://www.nlm.nih.gov/bsd/licensee/elements_descriptions.html#articleidlist
	//Authors          []*Author    `gorm:"many2many:article_author;"`
	//Chemicals        []*Chemical  `gorm:"many2many:article_chemical;"`
	//Citations        []*Citation  `gorm:"many2many:article_citation;"`
	CoiStatement string
	//DataBanks        []*DataBank `gorm:"many2many:article_databank;"`
	CopyrightInformation string
	DateRevised          uint64 //YYYYMMDD 20140216

	ID            uint32        `gorm:"primary_key" sql:"type:int4"` // PMID
	ELocationID   []ELocationID `gorm:"-"`
	Issue         string        `sql:"size:32" sql:"type:int4"`
	JournalID     uint32
	Keywords      []*Keyword `gorm:"-"`
	KeywordsOwner string     `sql:"size:16"`
	Language      string     `sql:"size:3"`
	//MeshDescriptors  []*MeshDescriptor  `gorm:"many2many:article_meshdescriptor;"`
	PubDay   uint8  `sql:"type:int1"`
	PubMonth string `sql:"size:8"`
	PubYear  uint16 `sql:"type:int2"`
	//OtherIds         []*Other           `gorm:"many2many:article_other;"` // https://www.nlm.nih.gov/bsd/licensee/elements_descriptions.html#otherid
	//PublicationTypes []*PublicationType `gorm:"many2many:article_publicationtype;"`
	Pagination string
	Retracted  bool
	//SupplMesh        []*SupplMesh `gorm:"many2many:article_supplmesh;"`
	Title   string
	Version uint8 `sql:"type:int1"`
	Volume  string

	SourceXMLFilename string `gorm:"-"`
}

type AbstractText struct {
	ID             uint32
	Order          int8
	Label          string
	NlmCategory    string
	Text           string
	TruncatedAt250 bool
}

type PublicationType struct {
	gorm.Model
	UI   string `sql:"size:8"`
	Name string `sql:"size:64"`
}

type ArticleID struct {
	gorm.Model
	OtherArticleID string `sql:"size:64"`
	Type           string `sql:"size:12"`
	ArticleID      uint32 `sql:"type:int4"`
	//Articles       []*Article `gorm:"many2many:article_articleid;"`
}

type DataBank struct {
	gorm.Model
	Name             string `sql:"size:32"`
	AccessionNumbers []*AccessionNumber
	//Articles         []*Article `gorm:"many2many:article_databank;"`
	//ArticleID        uint32 `sql:"type:int4"`
}

type AccessionNumber struct {
	gorm.Model
	Number     string `sql:"size:32"`
	DataBankID uint32 `sql:"type:int2"`
}

type Other struct {
	gorm.Model
	Source  string
	OtherID string
}

type Journal struct {
	ID              uint32 `gorm:"primaryKey" sql:"type:int4"`
	IsoAbbreviation string `sql:"size:128"`
	Issn            string `sql:"size:10"`
	Title           string
	Country         string
	MedlineTA       string
	NlmUniqueID     string

	//Articles        []*Article
}

type SourceXMLFile struct {
	ID       uint32
	Filename string
}

type Author struct {
	gorm.Model
	LastName         string `sql:"size:48"`
	FirstName        string `sql:"size:16"`
	MiddleName       string `sql:"size:16"`
	Affiliation      string
	CollectiveName   string
	Identifier       string `sql:"size:16"`
	IdentifierSource string `sql:"size:48"`
	//Articles         []*Article `gorm:"many2many:article_author;"`
}

type Keyword struct {
	ID         uint32 `gorm:"primaryKey"`
	Owner      string
	MajorTopic bool
	Name       string `sql:"size:128"`
}

type SupplMesh struct {
	gorm.Model
	Type string `xml:"Type,attr"  json:",omitempty"`
	UI   string `xml:"UI,attr"  json:",omitempty"`
	Name string `xml:",chardata" json:",omitempty"`
}

type MeshQualifier struct {
	gorm.Model
	MajorTopic       bool
	Name             string `sql:"size:128"`
	MeshDescriptorID uint32 `sql:"type:int4"`
	UI               string
}

type Gene struct {
	gorm.Model
	Name string
}

type Chemical struct {
	gorm.Model
	Name     string
	Registry string `sql:"size:32"`
	UI       string `sql:"size:7"`
	//Articles []*Article `gorm:"many2many:article_chemical;"`
}

type Citation struct {
	gorm.Model
	PMID uint32
}
type ELocationID struct {
	EIdType string
	ValidYN bool
	Text    string
}
