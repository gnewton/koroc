package main

import (
	"testing"
)

//////////////////////////////////////////////////////////////////////
// Failing tests

func TestField_CheckValue_WantUint32GotText(t *testing.T) {
	_, f0, _, _ := articleTable(new(DialectSqlite3))

	if err := f0.CheckValueType("foo"); err == nil {
		t.Fatal(err)
	}
}

func TestField_CheckValue_WantUint32GotBool(t *testing.T) {
	_, f0, _, _ := articleTable(new(DialectSqlite3))

	if err := f0.CheckValueType(true); err == nil {
		t.Fatal(err)
	}
}

func TestField_CheckValue_WantUint32GotInt(t *testing.T) {
	_, f0, _, _ := articleTable(new(DialectSqlite3))

	if err := f0.CheckValueType(32); err == nil {
		t.Fatal(err)
	}
}

func TestField_CheckValue_WantTextGotInt(t *testing.T) {
	_, _, f1, _ := articleTable(new(DialectSqlite3))

	if err := f1.CheckValueType(32); err == nil {
		t.Fatal(err)
	}
}

func TestField_CheckValue_NilValue(t *testing.T) {
	_, f0, _, _ := articleTable(new(DialectSqlite3))

	if err := f0.CheckValueType(nil); err == nil {
		t.Fatal(err)
	}
}

//////////////////////////////////////////////////////////////////////
// Positive tests
func TestField_CheckValue_WantUint32GotUint(t *testing.T) {
	_, f0, _, _ := articleTable(new(DialectSqlite3))

	if err := f0.CheckValueType(uint32(32)); err != nil {
		t.Fatal(err)
	}
}

func TestField_CheckValue_WantTextGotText(t *testing.T) {
	_, _, f1, _ := articleTable(new(DialectSqlite3))

	if err := f1.CheckValueType("foo"); err != nil {
		t.Fatal(err)
	}
}

func TestField_CheckValue_WantBoolGotBool(t *testing.T) {
	_, _, _, f2 := articleTable(new(DialectSqlite3))

	if err := f2.CheckValueType(true); err != nil {
		t.Fatal(err)
	}
}
