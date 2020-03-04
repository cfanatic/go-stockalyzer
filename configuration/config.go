package configuration

import (
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type keys int

const (
	PATH_R             = "misc/config.toml"
	PATH_D             = "../../misc/config.toml"
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
	path, _ := filepath.Abs(PATH_R)
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
