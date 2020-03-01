package finance

import (
	"fmt"
	"os"
	"time"

	"github.com/wcharczuk/go-chart"
)

type IFinance interface {
	GetProfile() *Profile
	GetQuote() *Quote
	GetCandle(from, to string) *Candle

	GetError() error

	Ticker() *string
	XValues() *[]time.Time
	YValues() *[]float64
}

type Finance struct {
	Ticker string
	*Profile
	*Quote
	*Candle
}

type Profile struct {
	Country      string
	Currency     string
	Description  string
	Exchange     string
	GICSIndustry string
	GICSSector   string
	ISIN         string
	Name         string
}

type Quote struct {
	Open      float64
	High      float64
	Low       float64
	Current   float64
	PrevClose float64
}

type Candle struct {
	Close  []float64
	High   []float64
	Low    []float64
	Open   []float64
	Times  []time.Time
	Volume []float64
}

func Plot(stock IFinance) {
	if err := stock.GetError(); err == nil {
		quotes := chart.TimeSeries{
			Name: *stock.Ticker(),
			Style: chart.Style{
				StrokeColor: chart.GetDefaultColor(0),
			},
			XValues: *stock.XValues(),
			YValues: *stock.YValues(),
		}

		graph := chart.Chart{
			XAxis: chart.XAxis{
				ValueFormatter: chart.TimeHourValueFormatter,
				TickPosition:   chart.TickPositionBetweenTicks,
			},
			// YAxis: chart.YAxis{
			// 	Range: &chart.ContinuousRange{},
			// },
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
	} else {
		panic(err)
	}
}

func Print(stock IFinance) {
	if err := stock.GetError(); err == nil {
		time := *stock.XValues()
		price := *stock.YValues()
		for i := range time {
			fmt.Printf("%3d | %+v | %+v | %v\n", i, time[i].Unix(), time[i], price[i])
		}
	} else {
		panic(err)
	}
}
