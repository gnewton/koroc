package main

import (
	"database/sql"
	"errors"
	"log"
)

type Base struct {
	counter    *CounterImpl
	db         *sql.DB
	tx         *sql.Tx
	insertPrep *sql.Stmt
	updatePrep *sql.Stmt
	cache      *Cache
}

type Persist struct {
	Base
	stringCache *StringCache
}

type ArticlePersist struct {
	Persist
	chemicalPersist *Persist
	commitSize      int
	commitCounter   int
}

const createArticlesTable = "CREATE TABLE \"articles\" (\"abstract\" varchar(255),\"day\" integer,\"id\" integer primary key autoincrement,\"issue\" varchar(255),\"journal_id\" bigint,\"keywords_owner\" varchar(255),\"language\" varchar(255),\"month\" varchar(8),\"title\" varchar(255),\"volume\" varchar(255),\"year\" integer,\"date_revised\" bigint );"

const prepInsertArticle = "INSERT INTO articles (abstract,day,id,issue,journal_id,keywords_owner,language,month,title,volume,year,date_revised) VALUES (?,?,?,?,?,?,?,?,?,?,?,?)"

const prepUpdateArticle = "UPDATE articles set abstract=?,day=?,id=?,issue=?,journal_id=?,keywords_owner=?,language=?,month=?,title=?,volume=?,year=?,date_revised=? where id=?"

func NewArticlePersist(db *sql.DB, commitSize int) (*ArticlePersist, error) {
	if db == nil {
		return nil, errors.New("DB is nil")
	}
	if commitSize < 0 {
		return nil, errors.New("Commit size cannot be less than zero")
	}

	if commitSize == 0 {
		// default
		commitSize = 5000
	}

	var err error
	ap := new(ArticlePersist)
	ap.db = db

	ap.tx, ap.insertPrep, ap.updatePrep, err = newTx(ap.db, prepInsertArticle, prepUpdateArticle)
	if err != nil {
		return nil, err
	}

	ap.counter = NewCounter()
	ap.commitSize = commitSize
	ap.cache = NewCache()
	ap.chemicalPersist, err = NewChemicalPersist(ap.tx)
	if err != nil {
		return nil, err
	}
	return ap, nil
}

func (ap *ArticlePersist) Save(article *Article) error {
	var err error = nil

	if article == nil {
		return nil
	}
	article.ID = ap.counter.Next()

	if ap.cache.Exists(article.ID) {
		//delete articled
	}

	var abstract *string
	if len(article.Abstract) == 0 {
		abstract = nil
	} else {
		abstract = &article.Abstract
	}

	var title *string
	if len(article.Title) == 0 {
		title = nil
	} else {
		title = &article.Title
	}

	vars := []interface{}{abstract, article.Day, article.ID, article.Issue, article.JournalID, article.KeywordsOwner, article.Language, article.Month, title, article.Volume, article.Year, article.DateRevised}

	_, err = ap.tx.Stmt(ap.insertPrep).Exec(vars)
	//_, err = ap.tx.Stmt(ap.insertPrep).Exec(abstract, article.Day, article.ID, article.Issue, article.JournalID, article.KeywordsOwner, article.Language, article.Month, title, article.Volume, article.Year, article.DateRevised)
	if err != nil {
		log.Fatal(err)
		//return err
	}

	if article.Authors != nil && len(article.Authors) > 0 {
		for _ = range article.Authors {
			//log.Println(auth)
			//auth.Save(db, article.ID)
		}
	}

	for _, chem := range article.Chemicals {
		//log.Println(chem)
		key := chem.Name + chem.Registry
		//if value, ok := ap.chemicalPersist.stringCache.Exists(key); ok {
		if _, ok := ap.chemicalPersist.stringCache.Exists(key); ok {
			// use value to make join of chem.id and article.id
		} else {
			id := ap.chemicalPersist.counter.Next()
			_, err = ap.tx.Stmt(ap.chemicalPersist.insertPrep).Exec(id, chem.Name, chem.Registry)
			ap.chemicalPersist.stringCache.Add(key, id)
		}
		if err != nil {
			log.Fatal(err)
			//return err
		}
		//chem.Save(db, article.ID)
	}

	// for citation := range article.Citations {
	// 	log.Println(citation)
	// 	//chem.Save(db, article.ID)
	// }

	// for desc := range article.MeshDescriptors {
	// 	log.Println(desc)
	// 	//chem.Save(db, article.ID)
	// }
	ap.commitCounter += 1
	if ap.commitCounter > ap.commitSize {
		log.Println("Starting commit")
		err = ap.tx.Commit()
		if err != nil {
			return err
		}
		log.Println("Done commit")
		ap.tx, ap.insertPrep, ap.updatePrep, err = newTx(ap.db, prepInsertArticle, prepUpdateArticle)
		ap.commitCounter = 0

	}
	return err
}

const createChemTable = "CREATE TABLE IF NOT EXISTS \"chemicals\" (\"id\" integer primary key autoincrement,\"name\" varchar(255),\"registry\" varchar(32) );"
const prepInsertChem = "INSERT INTO chemicals (id,name,registry) VALUES (?,?,?)"
const prepUpdateChem = "UPDATE chemicals set id=?, name=?,registry=?"

func NewChemicalPersist(tx *sql.Tx) (*Persist, error) {
	p := new(Persist)

	p.counter = NewCounter()
	p.stringCache = NewStringCache()
	var err error
	p.insertPrep, p.updatePrep, err = newPreparedStatements(tx, prepInsertChem, prepUpdateChem)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return p, nil
}

func newPreparedStatements(tx *sql.Tx, insertString, updateString string) (*sql.Stmt, *sql.Stmt, error) {
	insertPrep, err := tx.Prepare(insertString)
	if err != nil {
		return nil, nil, err
	}

	updatePrep, err := tx.Prepare(updateString)
	if err != nil {
		return nil, nil, err
	}
	return insertPrep, updatePrep, nil
}

func newTx(db *sql.DB, insertString, updateString string) (*sql.Tx, *sql.Stmt, *sql.Stmt, error) {
	log.Println("newTx")
	tx, err := db.Begin()
	if err != nil {
		return nil, nil, nil, err
	}
	insertPrep, updatePrep, err := newPreparedStatements(tx, insertString, updateString)

	return tx, insertPrep, updatePrep, err

}
