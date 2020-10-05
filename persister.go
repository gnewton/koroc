package main

import (
	"github.com/gnewton/pubmedSqlStructs"
	"sync"
)

type Persister interface {
	Persist(chan []*pubmedSqlStructs.Article, sync.WaitGroup) error
}
