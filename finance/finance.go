package finance

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ryanuber/columnize"
	"github.com/wcharczuk/go-chart"
)

type Duration int

const (
	Intraday Duration = iota
	D5
	D10
	M1
	M3
	M6
	Y1
	Y3
	Y5
	Max
)

var Durations = []string{"Intraday", "D5", "D10", "M1", "M3", "M6", "Y1", "Y3", "Y5", "Max"}

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
	case Intraday:
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
	case M3, M6, Y1:
		time = *stock.XValues()
		tick = append([]chart.Tick{}, chart.Tick{Value: float64(0), Label: ""})
		grid = []chart.GridLine{}
		for i := 1; i < len(time)-1; i++ {
			_, month1, _ := time[i].Date()
			_, month2, _ := time[i+1].Date()
			if month1 != month2 {
				tick = append(tick, chart.Tick{Value: float64(i + 1), Label: fmt.Sprintf("%s", time[i+1].Format("Jan"))})
				grid = append(grid, chart.GridLine{Value: float64(i + 1)})
			}
		}
		tick = append(tick, chart.Tick{Value: float64(len(time)), Label: ""})
	case Y5, Max:
		time = *stock.XValues()
		tick = append([]chart.Tick{}, chart.Tick{Value: float64(0), Label: ""})
		grid = []chart.GridLine{}
		for i := 1; i < len(time)-1; i++ {
			year1, _, _ := time[i].Date()
			year2, _, _ := time[i+1].Date()
			if year1 != year2 {
				tick = append(tick, chart.Tick{Value: float64(i + 1), Label: fmt.Sprintf("%s", time[i+1].Format("2006"))})
				grid = append(grid, chart.GridLine{Value: float64(i + 1)})
			}
		}
		tick = append(tick, chart.Tick{Value: float64(len(time)), Label: ""})
	default:
		panic("Unkown duration parameter to plot stock chart")
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
				Name: stock.GetProfile().Name,
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
	path, _ := filepath.Abs("misc/plot")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, os.ModePerm)
	}
	name := stock.GetProfile().Name
	name = strings.Replace(name, " ", "_", -1)
	name = name + "_" + Durations[*stock.Duration()]
	f, _ := os.Create(fmt.Sprintf("misc/plot/%s.png", name))
	defer f.Close()
	graph.Render(chart.PNG, f)
}

func Print(stock IFinance) {
	time := *stock.XValues()
	price := *stock.YValues()
	for i := range time {
		fmt.Printf("%3d | %+v | %+v | %v\n", i, time[i].Unix(), time[i], price[i])
	}
	fmt.Println()
}

func Performance(stock IFinance) {
	var row strings.Builder
	var out []string
	var quotes [](*[]float64)
	var categories = []string{"Performance", "High", "Low"}
	var durations = []Duration{Intraday, D10, M1, M3, Y1, Y3, Y5, Max}

	getMax := func(values []float64) float64 {
		max := 0.0
		for _, value := range values {
			if value > max {
				max = value
			}
		}
		return max
	}
	getMin := func(values []float64) float64 {
		min := values[0]
		for _, value := range values {
			if value < min {
				min = value
			}
		}
		return min
	}

	for _, duration := range durations {
		stock.GetChart(duration)
		quotes = append(quotes, stock.YValues())
	}

	config := columnize.DefaultConfig()
	config.Glue = "      "

	row.WriteString(fmt.Sprintf("%s | Intraday | D10 | M1 | M3 | Y1 | Y3 | Y5 | Max", stock.GetProfile().Name))
	out = append(out, row.String())
	out = append(out, "")

	for _, category := range categories {
		row.Reset()
		row.WriteString(fmt.Sprintf("%s |", category))
		switch category {
		case "Performance":
			for i, tmp := range quotes {
				if i > 0 {
					quote := *tmp
					row.WriteString(fmt.Sprintf("%.2f%% |", ((quote[len(quote)-1]-quote[0])/quote[0])*100))
				} else {
					quote := stock.GetQuote()
					row.WriteString(fmt.Sprintf("%.2f%% |", ((quote.Current-quote.PrevClose)/quote.PrevClose)*100))
				}
			}
		case "High":
			for _, tmp := range quotes {
				quote := *tmp
				row.WriteString(fmt.Sprintf("%.2f |", getMax(quote)))
			}
		case "Low":
			for _, tmp := range quotes {
				quote := *tmp
				row.WriteString(fmt.Sprintf("%.2f |", getMin(quote)))
			}
		}
		out = append(out, row.String())
		out = append(out, "")
	}

	result := columnize.Format(out, config)
	fmt.Println(result)
}
