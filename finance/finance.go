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
	graph := chart.Chart{}
	switch *stock.Duration() {
	case D1:
		time := *stock.XValues()
		grids := []chart.GridLine{}
		for i := 1; i < len(time)-1; i++ {
			hour1, _, _ := time[i].Clock()
			hour2, _, _ := time[i+1].Clock()
			if hour1 != hour2 {
				grids = append(grids, chart.GridLine{Value: chart.TimeToFloat64(time[i+1])})
			}
		}
		graph = chart.Chart{
			XAxis: chart.XAxis{
				ValueFormatter: chart.TimeHourValueFormatter,
				TickPosition:   chart.TickPositionUnderTick,
				GridMajorStyle: chart.Style{
					StrokeColor: chart.ColorAlternateLightGray,
					StrokeWidth: 1.0,
				},
				GridLines: grids,
			},
			YAxis: chart.YAxis{
				Range: &chart.ContinuousRange{},
			},
			Series: []chart.Series{
				chart.TimeSeries{
					Name: *stock.Ticker(),
					Style: chart.Style{
						StrokeColor: chart.GetDefaultColor(0),
					},
					XValues: *stock.XValues(),
					YValues: *stock.YValues(),
				},
			},
		}
	case D5, D10:
		time := *stock.XValues()
		ticks := append([]chart.Tick{}, chart.Tick{Value: float64(0), Label: fmt.Sprintf("%s", time[0].Format("2006-01-02"))})
		grids := []chart.GridLine{}
		for i := 1; i < len(time)-1; i++ {
			_, _, day1 := time[i].Date()
			_, _, day2 := time[i+1].Date()
			if day1 != day2 {
				ticks = append(ticks, chart.Tick{Value: float64(i + 1), Label: fmt.Sprintf("%s", time[i+1].Format("2006-01-02"))})
				grids = append(grids, chart.GridLine{Value: float64(i + 1)})
			}
		}
		ticks = append(ticks, chart.Tick{Value: float64(len(time)), Label: ""})
		graph = chart.Chart{
			XAxis: chart.XAxis{
				ValueFormatter: chart.TimeDateValueFormatter,
				Ticks:          ticks,
				TickPosition:   chart.TickPositionUnderTick,
				GridMajorStyle: chart.Style{
					StrokeColor: chart.ColorAlternateLightGray,
					StrokeWidth: 1.0,
				},
				GridLines: grids,
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
	default:
		panic("Unkown chart duration parameter during plot")
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
