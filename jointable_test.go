package main

import (
	"github.com/gnewton/pubmedSqlStructs"
	"testing"
)

func TestCreateSql(t *testing.T) {
	jt, err := NewJoinTable("article_keyword", "article_id", "keyword_id", kwjtCreateSql, kwjtInsert)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(jt.CreateSql())
}

func TestInsertSql(t *testing.T) {
	jt, err := NewJoinTable("article_keyword", "article_id", "keyword_id", kwjtCreateSql, kwjtInsert)
	if err != nil {
		t.Fatal(err)
	}
	kw := pubmedSqlStructs.Keyword{
		MajorTopic: true,
		Name:       "blood"}

	rightId, newItem, sql, err := jt.AddJoinItem(42, &kw)
	t.Log(newItem)
	if err != nil {
		t.Fatal(err)
	}
	if rightId != 0 {
		t.Fatal()
	}
	if !newItem {
		t.Fatal()
	}
	t.Log(sql)

	kw = pubmedSqlStructs.Keyword{
		MajorTopic: false,
		Name:       "liver"}

	rightId, newItem, sql, err = jt.AddJoinItem(934, &kw)
	t.Log(newItem)
	if err != nil {
		t.Fatal(err)
	}
	if !newItem {
		t.Fatal()
	}
	t.Log(sql)

	kw = pubmedSqlStructs.Keyword{
		MajorTopic: true,
		Name:       "heart"}
	rightId, newItem, sql, err = jt.AddJoinItem(94, &kw)
	t.Log(newItem)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(sql)

	rightId, newItem, sql, err = jt.AddJoinItem(42, &kw)
	t.Log(newItem)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(sql)

	rightId, newItem, sql, err = jt.AddJoinItem(48, &kw)
	t.Log(newItem)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(sql)
}
