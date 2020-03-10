package finance

import (
	"fmt"
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
	token := configuration.Get(configuration.FINNHUB_TOKEN).(string)
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
	if fh.dateEqual(t1, t2) {
		if t1.Weekday() == time.Saturday || t2.Weekday() == time.Sunday {
			panic("Stock market is closed on weekends")
		}
	}
	var resolution finnhub.CandleResolution
	switch fh.Finance.Duration {
	case D1:
		resolution = finnhub.CandleResolutionSecond
	case D5:
		resolution = finnhub.CandleResolution5Second
	case D10:
		resolution = finnhub.CandleResolution15Second
	case M3:
	case M6:
	case Y1:
	case Y5:
	case Max:
	default:
		panic("Unknown resolution parameter")
	}
	if candle, err := fh.client.Stock.GetCandle(fh.Finance.Ticker, resolution, param); err == nil {
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

func (fh *Finnhub) GetChart(duration Duration) *Candle {
	var from, to string
	now := time.Now()
	switch duration {
	case D1:
		from = fmt.Sprintf("%s 08:00:00", now.Format("2006-01-02"))
		to = fmt.Sprintf("%s 22:00:00", now.Format("2006-01-02"))
	case D5:
		then := now.AddDate(0, 0, fh.dateShift(now, 5))
		from = fmt.Sprintf("%s 08:00:00", then.Format("2006-01-02"))
		to = fmt.Sprintf("%s 22:00:00", now.Format("2006-01-02"))
	case D10:
		then := now.AddDate(0, 0, fh.dateShift(now, 10))
		from = fmt.Sprintf("%s 08:00:00", then.Format("2006-01-02"))
		to = fmt.Sprintf("%s 22:00:00", now.Format("2006-01-02"))
	case M3:
	case M6:
	case Y1:
	case Y5:
	case Max:
	default:
		panic("Unknown chart duration parameter")
	}
	fh.Finance.Duration = duration
	return fh.GetCandle(from, to)
}

func (fh *Finnhub) Ticker() *string {
	if fh.err != nil {
		panic(fh.err)
	}
	return &fh.Finance.Ticker
}

func (fh *Finnhub) Duration() *Duration {
	if fh.err != nil {
		panic(fh.err)
	}
	return &fh.Finance.Duration
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

func (fh *Finnhub) dateEqual(date1, date2 time.Time) bool {
	y1, m1, d1 := date1.Date()
	y2, m2, d2 := date2.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}

func (fh *Finnhub) dateShift(start time.Time, days int) int {
	delta := 0
	for i := 0; i < days; i++ {
		if (start.AddDate(0, 0, -i)).Weekday() == time.Saturday ||
			start.AddDate(0, 0, -i).Weekday() == time.Sunday {
			delta++
		}
	}
	// Adjust shift based on whether D5 or D10 is requested
	if days%2 == 0 {
		return -days - delta
	} else {
		return -days - delta + 1
	}
}
