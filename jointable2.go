package main

import (
//"errors"
)

func NewJoinTable2(leftTable, rightTable *Table) (*Table, error) {
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

	jt.leftTable = leftTable
	jt.rightTable = rightTable

	return jt, nil
}
