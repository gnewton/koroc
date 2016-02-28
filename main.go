package main

import (
	"encoding/xml"
	"flag"
	"github.com/gnewton/pubmedstruct"
	"io"
	"io/ioutil"
	//_ "net/http/pprof"
	"strings"
	//"github.com/davecheney/profile"
	"github.com/jinzhu/gorm"
	"net/http"

	"log"
	//	"net/http"
	"os"
	"strconv"
	//"strings"
	"math"
	"time"
)

var TransactionSize = 5000
var chunkSize = 1000
var CloseOpenSize int64 = 99950000
var chunkChannelSize = 3
var dbFileName = "./pubmed_sqlite.db"
var sqliteLogFlag = false
var LoadNRecordsPerFile int64 = math.MaxInt64
var recordPerFileCounter int64 = 0
var doNotWriteToDbFlag = false

const CommentsCorrections_RefType = "Cites"
const PUBMED_ARTICLE = "PubmedArticle"

var out int = -1
var JournalIdCounter int64 = 0
var counters map[string]*int
var closeOpenCount int64 = 0

func init() {
	//defer profile.Start(profile.CPUProfile).Stop()
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	flag.BoolVar(&sqliteLogFlag, "L", sqliteLogFlag, "Turn on sqlite logging")
	flag.StringVar(&dbFileName, "f", dbFileName, "SQLite output filename")

	flag.IntVar(&TransactionSize, "t", TransactionSize, "Size of transactions")
	flag.IntVar(&chunkSize, "C", chunkSize, "Size of chunks")
	flag.Int64Var(&CloseOpenSize, "z", CloseOpenSize, "Num of records before sqlite connection is closed then reopened")
	flag.Int64Var(&LoadNRecordsPerFile, "N", LoadNRecordsPerFile, "Load only N records from each file")
	flag.BoolVar(&sqliteLogFlag, "V", sqliteLogFlag, "Turn on sqlite logging")

	flag.BoolVar(&doNotWriteToDbFlag, "X", doNotWriteToDbFlag, "Do not write to db. Rolls back transaction. For debugging")

	flag.Parse()

	if len(flag.Args()) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	logInit(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)
}

func main() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	//defer profile.Start(profile.CPUProfile).Stop()

	db, err := dbInit()
	if err != nil {
		Error.Fatal(err)
		return
	}
	defer func() {
		err = db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	articleChannel := make(chan []*Article, chunkChannelSize)

	done := make(chan bool)

	go articleAdder(articleChannel, done, db, TransactionSize)
	var count int64 = 0
	chunkCount := 0
	arrayIndex := 0

	var articleArray []*Article

	for i, filename := range flag.Args() {
		log.Println(i, " -- Input file: "+filename)
	}

	// Loop through files
	for i, filename := range flag.Args() {

		log.Println("Opening: "+filename, " ", i+1, " of ", len(flag.Args()))
		//log.Println(strconv.Itoa(i) + " of " + strconv.Itoa(len(flag.Args)-1))
		reader, _, err := genericReader(filename)

		if err != nil {
			log.Fatal(err)
			return
		}
		arrayIndex = 0
		articleArray = make([]*Article, chunkSize)

		decoder := xml.NewDecoder(reader)
		counters = make(map[string]*int)

		// Loop through XML
		for {
			if recordPerFileCounter == LoadNRecordsPerFile {
				log.Println("break file load. LoadNRecordsPerFile", count, LoadNRecordsPerFile)
				recordPerFileCounter = 0
				break
			}
			token, _ := decoder.Token()
			if token == nil {
				break
			}
			switch se := token.(type) {
			case xml.StartElement:
				if se.Name.Local == PUBMED_ARTICLE && se.Name.Space == "" {
					if count%10000 == 0 && count != 0 {
						log.Println("------------")
						log.Printf("count=%d\n", count)
						log.Printf("arrayIndex=%d\n", arrayIndex)
						log.Println("------------")
					}

					count = count + 1
					recordPerFileCounter = recordPerFileCounter + 1
					var pubmedArticle pubmedstruct.PubmedArticle
					decoder.DecodeElement(&pubmedArticle, &se)
					article := pubmedArticleToDbArticle(&pubmedArticle)
					if article == nil {
						log.Println("-----------------nil")
						continue
					}
					articleArray[arrayIndex] = article
					arrayIndex = arrayIndex + 1
					if arrayIndex >= chunkSize {
						//log.Printf("Sending chunk %d", chunkCount)
						chunkCount = chunkCount + 1
						//pubmedArticleChannel <- &pubmedArticle
						//log.Printf("%v\n", articleArray)
						articleChannel <- articleArray
						log.Println("Sent")
						articleArray = make([]*Article, chunkSize)
						arrayIndex = 0
					}
				}
			}

		}
		if arrayIndex > 0 && arrayIndex < chunkSize {
			articleChannel <- articleArray
			chunkCount = chunkCount + 1
		}
	}

	close(articleChannel)
	_ = <-done

}

func pubmedArticleToDbArticle(p *pubmedstruct.PubmedArticle) *Article {
	medlineCitation := p.MedlineCitation
	pArticle := medlineCitation.Article
	if pArticle == nil {
		log.Println("nil-----------")
		return nil
	}
	var err error
	dbArticle := new(Article)
	dbArticle.ID, err = strconv.ParseInt(p.MedlineCitation.PMID.Text, 10, 64)
	if err != nil {
		log.Println(err)
	}
	dbArticle.Abstract = ""
	//if pArticle !=pArticle.Abstract != nil && pArticle.Abstract.AbstractText != nil {
	if pArticle.Abstract != nil && pArticle.Abstract.AbstractText != nil {
		for i, _ := range pArticle.Abstract.AbstractText {
			dbArticle.Abstract = dbArticle.Abstract + " " + pArticle.Abstract.AbstractText[i].Text
		}
	}

	dbArticle.Title = pArticle.ArticleTitle.Text

	if pArticle.Journal != nil {
		if pArticle.Journal.JournalIssue != nil {
			if pArticle.Journal.JournalIssue.PubDate != nil {
				if pArticle.Journal.JournalIssue.PubDate.Year != nil {
					dbArticle.Year, err = strconv.Atoi(pArticle.Journal.JournalIssue.PubDate.Year.Text)

					if err != nil {
						log.Println(err)
					}

				} else {
					if pArticle.Journal.JournalIssue.PubDate.MedlineDate == nil || pArticle.Journal.JournalIssue.PubDate.MedlineDate.Text == "" {
						log.Println("MedlineDate is nil? ", pArticle.Journal.JournalIssue.PubDate.MedlineDate)
					} else {
						dbArticle.Year = medlineDate2Year(pArticle.Journal.JournalIssue.PubDate.MedlineDate.Text)
					}
				}
				if pArticle.Journal.JournalIssue.PubDate.Month != nil {
					dbArticle.Month = pArticle.Journal.JournalIssue.PubDate.Month.Text
				}
				if pArticle.Journal.JournalIssue.PubDate.Day != nil {
					dbArticle.Day, err = strconv.Atoi(pArticle.Journal.JournalIssue.PubDate.Day.Text)
					if err != nil {
						log.Println(err)
					}
				}
			} else {
				log.Println("Journal.JournalIssue.PubDate=nil pmid=", dbArticle.ID)
			}
			if dbArticle.Year < 1000 {
				log.Println("*******************************************")
				log.Println("Year=Error ", dbArticle.ID)
				log.Println(dbArticle.Year)
				log.Printf("%+v\n", pArticle.Journal.JournalIssue)

				log.Printf("%+v\n", pArticle.Journal.JournalIssue.PubDate.Year)
				if pArticle.Journal.JournalIssue.PubDate.MedlineDate != nil {
					log.Printf("%+v\n", pArticle.Journal.JournalIssue.PubDate.MedlineDate.Text)
				}
				log.Println("*******************************************")
			}
		}
	}

	if medlineCitation.CommentsCorrectionsList != nil {
		actualCitationCount := 0
		for _, commentsCorrection := range medlineCitation.CommentsCorrectionsList.CommentsCorrections {
			if commentsCorrection.Attr_RefType == CommentsCorrections_RefType {
				actualCitationCount = actualCitationCount + 1
			}
		}

		dbArticle.Citations = make([]Citation, actualCitationCount)
		counter := 0
		var err error
		for _, commentsCorrection := range medlineCitation.CommentsCorrectionsList.CommentsCorrections {
			if commentsCorrection.Attr_RefType == CommentsCorrections_RefType {
				citation := new(Citation)
				citation.Pmid, err = strconv.ParseInt(commentsCorrection.PMID.Text, 10, 64)
				if err != nil {
					log.Println(err)
				}
				citation.RefSource = commentsCorrection.RefSource.Text
				dbArticle.Citations[counter] = *citation
				counter = counter + 1
			}
		}

	}

	if medlineCitation.ChemicalList != nil {
		dbArticle.Chemicals = make([]Chemical, len(medlineCitation.ChemicalList.Chemical))
		for i, chemical := range medlineCitation.ChemicalList.Chemical {
			dbChemical := new(Chemical)
			dbChemical.Name = chemical.NameOfSubstance.Text
			dbChemical.Registry = chemical.RegistryNumber.Text
			dbArticle.Chemicals[i] = *dbChemical
		}

	}

	// if medlineCitation.MeshHeadingList != nil {
	// 	dbArticle.MeshTerms = make([]MeshTerm, len(medlineCitation.MeshHeadingList.MeshHeading))
	// 	for i, mesh := range medlineCitation.MeshHeadingList.MeshHeading {
	// 		dbMesh := new(MeshTerm)
	// 		dbMesh.Descriptor = mesh.DescriptorName.Text
	// 		//dbMesh.Qualifier = mesh.QualifierName.Text
	// 		dbArticle.MeshTerms[i] = *dbMesh
	// 	}
	// }

	if pArticle.Journal != nil {
		//journal := Journal{}
		//db.First(&journal, 10)
		//db.First(&user, 10)
		//db.Where("name = ?", "hello world").First(&User{}).Error == gorm.RecordNotFound
		//fmt.Println(pArticle.Journal.Title.Text)
		journal := Journal{
			Title: pArticle.Journal.Title.Text,
		}
		//journal := new(Journal)
		//journal.Id = JournalIdCounter
		//journal.Title = pArticle.Journal.Title.Text
		if pArticle.Journal.ISSN != nil {
			journal.Issn = pArticle.Journal.ISSN.Text
		}
		dbArticle.Journal = journal
		//dbArticle.journal_id.Int64 = journal.Id
		//dbArticle.journal_id.Valid = true
		//JournalIdCounter = JournalIdCounter + 1

	}

	if pArticle.AuthorList != nil {
		dbArticle.Authors = make([]Author, len(pArticle.AuthorList.Author))
		for i, author := range pArticle.AuthorList.Author {
			dbAuthor := new(Author)
			if author.Identifier != nil {
				//dbAuthor.Id = author.Identifier.Text
			}
			if author.LastName != nil {
				dbAuthor.LastName = author.LastName.Text
			}
			if author.ForeName != nil {
				dbAuthor.FirstName = author.ForeName.Text
			}
			if author.Affiliation != nil {
				dbAuthor.Affiliation = author.Affiliation.Text
			}
			dbArticle.Authors[i] = *dbAuthor
		}
	}

	return dbArticle
}

func articleAdder(articleChannel chan []*Article, done chan bool, db *gorm.DB, commitSize int) {
	log.Println("Start articleAdder")
	tx := db.Begin()
	t0 := time.Now()
	var totalCount int64 = 0
	counter := 0
	chunkCount := 0
	for articleArray := range articleChannel {
		log.Println("-- Consuming chunk ", chunkCount)

		log.Printf("articleAdder counter=%d", counter)
		log.Printf("TOTAL counter=%d", totalCount)

		log.Println(commitSize)
		if doNotWriteToDbFlag {
			counter = counter + len(articleArray)
			totalCount = totalCount + int64(len(articleArray))
			continue
		}

		tmp := articleArray
		for i := 0; i < len(tmp); i++ {
			article := tmp[i]
			if article == nil {
				//log.Println(i, " ******** Article is nil")
				continue
			}

			counter = counter + 1
			totalCount = totalCount + 1
			closeOpenCount = closeOpenCount + 1
			if counter == commitSize {
				tc0 := time.Now()
				tx.Commit()
				t1 := time.Now()
				log.Printf("The commit took %v to run.\n", t1.Sub(tc0))
				log.Printf("The call took %v to run.\n", t1.Sub(t0))
				t0 = time.Now()
				counter = 0

				if closeOpenCount >= CloseOpenSize {
					log.Println("CLOSEOPEN $$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$ ", closeOpenCount)
					var err error
					db, err = dbCloseOpen(db)
					if err != nil {
						log.Fatal(err)
					}
					closeOpenCount = 0
				}

				tx = db.Begin()
			}
			if err := tx.Create(article).Error; err != nil {
				tx.Rollback()
				log.Println("\\\\\\\\\\\\\\\\")
				log.Println("[", err, "]")
				log.Printf("PMID=%d", article.ID)
				//if !strings.HasSuffix(err.Error(), "PRIMARY KEY must be unique") {
				//continue
				//}
				//log.Println("Returning from articleAdder")
				//log.Fatal(" Fatal\\\\\\\\\\\\\\\\")
				//return
				tx = db.Begin()
			}

		}
		log.Println("-- END chunk ", chunkCount)
	}
	if !doNotWriteToDbFlag {
		tx.Commit()
		makeIndexes(db)
	}
	db.Close()
	done <- true
}

// From: http://www.goinggo.net/2013/11/using-log-package-in-go.html
var (
	Trace   *log.Logger
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
)

func logInit(
	traceHandle io.Writer,
	infoHandle io.Writer,
	warningHandle io.Writer,
	errorHandle io.Writer) {

	Trace = log.New(traceHandle,
		"TRACE: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Info = log.New(infoHandle,
		"INFO: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Warning = log.New(warningHandle,
		"WARNING: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Error = log.New(errorHandle,
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)
}

func medlineDate2Year(md string) int {
	// case <MedlineDate>1952 Mar-Apr</MedlineDate>
	var year int
	var err error
	// case 2000-2001
	if string(md[4]) == string('-') {
		yearStrings := strings.Split(md, "-")
		//case 1999-00
		if len(yearStrings[1]) != 4 {
			year, err = strconv.Atoi(yearStrings[0])
		} else {
			year, err = strconv.Atoi(yearStrings[1])
		}
		if err != nil {
			log.Println("error!! ", err)
		}
	} else {
		// case 1999 June 6
		yearString := strings.TrimSpace(strings.Split(md, " ")[0])
		yearString = yearString[0:4]
		//year, err = strconv.Atoi(strings.TrimSpace(strings.Split(md, " ")[0]))
		year, err = strconv.Atoi(yearString)
		if err != nil {
			log.Println("error!! ", err)
		}
	}
	if year == 0 {
		log.Println("medlineDate2Year ", year, md, " [", strings.TrimSpace(string(md[4])), "]")
	}
	return year

}
