package main

import (
	"database/sql"
	"errors"
	"log"
)

//const P_ARTICLE = "insert into articles (full_abstract,abstract_copyright_information,coi_statement,copyright_information,date_revised, id, issue, journal_id, language, pub_day, pub_month, pub_year, pagination, retracted, title, version, volume) values (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"

const P_ARTICLE = "insert into articles (id, title) values(?,?)"

const P_ARTICLE_DELETE = "delete from articles where id=?"

func (a *Article) Save(stmt *sql.Stmt) error {
	if stmt == nil {
		return errors.New("Statement is nil")
	}
	_, err := stmt.Exec(a.ID, a.Title)
	//log.Println("Inserting article:", a.ID, a.Title)
	if err != nil {
		log.Fatal(err)
	}
	return err
}

func (a Article) MakePreparedStatements(tx *sql.Tx) (*sql.Stmt, *sql.Stmt, error) {
	if tx == nil {
		log.Println(errors.New("Database is nil"))
	}
	insert, err := tx.Prepare(P_ARTICLE)
	if err != nil {
		return nil, nil, err
	}
	delete, err := tx.Prepare(P_ARTICLE_DELETE)
	if err != nil {
		return nil, nil, err
	}
	return insert, delete, nil
}

func (a *Article) Delete(stmt *sql.Stmt) error {
	if stmt == nil {
		return errors.New("Statement is nil")
	}
	r, err := stmt.Exec(a.ID)
	if err != nil {
		return err
	}
	n, err := r.RowsAffected()
	if err != nil {
		return err
	}
	if n != 1 && err == nil {
		err = errors.New("Wrong number of articles")
		log.Println(stmt)
		log.Println("id=", a.ID)
		log.Println("n=", n)
		log.Println(err)
		return err
	}
	return err

}
