package finance

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
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

var Durations = [...]string{
	"Intraday",
	"D5",
	"D10",
	"M1",
	"M3",
	"M6",
	"Y1",
	"Y3",
	"Y5",
	"Max",
}

type IFinance interface {
	GetQuote() *Quote
	GetCandle(from, to string) *Candle
	GetChart(duration Duration) *Candle

	Profile() *Profile
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
	var (
		times []time.Time
		tick  []chart.Tick
		grid  []chart.GridLine
	)
	minIndex, minValue := minValue(stock.YValues())
	maxIndex, maxValue := maxValue(stock.YValues())
	switch *stock.Duration() {
	case Intraday:
		times = *stock.XValues()
		tick = append([]chart.Tick{}, chart.Tick{Value: float64(0), Label: fmt.Sprintf("%s", times[0].Format("01-02 3PM"))})
		grid = []chart.GridLine{}
		for i := 1; i < len(times)-1; i++ {
			hour1, _, _ := times[i].Clock()
			hour2, _, _ := times[i+1].Clock()
			if hour1 != hour2 {
				tick = append(tick, chart.Tick{Value: float64(i + 1), Label: fmt.Sprintf("%s", times[i+1].Format("01-02 3PM"))})
				grid = append(grid, chart.GridLine{Value: float64(i + 1)})
			}
		}
		tick = append(tick, chart.Tick{Value: float64(len(times)), Label: ""})
	case D5, D10:
		times = *stock.XValues()
		tick = append([]chart.Tick{}, chart.Tick{Value: float64(0), Label: fmt.Sprintf("%s", times[0].Format("2006-01-02"))})
		grid = []chart.GridLine{}
		for i := 1; i < len(times)-1; i++ {
			_, _, day1 := times[i].Date()
			_, _, day2 := times[i+1].Date()
			if day1 != day2 {
				tick = append(tick, chart.Tick{Value: float64(i + 1), Label: fmt.Sprintf("%s", times[i+1].Format("2006-01-02"))})
				grid = append(grid, chart.GridLine{Value: float64(i + 1)})
			}
		}
		tick = append(tick, chart.Tick{Value: float64(len(times)), Label: ""})
	case M3, M6, Y1:
		times = *stock.XValues()
		tick = append([]chart.Tick{}, chart.Tick{Value: float64(0), Label: ""})
		grid = []chart.GridLine{}
		for i := 1; i < len(times)-1; i++ {
			_, month1, _ := times[i].Date()
			_, month2, _ := times[i+1].Date()
			if month1 != month2 {
				tick = append(tick, chart.Tick{Value: float64(i + 1), Label: fmt.Sprintf("%s", times[i+1].Format("Jan"))})
				grid = append(grid, chart.GridLine{Value: float64(i + 1)})
			}
		}
		tick = append(tick, chart.Tick{Value: float64(len(times)), Label: ""})
	case Y3, Y5, Max:
		times = *stock.XValues()
		tick = append([]chart.Tick{}, chart.Tick{Value: float64(0), Label: ""})
		grid = []chart.GridLine{}
		for i := 1; i < len(times)-1; i++ {
			year1, _, _ := times[i].Date()
			year2, _, _ := times[i+1].Date()
			if year1 != year2 {
				tick = append(tick, chart.Tick{Value: float64(i + 1), Label: fmt.Sprintf("%s", times[i+1].Format("2006"))})
				grid = append(grid, chart.GridLine{Value: float64(i + 1)})
			}
		}
		tick = append(tick, chart.Tick{Value: float64(len(times)), Label: ""})
	default:
		panic("Unsupported duration parameter to plot stock chart")
	}
	graph := chart.Chart{
		Width:  1280,
		Height: 720,
		Background: chart.Style{
			Padding: chart.Box{
				Top: 75,
			},
		},
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
				Name: stock.Profile().Name + " - " + fmt.Sprintf("%s", Durations[*stock.Duration()]),
				Style: chart.Style{
					StrokeColor: chart.GetDefaultColor(0),
				},
				XValues: func() []float64 {
					xvalues := make([]float64, len(times))
					for i := 0; i < len(times); i++ {
						xvalues[i] = float64(i)
					}
					return xvalues
				}(),
				YValues: *stock.YValues(),
			},
			chart.AnnotationSeries{
				Annotations: []chart.Value2{
					{XValue: float64(minIndex), YValue: minValue, Label: strconv.FormatInt(int64(minValue), 10)},
					{XValue: float64(maxIndex), YValue: maxValue, Label: strconv.FormatInt(int64(maxValue), 10)},
				},
				Style: chart.Style{
					StrokeColor: chart.GetDefaultColor(0),
				},
			},
		},
	}
	graph.Elements = []chart.Renderable{
		chart.LegendThin(&graph),
	}
	path, _ := filepath.Abs("misc/plot")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, os.ModePerm)
	}
	name := stock.Profile().Name
	name = strings.Replace(name, " ", "_", -1)
	name = name + "_" + Durations[*stock.Duration()]
	f, _ := os.Create(fmt.Sprintf("misc/plot/%s.png", name))
	defer f.Close()
	graph.Render(chart.PNG, f)
}

func Performance(stock IFinance) {
	var (
		row        strings.Builder
		out        []string
		candles    [](*[]float64)
		categories = []string{"Performance", "High", "Low"}
		durations  = []Duration{Intraday, D10, M1, M3, Y1, Y3, Y5, Max}
	)
	stopwatch := func() func() {
		start := time.Now()
		return func() {
			fmt.Printf("Time: %.2f seconds\n", time.Since(start).Seconds())
		}
	}
	defer stopwatch()()
	row.WriteString(fmt.Sprintf("%s | Intraday | D10 | M1 | M3 | Y1 | Y3 | Y5 | Max", stock.Profile().Name))
	out = append(out, row.String())
	out = append(out, "")
	for _, duration := range durations {
		stock.GetChart(duration)
		candles = append(candles, stock.YValues())
	}
	for _, category := range categories {
		row.Reset()
		row.WriteString(fmt.Sprintf("%s |", category))
		switch category {
		case "Performance":
			quote := stock.GetQuote()
			for i, tmp := range candles {
				if i > 0 {
					candle := *tmp
					if openMarket() == true {
						row.WriteString(fmt.Sprintf("%.2f%% |", ((quote.PrevClose-candle[0])/candle[0])*100))
					} else {
						row.WriteString(fmt.Sprintf("%.2f%% |", ((quote.Current-candle[0])/candle[0])*100))
					}
				} else {
					row.WriteString(fmt.Sprintf("%.2f%% |", ((quote.Current-quote.PrevClose)/quote.PrevClose)*100))
				}
			}
		case "High":
			for _, candle := range candles {
				_, maxValue := maxValue(candle)
				row.WriteString(fmt.Sprintf("%.2f |", maxValue))
			}
		case "Low":
			for _, candle := range candles {
				_, minValue := minValue(candle)
				row.WriteString(fmt.Sprintf("%.2f |", minValue))
			}
		}
		out = append(out, row.String())
		out = append(out, "")
	}
	config := columnize.DefaultConfig()
	config.Glue = "      "
	result := columnize.Format(out, config)
	fmt.Println(result)
}

func Print(stock IFinance) {
	time := *stock.XValues()
	price := *stock.YValues()
	for i := range time {
		fmt.Printf("%3d | %+v | %+v | %v\n", i, time[i].Unix(), time[i], price[i])
	}
	fmt.Println()
}

func openMarket() bool {
	var open bool
	current := time.Now()
	day := current.Weekday()
	hour := current.Hour()
	min := current.Minute()
	if day == time.Saturday || day == time.Sunday {
		open = false
	} else {
		if hour >= 17 && min >= 30 {
			open = false
		} else {
			open = true
		}
	}
	return open
}

func maxValue(values *[]float64) (int, float64) {
	idx, max := 0, 0.0
	for i, value := range *values {
		if value > max {
			idx = i
			max = value
		}
	}
	return idx, max
}

func minValue(values *[]float64) (int, float64) {
	idx, min := 0, (*values)[0]
	for i, value := range *values {
		if value < min {
			idx = i
			min = value
		}
	}
	return idx, min
}
