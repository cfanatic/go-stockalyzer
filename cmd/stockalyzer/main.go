package main

import (
	"github.com/cfanatic/stockalyzer/finance"
)

func main() {
	var stock finance.IFinance

	stock = finance.NewFinnhub("ADS.DE")
	stock.GetCandle("2020-03-04 08:00:00", "2020-03-04 22:00:00")
	stock.GetChart(finance.M6)

	finance.Print(stock)
	finance.Plot(stock)
}
