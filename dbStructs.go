package main

type Article struct {
	Abstract  string
	Authors   []Author   `gorm:"many2many:Article_Author;"`
	Chemicals []Chemical `gorm:"many2many:Article_Chemical;"`
	Citations []Citation `gorm:"many2many:Article_Citation;"`
	Day       int
	Genes     []Gene `gorm:"many2many:Article_Gene;"`
	ID        int64  `gorm:"primary_key"`
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
	ID    int
	Title string
	Issn  string
}

type Author struct {
	ID          string
	LastName    string
	FirstName   string
	MiddleName  string
	Affiliation string
}

type MeshTerm struct {
	ID         int `gorm:"primary_key"`
	Descriptor string
	Qualifier  string
}

type Gene struct {
	ID   int
	Name string
}

type Chemical struct {
	ID       int
	Name     string
	Registry string
}

type Citation struct {
	ID           int64
	SourcePmid   int64
	CitationPmid int64
}
