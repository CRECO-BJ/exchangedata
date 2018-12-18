package database

import (
	"log"
	"testing"

	"github.com/exchangedata/common"
)

type currencyTest struct {
	in   common.Currency
	want bool
	info string
}

var tc = []currencyTest{
	{common.Currency{}, false, "empty name"},
	{common.Currency{Name: "bitcoin", Abbr: "BTC"}, true, "normal one"},
	{common.Currency{Name: "bitcoin1", Abbr: "BTC"}, true, "abbrevation the same, full name different"},
	{common.Currency{Name: "bitcoin cash", Abbr: "BCH"}, true, "another good one"},
	{common.Currency{Name: "bitcoin cash", Abbr: "bcc"}, true, "another duplicated name"},
}

func TestSaveCurrency(t *testing.T) {
	ds := NewDataStore("mysql")
	if err := ds.OpenDB(); err != nil {
		t.Fatal("open db failed")
	}
	defer ds.CloseDB()

	ds.AutoMigrate()
	for _, m := range tc {
		result := true
		err := ds.UpdateCurrency(&m.in)
		if err != nil {
			log.Println("Error ", err)
			result = false
		}
		if result != m.want {
			t.Fatalf("%s - failed", m.info)
		}
	}
}

type exchangerTest struct {
	in   common.Exchanger
	want bool
	info string
}

var te = []exchangerTest{
	{common.Exchanger{}, false, "empty name"},
	{common.Exchanger{Name: "Poloniex"}, true, "normal one"},
	{common.Exchanger{Name: "Poloniex", Info: "这是一个不错的交易所"}, true, "duplicate name, more info"},
	{common.Exchanger{Name: "binance"}, true, "another one"},
	{common.Exchanger{Name: "okcoin", Info: "Okex is located in japan"}, true, "another info"},
}

func TestSaveExchager(t *testing.T) {
	ds := NewDataStore("mysql")
	if err := ds.OpenDB(); err != nil {
		t.Fatal("open db failed")
	}
	defer ds.CloseDB()

	ds.AutoMigrate()
	for _, m := range te {
		result := true
		err := ds.UpdateExchanger(&m.in)
		if err != nil {
			log.Println("Error ", err)
			result = false
		}
		if result != m.want {
			log.Fatalf("%s - failed", m.info)
		}
	}
}

// drop table access_secrets,communication_apis,currencies,currency_exchangers,exchangers,markets,symbols,tickers,ask_pricevols,bid_pricevols,order_books,price_vols,trades;
