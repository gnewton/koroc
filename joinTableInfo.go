package main

import (
	"errors"
	"strconv"
)

type JoinTableInfo struct {
	leftTable, rightTable      *Table
	rightTableIDCache          map[string]uint64
	rightTableIDCacheKeyFields []*Field
	rightTableIDCounter        uint64
}

const SEPARATOR = "|"

func (jt *JoinTableInfo) makeKey(rec *Record) (string, error) {
	if jt.rightTableIDCacheKeyFields == nil {
		return "", errors.New("rightTableIDCacheKeyFields is nil")
	}
	if len(jt.rightTableIDCacheKeyFields) == 0 {
		return "", errors.New("rightTableIDCacheKeyFields is len 0")
	}

	var key string
	for i, _ := range jt.rightTableIDCacheKeyFields {
		if jt.rightTableIDCacheKeyFields[i] == nil {
			return "", errors.New("rightTableIDCacheKeyFields field is nil")
		}
		field := jt.rightTableIDCacheKeyFields[i]
		fieldType := field.typ
		fieldValue := rec.values[field.positionInTable]
		if i != 0 {
			key += SEPARATOR
		}
		valueString, err := fieldType.ValueToString(field, fieldValue)
		if err != nil {
			return "", err
		}
		key += (strconv.Itoa(i) + ":" + valueString)
	}
	return key, nil
}
