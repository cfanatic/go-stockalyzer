package main

import (
	"log"
	"os"
	"time"

	"github.com/cfanatic/stockalyzer/finance"
)

const (
	mode     = "plots"
	duration = finance.Max
	sleep    = 500
	company  = "ADS.DE"
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

	log.SetOutput(os.Stdout)

	switch stock = finance.NewFinnhub(company); mode {
	case "print":
		stock.GetCandle("2020-03-01 08:00:00", "2020-03-10 22:00:00")
		finance.Print(stock)
	case "plot":
		stock.GetChart(duration)
		finance.Plot(stock)
	case "plots":
		for _, duration := range durations {
			log.Println("Plotting chart for " + finance.Durations[duration])
			stock.GetChart(duration)
			finance.Plot(stock)
			time.Sleep(sleep * time.Millisecond)
		}
	case "performance":
		finance.Performance(stock)
	}
}
