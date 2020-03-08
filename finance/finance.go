package finance

import (
	"fmt"
	"os"
	"time"

	"github.com/wcharczuk/go-chart"
)

type Duration int

const (
	D1 Duration = iota
	D5
	D10
	M3
	M6
	Y1
	Y5
	Max
)

type IFinance interface {
	GetProfile() *Profile
	GetQuote() *Quote
	GetCandle(from, to string) *Candle
	GetChart(period Duration) *Candle

	Ticker() *string
	Times() *[]time.Time
	XValues() *[]time.Time
	YValues() *[]float64
}

type Finance struct {
	Ticker   string
	Times    []time.Time
	Duration Duration
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
	times := stock.Times()
	ticks := []chart.Tick{
		chart.Tick{Value: float64(0), Label: fmt.Sprintf("%s", (*times)[0].Format("2006-01-02"))},
		chart.Tick{Value: float64(35), Label: fmt.Sprintf("%s", (*times)[1].Format("2006-01-02"))},
		chart.Tick{Value: float64(70), Label: fmt.Sprintf("%s", (*times)[2].Format("2006-01-02"))},
		chart.Tick{Value: float64(105), Label: fmt.Sprintf("%s", (*times)[3].Format("2006-01-02"))},
		chart.Tick{Value: float64(140), Label: fmt.Sprintf("%s", (*times)[4].Format("2006-01-02"))},
		chart.Tick{Value: float64(174), Label: ""},
	}
	xval := make([]float64, len(*stock.XValues()))
	for i := 0; i < len(*stock.XValues()); i++ {
		xval[i] = float64(i)
	}
	quotes := chart.ContinuousSeries{
		Name: *stock.Ticker(),
		Style: chart.Style{
			StrokeColor: chart.GetDefaultColor(0),
		},
		XValues: xval,
		YValues: *stock.YValues(),
	}
	graph := chart.Chart{
		XAxis: chart.XAxis{
			Ticks:          ticks,
			TickPosition:   chart.TickPositionUnderTick,
			ValueFormatter: chart.TimeDateValueFormatter,
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

func Print(stock IFinance) {
	time := *stock.XValues()
	price := *stock.YValues()
	for i := range time {
		fmt.Printf("%3d | %+v | %+v | %v\n", i, time[i].Unix(), time[i], price[i])
	}
}
