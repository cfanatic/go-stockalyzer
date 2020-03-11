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
	GetChart(duration Duration) *Candle

	Ticker() *string
	Duration() *Duration
	XValues() *[]time.Time
	YValues() *[]float64
}

type Finance struct {
	Ticker   string
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
	var time []time.Time
	var tick []chart.Tick
	var grid []chart.GridLine
	switch *stock.Duration() {
	case D1:
		time = *stock.XValues()
		tick = append([]chart.Tick{}, chart.Tick{Value: float64(0), Label: fmt.Sprintf("%s", time[0].Format("01-02 3PM"))})
		grid = []chart.GridLine{}
		for i := 1; i < len(time)-1; i++ {
			hour1, _, _ := time[i].Clock()
			hour2, _, _ := time[i+1].Clock()
			if hour1 != hour2 {
				tick = append(tick, chart.Tick{Value: float64(i + 1), Label: fmt.Sprintf("%s", time[i+1].Format("01-02 3PM"))})
				grid = append(grid, chart.GridLine{Value: float64(i + 1)})
			}
		}
		tick = append(tick, chart.Tick{Value: float64(len(time)), Label: ""})
	case D5, D10:
		time = *stock.XValues()
		tick = append([]chart.Tick{}, chart.Tick{Value: float64(0), Label: fmt.Sprintf("%s", time[0].Format("2006-01-02"))})
		grid = []chart.GridLine{}
		for i := 1; i < len(time)-1; i++ {
			_, _, day1 := time[i].Date()
			_, _, day2 := time[i+1].Date()
			if day1 != day2 {
				tick = append(tick, chart.Tick{Value: float64(i + 1), Label: fmt.Sprintf("%s", time[i+1].Format("2006-01-02"))})
				grid = append(grid, chart.GridLine{Value: float64(i + 1)})
			}
		}
		tick = append(tick, chart.Tick{Value: float64(len(time)), Label: ""})
	default:
		panic("Unkown chart duration parameter during plot")
	}
	graph := chart.Chart{
		XAxis: chart.XAxis{
			GridMajorStyle: chart.Style{
				StrokeColor: chart.ColorAlternateLightGray,
				StrokeWidth: 1.0,
			},
			GridLines:      grid,
			Ticks:          tick,
			TickPosition:   chart.TickPositionUnderTick,
			ValueFormatter: chart.TimeDateValueFormatter,
		},
		YAxis: chart.YAxis{
			GridMinorStyle: chart.Style{
				StrokeColor: chart.ColorAlternateLightGray,
				StrokeWidth: 1.0,
			},
			Range: &chart.ContinuousRange{},
			ValueFormatter: func(v interface{}) string {
				if v, isFloat := v.(float64); isFloat {
					return fmt.Sprintf("%0.f", v)
				}
				return ""
			},
		},
		Series: []chart.Series{
			chart.ContinuousSeries{
				Name: *stock.Ticker(),
				Style: chart.Style{
					StrokeColor: chart.GetDefaultColor(0),
				},
				XValues: func() []float64 {
					xvalues := make([]float64, len(time))
					for i := 0; i < len(time); i++ {
						xvalues[i] = float64(i)
					}
					return xvalues
				}(),
				YValues: *stock.YValues(),
			},
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
