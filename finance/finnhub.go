package finance

import "fmt"

type Finnhub struct {
	key string
}

func NewFinnhub(key string) *Finnhub {
	return &Finnhub{key: key}
}

func (fh *Finnhub) Get() {
}

func (fh *Finnhub) Print() {
	fmt.Println("finance/finnhub")
}
