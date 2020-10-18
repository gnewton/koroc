package main

import (
	"testing"
)

func TestCreateSql(t *testing.T) {
	_, err := NewJoinTable("article_keyword", "article_id", "keyword_id", kwjtCreateSql, kwjtInsert, kwtSavePreparedSql)
	if err != nil {
		t.Fatal(err)
	}
}

func TestInsertSql(t *testing.T) {
	jt, err := NewJoinTable("article_keyword", "article_id", "keyword_id", kwjtCreateSql, kwjtInsert, kwtSavePreparedSql)
	if err != nil {
		t.Fatal(err)
	}
	kw := Keyword{
		MajorTopic: true,
		Name:       "blood"}

	rightId, newItem, _, err := jt.AddJoinItem(42, &kw)

	if err != nil {
		t.Fatal(err)
	}
	if rightId != 0 {
		t.Fatal()
	}
	if !newItem {
		t.Fatal()
	}

	kw = Keyword{
		MajorTopic: false,
		Name:       "liver"}

	rightId, newItem, _, err = jt.AddJoinItem(934, &kw)
	if err != nil {
		t.Fatal(err)
	}
	if !newItem {
		t.Fatal()
	}

	kw = Keyword{
		MajorTopic: true,
		Name:       "heart"}
	rightId, newItem, _, err = jt.AddJoinItem(94, &kw)
	if err != nil {
		t.Fatal(err)
	}

	rightId, newItem, _, err = jt.AddJoinItem(42, &kw)
	if err != nil {
		t.Fatal(err)
	}

	rightId, newItem, _, err = jt.AddJoinItem(48, &kw)
	if err != nil {
		t.Fatal(err)
	}
}
