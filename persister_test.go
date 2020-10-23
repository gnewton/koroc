package main

import (
	"errors"
	"testing"
)

//////////////////////////////////////////////////////////////////////
// Failing tests
func TestPersist_InsertJoin_JTNil(t *testing.T) {
	//t.Fatal("TODO")
}

func TestPersist_InsertJoin_JTInfoNil(t *testing.T) {
	//	t.Fatal("TODO")
}

func TestPersist_InsertJoin_LeftRecNil(t *testing.T) {
	//	t.Fatal("TODO")
}

func TestPersist_InsertJoin_RightRecNil(t *testing.T) {
	//	t.Fatal("TODO")
}

func TestPersist_InsertJoin_LeftRecNEJTInfoLeftRec(t *testing.T) {
	//	t.Fatal("TODO")
}

func TestPersist_InsertJoin_RightRecNEJTInfoRightRec(t *testing.T) {
	//	t.Fatal("TODO")
}

//////////////////////////////////////////////////////////////////////
// Positive tests
func TestPersist_positive(t *testing.T) {
	//t.Fatal("TODO")
}

func TestPersist_Insert(t *testing.T) {
	dialect, err := NewDialectSqlite3()
	if err != nil {
		t.Fatal(err)
	}
	db, err := dialect.OpenDB(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	p := Persister{
		dialect: dialect,
		db:      db,
	}
	personTable, _, _, _, err := personTableFull(dialect)
	if err != nil {
		t.Fatal(err)
	}

	err = p.CreateTables(personTable)
	if err != nil {
		t.Fatal(err)
	}
	p.tx, err = p.db.Begin()
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		err := p.tx.Commit()
		if err != nil {
			t.Fatal(err)
		}
	}()
	err = makeNewPreparedStatements(dialect, personTable, p.tx)
	if err != nil {
		t.Fatal(err)
	}
	rec, err := personTable.Record()
	if err != nil {
		t.Fatal(err)
	}
	if err := rec.AddN(0, uint32(42)); err != nil {
		t.Fatal(err)
	}
	if err := rec.AddN(1, "Bill"); err != nil {
		t.Fatal(err)
	}
	if err := rec.AddN(2, true); err != nil {
		t.Fatal(err)
	}

	// rec := Record{
	// 	table:  personTable,
	// 	values: []*interface{}{uint32(42), "Bill", true},
	// }

	if err := rec.Initialize(false); err != nil {
		t.Fatal(err)
	}

	err = p.Insert(rec)
	if err != nil {
		t.Fatal(err)
	}
	if err = rec.Reset(); err != nil {
		t.Fatal(err)
	}
	if err := rec.AddN(0, uint32(49)); err != nil {
		t.Fatal(err)
	}
	if err := rec.AddN(1, "Harry"); err != nil {
		t.Fatal(err)
	}
	if err := rec.AddN(2, false); err != nil {
		t.Fatal(err)
	}
	err = p.Insert(rec)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(personTable.selectOneRecordByPKPreparedStatementSql)
	// Added: 42; Select: 32: should fail
	newRec, err := personTable.Record()
	if err != nil {
		t.Fatal(err)
	}
	err = p.SelectOneRecordByPK(personTable, uint32(32), newRec)
	if err == nil {
		t.Fatal(err)
	}

	// Should succeed
	t.Log(newRec)
	newRec, err = personTable.Record()
	if err != nil {
		t.Fatal(err)
	}
	err = p.SelectOneRecordByPK(personTable, uint32(42), newRec)
	if err != nil {
		t.Fatal(err)
	}

	if newRec == nil {
		t.Fatal(errors.New("selected record is nil"))
	}
	t.Log(newRec)
	newRec, err = personTable.Record()
	if err != nil {
		t.Fatal(err)
	}
	err = p.SelectOneRecordByPK(personTable, uint32(49), newRec)
	if err != nil {
		t.Fatal(err)
	}

	if newRec == nil {
		t.Fatal(errors.New("selected record is nil"))
	}
	t.Log(newRec)
}
