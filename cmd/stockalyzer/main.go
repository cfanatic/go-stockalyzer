package main

import (
	"fmt"

	"github.com/cfanatic/stockalyzer/finance"
)

const (
	FINNHUB_KEY = ""
)

func main() {
	var stock finance.IFinance

	stock = finance.NewFinnhub(FINNHUB_KEY)

	profile := stock.Profile("AAPL")
	candle := stock.Candle("AAPL")
	quote := stock.Quote("AAPL")

	fmt.Printf("%+v\n", profile)
	fmt.Printf("%+v\n", candle.Times)
	fmt.Printf("%+v\n", quote)
}
