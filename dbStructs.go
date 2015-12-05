package main

type Article struct {
	Abstract  string
	Authors   []Author   `gorm:"many2many:Article_Author;"`
	Chemicals []Chemical `gorm:"many2many:Article_Chemical;"`
	Citations []Citation `gorm:"many2many:Article_Citation;"`
	Day       int
	Genes     []Gene `gorm:"many2many:Article_Gene;"`
	Id        int64  `gorm:"primary_key"`
	Issue     string
	Journal   Journal
	Language  string
	MeshTerms []MeshTerm `gorm:"many2many:Article_MeshTerm;"`
	Month     int
	Title     string
	Volume    string
	Year      int
}

type Journal struct {
	Id    int `gorm:"primary_key"`
	Title string
	Issn  string
}

type Author struct {
	Id          int `gorm:"primary_key"`
	LastName    string
	FirstName   string
	MiddleName  string
	Affiliation string
}

type MeshTerm struct {
	Id         int `gorm:"primary_key"`
	Descriptor string
	Qualifier  string
}

type Gene struct {
	Id   int `gorm:"primary_key"`
	Name string
}

type Chemical struct {
	Id       int `gorm:"primary_key"`
	Name     string
	Registry string
}

type Citation struct {
	Id           int64 `gorm:"primary_key"`
	SourcePmid   int64
	CitationPmid int64
}
