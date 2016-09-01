package main

import (
	"database/sql"
	"flag"
	"github.com/gnewton/pubmedSqlStructs"
	"github.com/gnewton/pubmedstruct"
	"github.com/jinzhu/gorm"
	"io"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
	//	"net/http"
	//"github.com/davecheney/profile"
	//"strings"
	//_ "net/http/pprof"
	"encoding/xml"

	"runtime/pprof"
)

var TransactionSize = 50000

var chunkSize = 10000
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
	f, err := os.Create("cpuprofile")
	if err != nil {
		log.Fatal(err)
	}
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

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

	articleChannel := make(chan []*pubmedSqlStructs.Article, chunkChannelSize)

	done := make(chan bool)

	//go articleAdder(articleChannel, done, db, TransactionSize)

	db2, err := sql.Open("sqlite3",
		dbFileName+"M")
	if err != nil {
		log.Fatal(err)
	}
	defer db2.Close()
	sqlite3Config(db2)

	_, err = db2.Exec(createArticlesTable)
	if err != nil {
		panic(err)
	}

	go articleAdder2(articleChannel, done, db2, TransactionSize)

	var count int64 = 0
	chunkCount := 0
	arrayIndex := 0

	var articleArray []*pubmedSqlStructs.Article

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
		articleArray = make([]*pubmedSqlStructs.Article, chunkSize)

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
						articleArray = make([]*pubmedSqlStructs.Article, chunkSize)
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

	//log.Println(journalMap)

	f, err = os.Create("memProfile")
	if err != nil {
		log.Fatal(err)
	}
	pprof.WriteHeapProfile(f)
	f.Close()
}

func pubmedArticleToDbArticle(p *pubmedstruct.PubmedArticle) *pubmedSqlStructs.Article {
	medlineCitation := p.MedlineCitation
	pArticle := medlineCitation.Article
	if pArticle == nil {
		log.Println("nil-----------")
		return nil
	}
	var err error
	dbArticle := new(pubmedSqlStructs.Article)
	dbArticle.ID, err = strconv.ParseInt(p.MedlineCitation.PMID.Text, 10, 64)
	if err != nil {
		log.Println(err)
	}

	// Abstract
	dbArticle.Abstract = ""
	if pArticle.Abstract != nil && pArticle.Abstract.AbstractText != nil {
		for i, _ := range pArticle.Abstract.AbstractText {
			dbArticle.Abstract = dbArticle.Abstract + " " + pArticle.Abstract.AbstractText[i].Text
		}
	}

	// Title
	dbArticle.Title = pArticle.ArticleTitle.Text

	// DateRevised
	if p.MedlineCitation.DateRevised != nil {
		d := p.MedlineCitation.DateRevised.Year.Text + p.MedlineCitation.DateRevised.Month.Text + p.MedlineCitation.DateRevised.Day.Text
		dbArticle.DateRevised, err = strconv.ParseInt(d, 10, 64)
	} else {
		if p.MedlineCitation.DateCompleted != nil {
			d := p.MedlineCitation.DateCompleted.Year.Text + p.MedlineCitation.DateCompleted.Month.Text + p.MedlineCitation.DateCompleted.Day.Text
			dbArticle.DateRevised, err = strconv.ParseInt(d, 10, 64)
		} else {
			if p.MedlineCitation.DateCreated != nil {
				d := p.MedlineCitation.DateCreated.Year.Text + p.MedlineCitation.DateCreated.Month.Text + p.MedlineCitation.DateCreated.Day.Text
				dbArticle.DateRevised, err = strconv.ParseInt(d, 10, 64)
			} else {
				log.Println(p.MedlineCitation)
			}
		}
	}

	// Date
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
		}

	}

	//if medlineCitation.OtherID != nil {
	//dbArticle.OtherId = medlineCitation.OtherID
	//}

	if medlineCitation.KeywordList != nil && medlineCitation.KeywordList.Keyword != nil && len(medlineCitation.KeywordList.Keyword) > 0 {
		dbArticle.Keywords = makeKeywords(medlineCitation.KeywordList.Attr_Owner, medlineCitation.KeywordList.Keyword)
		dbArticle.KeywordsOwner = medlineCitation.KeywordList.Attr_Owner
	}

	// Citations
	if medlineCitation.CommentsCorrectionsList != nil {

		actualCitationCount := 0
		for i, _ := range medlineCitation.CommentsCorrectionsList.CommentsCorrections {
			commentsCorrection := medlineCitation.CommentsCorrectionsList.CommentsCorrections[i]
			//log.Printf("%+v\n", *commentsCorrection)

			if commentsCorrection.Attr_RefType == CommentsCorrections_RefType {
				actualCitationCount = actualCitationCount + 1
				//log.Println(commentsCorrection.PMID.Text)
			}
		}

		dbArticle.Citations = make([]*pubmedSqlStructs.Citation, actualCitationCount)
		counter := 0
		var err error
		for i, _ := range medlineCitation.CommentsCorrectionsList.CommentsCorrections {
			commentsCorrection := medlineCitation.CommentsCorrectionsList.CommentsCorrections[i]
			if commentsCorrection.Attr_RefType == CommentsCorrections_RefType {
				citation := new(pubmedSqlStructs.Citation)
				citation.ID, err = strconv.ParseInt(commentsCorrection.PMID.Text, 10, 64)
				if err != nil {
					log.Println(err)
				}
				//citation.RefSource = commentsCorrection.RefSource.Text
				//citation.ID = commentsCorrection.RefSource.Text
				dbArticle.Citations[counter] = citation
				//log.Println("---", dbArticle.Citations[counter].ID)
				counter = counter + 1
			}
		}

	}

	// Chemicals
	if medlineCitation.ChemicalList != nil {
		dbArticle.Chemicals = makeChemicals(medlineCitation.ChemicalList.Chemical)
	}

	//mesh headings
	if medlineCitation.MeshHeadingList != nil {
		dbArticle.MeshDescriptors = makeMeshDescriptors(medlineCitation.MeshHeadingList.MeshHeading)
	}

	if pArticle.Journal != nil {
		defer func() {
			// recover from panic if one occured. Set err to nil otherwise.
			err := recover()
			if err != nil {
				log.Println("@@@@@@@@@@@@@@@@@@@@@@@ ", dbArticle.ID)
				log.Panic(err)
			}

		}()

		foo := makeJournal(pArticle.Journal)
		dbArticle.Journal = foo
	}

	if pArticle.AuthorList != nil {
		dbArticle.Authors = make([]pubmedSqlStructs.Author, len(pArticle.AuthorList.Author))
		for i, _ := range pArticle.AuthorList.Author {
			author := pArticle.AuthorList.Author[i]
			dbAuthor := new(pubmedSqlStructs.Author)
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

const prepArticle = "INSERT INTO articles (abstract,day,id,issue,journal_id,keywords_owner,language,month,title,volume,year,date_revised) VALUES (?,?,?,?,?,?,?,?,?,?,?,?)"
const createArticlesTable = "CREATE TABLE \"articles\" (\"abstract\" varchar(255),\"day\" integer,\"id\" integer primary key autoincrement,\"issue\" varchar(255),\"journal_id\" bigint,\"keywords_owner\" varchar(255),\"language\" varchar(255),\"month\" varchar(8),\"title\" varchar(255),\"volume\" varchar(255),\"year\" integer,\"date_revised\" bigint );"

func articleAdder2(articleChannel chan []*pubmedSqlStructs.Article, done chan bool, db *sql.DB, commitSize int) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := tx.Prepare(prepArticle)
	if err != nil {
		log.Fatal(err)
	}

	chunkCount := 0
	var totalCount int64 = 0
	counter := 0
	for articleArray := range articleChannel {
		log.Println("-- Consuming chunk ", chunkCount)
		chunkCount += 1

		for i := 0; i < len(articleArray); i++ {
			a := articleArray[i]
			if a == nil {
				//log.Println(i, " ******** Article is nil")
				continue
			}
			//log.Println(article.Title)
			_, err = stmt.Exec(a.Abstract, a.Day, a.ID, a.Issue, a.JournalID, a.KeywordsOwner, a.Language, a.Month, a.Title, a.Volume, a.Year, a.DateRevised)
			if a.ID == 20029614 {
				log.Println(a.ID, "|||", a.Abstract, a.Day, a.Issue, a.JournalID, a.KeywordsOwner, a.Language, a.Month, a.Title, a.Volume, a.Year, a.DateRevised)
			}
			if err != nil {
				if err.Error() == "UNIQUE constraint failed: articles.id" {
					log.Println("*** ", a.ID, "|||", a.Abstract, a.Day, a.Issue, a.JournalID, a.KeywordsOwner, a.Language, a.Month, a.Title, a.Volume, a.Year, a.DateRevised)
					continue
				}
				log.Println(a.ID, "|||", a.Abstract, a.Day, a.Issue, a.JournalID, a.KeywordsOwner, a.Language, a.Month, a.Title, a.Volume, a.Year, a.DateRevised)
				log.Println(err)
				log.Fatal(err)
			}

			counter = counter + 1
			totalCount = totalCount + 1
			if counter == commitSize {
				counter = 0
				log.Println("************ committing", totalCount)
				tx.Commit()
				stmt.Close()
				tx, err = db.Begin()
				if err != nil {
					log.Fatal(err)
				}
				stmt, err = tx.Prepare(prepArticle)
				if err != nil {
					log.Fatal(err)
				}
			}
		}
	}
	tx.Commit()
	stmt.Close()
	db.Close()
	done <- true
}

func articleAdder(articleChannel chan []*pubmedSqlStructs.Article, done chan bool, db *gorm.DB, commitSize int) {
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
	//log.Println(md)
	if len(md) == 5 {
		year, err = strconv.Atoi(md)
		if err != nil {
			log.Println("error!! ", err)
			year = 0
		}
		return year
	}
	if len(md) >= 5 && string(md[4]) == string('-') {
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
