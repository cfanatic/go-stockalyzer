package main

import (
	"fmt"

	"github.com/cfanatic/stockalyzer/finance"
)

func main() {
	var stock finance.IFinance

	stock = finance.NewFinnhub("ADS.DE")
	candle := stock.GetCandle("2020-02-27 08:00:00", "2020-02-27 22:00:00")

	if err := stock.GetError(); err == nil {
		for i := range candle.Times {
			fmt.Printf("%3d | %+v | %+v | %v\n", i, candle.Times[i].Unix(), candle.Times[i], candle.Open[i])
		}
	} else {
		panic(err)
	}

	finance.Plot(stock)
}
