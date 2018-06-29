package main

type Counter interface {
	Next() int64
}

type CounterImpl struct {
	n int64
}

func NewCounter() *CounterImpl {
	counter := new(CounterImpl)
	counter.n = 0
	return counter
}

func (counter *CounterImpl) Next() int64 {
	counter.n = counter.n + 1
	return counter.n
}
