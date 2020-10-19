package main

import (
	"errors"
	"strconv"
)

type DialectSqlite3 struct {
}

//func (d *DialectSqlite3) PreparedValuePlaceHolder(counter int) (string, error) {
//	return "?", nil
//}

func (d *DialectSqlite3) DeleteByPKPreparedStatementSql(table string, pk string) (string, error) {
	if table == "" {
		return "", errors.New("table is empty")
	}
	if pk == "" {
		return "", errors.New("pk is empty")
	}
	return "DELETE FROM " + table + " where " + pk + "=?", nil
}

func (d *DialectSqlite3) InsertPreparedStatementSql(table string, fields []*Field) (string, error) {
	if table == "" {
		return "", errors.New("table is empty")
	}
	if len(fields) == 0 {
		return "", errors.New("fields zero length")
	}

	sql := "INSERT INTO " + table + "("

	var holders string
	for i, _ := range fields {
		if fields[i] == nil {
			return "", errors.New("field is nil: " + strconv.Itoa(i))
		}
		if i != 0 {
			holders += ","
			sql += ","
		}
		holders += "?"
		sql += fields[i].name
	}

	sql += ") VALUES (" + holders + ")"

	return sql, nil
}

func (d *DialectSqlite3) CreateTableSql(table string, fields []*Field, pk string) (string, error) {
	if table == "" {
		return "", errors.New("table is empty")
	}
	if len(fields) == 0 {
		return "", errors.New("fields zero length")
	}
	if pk == "" {
		return "", errors.New("pk is empty")
	}

	sql := "CREATE TABLE " + table + " ("

	fieldsSql, err := d.fieldCreates(fields, pk)
	if err != nil {
		return "", err
	}

	//sql += fieldsSql + ", PRIMARY KEY " + pk.name + ")"
	sql += fieldsSql + ")"

	return sql, nil
}

/////////////////
func (t *DialectSqlite3) createSql(f *Field, pk string) (string, error) {
	typ, err := f.makeSqlType()
	if err != nil {
		return "", err
	}
	sql := f.name + " " + typ
	if f.name == pk {
		sql += " PRIMARY KEY"
	}
	return sql, nil
}

func (d *DialectSqlite3) fieldCreates(fields []*Field, pk string) (string, error) {
	if fields == nil {
		return "", errors.New("fields is nil")
	}
	if len(fields) == 0 {
		return "", errors.New("fields zero length")
	}
	var s string
	for i, _ := range fields {
		if i != 0 {
			s += ", "
		}
		fs, err := d.createSql(fields[i], pk)
		if err != nil {
			return "", err
		}
		s += fs
	}
	return s, nil
}

func (f *Field) makeSqlType() (string, error) {
	switch f.typ {
	case Text:
		return "TEXT", nil
	case Uint32:
		//return "INT4", nil
		return "INTEGER", nil
	case Uint64:
		return "BIGINT", nil
	case Boolean:
		return "BOOLEAN", nil
	}
	return "", errors.New("Unknown type")

}
