package main

import (
	//"github.com/gnewton/
	"sync"
)

type Persister interface {
	Persist(chan []*Article, sync.WaitGroup) error
}
