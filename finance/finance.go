package finance

import "github.com/m1/go-finnhub"

type Company = finnhub.Company
type Candle = finnhub.Candle
type Quote = finnhub.Quote

type IFinance interface {
	GetProfile(symbol string) *Company
	GetQuote(symbol string) *Quote
	GetCandle(symbol, from, to string) *Candle

	GetName() string
	GetError() error
}
