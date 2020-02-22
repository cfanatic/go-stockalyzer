package finance

import (
	"time"

	"github.com/m1/go-finnhub"
	"github.com/m1/go-finnhub/client"
)

type Finnhub struct {
	cli client.Client
	err error
}

func NewFinnhub(key string) *Finnhub {
	cli := client.New(key)
	return &Finnhub{cli: *cli}
}

func (fh *Finnhub) Profile(symbol string) *finnhub.Company {
	var profile *finnhub.Company
	profile, fh.err = fh.cli.Stock.GetProfile(symbol)
	return profile
}

func (fh *Finnhub) Candle(symbol string) *finnhub.Candle {
	var candle *finnhub.Candle

	// count := finnhub.CandleDefaultCount
	layout := "01/02/2006 3:04:05 PM"
	from, _ := time.Parse(layout, "02/21/2020 07:00:00 AM")
	to, _ := time.Parse(layout, "02/21/2020 09:00:00 PM")

	param := &finnhub.CandleParams{
		Count: nil,
		From:  &from,
		To:    &to,
	}
	candle, fh.err = fh.cli.Stock.GetCandle(symbol, finnhub.CandleResolutionSecond, param)
	return candle
}

func (fh *Finnhub) Quote(symbol string) *finnhub.Quote {
	var quote *finnhub.Quote
	quote, fh.err = fh.cli.Stock.GetQuote(symbol)
	return quote
}
