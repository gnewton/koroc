package main

import (
	"database/sql"
	"errors"
	"fmt"
	//"github.com/gnewton/
	"log"
	"strconv"
)

type JTFieldsCreateSql func() string
type InsertValuesSql func(interface{}) (names string, fields string, values string, err error)
type SavePreparedSql func() (fields string, numFields string)

const MajorTopicField = "major_topic"

//SavePreparedSql
func kwtSavePreparedSql() (fields string, numFields string) {
	return ", " + MajorTopicField, ",?"
}

// JTFieldsCreateSql
func kwjtCreateSql() string {
	return MajorTopicField + " boolean"
}

// InsertValuesSql
func kwjtInsert(i interface{}) (string, string, string, error) {
	keyword, ok := i.(*Keyword)
	if !ok {
		log.Println("Unable to type assert Keyword")
		return "", "", "", errors.New("Unable to type assert Keyword")
	}
	fields := ", " + MajorTopicField
	var values string
	if keyword.MajorTopic {
		values = ", true"
	} else {
		values = ", false"
	}
	return keyword.Name, fields, values, nil

}

func NewJoinTable(joinTableName string, leftJoinField string, rightJoinField string, fcreate JTFieldsCreateSql, finsert InsertValuesSql, fSave SavePreparedSql) (*JoinTable, error) {
	jt := new(JoinTable)
	jt.ids = make(map[string]uint32)
	jt.leftJoinField = leftJoinField
	jt.rightJoinField = rightJoinField
	jt.joinTableName = joinTableName
	jt.fCreate = fcreate
	jt.fInsert = finsert
	jt.fSave = fSave
	return jt, nil
}

type JoinTable struct {
	ids                           map[string]uint32
	counter                       uint32
	leftJoinField, rightJoinField string
	joinTableName                 string
	fCreate                       JTFieldsCreateSql
	fInsert                       InsertValuesSql
	fSave                         SavePreparedSql
}

func (jt *JoinTable) Delete(leftJoinField uint32, stmt *sql.Stmt) error {
	if stmt == nil {
		return errors.New("Statement is nil")
	}
	_, err := stmt.Exec(leftJoinField)
	return err
}

func (jt *JoinTable) SavePreparedStatement(tx *sql.Tx) (*sql.Stmt, error) {
	var fields, valuePlaceHolders string
	if jt.fSave != nil {
		fields, valuePlaceHolders = jt.fSave()
	}

	sql := "INSERT INTO " + jt.joinTableName + "(" + jt.leftJoinField + "," + jt.rightJoinField + fields + ") VALUES (?,?" + valuePlaceHolders + ")"
	return tx.Prepare(sql)
}

func (jt *JoinTable) DeletePreparedStatement(tx *sql.Tx) (*sql.Stmt, error) {
	if tx == nil {
		log.Println(errors.New("Database is nil"))
	}
	sql := fmt.Sprintln("DELETE FROM", jt.joinTableName, "WHERE", jt.leftJoinField, "=?")
	return tx.Prepare(sql)
}

func (jt *JoinTable) CreateSql() (string, string) {
	fc := ""
	if jt.fCreate != nil {
		fc = "," + jt.fCreate()
	}
	//return "CREATE TABLE " + jt.joinTableName + " (" + jt.leftJoinField + " integer, " + jt.rightJoinField + " integer, " + fc + " PRIMARY KEY (" + jt.leftJoinField + "," + jt.rightJoinField + "))"
	return "CREATE TABLE " + jt.joinTableName + " (" + jt.leftJoinField + " integer(4), " + jt.rightJoinField + " integer(4)" + fc + ")", "CREATE UNIQUE INDEX idx_" + jt.joinTableName + " ON " + jt.joinTableName + "(" + jt.leftJoinField + ", " + jt.rightJoinField + ")"

}

// Returns the id of the Joined item
// If new, returns true indicating the Keyword item should be saved, with the new ID (which is returned)
// Always, this means a
func (jt *JoinTable) AddJoinItem(leftId uint32, i interface{}) (rightId uint32, newItem bool, joinSql string, err error) {
	if i == nil || leftId < 0 || jt.ids == nil || len(jt.joinTableName) == 0 || len(jt.leftJoinField) == 0 || len(jt.rightJoinField) == 0 {
		return 0, false, "", errors.New("Bad input data")
	} else {
		// make key
		key, fields, values := "", "", ""
		if jt.fInsert != nil {
			key, fields, values, err = jt.fInsert(i)
			if err != nil {
				log.Printf("%T\n", i)
				log.Printf("%v\n", i)
				log.Fatal(err)
			}
		}

		var ok bool
		// Item already exists?
		if rightId, ok = jt.ids[key]; ok {
			// Already exists
			newItem = false
		} else {
			newItem = true
			// New join
			rightId = jt.counter
			jt.ids[key] = rightId
			jt.counter++
		}

		// insert sql
		joinSql := "INSERT INTO " + jt.joinTableName + " (" + jt.leftJoinField + "," + jt.rightJoinField + fields + ") VALUES (" + strconv.FormatUint(uint64(leftId), 10) + "," + strconv.FormatUint(uint64(rightId), 10) + values + ")"
		//log.Println("Inserting join: " + joinSql)
		return rightId, newItem, joinSql, nil

	}
}
