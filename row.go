package main

import (
	"errors"
	"fmt"
	"log"
)

type Row struct {
	table                  *Table
	ints                   map[string]int64
	strings                map[string]string
	blobs                  map[string][]byte
	fieldHasValue          map[string]struct{}
	relationRows           []*Row
	thisOneToManyRelation  *OneToMany
	thisManyToManyRelation *ManyToMany
	isInited               bool
}

func (row *Row) AddRelationRow(relRow *Row) error {
	if relRow == nil {
		return errors.New("Relation row is nil")
	}
	if !row.isInited {
		return errors.New("Row is not inited. Need to run row.Init()")
	}
	row.relationRows = append(row.relationRows, relRow)

	// loop through this row's table's relations
	for i, _ := range row.table.oneToMany {

		rel := row.table.oneToMany[i]
		fmt.Println("-------------- one2N", *rel)
		if rel.rightTable == relRow.table {
			relRow.thisOneToManyRelation = rel
		}
	}

	for i, _ := range row.table.manyToMany {
		rel := row.table.manyToMany[i]
		fmt.Println("--------------", *rel)
		if rel.rightTable == relRow.table {
			relRow.thisManyToManyRelation = rel
		}
	}

	return nil
}

func (row *Row) Init(t *Table) {
	row.table = t

	row.ints = make(map[string]int64, 0)
	row.strings = make(map[string]string, 0)
	row.blobs = make(map[string][]byte, 0)
	row.fieldHasValue = make(map[string]struct{}, 0)

	row.isInited = true
}

func (row *Row) setFieldHasValue(fieldName string) {
	row.fieldHasValue[fieldName] = struct{}{}
}

func (row *Row) hasValue(fieldName string) bool {
	_, hasValue := row.fieldHasValue[fieldName]
	return hasValue
}

func (r *Row) Print() {
	fmt.Println("Row: Table:", r.table.name)
	fmt.Println("Row: fields:")
	for k, v := range r.ints {
		fmt.Println("\tints", k, "=", v)
	}
	for k, v := range r.strings {
		fmt.Println("\tstring", k, "=", v)
	}

	for i, _ := range r.relationRows {
		fmt.Println("\tRelation: table=", r.relationRows[i].table.name, " relation:")
		if r.relationRows[i].thisOneToManyRelation != nil {
			fmt.Println("\t", r.relationRows[i].thisOneToManyRelation.leftKeyField.GetName())
		}

	}
}

func (r *Row) SetInt64(f string, v int64) error {
	if !r.isInited {
		return errors.New("Row is not inited. Need to run row.Init()")
	}
	log.Println("---", f, v)
	log.Println("--- r.table.fieldsMap", r.table.fieldsMap)
	log.Println("--- r.table.fieldsMap", r.table.fieldsMap)
	if _, ok := r.table.fieldsMap[f]; ok {
		log.Println("---********************", f, v)
		r.ints[f] = v
		r.setFieldHasValue(f)
		return nil
	}
	return errors.New("Field of type string named [" + f + "] does not exist")
}

func (r *Row) SetString(f string, v string) error {
	if !r.isInited {
		return errors.New("Row is not inited. Need to run row.Init()")
	}
	log.Println("---", f, v)
	if _, ok := r.table.fieldsMap[f]; ok {
		log.Println("---********************", f, v)
		r.strings[f] = v
		r.setFieldHasValue(f)
		return nil
	}
	err := errors.New("Field of type string named [" + f + "] does not exist")
	return err
}

func (r *Row) fieldValuesForPreparedStatement() []interface{} {
	fields := make([]interface{}, 0)
	for i, _ := range r.table.Fields {
		field := r.table.Fields[i]
		fieldName := field.GetName()
		if r.hasValue(fieldName) {
			if v, ok := r.ints[fieldName]; ok {
				fields = append(fields, v)
			}

			if v, ok := r.strings[fieldName]; ok {
				fields = append(fields, v)
			}
			//log.Println("fieldsForPreparedStatement", i, field, fields)
		}
	}
	return fields
}

func (r *Row) fieldNamesForPreparedStatement() []string {
	fields := make([]string, 0)
	for i, _ := range r.table.Fields {
		field := r.table.Fields[i]
		fieldName := field.GetName()
		if r.hasValue(fieldName) {
			if _, ok := r.ints[fieldName]; ok {
				fields = append(fields, fieldName)
			}

			if _, ok := r.strings[fieldName]; ok {
				fields = append(fields, fieldName)
			}
			//log.Println("fieldsForPreparedStatement", i, field, fields)
		}
	}
	return fields
}
