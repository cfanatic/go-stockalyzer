package finance

import (
	"time"

	"github.com/m1/go-finnhub"
	"github.com/m1/go-finnhub/client"
)

type Finnhub struct {
	client client.Client
	name   string
	err    error
}

func NewFinnhub(key string) *Finnhub {
	c := client.New(key)
	return &Finnhub{client: *c, name: "", err: nil}
}

func (fh *Finnhub) GetProfile(symbol string) *Company {
	var p *Company
	if p, fh.err = fh.client.Stock.GetProfile(symbol); fh.err == nil {
		fh.name = symbol
	} else {
		fh.name = ""
	}
	return p
}

func (fh *Finnhub) GetQuote(symbol string) *Quote {
	var q *Quote
	if q, fh.err = fh.client.Stock.GetQuote(symbol); fh.err == nil {
		fh.name = symbol
	} else {
		fh.name = ""
	}
	return q
}

func (fh *Finnhub) GetCandle(symbol, from, to string) *Candle {
	var c *Candle
	layout := "2006-01-02 15:04:05"
	t1, _ := time.Parse(layout, from)
	t2, _ := time.Parse(layout, to)
	param := &finnhub.CandleParams{
		Count: nil,
		From:  &t1,
		To:    &t2,
	}
	if c, fh.err = fh.client.Stock.GetCandle(symbol, finnhub.CandleResolutionSecond, param); fh.err == nil {
		fh.name = symbol
	} else {
		fh.name = ""
	}
	return c
}

func (fh *Finnhub) GetName() string {
	return fh.name
}

func (fh *Finnhub) GetError() error {
	return fh.err
}
