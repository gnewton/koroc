package main

import (
	"database/sql"
	"errors"
	"log"
)

type Persister struct {
	dialect Dialect
	db      *sql.DB
	tx      *sql.Tx
}

func (p *Persister) CreateTables(tables ...*Table) error {
	if p.tx != nil {
		return errors.New("Existing transaction must be nil")
	}
	var err error
	p.tx, err = p.db.Begin()
	if err != nil {
		log.Println(err)
		return err
	}
	for i, _ := range tables {
		tab := tables[i]
		createSql, err := p.dialect.CreateTableSql(tab.name, tab.fields, tab.pk.name)
		if err != nil {
			return err
		}
		result, err := p.tx.Exec(createSql)
		if err != nil {
			return err
		}
		if result == nil {
			return errors.New("result is nil")
		}

	}
	err = p.tx.Commit()
	if err != nil {
		return err
	}
	p.tx = nil
	return nil
}
func (p *Persister) ExistsByPK(tab *Table, v interface{}) (bool, error) {
	return false, errors.New("TODO")
}

func (p *Persister) SelectOneRecordByPK(tab *Table, v interface{}, rec *Record) error {

	if tab == nil {
		return errors.New("Table is nil")
	}
	if rec == nil {
		return errors.New("Record is nil")
	}
	if tab.selectOneRecordByPKPreparedStatement == nil {
		return errors.New("Table.selectOneRecordByPKPreparedStatement is nil; table:" + tab.name)
	}

	rows, err := tab.selectOneRecordByPKPreparedStatement.Query(v)
	if err != nil {
		return err
	}
	defer rows.Close()
	if !rows.Next() {
		return errors.New("No records returned" + tab.name)
	}

	valuesPointers := make([]interface{}, len(rec.values))
	for i, _ := range rec.values {
		valuesPointers[i] = &rec.values[i]
	}
	if err := rows.Scan(valuesPointers...); err != nil {
		return err
	}
	if rows.Next() {
		return errors.New(">1 records returned" + tab.name)
	}
	return nil

}

func (p *Persister) DeleteByPK(tab *Table, v interface{}) error {
	if tab == nil {
		return errors.New("Table is nil")
	}
	if tab.deleteByPKPreparedStatement == nil {
		return errors.New("Table.deleteByPKPreparedStatement is nil; table:" + tab.name)
	}
	results, err := tab.deleteByPKPreparedStatement.Exec(v)
	if err != nil {
		return err
	}
	if results == nil {
		return err
	}
	rowsAffected, err := results.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected != 1 {
		return err
	}
	return nil
}

func (p *Persister) Insert(rec *Record) error {
	if rec == nil {
		err := errors.New("Record is nil")
		log.Println(err)
		return err
	}

	if rec.table == nil {
		err := errors.New("Record.table is nil")
		log.Println(err)
		return err
	}

	if rec.table.insertPreparedStatement == nil {
		err := errors.New("Prepared statement is nil: table:" + rec.table.name)
		return err
	}
	result, err := rec.table.insertPreparedStatement.Exec(rec.values...)

	if err != nil {
		log.Println(err)
		return err
	}
	if result == nil {
		log.Println(err)
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Println(err)
		return err
	}
	if rowsAffected != 1 {
		log.Println(err)
		return err
	}
	return nil
}

func (p *Persister) JoinTableInsert(joinTable *Table, leftRec, rightRec *Record) error {
	// TODO: leftRec is saved elsewhere; cache rightRec ids; save rightRec; joinTable should contain fields used for cache key
	if leftRec.table != joinTable.joinTableInfo.leftTable {
		return errors.New("Record left table does not match join table left table")
	}

	if rightRec.table != joinTable.joinTableInfo.rightTable {
		return errors.New("Record right table does not match join table right table")
	}

	jrec, err := joinTable.Record()
	if err != nil {
		return err
	}
	// left table id value
	jrec.AddN(0, leftRec.values[leftRec.table.pk.positionInTable])
	// left table id value
	jrec.AddN(1, rightRec.values[rightRec.table.pk.positionInTable])

	if err := p.Insert(jrec); err != nil {
		return err
	}

	return nil

}

func (p *Persister) TxCommit(dialect Dialect, tables ...*Table) error {
	var err error
	for i, _ := range tables {
		tab := tables[i]
		err = closePreparedStatements(tab.insertPreparedStatement, tab.deleteByPKPreparedStatement, tab.selectOneRecordByPKPreparedStatement)

		if err != nil {
			return err
		}
		tab.insertPreparedStatement = nil
		tab.deleteByPKPreparedStatement = nil
		tab.selectOneRecordByPKPreparedStatement = nil
	}

	err = p.tx.Commit()
	if err != nil {
		return err
	}
	p.tx = nil
	return nil
}

func (p *Persister) TxBegin(dialect Dialect, tables ...*Table) error {
	var err error

	p.tx, err = p.db.Begin()
	if err != nil {
		return err
	}

	for i, _ := range tables {
		tab := tables[i]
		err = makeNewPreparedStatements(dialect, tab, p.tx)
		if err != nil {
			return err
		}
	}

	return err
}

func closePreparedStatements(stmts ...*sql.Stmt) error {
	for i, _ := range stmts {
		stmt := stmts[i]
		if stmt != nil {
			if err := stmt.Close(); err != nil {
				return err
			}
		}
	}
	return nil

}

func makeNewPreparedStatements(dialect Dialect, tab *Table, tx *sql.Tx) error {
	var err error

	// INSERT
	if tab.insertPreparedStatementSql == "" {
		tab.insertPreparedStatementSql, err = dialect.InsertPreparedStatementSql(tab.name, tab.fields)
		if err != nil {
			return err
		}
	}
	tab.insertPreparedStatement, err = tx.Prepare(tab.insertPreparedStatementSql)
	if err != nil {
		return err
	}

	// DELETE BY PK
	if tab.deleteByPKPreparedStatementSql == "" {
		tab.deleteByPKPreparedStatementSql, err = dialect.DeleteByPKPreparedStatementSql(tab.name, tab.pk.name)
		if err != nil {
			return err
		}
	}
	tab.deleteByPKPreparedStatement, err = tx.Prepare(tab.deleteByPKPreparedStatementSql)
	if err != nil {
		return err
	}

	// SELECT BY PK
	if tab.selectOneRecordByPKPreparedStatementSql == "" {
		tab.selectOneRecordByPKPreparedStatementSql, err = dialect.SelectOneRecordByPKPreparedStatementSql(tab.name, tab.fields, tab.pk.name)

		if err != nil {
			return err
		}
	}
	tab.selectOneRecordByPKPreparedStatement, err = tx.Prepare(tab.selectOneRecordByPKPreparedStatementSql)
	if err != nil {
		return err
	}

	return nil
}
