package main

import (
	"errors"
	"fmt"
	"log"
)

type Table struct {
	name       string
	Fields     []Field
	fieldsMap  map[string]Field
	inited     bool
	oneToMany  []*OneToMany
	manyToMany []*ManyToMany
}

func NewTable(tableName string) (*Table, error) {
	if len(tableName) == 0 {
		return nil, errors.New("Table name is zero length")
	}
	t := new(Table)
	t.name = tableName
	t.Fields = make([]Field, 0)
	t.fieldsMap = make(map[string]Field, 0)
	t.oneToMany = make([]*OneToMany, 0)
	t.manyToMany = make([]*ManyToMany, 0)
	t.inited = true

	return t, nil
}

func (t *Table) Print() {
	fmt.Println("---------------------------------------------------------------------------------------------------------------")
	fmt.Println("Table:", t.name)
	for _, f := range t.Fields {
		fmt.Println("\t", f.String())
	}
	for _, r := range t.oneToMany {
		r.Print()
	}

	for _, r := range t.manyToMany {
		r.Print()
	}
	fmt.Println("---------------------------------------------------------------------------------------------------------------")
}

// id
// ADD timestamp
func (t *Table) SetDefaultFields() error {
	id, err := t.NewFieldInt64("id")
	if err != nil {
		log.Println(err)
		return err
	}
	id.PrimaryKey()

	return nil
}

func (t *Table) NewOneToMany(rightTable *Table, leftKeyField Field, uniqueRightFields ...Field) (*OneToMany, error) {
	if rightTable == nil {
		return nil, errors.New("Right table must be not nil")
	}
	if leftKeyField == nil {
		return nil, errors.New("leftFieldKey cannt be nil")
	}
	if len(uniqueRightFields) == 0 {
		return nil, errors.New("Need at least one uniqueRightField")
	}

	r := new(OneToMany)
	r.leftTable = t
	r.rightTable = rightTable
	r.leftKeyField = leftKeyField
	r.rightKeyField = "id"
	r.rightTableUniqueFields = uniqueRightFields
	t.oneToMany = append(t.oneToMany, r)
	return r, nil
}

func (t *Table) NewManyToMany(rightTable *Table, uniqueRightFields ...Field) *OneToMany {
	r := new(OneToMany)
	r.rightTable = rightTable
	r.rightKeyField = "id"
	r.rightTableUniqueFields = uniqueRightFields
	t.oneToMany = append(t.oneToMany, r)
	return r
}

func (t *Table) initField(f Field) {
	//f.Table(t)
}

func (t *Table) NewFieldString(name string) (FieldInt64, error) {
	f := new(FieldStringImpl)
	return f, t.add(f, name)
}

func (t *Table) AddOneToMany(r *OneToMany) error {
	if r == nil {
		return errors.New("OneToMany is nil")
	}
	t.oneToMany = append(t.oneToMany, r)
	return nil
}

func (t *Table) AddManyToMany(r *ManyToMany) error {
	if r == nil {
		return errors.New("ManyToMany is nil")
	}
	t.manyToMany = append(t.manyToMany, r)
	return nil
}

func (t *Table) NewFieldInt64(name string) (FieldInt64, error) {
	f := new(FieldInt64Impl)
	return f, t.add(f, name)
}

func (t *Table) add(f Field, name string) error {
	if !t.inited {
		return errors.New("Table not inited")
	}
	if len(name) == 0 {
		return errors.New("Field name is zero length")
	}
	f.Name(name)
	t.Fields = append(t.Fields, f)
	t.fieldsMap[f.GetName()] = f
	log.Println(f, t.Fields)
	return nil
}

func (t *Table) SetName(n string) {
	t.name = n
}

func (t *Table) NewRow() *Row {
	row := new(Row)
	row.Init(t)
	return row
}
