package main

import (
	"database/sql"
	"errors"
	//"fmt"
	//"log"
)

type JoinTableInfo struct {
	leftTable, rightTable      *Table
	rightTableIDCache          map[string]uint64
	rightTableIDCacheKeyFields []*Field
}

type Table struct {
	fields                                  []*Field
	name                                    string
	pk                                      *Field
	fmap                                    map[string]struct{}
	fieldCounter                            int
	dialect                                 Dialect
	insertPreparedStatement                 *sql.Stmt
	insertPreparedStatementSql              string
	deleteByPKPreparedStatement             *sql.Stmt
	deleteByPKPreparedStatementSql          string
	selectOneRecordByPKPreparedStatement    *sql.Stmt
	selectOneRecordByPKPreparedStatementSql string
	// if this is a JoinTable
	joinTableInfo *JoinTableInfo
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

func (t *Table) Record() (*Record, error) {
	rec := Record{
		table: t,
	}
	err := rec.Initialize(true)

	return &rec, err
}
