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
	var resolution finnhub.CandleResolution
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
	switch *fh.Duration() {
	case D1:
		resolution = finnhub.CandleResolutionSecond
	case D5:
		resolution = finnhub.CandleResolution5Second
	case D10:
		resolution = finnhub.CandleResolution15Second
	case M3:
		resolution = finnhub.CandleResolutionMinute
	case M6:
		resolution = finnhub.CandleResolutionMinute
	case Y1:
		resolution = finnhub.CandleResolutionDay
	case Y5:
		resolution = finnhub.CandleResolutionDay
	case Y30:
		resolution = finnhub.CandleResolutionWeek
	default:
		panic("Unknown chart duration parameter")
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
	var now, then time.Time
	fh.Finance.Duration = duration
	now = time.Now()
	switch duration {
	case D1:
		then = now.AddDate(0, 0, fh.dateShift(now, 0))
	case D5:
		then = now.AddDate(0, 0, fh.dateShift(now, 5))
	case D10:
		then = now.AddDate(0, 0, fh.dateShift(now, 10))
	case M3:
		then = now.AddDate(0, -3, 0)
	case M6:
		then = now.AddDate(0, -6, 0)
	case Y1:
		then = now.AddDate(-1, 0, 0)
	case Y5:
		then = now.AddDate(-5, 0, 0)
	case Y30:
		then = now.AddDate(-30, 0, 0)
	default:
		panic("Unknown chart duration parameter")
	}
	from = fmt.Sprintf("%s 08:00:00", then.Format("2006-01-02"))
	to = fmt.Sprintf("%s 22:00:00", now.Format("2006-01-02"))
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
	delta, shift := 0, 0
	for i := 0; i < days; i++ {
		if (start.AddDate(0, 0, -i)).Weekday() == time.Saturday ||
			start.AddDate(0, 0, -i).Weekday() == time.Sunday {
			delta++
		}
	}
	switch *fh.Duration() {
	case D1:
		shift = -days - delta
	case D5:
		shift = -days - delta
	case D10:
		shift = -days - delta - 1
	default:
		panic("Unknown chart duration parameter")
	}
	return shift
}
