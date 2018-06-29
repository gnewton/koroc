package main

import (
	"database/sql"
	"flag"
	//"github.com/pkg/profile"
	"github.com/jinzhu/gorm"
	"io"
	"io/ioutil"
	"log"
	"math"
	"os"
	"sync"
	"time"
)

var TransactionSize = 50000

//var TransactionSize = 1000

var chunkSize = 1000

//var chunkSize = 100
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
	if true {
		foo()
		return
	}
	var wg sync.WaitGroup
	var wgExtract sync.WaitGroup

	//defer profile.Start().Stop()

	if meshFileName != "" {
		loadMesh(meshFileName)
	}

	dbc := DBConnector{dbFilename: dbFilename}

	gdb, err := dbc.Open()

	if err != nil {
		Error.Fatal(err)
		return
	}
	defer func() {
		err = gdb.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	dbInit(gdb)
	gdb.Close()

	db, err := dbOpen2(dbFilename)

	articleChannel := make(chan []*Article, 20)

	wg.Add(1)
	go articleAdder3(articleChannel, db, TransactionSize, &wg)

	// Loop through files

	n := len(flag.Args())
	filenameChannel := make(chan string, n)

	numExtractors := 8
	for i := 0; i < numExtractors; i++ {
		wgExtract.Add(1)
		go readFromFileAndExtractXML(filenameChannel, &dbc, articleChannel, &wgExtract)
	}

	for _, filename := range flag.Args() {
		log.Println(" -- Input file: " + filename)
		filenameChannel <- filename
	}

	close(filenameChannel)
	wgExtract.Wait()
	close(articleChannel)
	wg.Wait()
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
