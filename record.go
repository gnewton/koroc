package main

import (
	"errors"
	"strconv"
)

type Record struct {
	table  *Table
	values []interface{}
}

func (r *Record) Initialize() (*Record, error) {
	if r.table == nil {
		return nil, errors.New("Table is nil")
	}
	if r.table.fields == nil {
		return nil, errors.New("Table.fields is nil")
	}
	if len(r.table.fields) == 0 {
		return nil, errors.New("Table.fields is len 0")
	}
	r.values = make([]interface{}, len(r.table.fields))
	return r, nil
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

func (r *Record) AddN(i int, v interface{}) error {
	if r.table == nil {
		return errors.New("Table is nil")
	}

	if i < 0 {
		return errors.New("Index < 0")
	}

	if r.values == nil {
		r.Initialize()
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
		r, err = r.Initialize()
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
