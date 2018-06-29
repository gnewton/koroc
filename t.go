package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
)

type Named interface {
}

type Column struct {
	field Field
}

type HasDB struct {
	db *sql.DB
}

type TableWriter interface {
	SetDB(*sql.DB) error
	CreateTable(*Table) error
}

type RowWriter interface {
	Init(*sql.DB) error
	SetTxSize(int) error
	Close() error
	Write(*Row) error
}

type RowWriterImp struct {
	HasDB
	txSize int
}

func (w *RowWriterImp) Init(db *sql.DB) error {
	if db == nil {
		return errors.New("db is nil")
	}
	w.db = db
	w.txSize = 10000
	return nil
}

func (w *RowWriterImp) SetTxSize(s int) error {
	if s < 1 {
		return errors.New("tx size is < 1")
	}
	w.txSize = s
	return nil
}

func (w *RowWriterImp) Close() error {
	if w.db == nil {
		return errors.New("db is nil")
	}
	return w.db.Close()
}

func (w *RowWriterImp) Write(row *Row) error {
	t := row.table

	s := "insert into " + t.name + "("
	values := " values ("
	for i, _ := range t.Fields {
		if i != 0 {
			s += ","
			values += ","
		}
		s += t.Fields[i].GetName()

	}
	s += ")" + values
	fmt.Println("===========================", s)
	return nil
}

func foo() {
	db, err := dbOpen("test1.sqlite3")
	defer db.Close()

	if false {
		t, err := NewTable("articles")
		t.SetDefaultFields()

		title, err := t.NewFieldString("title")
		if err != nil {
			log.Fatal(err)
		}
		title.Length(64)
		title.NotNullable()

		journal_id, err := t.NewFieldInt64("journal_id")
		if err != nil {
			log.Fatal(err)
		}
		journal_id.NotNullable()
		journal_id.Indexed()

		articleRow := t.NewRow()
		err = articleRow.SetString("title", "this is the title")
		if err != nil {
			log.Fatal(err)
		}
		err = articleRow.SetInt64("id", 388)
		if err != nil {
			log.Fatal(err)
		}

		err = articleRow.SetInt64("journal_id", 76355)
		if err != nil {
			log.Fatal(err)
		}

		rw := new(RowWriterImp)
		rw.Write(articleRow)

		////////////////////////////////////////////////////////////////////////
	}
	journals, err := NewTable("journals")
	journals.SetDefaultFields()

	issn, err := journals.NewFieldString("issn")
	if err != nil {
		log.Fatal(err)
	}
	issn.Length(16)
	issn.NotNullable()

	name, err := journals.NewFieldString("name")
	if err != nil {
		log.Fatal(err)
	}
	name.Length(64)
	name.NotNullable()

	//journalRelation := t.NewOneToMany(journals, journal_id, issn)
	//_, err = t.NewOneToMany(journals, journal_id, issn)
	// _, err = t.NewOneToMany(journals, issn, issn)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	journalRow := journals.NewRow()
	err = journalRow.SetInt64("id", 1077577)
	if err != nil {
		log.Fatal(err)
	}

	err = journalRow.SetString("issn", "lkajbasdflkajsfbajksdfb")
	if err != nil {
		log.Fatal(err)
	}

	fields := journalRow.fieldNamesForPreparedStatement()
	log.Println("journalRow.fieldsForPreparedStatement", fields)

	fieldValues := journalRow.fieldValuesForPreparedStatement()
	log.Println("journalRow.fieldValuesForPreparedStatement", fieldValues)
	//fields := fieldsForPreparedStatement
	//

	//log.Println(row.fieldsForPreparedStatement())
	// log.Println(row)
	// log.Println(t)
	fmt.Println("-----------")
	// t.Print()
	// articleRow.AddRelationRow(journalRow)
	// articleRow.Print()

}
