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
	t1, _ := time.Parse("2006-01-02 15:04:05", from)
	t2, _ := time.Parse("2006-01-02 15:04:05", to)
	fh.dateDuration(t1, t2)
	switch *fh.Duration() {
	case Intraday:
		t1 = t1.AddDate(0, 0, fh.dateShift(t1, 0))
		resolution = finnhub.CandleResolutionSecond
	case D5:
		t1 = t1.AddDate(0, 0, fh.dateShift(t2, 5))
		resolution = finnhub.CandleResolution5Second
	case D10:
		t1 = t1.AddDate(0, 0, fh.dateShift(t2, 10))
		resolution = finnhub.CandleResolution15Second
	case M1, M3, M6:
		resolution = finnhub.CandleResolutionMinute
	case Y1, Y3, Y5:
		resolution = finnhub.CandleResolutionDay
	case Max:
		resolution = finnhub.CandleResolutionWeek
	default:
		panic("Unknown chart duration parameter")
	}
	param := &finnhub.CandleParams{
		Count: nil,
		From:  &t1,
		To:    &t2,
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
	now = time.Now()
	switch duration {
	case Intraday:
		then = now.AddDate(0, 0, 0)
	case D5:
		then = now.AddDate(0, 0, -5)
	case D10:
		then = now.AddDate(0, 0, -10)
	case M1:
		then = now.AddDate(0, -1, 0)
	case M3:
		then = now.AddDate(0, -3, 0)
	case M6:
		then = now.AddDate(0, -6, 0)
	case Y1:
		then = now.AddDate(-1, 0, 0)
	case Y3:
		then = now.AddDate(-3, 0, 0)
	case Y5:
		then = now.AddDate(-5, 0, 0)
	case Max:
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

func (fh *Finnhub) dateDuration(from, to time.Time) {
	t1, _ := time.Parse("2006-01-02", fmt.Sprintf("%d-%02d-%02d", from.Year(), from.Month(), from.Day()))
	t2, _ := time.Parse("2006-01-02", fmt.Sprintf("%d-%02d-%02d", to.Year(), to.Month(), to.Day()))
	switch diff := t2.Sub(t1).Hours(); {
	case diff <= 24:
		fh.Finance.Duration = Intraday
	case diff <= 24*5:
		fh.Finance.Duration = D5
	case diff <= 24*10:
		fh.Finance.Duration = D10
	case diff <= 24*31*1:
		fh.Finance.Duration = M1
	case diff <= 24*31*3:
		fh.Finance.Duration = M3
	case diff <= 24*31*6:
		fh.Finance.Duration = M6
	case diff <= 24*31*12:
		fh.Finance.Duration = Y1
	case diff <= 24*31*12*3:
		fh.Finance.Duration = Y3
	case diff <= 24*31*12*5:
		fh.Finance.Duration = Y5
	case diff > 24*31*12*5:
		fh.Finance.Duration = Max
	default:
		panic("Unknown date duration")
	}
}

func (fh *Finnhub) dateShift(start time.Time, days int) int {
	delta, shift := 0, 0
	switch *fh.Duration() {
	case Intraday:
		if start.Weekday() == time.Saturday {
			shift = -1
		} else if start.Weekday() == time.Sunday {
			shift = -2
		}
	case D5:
		for i := 0; i < days; i++ {
			if (start.AddDate(0, 0, -i)).Weekday() == time.Saturday ||
				start.AddDate(0, 0, -i).Weekday() == time.Sunday {
				delta++
			}
		}
		shift = -delta
	case D10:
		for i := 0; i < days; i++ {
			if (start.AddDate(0, 0, -i)).Weekday() == time.Saturday ||
				start.AddDate(0, 0, -i).Weekday() == time.Sunday {
				delta++
			}
		}
		shift = -delta - 1
	default:
		panic("Unknown chart duration parameter")
	}
	return shift
}
