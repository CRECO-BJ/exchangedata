package main

import (
	"github.com/exchangedata/common"
	"github.com/exchangedata/database"
)

func main() {
	ds := database.NewDataStore("mysql")
	if err := ds.OpenDB(); err != nil {
		panic("open db failed")
	}
	defer ds.CloseDB()

	ds.GetDB().Debug().AutoMigrate(&common.Exchanger{}, &common.Market{}, &common.Currency{},
		&common.CommunicationAPI{}, &common.AccessSecret{}, &common.Symbol{})
	ds.GetDB().Debug().AutoMigrate(&common.Ticker{}, &common.Trade{}, &common.OrderBook{}, &common.PriceVol{})
}
