package main

import (
	"errors"
	"strconv"
)

type Record struct {
	table     *Table
	values    []interface{}
	outValues []interface{}
}

func (r *Record) Initialize(initializeValues bool) error {
	if r.table == nil {
		return errors.New("Table is nil")
	}
	if r.table.fields == nil {
		return errors.New("Table.fields is nil")
	}
	if len(r.table.fields) == 0 {
		return errors.New("Table.fields is len 0")
	}
	if initializeValues {
		r.values = make([]interface{}, len(r.table.fields))
		r.outValues = make([]interface{}, len(r.table.fields))
		for i, _ := range r.values {
			r.outValues[i] = &r.values[i]
		}
	}
	return nil
}

func (r *Record) Reset() error {
	if r.values == nil {
		return errors.New("values is nil")
	}
	for i := 0; i < len(r.values); i++ {
		r.values[i] = nil
	}
	return nil
}
func (r *Record) Get(f *Field) (interface{}, error) {
	if f == nil {
		return nil, errors.New("field is nil")
	}
	positionInTable := f.positionInTable
	if positionInTable < 0 {
		return nil, errors.New("positionInTable index is < 0")
	}
	return r.values[positionInTable], nil

}
func (r *Record) AddN(i int, v interface{}) error {
	if r.table == nil {
		return errors.New("Table is nil")
	}

	if i < 0 {
		return errors.New("Index < 0")
	}

	if r.values == nil {
		r.Initialize(true)
	}

	if i >= len(r.values) {
		return errors.New("Out of bounds")
	}

	if err := r.table.fields[i].CheckValueType(v); err != nil {
		return err
	}

	r.values[i] = v
	return nil
}

func (r *Record) Add(f *Field, v interface{}) error {
	if r.table == nil {
		return errors.New("Table is nil")
	}
	if r.values == nil {
		var err error
		err = r.Initialize(true)
		if err != nil {
			return err
		}
	}
	if f.positionInTable > len(r.values) || f.positionInTable < 0 {
		return errors.New("Field positionInTable out of bounds:" + r.table.name + ":" + f.name + ":" + strconv.Itoa(f.positionInTable))
	}
	if err := r.table.fields[f.positionInTable].CheckValueType(v); err != nil {
		return err
	}
	r.values[f.positionInTable] = v
	return nil
}
