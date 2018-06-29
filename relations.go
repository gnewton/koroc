package main

import "fmt"

type Relation struct {
	leftTable              *Table
	rightTable             *Table
	rightTableUniqueFields []Field
}

type OneToMany struct {
	Relation
	leftKeyField  Field
	rightKeyField string
}

type ManyToMany struct {
	Relation
	leftKeyField, rightKeyField string
}

func (r *OneToMany) Print() {
	fmt.Print("\tRelation: OneToMany: ", r.leftKeyField.GetName(), "  Table: [", r.rightTable.name, ".", r.leftKeyField.GetName(), "] UniqueFields: [")
	for i, f := range r.rightTableUniqueFields {
		if i != 0 {
			fmt.Print(", ")
		}
		fmt.Print(f.GetName())
	}
	fmt.Print("]")
	fmt.Println()
}

func (r *ManyToMany) Print() {
	fmt.Print("\tRelation: ManyToMany", r.rightTable.name)
	for _, f := range r.rightTableUniqueFields {
		fmt.Print(", ", f.GetName())
	}
	fmt.Println()
}
