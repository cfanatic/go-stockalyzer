package main

import (
	"github.com/cfanatic/stockalyzer"
	"github.com/cfanatic/stockalyzer/database"
	"github.com/cfanatic/stockalyzer/finance"
)

func main() {
	stockalyzer.New()
	database.New("address", "database", "collection")
	finance.NewFinnhub()
	finance.NewYahoo()
}
