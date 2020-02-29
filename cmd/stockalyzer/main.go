package main

import (
	"fmt"
	"os"

	"github.com/cfanatic/stockalyzer/finance"
	"github.com/wcharczuk/go-chart"
)

const (
	FINNHUB_KEY = ""
)

func main() {
	var stock finance.IFinance

	stock = finance.NewFinnhub(FINNHUB_KEY)

	// profile := stock.GetProfile("ADS.DE")
	// quote := stock.GetQuote("ADS.DE")
	candle := stock.GetCandle("ADS.DE", "2020-02-27 08:00:00", "2020-02-27 22:00:00")

	if err := stock.GetError(); err == nil {
		for i := range candle.Times {
			fmt.Printf("%3d | %+v | %+v | %v\n", i, candle.Times[i].Unix(), candle.Times[i], candle.Open[i])
		}
	} else {
		panic(err)
	}

	quotes := chart.TimeSeries{
		Name: stock.GetName(),
		Style: chart.Style{
			StrokeColor: chart.GetDefaultColor(0),
		},
		XValues: candle.Times,
		YValues: candle.Open,
	}

	graph := chart.Chart{
		XAxis: chart.XAxis{
			ValueFormatter: chart.TimeHourValueFormatter,
			TickPosition:   chart.TickPositionBetweenTicks,
		},
		YAxis: chart.YAxis{
			// Range: &chart.ContinuousRange{
			// 	Max: 300.0,
			// 	Min: 200.0,
			// },
		},
		Series: []chart.Series{
			quotes,
		},
	}

	graph.Elements = []chart.Renderable{
		chart.Legend(&graph),
	}

	f, _ := os.Create("bin/output.png")
	defer f.Close()
	graph.Render(chart.PNG, f)
}
