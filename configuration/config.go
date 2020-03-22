package configuration

import (
	"flag"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type keys int

const (
	PATH             = "cmd/stockalyzer/config.toml"
	PATH_R           = "misc/config.toml"
	PATH_D           = "../../misc/config.toml"
	MARKETHOURS keys = iota
	CHARTSIZE
	TOKEN
)

var (
	mode = flag.String("mode", PATH, "define execution mode")
)

type config struct {
	General general
	Finnhub finnhub
}

type general struct {
	MarketHours []int
	ChartSize   []int
}

type finnhub struct {
	Token string
}

func Get(key keys) interface{} {
	var (
		conf config
		path string
	)
	if flag.Parse(); *mode == "release" {
		path, _ = filepath.Abs(PATH_R)
	} else if *mode == "debug" {
		path, _ = filepath.Abs(PATH_D)
	} else {
		path, _ = filepath.Abs(*mode)
	}
	if _, err := toml.DecodeFile(path, &conf); err != nil {
		panic(err)
	}
	switch key {
	case MARKETHOURS:
		return conf.General.MarketHours
	case CHARTSIZE:
		return conf.General.ChartSize
	case TOKEN:
		return conf.Finnhub.Token
	default:
		return nil
	}
}

func MarketHours() (int, int) {
	tmp := Get(MARKETHOURS).([]int)
	return tmp[0], tmp[1]
}

func ChartSize() (int, int) {
	tmp := Get(CHARTSIZE).([]int)
	return tmp[0], tmp[1]
}

func Token() string {
	return Get(TOKEN).(string)
}
