package main

import (
	"database/sql"
	//pm "github.com/gnewton/
	"errors"
	pmx "github.com/gnewton/pubmedstruct"
	"log"
)

const P_KW = "insert into keywords (id, owner, major_topic, name) values (?,?,?,?)"

func (kw *Keyword) Save(s *sql.Stmt) {
	//log.Println("Inserting keyword", kw.ID, kw.Owner, kw.MajorTopic, kw.Name)
	_, err := s.Exec(kw.ID, kw.Owner, kw.MajorTopic, kw.Name)
	if err != nil {
		log.Fatal(kw.ID, kw.Owner, kw.MajorTopic, kw.Name)
		log.Fatal(err)
	}
}

func (kw Keyword) MakePreparedStatement(tx *sql.Tx) (*sql.Stmt, error) {
	if tx == nil {
		log.Println(errors.New("Database is nil"))
	}
	return tx.Prepare(P_KW)
}

func ExtractKeywords(owner string, keywords []*pmx.Keyword) []*Keyword {
	newKeywords := make([]*Keyword, len(keywords))

	for i, _ := range keywords {
		keyword := keywords[i]
		newKeywords[i] = findKeyword(owner, keyword)
	}

	return newKeywords
}
