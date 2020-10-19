package main

import (
	"testing"
)

// Failing tests
func TestRecord_AddN_IndexTooLargeLimit(t *testing.T) {
	tab, _, _, _, err := personTableFull(new(DialectSqlite3))
	if err != nil {
		t.Fatal(err)
	}

	rec := tab.Record()
	if err := rec.AddN(3, uint32(45)); err == nil {
		t.Fatal("Should fail")
	}
}

func TestRecord_Add_FieldPositionInTableTooLarge(t *testing.T) {
	_, f0, _, _ := personTable(new(DialectSqlite3))
	f0.positionInTable = 999
	rec := new(Record)

	if err := rec.Add(f0, uint32(45)); err == nil {
		t.Fatal("Should fail")
	}
}

func TestRecord_AddN_IndexTooSmallLimit(t *testing.T) {
	tab, _, _, _, err := personTableFull(new(DialectSqlite3))
	if err != nil {
		t.Fatal(err)
	}

	rec := tab.Record()
	if err := rec.AddN(-1, uint32(45)); err == nil {
		t.Fatal("Should fail")
	}
}

func TestRecord_Add_FieldPositionInTableTooSmall(t *testing.T) {
	tab, f0, _, _ := personTable(new(DialectSqlite3))
	tab.AddField(f0)
	f0.positionInTable = -1
	rec := tab.Record()

	if err := rec.Add(f0, uint32(45)); err == nil {
		t.Fatal("Should fail")
	}
}
func TestRecord_AddN_NilRecValues(t *testing.T) {
	rec := new(Record)
	rec.values = nil

	if err := rec.AddN(0, uint32(45)); err == nil {
		t.Fatal("Should fail")
	}
}

func TestRecord_Add_NilRecValues(t *testing.T) {
	_, f0, _, _ := personTable(new(DialectSqlite3))
	rec := new(Record)
	rec.values = nil

	if err := rec.Add(f0, uint32(45)); err == nil {
		t.Fatal("Should fail")
	}
}

func TestRecord_Add_NonUint32Int(t *testing.T) {
	tab, f0, _, _ := personTable(new(DialectSqlite3))

	rec := tab.Record()
	if err := rec.Add(f0, 45); err == nil {
		t.Fatal("Should fail")
	}
}

func TestRecord_AddN_NonUint32Int(t *testing.T) {
	tab, _, _, _, err := personTableFull(new(DialectSqlite3))
	if err != nil {
		t.Fatal(err)
	}

	rec := tab.Record()
	if err := rec.AddN(0, 45); err == nil {
		t.Fatal("Should fail")
	}
}

func TestRecord_Add_NonText(t *testing.T) {
	tab, _, f1, _ := personTable(new(DialectSqlite3))

	rec := tab.Record()
	if err := rec.Add(f1, 45); err == nil {
		t.Fatal("Should fail")
	}
}

func TestRecord_AddN_NonText(t *testing.T) {
	tab, _, _, _, err := personTableFull(new(DialectSqlite3))
	if err != nil {
		t.Fatal(err)
	}

	rec := tab.Record()
	if err := rec.AddN(1, 45); err == nil {
		t.Fatal("Should fail")
	}
}

func TestRecord_AddN_NegIndex(t *testing.T) {
	tab, _, _, _, err := personTableFull(new(DialectSqlite3))
	if err != nil {
		t.Fatal(err)
	}

	rec := tab.Record()
	if err := rec.AddN(-1, 45); err == nil {
		t.Fatal("Should fail")
	}
}

func TestRecord_Add_NilTable(t *testing.T) {
	_, f0, _, _ := personTable(new(DialectSqlite3))
	rec := new(Record)

	if err := rec.Add(f0, uint32(45)); err == nil {
		t.Fatal("Should fail")
	}
}

func TestRecord_AddN_NilTable(t *testing.T) {
	rec := new(Record)

	if err := rec.AddN(0, uint32(45)); err == nil {
		t.Fatal("Should fail")
	}
}

func TestRecord_Initialize_NilTable(t *testing.T) {
	r := Record{
		table: nil,
	}
	if _, err := r.Initialize(); err == nil {
		t.Fatal()
	}
}

func TestRecord_Initialize_TableFieldsNil(t *testing.T) {
	r := Record{
		table: new(Table),
	}
	if _, err := r.Initialize(); err == nil {
		t.Fatal()
	}
}

func TestRecord_Initialize_TableFieldsZeroLen(t *testing.T) {
	r := Record{
		table: new(Table),
	}
	r.table.fields = make([]*Field, 0)
	if _, err := r.Initialize(); err == nil {
		t.Fatal()
	}
}

func TestRecord_Reset_RecordValuesNil(t *testing.T) {
	r := Record{
		table: new(Table),
	}
	r.values = nil
	if err := r.Reset(); err == nil {
		t.Fatal()
	}
}

//////////////////////////////////////////////////////////////////////
// Positive tests
