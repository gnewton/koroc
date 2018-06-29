package main

import (
	"errors"
	//"log"
)

type Join struct {
	leftKey   string
	rightKey  string
	tableName string
	cache     map[int64]struct{}
}

func NewJoin(leftKey, rightKey, tableName string) (*Join, error) {
	if len(leftKey) == 0 || len(rightKey) == 0 {
		return nil, errors.New("One of the keys is empty")
	}
	if len(tableName) == 0 {
		return nil, errors.New("Table name is empty")
	}
	//join := new(Join)
	return nil, nil
}

func (j *Join) Save(leftValue, rightValue int64) error {

	return nil
}

func (j *Join) Create() error {

	return nil
}
