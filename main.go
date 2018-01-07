package main

import (
	"database/sql"
	"encoding/xml"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gnewton/pubmedSqlStructs"
	"github.com/gnewton/pubmedstruct"
	"github.com/jinzhu/gorm"
)

var TransactionSize = 100000

var chunkSize = 10000
var CloseOpenSize int64 = 99950000
var chunkChannelSize = 3

var dbFilename = "./pubmed_sqlite.db"
var meshFileName = ""
var sqliteLogFlag = false
var LoadNRecordsPerFile int64 = math.MaxInt64
var recordPerFileCounter int64 = 0
var doNotWriteToDbFlag = false

const CommentsCorrections_RefType = "Cites"
const PUBMED_ARTICLE = "PubmedArticle"

var out int = -1
var JournalIdCounter int64 = 0
var counters map[string]*int
var articleIdsInDBCache map[int64]int = make(map[int64]int, 100000)
var closeOpenCount int64 = 0

var empty struct{}

func init() {

	//defer profile.Start(profile.CPUProfile).Stop()
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	flag.BoolVar(&sqliteLogFlag, "L", sqliteLogFlag, "Turn on sqlite logging")
	flag.StringVar(&dbFilename, "f", dbFilename, "SQLite output filename")
	flag.StringVar(&meshFileName, "m", meshFileName, "MeSH descriptor sqlite3 filename")

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

	if meshFileName != "" {
		loadMesh(meshFileName)
	}

	dbc := DBConnector{dbFilename: dbFilename}

	db, err := dbc.Open()

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

	dbInit(db)
	db.Close()

	articleChannel := make(chan []*pubmedSqlStructs.Article, 20)
	txChannel := make(chan *gorm.DB, 5)

	done := make(chan bool, 8)
	articleAdderDone := make(chan bool)
	go articleAdder(articleChannel, &dbc, db, txChannel, TransactionSize, articleAdderDone)
	//go articleAdder2(articleChannel, db, TransactionSize)
	go committor(txChannel, done)

	for i, filename := range flag.Args() {
		log.Println(i, " -- Input file: "+filename)
	}

	// Loop through files

	n := len(flag.Args())
	filenameChannel := make(chan string, n)

	numExtractors := 5
	for i := 0; i < numExtractors; i++ {
		go readFromFileAndExtractXML(filenameChannel, &dbc, articleChannel, done)
	}

	for _, filename := range flag.Args() {
		filenameChannel <- filename
	}
	close(filenameChannel)

	//for _, _ = range flag.Args() {
	for i := 0; i < numExtractors; i++ {

		_ = <-done
	}

	close(articleChannel)
	_ = <-articleAdderDone
}

// Reads from the channel arrays of pubmedarticles and
func readFromFileAndExtractXML(c chan string, dbc *DBConnector, articleChannel chan []*pubmedSqlStructs.Article, done chan bool) {
	articleArray := make([]*pubmedSqlStructs.Article, chunkSize)
	for filename := range c {
		if filename == "" {
			break
		}
		log.Println("Opening: " + filename)
		reader, _, err := genericReader(filename)

		if err != nil {
			log.Fatal(err)
			return
		}

		decoder := xml.NewDecoder(reader)
		counters = make(map[string]*int)

		count := 0

		// Loop through XML
		for {

			if recordPerFileCounter == LoadNRecordsPerFile {
				log.Println("break file load. LoadNRecordsPerFile", count, LoadNRecordsPerFile)
				recordPerFileCounter = 0
				break
			}
			token, err := decoder.Token()
			if err != nil {
				if err == io.EOF {
					break
				}
				log.Println("Fatal error in file:", filename)
				//log.Fatal(err)
			}
			if token == nil {
				break
			}
			switch se := token.(type) {
			case xml.StartElement:
				if se.Name.Local == PUBMED_ARTICLE && se.Name.Space == "" {
					if count%10000 == 0 && count != 0 {
						log.Printf("count=%d\n", count)
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
					articleArray[count] = article
					if count == chunkSize-1 {
						articleChannel <- articleArray
						articleArray = make([]*pubmedSqlStructs.Article, chunkSize)
						count = 0
					}

				}
			}
		}
	}

	// Is there something not yet sent in the articleArray?
	if len(articleArray) > 0 {
		articleChannel <- articleArray
	}

	done <- true
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

	dbArticle.Version, err = strconv.Atoi(p.MedlineCitation.PMID.Attr_Version)
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

	makeDataBanks(pArticle, dbArticle)

	if p.PubmedData.ArticleIdList != nil {
		dbArticle.ArticleIDs = makeArticleIdList(p.PubmedData.ArticleIdList)
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
	//log.Println("mesh")
	//log.Println(medlineCitation.MeshHeadingList)
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

func makeArticleIdList(alist *pubmedstruct.ArticleIdList) []*pubmedSqlStructs.ArticleID {
	if alist.ArticleId == nil {
		return nil
	}
	arts := make([]*pubmedSqlStructs.ArticleID, len(alist.ArticleId))

	for i, _ := range alist.ArticleId {
		aid := new(pubmedSqlStructs.ArticleID)
		aid.Type = alist.ArticleId[i].Attr_IdType
		aid.ID = alist.ArticleId[i].Text
		arts[i] = aid
	}

	return arts
}

var databankCounter int64 = 0

func makeDataBanks(src *pubmedstruct.Article, dest *pubmedSqlStructs.Article) {
	if src.DataBankList == nil || src.DataBankList.DataBank == nil {
		return
	}

	if dest.DataBanks == nil {
		dest.DataBanks = make([]*pubmedSqlStructs.DataBank, len(src.DataBankList.DataBank))
	}
	for i, _ := range src.DataBankList.DataBank {
		bank := src.DataBankList.DataBank[i]
		dest.DataBanks[i] = new(pubmedSqlStructs.DataBank)
		dest.DataBanks[i].Name = bank.DataBankName.Text
		log.Println(bank.DataBankName.Text)
		dest.DataBanks[i].ID = databankCounter
		databankCounter = databankCounter + 1

		if dest.DataBanks[i].AccessionNumbers == nil {
			dest.DataBanks[i].AccessionNumbers = make([]*pubmedSqlStructs.AccessionNumber, len(bank.AccessionNumberList.AccessionNumber))
		}
		for j, _ := range bank.AccessionNumberList.AccessionNumber {
			dest.DataBanks[i].AccessionNumbers[j] = new(pubmedSqlStructs.AccessionNumber)
			dest.DataBanks[i].AccessionNumbers[j].Number = bank.AccessionNumberList.AccessionNumber[j].Text
			dest.DataBanks[i].AccessionNumbers[j].ID = databankCounter
			databankCounter = databankCounter + 1
			log.Println("\t", bank.AccessionNumberList.AccessionNumber[j].Text)
		}

	}
}

const createArticlesTable = "CREATE TABLE \"articles\" (\"abstract\" varchar(255),\"day\" integer,\"id\" integer primary key autoincrement,\"issue\" varchar(255),\"journal_id\" bigint,\"keywords_owner\" varchar(255),\"language\" varchar(255),\"month\" varchar(8),\"title\" varchar(255),\"volume\" varchar(255),\"year\" integer,\"date_revised\" bigint );"

const prepInsertArticle = "INSERT INTO articles (abstract,day,id,issue,journal_id,keywords_owner,language,month,title,volume,year,date_revised) VALUES (?,?,?,?,?,?,?,?,?,?,?,?)"

const prepUpdateArticle = "UPDATE articles set abstract=?,day=?,id=?,issue=?,journal_id=?,keywords_owner=?,language=?,month=?,title=?,volume=?,year=?,date_revised=? where id=?"

func articleAdder2(articleChannel chan []*pubmedSqlStructs.Article, db *sql.DB, commitSize int) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	stmtInsert, err := tx.Prepare(prepInsertArticle)
	if err != nil {
		log.Fatal(err)
	}

	stmtUpdate, err := tx.Prepare(prepUpdateArticle)
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
				continue
			}

			// Have we already inserted this article (i.e. is this an update?)
			if _, ok := articleIdsInDBCache[a.ID]; ok {
				log.Println("Updating article:", a.ID)
				_, err = stmtUpdate.Exec(a.Abstract, a.Day, a.ID, a.Issue, a.JournalID, a.KeywordsOwner, a.Language, a.Month, a.Title, a.Volume, a.Year, a.DateRevised, a.ID)
			} else {
				articleIdsInDBCache[a.ID] = a.Version
				_, err = stmtInsert.Exec(a.Abstract, a.Day, a.ID, a.Issue, a.JournalID, a.KeywordsOwner, a.Language, a.Month, a.Title, a.Volume, a.Year, a.DateRevised)

			}

			if err != nil {
				log.Println(err)
				if err.Error() == "UNIQUE constraint failed: articles.id" {
					log.Println("*** ", a.ID, "|||", a.Abstract, a.Day, a.Issue, a.JournalID, a.KeywordsOwner, a.Language, a.Month, a.Title, a.Volume, a.Year, a.DateRevised)
					continue
				}
				log.Println(a.ID, "|||", a.Abstract, a.Day, a.Issue, a.JournalID, a.KeywordsOwner, a.Language, a.Month, a.Title, a.Volume, a.Year, a.DateRevised)
				log.Fatal(err)
			}

			counter = counter + 1
			totalCount = totalCount + 1
			if counter == commitSize {
				counter = 0
				log.Println("************ committing", totalCount)
				tx.Commit()
				stmtInsert.Close()
				stmtUpdate.Close()
				tx, err = db.Begin()
				if err != nil {
					log.Fatal(err)
				}
				stmtInsert, err = tx.Prepare(prepInsertArticle)
				if err != nil {
					log.Fatal(err)
				}
				stmtUpdate, err = tx.Prepare(prepUpdateArticle)
				if err != nil {
					log.Fatal(err)
				}
			}
		}
	}
	tx.Commit()
	stmtInsert.Close()
	stmtUpdate.Close()
	db.Close()

}

func updateArticle(article *pubmedSqlStructs.Article) (sql.Result, error) {
	return nil, nil
}

func insertArticle(article *pubmedSqlStructs.Article) (sql.Result, error) {
	return nil, nil
}

func committor(transactionChannel chan *gorm.DB, done chan bool) {
	for tx := range transactionChannel {
		log.Println("COMMIT starting")
		tx.Commit()
		tx.Close()
		log.Println("COMMIT done")
	}
	done <- true
}

func articleAdder(articleChannel chan []*pubmedSqlStructs.Article, dbc *DBConnector, db *gorm.DB, txChannel chan *gorm.DB, commitSize int, done chan bool) {
	log.Println("Start articleAdder")
	var err error
	db, err = dbc.Open()
	if err != nil {
		log.Fatal(err)
	}
	tx := db.Begin()
	//t0 := time.Now()
	var totalCount int64 = 0
	counter := 0
	chunkCount := 0
	for articleArray := range articleChannel {
		log.Println("-- Consuming chunk ", chunkCount)

		log.Printf("articleAdder counter=%d", counter)
		log.Printf("TOTAL counter=%d", totalCount)

		log.Println("Commit size=", commitSize)
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
				//tc0 := time.Now()
				//tx.Commit()
				log.Printf("Transaction channel length=%d", len(txChannel))
				txChannel <- tx
				//log.Println("transaction")
				//log.Println(tx)
				var err error
				tx, err = dbc.Open()
				if err != nil {
					log.Fatal(err)
				}
				//t1 := time.Now()
				//log.Printf("The commit took %v to run.\n", t1.Sub(tc0))
				//log.Printf("The call took %v to run.\n", t1.Sub(t0))
				//t0 = time.Now()
				counter = 0
				tx = tx.Begin()
				//log.Println("transaction")
				//log.Println(tx)
			}
			var err error
			if _, ok := articleIdsInDBCache[article.ID]; ok {
				log.Println("Updating article:", article.ID, article.Version)

				var oldArticle pubmedSqlStructs.Article
				tx.Where("ID = ?", article.ID).First(&oldArticle)
				tx.Delete(oldArticle)

				//err = tx.Update(article).Error
				err = tx.Create(article).Error
			} else {
				err = tx.Create(article).Error
				articleIdsInDBCache[article.ID] = article.Version
			}

			if err != nil {
				//log.Println("transaction")
				//log.Println(tx)
				log.Println(err)
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
		var err error
		tx, err = dbc.Open()
		if err != nil {
			log.Fatal(err)
		}
		makeIndexes(tx)
	}
	close(txChannel)
	db.Close()
	log.Println("-- END articleAdder")
	done <- true
	log.Println("++ END articleAdder")
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

var dateSeasonYear = []string{"Summer", "Winter", "Spring", "Fall"}

func medlineDate2Year(md string) int {
	// Other cases:
	//   Fall 2017; 8/15/12; Spring 2017; Summer 2017; Fall 2017;

	// case <MedlineDate>1952 Mar-Apr</MedlineDate>
	var year int
	var err error

	// case 2000-2001
	//log.Println(md)

	for i, _ := range dateSeasonYear {
		if strings.HasPrefix(md, dateSeasonYear[i]) {
			return seasonYear(md)
		}
	}

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
			log.Println("error!! yearString=[", yearString, "]", err)
		}
	}
	if year == 0 {
		log.Println("medlineDate2Year [", md, "] [", strings.TrimSpace(string(md[4])), "]")
	}
	return year

}

func seasonYear(md string) int {
	parts := strings.Split(md, " ")
	year, err := strconv.Atoi(parts[1])
	if err != nil {
		log.Fatal(err)
	}
	return year
}

const selectArticleIDs = "select id from articles"

type DBConnector struct {
	dbFilename string
}

func (dbc *DBConnector) Open() (*gorm.DB, error) {
	return gorm.Open("sqlite3", dbc.dbFilename)
}

func makeArticleIdsInDBCache(db *sql.DB) (map[int64]struct{}, error) {
	tx, err := db.Begin()
	defer tx.Commit()
	if err != nil {
		return nil, err
	}
	t0 := time.Now()

	rows, err := tx.Query(selectArticleIDs)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	articleIdsInDB := make(map[int64]struct{}, 10000)
	count := 0
	for rows.Next() {
		count += 1
		var id int64
		err = rows.Scan(&id)
		if err != nil {
			return nil, err
		}
		articleIdsInDB[id] = empty
	}
	log.Printf("The database took %v to load cache. Size:%d\n", time.Now().Sub(t0), count)
	return articleIdsInDB, err
}
