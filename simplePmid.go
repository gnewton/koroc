package main

type Article struct {
	Abstract string
	Day      int
	Issue    string
	Language string
	Month    int
	PMID     int64
	Title    string
	Volume   string
	Year     int
	Journal  int
}

type Journal struct {
	ID    int
	Title string
	ISSN  string
}

type Author struct {
	ID          int
	LastName    string
	FirstName   string
	MiddleName  string
	Affiliation string
}

type MeSHTerm struct {
	MeSHID     int
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
	SourcePMID   int64
	CitationPMID int64
}

////// joins
type ArticleGene struct {
	PMID   int64
	GeneID int
}

type ArticleMeSH struct {
	PMID   int64
	MeSHID int
}

type ArticleAuthor struct {
	PMID     int64
	AuthorID int
}

type ArticleChemical struct {
	PMID       int64
	ChemicalID int
}
