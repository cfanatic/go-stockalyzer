package finance

import (
	"time"

	"github.com/m1/go-finnhub"
	"github.com/m1/go-finnhub/client"
)

type Finnhub struct {
	Client client.Client
	Name   string
	Error  error
}

func NewFinnhub(key string) *Finnhub {
	c := client.New(key)
	return &Finnhub{Client: *c, Name: "", Error: nil}
}

func (fh *Finnhub) GetProfile(symbol string) *Company {
	var p *Company
	if p, fh.Error = fh.Client.Stock.GetProfile(symbol); fh.Error == nil {
		fh.Name = symbol
	} else {
		fh.Name = ""
	}
	return p
}

func (fh *Finnhub) GetQuote(symbol string) *Quote {
	var q *Quote
	if q, fh.Error = fh.Client.Stock.GetQuote(symbol); fh.Error == nil {
		fh.Name = symbol
	} else {
		fh.Name = ""
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
	if c, fh.Error = fh.Client.Stock.GetCandle(symbol, finnhub.CandleResolutionSecond, param); fh.Error == nil {
		fh.Name = symbol
	} else {
		fh.Name = ""
	}
	return c
}

func (fh *Finnhub) GetName() string {
	return fh.Name
}

func (fh *Finnhub) GetError() error {
	return fh.Error
}
