package finance

import "github.com/m1/go-finnhub"

type IFinance interface {
	Profile(symbol string) *finnhub.Company
	Candle(symbol string) *finnhub.Candle
	Quote(symbol string) *finnhub.Quote
}
