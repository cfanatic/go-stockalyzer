package main

import (
	"github.com/cfanatic/stockalyzer/finance"
)

const (
	FINNHUB_KEY = ""
)

func main() {
	var stock finance.IFinance

	stock = finance.NewFinnhub(FINNHUB_KEY)
	stock.Print()
}
