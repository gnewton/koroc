package main

// borrowed from otira https://github.com/gnewton/otira/blob/master/dialect.go
type Dialect interface {
	InsertPreparedStatementSql(table string, fields []*Field) (string, error)
	DeleteByPKPreparedStatementSql(table string, pk string) (string, error)
	CreateTableSql(table string, fields []*Field, pk string) (string, error)
}
