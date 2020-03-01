package finance

import (
	"time"

	"github.com/cfanatic/stockalyzer/configuration"

	"github.com/m1/go-finnhub"
	"github.com/m1/go-finnhub/client"
)

type Finnhub struct {
	Finance Finance
	client  client.Client
	err     error
}

func NewFinnhub(symbol string) *Finnhub {
	key := configuration.FINNHUB_TOKEN
	token := configuration.Get(key).(string)
	c := client.New(token)
	return &Finnhub{client: *c, Finance: Finance{Ticker: symbol}, err: nil}
}

func (fh *Finnhub) GetProfile() *Profile {
	p := Profile{}
	if profile, err := fh.client.Stock.GetProfile(fh.Finance.Ticker); err == nil {
		p.Country = profile.Country
		p.Currency = profile.Currency
		p.Description = profile.Description
		p.Exchange = profile.Exchange
		p.GICSIndustry = profile.GICSIndustry
		p.GICSSector = profile.GICSSector
		p.ISIN = profile.ISIN
		p.Name = profile.Name
		fh.Finance.Profile = &p
	} else {
		fh.err = err
	}
	return &p
}

func (fh *Finnhub) GetQuote() *Quote {
	q := Quote{}
	if quote, err := fh.client.Stock.GetQuote(fh.Finance.Ticker); err == nil {
		q.Open = quote.Open
		q.High = quote.High
		q.Low = quote.Low
		q.Current = quote.Current
		q.PrevClose = quote.PreviousClose
		fh.Finance.Quote = &q
	} else {
		fh.err = err
	}
	return &q
}

func (fh *Finnhub) GetCandle(from, to string) *Candle {
	c := Candle{}
	layout := "2006-01-02 15:04:05"
	t1, _ := time.Parse(layout, from)
	t2, _ := time.Parse(layout, to)
	param := &finnhub.CandleParams{
		Count: nil,
		From:  &t1,
		To:    &t2,
	}
	if candle, err := fh.client.Stock.GetCandle(fh.Finance.Ticker, finnhub.CandleResolutionSecond, param); err == nil {
		c.Close = candle.Close
		c.High = candle.High
		c.Low = candle.Low
		c.Open = candle.Open
		c.Times = candle.Times
		c.Volume = candle.Volume
		fh.Finance.Candle = &c
	} else {
		fh.err = err
	}
	return &c
}

func (fh *Finnhub) Ticker() *string {
	if fh.err != nil {
		panic(fh.err)
	}
	return &fh.Finance.Ticker
}

func (fh *Finnhub) XValues() *[]time.Time {
	if fh.err != nil {
		panic(fh.err)
	}
	return &fh.Finance.Candle.Times
}

func (fh *Finnhub) YValues() *[]float64 {
	if fh.err != nil {
		panic(fh.err)
	}
	return &fh.Finance.Candle.Open
}
