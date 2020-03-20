package main

import (
	"time"

	"github.com/cfanatic/stockalyzer/finance"
)

const (
	mode     = "plot"
	duration = finance.D5
	sleep    = 5
)

var (
	durations = [...]finance.Duration{
		finance.Intraday,
		finance.D5,
		finance.D10,
		finance.M3,
		finance.M6,
		finance.Y1,
		finance.Y3,
		finance.Y5,
		finance.Max,
	}
)

func main() {
	var stock finance.IFinance

	switch stock = finance.NewFinnhub("ADS.DE"); mode {
	case "print":
		stock.GetCandle("2020-03-01 08:00:00", "2020-03-10 22:00:00")
		finance.Print(stock)
	case "plot":
		stock.GetChart(duration)
		finance.Plot(stock)
	case "plots": // buggy: panic is thrown due to request limit
		for _, duration := range durations {
			stock.GetChart(duration)
			finance.Plot(stock)
			time.Sleep(sleep * time.Second)
		}
	case "performance":
		finance.Performance(stock)
	}
}
