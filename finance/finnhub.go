package finance

import (
	"time"

	"github.com/m1/go-finnhub"
	"github.com/m1/go-finnhub/client"
)

type Finnhub struct {
	Client client.Client
	Error  error
}

func NewFinnhub(key string) *Finnhub {
	c := client.New(key)
	return &Finnhub{Client: *c, Error: nil}
}

func (fh *Finnhub) GetProfile(symbol string) *Company {
	var p *Company
	p, fh.Error = fh.Client.Stock.GetProfile(symbol)
	return p
}

func (fh *Finnhub) GetQuote(symbol string) *Quote {
	var q *Quote
	q, fh.Error = fh.Client.Stock.GetQuote(symbol)
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
	c, fh.Error = fh.Client.Stock.GetCandle(symbol, finnhub.CandleResolutionSecond, param)
	return c
}

func (fh *Finnhub) GetError() error {
	return fh.Error
}
