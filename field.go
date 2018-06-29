package main

const quote = "'"

type Field interface {
	Name(string)
	GetName() string
	Length(int)
	Unique()
	Indexed()
	NotNullable()
	PrimaryKey()
	String() string
	//Table(*Table)
}

func (f *FieldImpl) Table(t *Table) {
	f.table = t
}

func (f *FieldImpl) String() string {
	return f.name
}

func (f *FieldImpl) Length(l int) {
	f.length = l
}

func (f *FieldImpl) PrimaryKey() {
	f.primaryKey = true
}

func (f *FieldImpl) Name(n string) {
	f.name = n
}

func (f *FieldImpl) GetName() string {
	return f.name
}

func (f *FieldImpl) Unique() {
	f.unique = true
}

func (f *FieldImpl) Indexed() {
	f.indexed = true
}

func (f *FieldImpl) NotNullable() {
	f.notNullable = true
}

type FieldBlobImpl struct {
	FieldImpl
}

type FieldInt64Impl struct {
	FieldImpl
	value int64
}

type FieldStringImpl struct {
	FieldImpl
	value string
}

func (f *FieldStringImpl) CreateSql() string {
	return quote + f.name + quote + " varchar(" // + f.length + ")"
}

func (f *FieldStringImpl) Value() interface{} {
	return f.value
}

func (f *FieldStringImpl) SetValue(v string) {
	f.value = v
}

func (f *FieldInt64Impl) CreateSql() string {
	return quote + f.name + quote + " int"
}

func (f *FieldInt64Impl) Value() interface{} {
	return f.value
}

func (f *FieldInt64Impl) SetValue(v int64) {
	f.value = v
}

type FieldBlob interface {
	Field
}

type FieldInt64 interface {
	Field
}

type FieldString interface {
	Field
}

type FieldImpl struct {
	name        string
	length      int
	unique      bool
	indexed     bool
	primaryKey  bool
	notNullable bool
	table       *Table
}
