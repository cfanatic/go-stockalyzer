package finance

import (
	"os"
	"time"

	"github.com/wcharczuk/go-chart"
)

type IFinance interface {
	GetProfile() *Profile
	GetQuote() *Quote
	GetCandle(from, to string) *Candle

	GetError() error
}

type Finance struct {
	Ticker  string
	Profile Profile
	Quote   Quote
	Candle  Candle
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
	quotes := chart.TimeSeries{
		Name: stock.GetName(),
		Style: chart.Style{
			StrokeColor: chart.GetDefaultColor(0),
		},
		XValues: nil,
		YValues: nil,
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
}

func Print() {
}
