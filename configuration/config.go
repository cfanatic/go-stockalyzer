package configuration

import (
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type list int

const (
	CONFIG             = "misc/config.toml"
	FINNHUB_TOKEN list = iota
)

type config struct {
	Finnhub finnhub
}

type finnhub struct {
	Token string
}

func Get(key list) interface{} {
	var conf config
	path, _ := filepath.Abs(CONFIG)
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
