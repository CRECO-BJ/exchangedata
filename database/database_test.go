package database

import (
	"log"
	"testing"
	"time"

	"github.com/exchangedata/common"
)

var testCurrencies = []common.Currency{
	common.Currency{Name: "bitcoin", Abbr: "BTC"},
	common.Currency{Name: "Litecoin", Abbr: "LTC"},
	common.Currency{Name: "bitcoin cash", Abbr: "BCH"},
	common.Currency{Name: "Dogecoin", Abbr: "DOGE"},
}

var testExchangers = []common.Exchanger{
	common.Exchanger{Name: "Poloniex"},
	common.Exchanger{Name: "bittrex"},
	common.Exchanger{Name: "okcoin", Info: "Okex is located in japan"},
}

var testMarkets = []*common.Market{
	&common.Market{Name: "LTC-BTC",
		Symbol: &common.Symbol{Base: &common.Currency{Name: "Litecoin", Abbr: "LTC"},
			Quote: &common.Currency{Name: "bitcoin", Abbr: "BTC"}},
		Active:     true,
		Info:       "common market test",
		Precision:  5,
		Limitation: common.Limitation{Min: 0.01, Max: 10000},
		MinStep:    0.000001,
		Exchanger:  &testExchangers[1]},
	&common.Market{Name: "DOGE-BTC",
		Symbol: &common.Symbol{
			Base:  &testCurrencies[3],
			Quote: &testCurrencies[0]},
		Active:     true,
		Info:       "extiguish test",
		Precision:  8,
		Limitation: common.Limitation{Min: 100, Max: 10000000},
		MinStep:    0.00000000001,
		Exchanger:  &testExchangers[1]},
	&common.Market{Name: "DOGE-BTC",
		Symbol: &common.Symbol{
			Base:  &testCurrencies[3],
			Quote: &testCurrencies[0]},
		Active:     true,
		Info:       "extiguish test",
		Precision:  8,
		Limitation: common.Limitation{Min: 100, Max: 10000000},
		MinStep:    0.00000000001,
		Exchanger:  &testExchangers[0]},
}

var testTickers = []*common.Ticker{
	&common.Ticker{
		Time:   time.Date(2018, 11, 22, 3, 4, 5, 6, time.Local),
		Market: testMarkets[1], Last: 382.98901522, Ask: 381.99755898, Bid: 379.41296309, High: 412.25844455,
		Percentage: -0.04312950, Low: 364.56122072, BaseVolume: 14969820.94951828, QuoteVolume: 38859.58435407,
	}, // poloniex model
	&common.Ticker{
		Time:   time.Date(2018, 11, 22, 3, 4, 5, 6, time.Local),
		Market: testMarkets[0], Last: 3.35579531, Bid: 2.05670368, Ask: 3.35579531,
	}, // bittrex model
	&common.Ticker{
		Time:   time.Date(2018, 11, 22, 3, 4, 5, 6, time.Local),
		Market: testMarkets[2], High: 0.00846390, BaseVolume: 1135176.4290665, Last: 0.00809068, Low: 0.00801497, Bid: 0.00808481, Ask: 0.00809001,
	}, //okex model
}

var testTrades = []*common.Trade{
	&common.Trade{
		Time:    time.Now(),
		OrderID: "4688134",
		Type:    "fill",
		Side:    "buy",
		Price:   3304.51,
		Amount:  23.1000008,
		Market:  testMarkets[1],
	},
	&common.Trade{
		Time:    time.Date(2018, 11, 12, 23, 54, 33, 00, time.UTC),
		OrderID: "1832441212123",
		Side:    "buy",
		Price:   33.3,
		Amount:  2008,
		Market:  testMarkets[2],
	},
}

var testOrderBooks = []*common.OrderBook{
	&common.OrderBook{
		Market: testMarkets[0],
		Time:   time.Now(),
		Asks: []*common.PriceVol{
			&common.PriceVol{Price: 0.00001853, Volume: 2537.5637},
			&common.PriceVol{Price: 0.00001854, Volume: 1567238.172367},
		},
		Bids: []*common.PriceVol{
			&common.PriceVol{Price: 0.00001841, Volume: 3645.3647},
			&common.PriceVol{Price: 0.00001840, Volume: 1637.3647},
		},
	},
}

type currencyTest struct {
	in   common.Currency
	want bool
	info string
}

var tc = []currencyTest{
	{common.Currency{}, true, "empty name"},
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
		if ds.UpdateCurrency(&m.in).Error != nil {
			log.Printf("Error %s", ds.GetDB().Error)
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

var (
	Currencs  map[string]*common.Currency
	Exchanges map[string]*common.Exchanger
	Markes    map[uint]*common.Market
)

func init() {
	Currencs = make(map[string]*common.Currency)
	Exchanges = make(map[string]*common.Exchanger)
	Markes = make(map[uint]*common.Market)
	for k, m := range testMarkets {
		if Currencs[m.Symbol.Base.Name] == nil || Currencs[m.Symbol.Base.Name].Name != m.Symbol.Base.Name {
			Currencs[m.Symbol.Base.Name] = m.Symbol.Base
		} else {
			m.Symbol.Base = Currencs[m.Symbol.Base.Name]
		}
		if Currencs[m.Symbol.Quote.Name] == nil || Currencs[m.Symbol.Quote.Name].Name != m.Symbol.Quote.Name {
			Currencs[m.Symbol.Quote.Name] = m.Symbol.Quote
		} else {
			m.Symbol.Quote = Currencs[m.Symbol.Quote.Name]
		}
		if Exchanges[m.Exchanger.Name] == nil || Exchanges[m.Exchanger.Name].Name != m.Exchanger.Name {
			Exchanges[m.Exchanger.Name] = m.Exchanger
		} else {
			m.Exchanger = Exchanges[m.Exchanger.Name]
		}
		Markes[uint(k)] = m
	}
}

func TestUpdateTicker(t *testing.T) {
	var err error
	ds := NewDataStore("mysql")
	if err = ds.OpenDB(); err != nil {
		panic("open db failed")
	}
	defer ds.CloseDB()

	ds.AutoMigrate()
	for _, c := range Currencs {
		if ds.UpdateCurrency(c).Error != nil {
			log.Println("error save currency", ds.GetDB().Error)
		}
	}
	for _, c := range Exchanges {
		if ds.UpdateExchanger(c).Error != nil {
			log.Println("error save exchanger", ds.GetDB().Error)
		}
	}
	for _, c := range Markes {
		if ds.UpdateMarket(c).Error != nil {
			log.Println("error save market", ds.GetDB().Error)
		}
	}
	for _, c := range testTickers {
		if ds.UpdateTicker(c).Error != nil {
			log.Println("error save ticker", ds.GetDB().Error)
		}
	}
	for _, c := range testTrades {
		if ds.UpdateTrade(c).Error != nil {
			log.Println("error save trade", ds.GetDB().Error)
		}
	}
	for _, c := range testOrderBooks {
		if ds.UpdateOrderBook(c).Error != nil {
			log.Println("error save orderbook", ds.GetDB().Error)
		}
	}
}

// drop table access_secrets,communication_apis,currencies,currency_exchangers,exchangers,markets,symbols,tickers,ask_pricevols,bid_pricevols,order_books,price_vols,trades;
