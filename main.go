package main

import (
	"database/sql"

	"flag"
	"io"
	"io/ioutil"
	"log"
	"math"
	"os"
	"sync"
	"time"

	"github.com/gnewton/pubmedSqlStructs"
	"github.com/pkg/profile"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
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
var loggingFlag = false
var sanitizeStringsFlag = false

const CommentsCorrections_RefType = "Cites"
const PUBMED_ARTICLE = "PubmedArticle"

var out int = -1
var JournalIdCounter int64 = 0
var counters map[string]*int
var articleIdsInDBCache map[uint32]uint8 = make(map[uint32]uint8, 100000)
var closeOpenCount int64 = 0

var empty struct{}

type foo struct{}

type ArticlesEnvelope struct {
	articles []*pubmedSqlStructs.Article
	n        int
}

func init() {

}

func main() {

	//defer profile.Start(profile.CPUProfile).Stop()
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	flag.BoolVar(&sqliteLogFlag, "L", sqliteLogFlag, "Turn on sqlite logging")
	flag.BoolVar(&loggingFlag, "l", loggingFlag, "Turn on verbose logging")
	flag.StringVar(&dbFilename, "f", dbFilename, "SQLite output filename")
	flag.StringVar(&meshFileName, "m", meshFileName, "MeSH descriptor sqlite3 filename")

	flag.IntVar(&TransactionSize, "t", TransactionSize, "Size of transactions")
	flag.IntVar(&chunkSize, "C", chunkSize, "Size of chunks")
	flag.Int64Var(&CloseOpenSize, "z", CloseOpenSize, "Num of records before sqlite connection is closed then reopened")
	flag.Int64Var(&LoadNRecordsPerFile, "N", LoadNRecordsPerFile, "Load only N records from each file")
	flag.BoolVar(&sqliteLogFlag, "V", sqliteLogFlag, "Turn on sqlite logging")
	flag.BoolVar(&sanitizeStringsFlag, "s", sanitizeStringsFlag, "Removes xml tags from strings")

	flag.BoolVar(&doNotWriteToDbFlag, "X", doNotWriteToDbFlag, "Do not write to db. Rolls back transaction. For debugging")

	flag.Parse()

	if len(flag.Args()) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	logInit(loggingFlag, ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)
	defer profile.Start(profile.MemProfile).Stop()
	if meshFileName != "" {
		loadMesh(meshFileName)
	}

	dbc := DBConnector{dbFilename: dbFilename}

	db, err := dbc.Open()

	if err != nil {
		Error.Fatal(err)
		return
	}

	dbInit(db)

	numExtractors := 12
	//articleChannel := make(chan []*pubmedSqlStructs.Article, numExtractors*3)
	articleChannel := make(chan ArticlesEnvelope, numExtractors*3)

	var addWg sync.WaitGroup
	var extractWg sync.WaitGroup

	addWg.Add(1)
	go articleAdder(articleChannel, &dbc, db, TransactionSize, &addWg)

	for i, filename := range flag.Args() {
		log.Println(i, " -- Input file: "+filename)
	}

	// Loop through files
	n := len(flag.Args())
	readFileChannel := make(chan string, n)

	for i := 0; i < numExtractors; i++ {
		extractWg.Add(1)
		go readFromFileAndExtractXML(i, readFileChannel, &dbc, articleChannel, &extractWg)

	}

	for _, filename := range flag.Args() {
		log.Println("Pushing file into channel", filename)
		readFileChannel <- filename
	}
	log.Println("Done pushing files")
	close(readFileChannel)
	log.Println("readFileChannel closed")
	extractWg.Wait()
	log.Println("Post 	extractWg.Wait")
	close(articleChannel)
	log.Println("articleChannel closed")
	addWg.Wait()
	log.Println("Post 	addWg.Wait")

}

// From: http://www.goinggo.net/2013/11/using-log-package-in-go.html
var (
	Trace   *log.Logger
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
)

func logInit(
	loggingFlag bool,
	traceHandle io.Writer,
	infoHandle io.Writer,
	warningHandle io.Writer,
	errorHandle io.Writer) {

	if !loggingFlag {
		log.SetOutput(ioutil.Discard)
		return
	}

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

const selectArticleIDs = "select id from articles"

type DBConnector struct {
	dbFilename string
	gdb        *gorm.DB
}

func (dbc *DBConnector) Open() (*gorm.DB, error) {
	var err error
	dbc.gdb, err = gorm.Open(sqlite.Open(dbc.dbFilename), &gorm.Config{})
	return dbc.gdb, err
}

func (dbc *DBConnector) DB() *gorm.DB {
	return dbc.gdb
}

func makeArticleIdsInDBCache(db *sql.DB) (map[uint32]struct{}, error) {
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
	articleIdsInDB := make(map[uint32]struct{}, 10000)
	count := 0
	for rows.Next() {
		count += 1
		var id uint32
		err = rows.Scan(&id)
		if err != nil {
			return nil, err
		}
		articleIdsInDB[id] = empty
	}
	log.Printf("The database took %v to load cache. Size:%d\n", time.Now().Sub(t0), count)
	return articleIdsInDB, err
}
