package main

import (
	"errors"
	"github.com/gnewton/pubmedSqlStructs"
	"log"
	"strconv"
)

type JTFieldsCreateSql func() string
type InsertValuesSql func(interface{}) (string, string, string)

const MajorTopicField = "major_topic"

func kwjtCreateSql() string {
	return MajorTopicField + " boolean,"
}

func kwjtInsert(i interface{}) (string, string, string) {
	keyword, ok := i.(*pubmedSqlStructs.Keyword)
	if !ok {
		log.Fatal("Unable to type assert Keyword")
	}
	fields := ", " + MajorTopicField
	var values string
	if keyword.MajorTopic {
		values = ", true"
	} else {
		values = ", false"
	}
	return keyword.Name, fields, values

}

func NewJoinTable(joinTableName string, leftJoinField string, rightJoinField string, fcreate JTFieldsCreateSql, finsert InsertValuesSql) (*JoinTable, error) {
	jt := new(JoinTable)
	jt.ids = make(map[string]uint32)
	jt.leftJoinField = leftJoinField
	jt.rightJoinField = rightJoinField
	jt.joinTableName = joinTableName
	jt.fcreate = fcreate
	jt.finsert = finsert
	return jt, nil
}

type JoinTable struct {
	ids                           map[string]uint32
	counter                       uint32
	leftJoinField, rightJoinField string
	joinTableName                 string
	fcreate                       JTFieldsCreateSql
	finsert                       InsertValuesSql
}

func (jt *JoinTable) CreateSql() string {
	fc := ""
	if jt.fcreate != nil {
		fc = jt.fcreate()
	}
	return "CREATE TABLE " + jt.joinTableName + " (" + jt.leftJoinField + " integer, " + jt.rightJoinField + " integer, " + fc + " PRIMARY KEY (" + jt.leftJoinField + "," + jt.rightJoinField + "))"
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
		if jt.finsert != nil {
			key, fields, values = jt.finsert(i)
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
		log.Println(jt.ids)
		// insert sql
		joinSql := "INSERT INTO " + jt.joinTableName + " (" + jt.leftJoinField + "," + jt.rightJoinField + fields + ") VALUES (" + strconv.FormatUint(uint64(leftId), 10) + "," + strconv.FormatUint(uint64(rightId), 10) + values + ")"
		return rightId, newItem, joinSql, nil

	}
}
