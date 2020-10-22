package main

import (
	"database/sql"
	"errors"
	_ "github.com/mattn/go-sqlite3"
	"testing"
)

//////////////////////////////////////////////////////////////////////
// Failing tests

func TestTable_AddField_NullField(t *testing.T) {
	tab, _, _, _ := personTable(new(DialectSqlite3))

	if err := tab.AddField(nil); err == nil {
		t.Fatal("Should fail")
	}
}

func TestTable_AddField_EmptyFieldName(t *testing.T) {
	tab, f0, _, _ := personTable(new(DialectSqlite3))
	f0.name = ""
	if err := tab.AddField(f0); err == nil {
		t.Fatal("Should fail")
	}
}

func TestTable_AddField_RepeatFieldName(t *testing.T) {
	tab, f0, _, _ := personTable(new(DialectSqlite3))
	t.Log(tab.fmap)
	if err := tab.AddField(f0); err != nil {
		t.Fatal("Should not fail")
	}
	t.Log(tab.fmap)
	if err := tab.AddField(f0); err == nil {
		t.Log(tab.fmap)
		t.Fatal("Should fail")
	}
}

//////////////////////////////////////////////////////////////////////
// Positive tests
func TestTable_CreateTable(t *testing.T) {
	db, err := newDB()
	if err != nil {
		t.Fatal(err)
	}

	_, _, _, _, err = _CreateTable(t, db)
	if err != nil {
		t.Fatal(err)
	}
}

func TestTable_InsertRecord(t *testing.T) {
	db, err := newDB()
	if err != nil {
		t.Fatal(err)
	}

	var id uint32 = 10
	_, _, err = _InsertRecord(t, db, id, "foo", true)
	if err != nil {
		t.Fatal(err)
	}
}

func TestTable_DeleteRecord(t *testing.T) {
	db, err := newDB()
	if err != nil {
		t.Fatal(err)
	}

	var id uint32 = 10
	db, tab, err := _InsertRecord(t, db, id, "foo", true)
	if err != nil {
		t.Fatal(err)
	}

	err = _DeleteRecord(t, id, db, tab)
	if err != nil {
		t.Fatal(err)
	}
}

//////////////////////////////////////////////////////////////////////
//helpers

func _DeleteRecord(t *testing.T, v0 uint32, db *sql.DB, tab *Table) error {

	preparedDeleteSql, err := tab.dialect.DeleteByPKPreparedStatementSql(tab.name, tab.pk.name)
	if err != nil {
		t.Log(err)
		return err
	}
	tx, err := db.Begin()
	if err != nil {
		t.Log(err)
		return err
	}
	stmt, err := tx.Prepare(preparedDeleteSql)
	t.Log(preparedDeleteSql)
	if err != nil {
		t.Log(err)
		return err
	}
	result, err := stmt.Exec(v0)
	if err != nil {
		t.Log(err)
		return err
	}
	if result == nil {
		t.Log(err)
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		t.Log(err)
		return err
	}
	if rowsAffected != 1 {
		return errors.New("Should only effect one row")
	}
	return nil
}

func _InsertRecord(t *testing.T, db *sql.DB, v0 uint32, v1 string, v2 bool) (*sql.DB, *Table, error) {

	tab, f0, f1, f2, err := _CreateTable(t, db)

	if err != nil {
		t.Log(err)
		return nil, nil, err
	}
	preparedInsertSql, err := tab.dialect.InsertPreparedStatementSql(tab.name, tab.fields)
	t.Log(preparedInsertSql)
	if err != nil {
		t.Log(err)
		return nil, nil, err
	}

	tx, err := db.Begin()
	if err != nil {
		t.Log(err)
		return nil, nil, err
	}
	stmt, err := tx.Prepare(preparedInsertSql)
	if err != nil {
		t.Log(err)
		return nil, nil, err
	}

	rec, err := tab.Record()
	if err != nil {
		t.Fatal(err)
	}
	if err := rec.Add(f0, v0); err != nil {
		t.Log(err)
		return nil, nil, err
	}

	if err := rec.Add(f1, v1); err != nil {
		t.Log(err)
		return nil, nil, err
	}

	if err := rec.Add(f2, v2); err != nil {
		t.Log(err)
		return nil, nil, err
	}

	result, err := stmt.Exec(rec.values...)
	if err != nil {
		t.Log(err)
		return nil, nil, err
	}
	if result == nil {
		t.Log(err)
		return nil, nil, err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		t.Log(err)
		return nil, nil, err
	}
	if rowsAffected != 1 {
		t.Log(err)
		return nil, nil, errors.New("Should only effect one row")
	}
	err = tx.Commit()
	if err != nil {
		return nil, nil, err
	}
	return db, tab, nil
}

func _CreateTable(t *testing.T, db *sql.DB) (*Table, *Field, *Field, *Field, error) {
	tab, f0, f1, f2, err := personTableFull(new(DialectSqlite3)) // TODO: Dialect should be passed in....

	if err != nil {
		return nil, nil, nil, nil, err
	}

	createTableSql, err := tab.dialect.CreateTableSql(tab.name, tab.fields, tab.pk.name)
	if err != nil {
		return nil, nil, nil, nil, err

	}
	t.Log(createTableSql)
	tx, err := db.Begin()
	if err != nil {
		return nil, nil, nil, nil, err

	}

	result, err := tx.Exec(createTableSql)
	if result == nil {
		return nil, nil, nil, nil, err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, nil, nil, nil, err
	}

	if rowsAffected != 0 {
		return nil, nil, nil, nil, errors.New("More than zero row affected")
	}

	err = tx.Commit()
	if err != nil {
		return nil, nil, nil, nil, err
	}

	return tab, f0, f1, f2, nil
}

func newDB() (*sql.DB, error) {
	return sql.Open("sqlite3", ":memory:")
}

func personTableFull(dialect Dialect) (*Table, *Field, *Field, *Field, error) {
	tab, f0, f1, f2 := personTable(dialect)
	if err := tab.AddField(f0); err != nil {
		return nil, nil, nil, nil, err
	}

	if err := tab.AddField(f1); err != nil {
		return nil, nil, nil, nil, err
	}

	if err := tab.AddField(f2); err != nil {
		return nil, nil, nil, nil, err
	}
	return tab, f0, f1, f2, nil
}

func personTable(dialect Dialect) (*Table, *Field, *Field, *Field) {

	tab := Table{name: "person",
		dialect: dialect,
	}

	f0 := Field{
		name: "id",
		typ:  Uint32,
	}

	f1 := Field{
		name: "first_name",
		typ:  Text,
	}

	f2 := Field{
		name: "has_car",
		typ:  Boolean,
	}
	tab.pk = &f0
	return &tab, &f0, &f1, &f2
}
