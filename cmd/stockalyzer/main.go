package main

import (
	"github.com/cfanatic/stockalyzer/finance"
)

func main() {
	var stock finance.IFinance

	stock = finance.NewFinnhub("ADS.DE")
	stock.GetCandle("2020-02-27 08:00:00", "2020-02-27 22:00:00")

	finance.Print(stock)
	finance.Plot(stock)
}
