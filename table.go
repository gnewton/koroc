package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
)

type Table struct {
	fields                                  []*Field
	name                                    string
	pk                                      *Field
	fieldMap                                map[string]struct{}
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

func (t *Table) String() string {
	return t.name
}

func (t *Table) SetPrimaryKey(pk *Field) error {
	if pk == nil {
		err := fmt.Errorf("Table [%s]: pk is nil", t.name)
		log.Println(err)
		return err
	}
	if pk.typ != Uint64 {
		err := fmt.Errorf("Table [%s] Field [%s]: Primary key is not uint64; is %s", t.name, pk.name, pk.typ.String())
		log.Println(err)
		return err
	}

	t.pk = pk
	return nil
}

func (t *Table) AddField(f *Field) error {
	if f == nil {
		return errors.New("Field is nil; table is " + t.name)
	}
	if f.name == "" {
		return errors.New("Field is empty; table is " + t.name)
	}
	if t.fieldMap == nil {
		t.fieldMap = make(map[string]struct{})
	}
	if _, ok := t.fieldMap[f.name]; ok {
		return errors.New("Field with that name already exists: " + f.name)
	} else {
		t.fieldMap[f.name] = struct{}{}
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
