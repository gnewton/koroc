package main

import (
	"errors"
	//"fmt"
	//"log"
)

type Table struct {
	fields       []*Field
	name         string
	pk           *Field
	fmap         map[string]struct{}
	fieldCounter int
	dialect      Dialect
}

func (t *Table) AddField(f *Field) error {
	if f == nil {
		return errors.New("Field is nil; table is " + t.name)
	}
	if f.name == "" {
		return errors.New("Field is empty; table is " + t.name)
	}
	if t.fmap == nil {
		t.fmap = make(map[string]struct{})
	}
	if _, ok := t.fmap[f.name]; ok {
		return errors.New("Field with that name already exists: " + f.name)
	} else {
		t.fmap[f.name] = struct{}{}
	}

	t.fields = append(t.fields, f)
	f.positionInTable = t.fieldCounter
	t.fieldCounter++
	return nil
}

func (t *Table) Record() *Record {
	rec := Record{
		table: t,
	}
	return &rec
}

/*
func (t *Table) CreateSql() (string, error) {
	return t.dialect.CreateTableSql(t.name, t.fields, t.pk)
}

func (t *Table) DeletePreparedStatement() (string, error) {
	return t.dialect.DeleteByPKPreparedStatementSql(t.name, t.pk)
}

func (t *Table) InsertPreparedStatement() (string, error) {
	return t.dialect.InsertPreparedStatementSql(t.name, t.fields)
}
*/
