package main

import (
	"database/sql"
	"errors"
	_ "github.com/mattn/go-sqlite3"
	"testing"
)

//////////////////////////////////////////////////////////////////////
// Failing tests
func TestTableFoo(t *testing.T) {
	tab, f0, f1, _ := articleTable(new(DialectSqlite3))

	if err := tab.AddField(f0); err != nil {
		t.Fatal(err)
	}
	if err := tab.AddField(f1); err != nil {
		t.Fatal(err)
	}
	tab.pk = f0

	rec := tab.Record()

	var i uint32 = 42
	if err := rec.Add(f0, i); err != nil {
		t.Fatal(err)
	}

	if err := rec.Add(f1, "fred"); err != nil {
		t.Fatal(err)
	}

	t.Log(rec.values)
	rec.Reset()
	t.Log(rec.values)

	if err := rec.AddN(0, uint32(997)); err != nil {
		t.Fatal(err)
	}
	if err := rec.AddN(1, "bill"); err != nil {
		t.Fatal(err)
	}
	t.Log(rec.values)

	rec.Reset()
	if err := rec.Add(f1, "foo"); err != nil {
		t.Fatal(err)
	}

	// if err != nil {
	// 	t.Fatal(err)
	// }

}

func TestTable_AddField_NullField(t *testing.T) {
	tab, _, _, _ := articleTable(new(DialectSqlite3))

	if err := tab.AddField(nil); err == nil {
		t.Fatal("Should fail")
	}
}

func TestTable_AddField_EmptyFieldName(t *testing.T) {
	tab, f0, _, _ := articleTable(new(DialectSqlite3))
	f0.name = ""
	if err := tab.AddField(f0); err == nil {
		t.Fatal("Should fail")
	}
}

func TestTable_AddField_RepeatFieldName(t *testing.T) {
	tab, f0, _, _ := articleTable(new(DialectSqlite3))
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

func TestTable_InsertRecord(t *testing.T) {
	db, tab, f0, f1, f2 := _CreateTable(t)

	preparedSql, err := tab.InsertPreparedStatement()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(preparedSql)

	tx, err := db.Begin()
	if err != nil {
		t.Fatal(err)
	}
	stmt, err := tx.Prepare(preparedSql)
	if err != nil {
		t.Fatal(err)
	}

	rec := tab.Record()
	if err := rec.Add(f0, uint32(10)); err != nil {
		t.Fatal(err)
	}

	if err := rec.Add(f1, "foo"); err != nil {
		t.Fatal(err)
	}

	if err := rec.Add(f2, true); err != nil {
		t.Fatal(err)
	}

	result, err := stmt.Exec(rec.values...)
	if err != nil {
		t.Fatal(err)
	}
	if result == nil {
		t.Fatal(err)
	}
	rowsAffected, err := result.RowsAffected()
	if rowsAffected != 1 {
		t.Fatal(errors.New("Should only effect one row"))
	}
}

func TestTable_CreateTable(t *testing.T) {
	_, _, _, _, _ = _CreateTable(t)
}

func _CreateTable(t *testing.T) (*sql.DB, *Table, *Field, *Field, *Field) {
	tab, f0, f1, f2, err := articleTableFull(new(DialectSqlite3))
	if err != nil {
		t.Fatal(err)
	}

	createTableSql, err := tab.CreateSql()
	if err != nil {
		t.Fatal(err)
	}

	db, err := newDB()
	if err != nil {
		t.Fatal(err)
	}

	t.Log(createTableSql)
	_, err = db.Exec(createTableSql)
	if err != nil {
		t.Fatal(err)
	}
	return db, tab, f0, f1, f2
}

//////////////////////////////////////////////////////////////////////
//helpers
func newDB() (*sql.DB, error) {
	return sql.Open("sqlite3", ":memory:")
}

func articleTableFull(dialect Dialect) (*Table, *Field, *Field, *Field, error) {
	tab, f0, f1, f2 := articleTable(dialect)
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

func articleTable(dialect Dialect) (*Table, *Field, *Field, *Field) {

	tab := Table{name: "articles",
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
