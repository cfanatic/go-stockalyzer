package configuration

import (
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type keys int

const (
	PATH               = "cmd/stockalyzer/config.toml"
	FINNHUB_TOKEN keys = iota
)

type config struct {
	Finnhub finnhub
}

type finnhub struct {
	Token string
}

func Get(key keys) interface{} {
	var conf config
	path, _ := filepath.Abs(PATH)
	if _, err := toml.DecodeFile(path, &conf); err != nil {
		panic(err)
	}
	switch key {
	case FINNHUB_TOKEN:
		return conf.Finnhub.Token
	default:
		return nil
	}
}
