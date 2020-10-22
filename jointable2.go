package main

import (
	"errors"
)

func NewJoinTable2(leftTable, rightTable *Table, additionalFields []*Field, keyFields ...*Field) (*Table, error) {
	if err := errorsNewJoinTable2(leftTable, rightTable, additionalFields, keyFields...); err != nil {
		return nil, err
	}

	lf := new(Field)
	lf.typ = leftTable.pk.typ
	lf.name = leftTable.name + "_" + leftTable.pk.name

	rf := new(Field)
	rf.typ = rightTable.pk.typ
	rf.name = rightTable.name + "_" + rightTable.pk.name

	jt := new(Table)
	jt.name = "jt_" + leftTable.name + "_" + rightTable.name
	jt.AddField(lf)
	jt.pk = lf
	jt.AddField(rf)

	if additionalFields != nil {
		for i, _ := range additionalFields {
			af := additionalFields[i]
			if af == nil {
				return nil, errors.New("Additional field is nil")
			}
			err := jt.AddField(af)
			if err != nil {
				return nil, err
			}
		}
	}
	jtInfo := new(JoinTableInfo)
	jt.joinTableInfo = jtInfo
	jtInfo.leftTable = leftTable
	jtInfo.rightTable = rightTable
	jtInfo.rightTableIDCacheKeyFields = keyFields

	return jt, nil
}

func errorsNewJoinTable2(leftTable, rightTable *Table, additionalFields []*Field, keyFields ...*Field) error {
	if leftTable == nil {
		return errors.New("left table is nil")
	}
	if rightTable == nil {
		return errors.New("right table is nil")
	}
	if leftTable.pk == nil {
		return errors.New("left table pk is nil")
	}
	if rightTable.pk == nil {
		return errors.New("right table pk is nil")
	}
	if leftTable.pk.name == "" {
		return errors.New("left table pk name is empty")
	}
	if rightTable.pk.name == "" {
		return errors.New("right table pk name is empty")

	}
	return nil
}
